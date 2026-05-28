/*
Copyright (c) 2023-2026 Microbus LLC and various contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package showcase

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/microbus-io/bespa/form"
	"github.com/microbus-io/bespa/website/resources"
	"github.com/microbus-io/bespa/website/shared"
	"github.com/microbus-io/bespa/website/storage"
	"github.com/microbus-io/errors"
)

func HandleDirEdit(w http.ResponseWriter, r *http.Request) {
	directory := shared.SessionOf(w, r).DirectoryDB()

	state := wf.StateOf(r)
	person := &storage.Person{}
	pageTitle := "Register new person"
	actionBtn := wf.ButtonFilled("").Add("Register").WithDisabled(!directory.CanInsert(r))

	p := directory.RandomPerson()
	href := fmt.Sprintf("?fname=%s&lname=%s&email=%s&phone=%s&birthday=%s&&city=%s&state=%s&zip=%s",
		p.FirstName, p.LastName, p.Email, p.Phone, p.Birthday.Format("2006-01-02"), p.City, p.State, p.Zip,
	)
	quickFill := wf.Link(href).Add(wf.Icon("create"), "Quick fill")

	if state.Has("id") {
		var ok bool
		person, ok = directory.Lookup(state.Get("id"))
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Not found"))
			return
		}
		pageTitle = person.FullName()
		actionBtn = wf.ButtonFilled("").Add("Update")
		quickFill.HideIf(true)
	}

	dropdown, _ := createStatesDropdown(person.State)

	form := wf.Form().Add(
		wf.Field().
			AddLeft("First name").
			AddRight(
				wf.InputText("fname", person.FirstName).
					WithRequired(true).
					WithLength(0, 20).
					RedrawIfChanged(r, "fname"),
			),
		wf.Field().
			AddLeft("Last name").
			AddRight(
				wf.InputText("lname", person.LastName).
					WithRequired(true).
					WithLength(0, 20).
					RedrawIfChanged(r, "lname"),
			),
		wf.Field().
			AddLeft("Email").
			AddRight(
				wf.InputEmail("email", person.Email).
					WithRequired(true).
					WithLength(0, 40).
					WithPredicate(func(value string) (bool, string) {
						// Disallow duplicate emails for different persons
						p, ok := directory.LookupByEmail(value)
						if ok && p.ID != person.ID {
							return false, "Email already in use"
						}
						return true, ""
					}).
					RedrawIfChanged(r, "email"),
			),
		wf.Field().
			AddLeft("Phone").
			AddRight(
				wf.InputPhone("phone", person.Phone).
					WithRequired(false).
					WithPlaceholder("555 555 5555").
					WithPattern(`^[0-9]{3} [0-9]{3} [0-9]{4}$`).
					WithLength(12, 12).
					RedrawIfChanged(r, "phone"),
			),
		wf.Field().
			AddLeft("Birthday", wf.HelpBubble("Must be 18+ years old")).
			AddRight(
				wf.InputDate("birthday", person.Birthday).
					WithRequired(true).
					WithMin(time.Now().AddDate(-120, 0, 0)).
					WithMax(time.Now().AddDate(-18, 0, 0)).
					RedrawIfChanged(r, "birthday"),
			),
		wf.Field().
			AddLeft("Address line 1").
			AddRight(
				wf.InputText("address1", person.Address1).
					WithLength(0, 40),
			),
		wf.Field().
			AddLeft("Address line 2").
			AddRight(
				wf.InputText("address2", person.Address2).
					WithLength(0, 40),
			),
		wf.Field().
			AddLeft("City, state, zip").
			AddRight(
				wf.Splitter(5, 3, 2).WithWrap(false).
					AddLeft(
						wf.InputText("city", person.City).
							WithRequired(true).
							WithLength(0, 20).
							RedrawIfChanged(r, "city"),
					).
					AddToCol(1,
						dropdown.RedrawIfChanged(r, "state"),
					).
					AddRight(
						wf.InputText("zip", person.Zip).
							WithLength(5, 5).
							WithPattern("^[0-9]{5}$").
							WithPlaceholder("00000").
							RedrawIfChanged(r, "zip"),
					),
			),

		wf.ButtonText("").Add("Cancel").WithHrefBack(),
		actionBtn,
	)

	snack := wf.Snackbar()
	if form.ReadyToCommit(r) && (person.ID != "" || directory.CanInsert(r)) {
		values := form.Values(r)
		person.FirstName = values.Get("fname")
		person.LastName = values.Get("lname")
		person.Email = values.Get("email")
		person.Phone = values.Get("phone")
		person.Birthday, _ = time.Parse("2006-01-02", values["birthday"])
		person.Address1 = values.Get("address1")
		person.Address2 = values.Get("address2")
		person.City = values.Get("city")
		person.State = values.Get("state")
		person.Zip = values.Get("zip")

		var err error
		if person.ID == "" {
			// Register new person
			person.ID, err = directory.Insert(person)
			if err == nil {
				backArg := ""
				if state.Has("_back") {
					backArg = "&_back=" + state.Get("_back")
				}
				snack = wf.Snackbar().Add(
					fmt.Sprintf("%s %s registered", person.FirstName, person.LastName),
					wf.SpacerBreak(),
					wf.TextAlignRight(wf.Link("dir-edit?id="+person.ID+backArg).Add("View")),
				)
				form.Reset(r)
			} else {
				snack = wf.Snackbar().Add(
					fmt.Sprintf("Failed to register %s %s: %s", person.FirstName, person.LastName, err.Error()),
				)
			}
		} else {
			// Update existing record
			err = directory.Update(person)
			if err == nil {
				snack = wf.Snackbar().Add(
					fmt.Sprintf("%s %s updated", person.FirstName, person.LastName),
					wf.SpacerBreak(),
					wf.TextAlignRight(wf.Link("dir-edit?id="+person.ID).Add("View")),
				)
				if wf.RedirectBack(w, r) {
					return
				}
			} else {
				snack = wf.Snackbar().Add(
					fmt.Sprintf("Failed to update %s %s: %s", person.FirstName, person.LastName, err.Error()),
				)
			}
		}

		actionBtn.WithDisabled(!directory.CanInsert(r))
	}

	page := wf.Page().Add(
		wf.AppBar(pageTitle).
			AddRight(quickFill.RedrawIf(wf.StateOf(r).HasChanges())),
		`The same form is used to add new data or edit existing data.
		It performs client-side validations to enforce required fields, minimum age and phone and zip code formats.
		Further validations are performed server-side to ensure the uniqueness of the email address
		as well as to enforce a cap of `, storage.MaxDirectoryRecords, ` records.`,
		wf.SpacerParagraph(),
		form,
		snack.RedrawIf(true),
	)

	shared.Render(w, r, page)
}

// HandleDirList shows a list of persons in a table.
func HandleDirList(w http.ResponseWriter, r *http.Request) {
	directory := shared.SessionOf(w, r).DirectoryDB()

	// Perform the delete action
	deleteSnackbar := wf.Snackbar()
	state := wf.StateOf(r)
	idToDelete := state.Get("delete")
	if idToDelete != "" {
		p, ok := directory.Lookup(idToDelete)
		if ok {
			directory.Delete(idToDelete)
			deleteSnackbar.Add(p.FullName() + " deleted").RedrawIf(true)
		}
		state.Del("name")
		state.Set("delete", "")
	}

	// Prepare the table
	tbl := wf.Table().
		WithDefaultSortOrder(r, "lname").
		WithDefaultPageRows(r, 10).
		RedrawIfChanged(r, "modal", "delete").
		Add(
			wf.Col("wx", 12, "left").Add(wf.Sorter("fname", "First name")),
			wf.Col("wx", 12, "left").Add(wf.Sorter("lname", "Last name")),
			wf.Col("n", 24, "left").Add(
				wf.Sorter("fname", "First"),
				" and ",
				wf.Sorter("lname", "last"),
				" name",
			),
			wf.Col("nwx", 3, "right").Add(wf.Sorter("age", "Age")),
			wf.Col("x", 32, "left").Add("email"),
			wf.Col("nwx", 3, "left").Add(wf.Sorter("state", "State")),
			wf.Col("nwx", 1, "right"),
		)

	// Prepare the data
	persons := directory.List()

	// Filter
	if strings.TrimSpace(tbl.Query(r)) != "" {
		q := strings.TrimSpace(strings.ToLower(tbl.Query(r)))
		filtered := []*storage.Person{}
		for _, p := range persons {
			if strings.Contains(strings.ToLower(p.LastName), q) ||
				strings.Contains(strings.ToLower(p.FirstName), q) ||
				strings.Contains(strings.ToLower(p.Email), q) ||
				strings.Contains(strings.ToLower(p.State), q) {
				filtered = append(filtered, p)
			}
		}
		persons = filtered
	}

	// Sort
	switch tbl.SortOrder(r) {
	case "lname", "":
		sort.Slice(persons, func(i, j int) bool {
			return strings.ToUpper(persons[i].LastName) < strings.ToUpper(persons[j].LastName)
		})
	case "-lname":
		sort.Slice(persons, func(j, i int) bool {
			return strings.ToUpper(persons[i].LastName) < strings.ToUpper(persons[j].LastName)
		})
	case "fname":
		sort.Slice(persons, func(i, j int) bool {
			return strings.ToUpper(persons[i].FirstName) < strings.ToUpper(persons[j].FirstName)
		})
	case "-fname":
		sort.Slice(persons, func(j, i int) bool {
			return strings.ToUpper(persons[i].FirstName) < strings.ToUpper(persons[j].FirstName)
		})
	case "state":
		sort.Slice(persons, func(i, j int) bool {
			return strings.ToUpper(persons[i].State) < strings.ToUpper(persons[j].State)
		})
	case "-state":
		sort.Slice(persons, func(j, i int) bool {
			return strings.ToUpper(persons[i].State) < strings.ToUpper(persons[j].State)
		})
	case "age":
		sort.Slice(persons, func(i, j int) bool {
			return persons[i].Birthday.After(persons[j].Birthday)
		})
	case "-age":
		sort.Slice(persons, func(j, i int) bool {
			return persons[i].Birthday.After(persons[j].Birthday)
		})
	}

	// Populate rows
	tbl.WithTotalRows(r, len(persons))
	rowFrom, rowTo := tbl.DisplayRange(r)
	for _, p := range persons[rowFrom:rowTo] {
		tbl.Add(wf.Row().Add(
			wf.QuickSearchUnderliner(p.FirstName),
			wf.QuickSearchUnderliner(p.LastName),
			wf.QuickSearchUnderliner(p.FirstName+" "+p.LastName),
			int(time.Since(p.Birthday).Hours()/24/365.25),
			wf.QuickSearchUnderliner(p.Email),
			wf.QuickSearchUnderliner(p.State),
			wf.MenuEllipsisH().Add(
				wf.Link("?modal="+url.QueryEscape("dir-edit?id="+p.ID+"&_back=^?modal=")).Add("Edit"),
				wf.Link("?alertdelete="+p.ID+"&name="+url.QueryEscape(p.FullName())).Add("Delete"),
			),
		).WithAction("?modal=" + url.QueryEscape("dir-edit?id="+p.ID+"&_back=^?modal=")))
	}

	page := wf.Page().Add(
		wf.AppBar("Directory"),
		`CRUD is a common pattern in most applications.
		The table shows the current data.
		The plus button allows adding new data.
		The kebab menu on each row includes the options to edit or delete existing data.`,

		wf.Toolbar().
			AddLeft(
				wf.ButtonTonal("").
					Add(wf.Icon("add")).
					WithHref("?modal=dir-edit?_back=^?modal=").
					WithDisabled(!directory.CanInsert(r)).
					RedrawIfChanged(r, "modal", "delete"),
				wf.QuickSearch(),
			).
			AddRight(
				wf.Paginator().
					RedrawIfChanged(r, "modal", "delete"),
			),
		tbl,
		wf.Toolbar().AddLeft(
			wf.PageSizer().
				RedrawIfChanged(r, "modal", "delete"),
		),

		wf.Modal("modal").Add(wf.EmbedHandler(mux.ServeHTTP, r, "GET", wf.StateOf(r).Get("modal"), nil)),
		wf.AlertError("alertdelete", "delete", "Delete?").Add(
			wf.Printf("Are you sure you want to delete {0}?", state.Get("name")),
			wf.ButtonText("").Add("Cancel"),
			wf.ButtonText("delete").Add("Delete"),
		),
		deleteSnackbar,

		wf.Debugger(),
	)

	shared.Render(w, r, page)
}

func createStatesDropdown(initialValue string) (*form.DropdownWidget, error) {
	dataFile, err := resources.Bundle.ReadFile("states.csv")
	if err != nil {
		return nil, errors.Trace(err)
	}
	records, err := csv.NewReader(bytes.NewReader(dataFile)).ReadAll()
	if err != nil {
		return nil, errors.Trace(err)
	}
	records = records[1:] // Discard header line

	dropdown := wf.Dropdown("state", initialValue).
		WithStretch(true).
		WithRequired(true).
		AddOption("", "")
	for _, rec := range records {
		dropdown.AddOption(rec[1], rec[0])
	}
	return dropdown, nil
}
