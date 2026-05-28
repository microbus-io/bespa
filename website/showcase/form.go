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
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/microbus-io/bespa/website/resources"
	"github.com/microbus-io/bespa/website/shared"
)

// handleFormInput demonstrates the use of a form with input fields.
func HandleReceiver(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"key":"123456"}`))
}

// handleFormInput demonstrates the use of a form with input fields.
func HandleFormInput(w http.ResponseWriter, r *http.Request) {
	dataFile, _ := resources.Bundle.ReadFile("states.csv")
	records, _ := csv.NewReader(bytes.NewReader(dataFile)).ReadAll()
	records = records[1:] // Discard header line
	richDropdown := wf.RichDropdown("states", "CA")
	for i := 0; i < len(records); i++ {
		abbrev := strings.ToLower(records[i][1])
		richDropdown.AddOption(records[i][1], wf.Avatar(records[i][0], "/images/stateflags/"+abbrev+".webp").WithNameLabel(true))
	}

	form := wf.Form().Add(
		wf.Field().
			AddLeft("Text").
			AddRight(
				wf.InputText("text", "").
					WithPlaceholder("Enter some text").
					WithAutoFocus(true),
			),
		wf.Field().
			AddLeft("Disabled").
			AddRight(
				wf.InputText("disabled", "Uneditable").WithDisabled(true),
			),
		wf.Field().
			AddLeft("Integer").
			AddRight(
				wf.InputInteger("int", "").WithMin(0),
			),
		wf.Field().
			AddLeft("Range", wf.HelpBubble("Zero to 100% at 10% increments")).
			AddRight(
				wf.InputRange("range", 50).WithMin(0).WithMax(100).WithStep(10),
			),
		wf.Field().
			AddLeft("Disabled range").
			AddRight(
				wf.InputRange("disabledrange", 90).WithMin(0).WithMax(100).WithStep(10).WithDisabled(true),
			),
		wf.Field().
			AddLeft("Email").
			AddRight(
				wf.InputEmail("email", "").WithPlaceholder("Email"),
			),
		wf.Field().
			AddLeft("URL").
			AddRight(
				wf.InputURL("url", ""),
			),
		wf.Field().
			AddLeft("Password").
			AddRight(
				wf.InputPassword("pw", "").WithLength(8, -1),
			),
		wf.Field().
			AddLeft("Numeric OTP").
			AddRight(
				wf.InputOneTimePassword("otp", "").WithPattern("[0-9]*").WithLength(6, 6).WithWidth(6),
			),
		wf.Field().
			AddLeft("Color").
			AddRight(
				wf.InputColor("color", "#426915"),
			),
		wf.Field().
			AddLeft("Multiple lines").
			AddRight(
				wf.InputText("multi", "Multiple\nlines of text\nare allowed").
					WithRows(3),
			),
		wf.Field().
			AddLeft("Disabled multiple lines").
			AddRight(
				wf.InputText("disabledmulti", "Uneditable\nmultiple\nlines").
					WithRows(3).
					WithDisabled(true),
			),
		wf.Field().
			AddLeft("Month").
			AddRight(
				wf.InputMonth("month", time.Now()),
			),
		wf.Field().
			AddLeft("Date and time").
			AddRight(
				wf.Gallery(
					wf.InputDate("date", time.Now().UTC().Truncate(time.Hour*24)),
					wf.InputTime("time", time.Now().Round(time.Minute*15)).WithStep(15*60),
				),
			),
		wf.Field().
			AddLeft("Time zone").
			AddRight(
				wf.InputTimeZone("tz", "US/Pacific"),
			),
		wf.Field().
			AddLeft("Host and port").
			AddRight(
				wf.Splitter(2, 0, 1).WithWrap(false).
					AddLeft(wf.InputText("host", "").WithPlaceholder("www.example.com")).
					AddRight(wf.InputText("port", "").WithPlaceholder("443")),
			),

		wf.Field().
			AddLeft("Dropdown selector").
			AddRight(
				wf.Dropdown("dropdown", "b").
					AddOption("r", "Red").
					AddOption("g", "Green").
					AddOption("b", "Blue"),
			),
		wf.Field().
			AddLeft("Disabled dropdown selector").
			AddRight(
				wf.Dropdown("disableddropdown", "b").WithDisabled(true).
					AddOption("r", "Red").
					AddOption("g", "Green").
					AddOption("b", "Blue"),
			),
		wf.Field().
			AddLeft("Rich drop down").
			AddRight(richDropdown),
		wf.Field().
			AddLeft("Disabled rich drop down").
			AddRight(
				wf.RichDropdown("disabledrich", "home").WithDisabled(true).
					AddOption("home",
						wf.TextStyle().WithColorPrimary().Add(wf.Icon("home").WithFill(true)),
						" is where the ",
						wf.TextStyle().WithColorError().Add(wf.Icon("favorite").WithFill(true)),
						" is",
					),
			),
		wf.Field().
			AddLeft("File").
			AddRight(
				wf.InputFile("file", "", "receiver"),
			),
		wf.Field().
			AddLeft("Checkboxes").
			AddRight(
				wf.Checkbox("cbon", true).Add("Checked"),
				wf.SpacerBreak(),
				wf.Checkbox("cboff", false).Add("Unchecked"),
				wf.SpacerBreak(),
				wf.Checkbox("disabledcbon", true).WithDisabled(true).Add("Can't uncheck"),
				wf.SpacerBreak(),
				wf.Checkbox("disabledcboff", false).WithDisabled(true).Add("Can't check"),
			),
		wf.Field().
			AddLeft("Radio buttons").
			AddRight(
				wf.Radio("country", "US").WithRequired(true).
					AddOption("US", "USA").
					AddOption("CA", "Canada").
					AddOption("MX", "Mexico"),
			),
		wf.Field().
			AddLeft("Disabled radio buttons").
			AddRight(
				wf.Radio("season", "1").WithHorizontal().WithDisabled(true).
					AddOption("1", "Winter").
					AddOption("2", "Spring").
					AddOption("3", "Summer").
					AddOption("4", "Fall"),
			),
		wf.Field().
			AddLeft("Filter chips").
			AddRight(wf.Gallery(
				wf.FilterChip("west", "West", true),
				wf.FilterChip("east", "East", false),
				wf.FilterChip("north", "North", true).WithDisabled(true),
				wf.FilterChip("south", "South", false).WithDisabled(true),
			)),
		wf.Field().
			AddLeft("Input chips").
			AddRight(
				wf.InputChips("chips", "/showcase/query-states").
					AddChip("CA", "California").
					AddChip("DE", "Delaware").
					WithMaxItems(4),
			),
		wf.Field().
			AddLeft("Switch toggles").
			AddRight(
				wf.Toggle("toggleoff", false),
				wf.Toggle("toggleon", true),
				wf.Toggle("disabledtoggleoff", false).WithDisabled(true),
				wf.Toggle("disabledtoggleon", true).WithDisabled(true),
			),
		wf.Field().
			AddLeft("Star rating").
			AddRight(
				wf.RatingStars("stars", 4),
			),
		wf.Field().
			AddLeft("Sentiment rating").
			AddRight(
				wf.RatingSentiment("sentiment", 4),
			),
		wf.Field().
			AddLeft("Star rating disabled").
			AddRight(
				wf.RatingStars("disabledstars", 2).WithDisabled(true),
			),
		wf.Field().
			AddLeft("Rich text editor").
			AddRight(
				wf.InputRichText("rich", "Hello <b>World</b>!").
					WithMentionFeed("@", 0, "Harry Potter", "Tom Riddle").
					WithToolbar(
						"bold", "italic", "underline", "strikethrough", "|",
						"fontColor", "fontBackgroundColor", "|",
						"numberedList", "bulletedList", "alignment", "|",
						"blockQuote", "link", "removeFormat",
					),
			),

		// Buttons grouped together are put into a toolbar automatically
		wf.ButtonText("").Add("Text").WithDisabled(true),
		wf.ButtonElevated("").Add("Elevated").WithDisabled(true),
		wf.ButtonTonal("").Add("Tonal").WithDisabled(true),
		wf.ButtonOutlined("").Add("Outlined").WithDisabled(true),
		wf.ButtonFilled("").Add("Filled").WithDisabled(true),

		wf.SpacerNewLine(),

		wf.ButtonText("").Add(wf.Icon("build"), "Build"),
		wf.ButtonElevated("").Add(wf.Icon("thumb down"), "Dislike"),
		wf.ButtonTonal("").Add(wf.Icon("electric bolt")),
		wf.ButtonOutlined("").Add(wf.Icon("settings")),
		wf.ButtonFilled("").Add(wf.Icon("visibility")),

		wf.SpacerNewLine(),

		wf.ButtonText("").Add("Cancel").WithHrefBack(),
		wf.ButtonElevated("").Add("Cancel").WithHrefBack(),
		wf.ButtonTonal("").Add("Cancel").WithHrefBack(),
		wf.ButtonOutlined("").Add("Cancel").WithHrefBack(),
		wf.ButtonFilled("save").Add("Save"),
	)

	if form.ReadyToCommit(r) {
		form.Reset(r)
	}

	page := wf.Page().Add(
		wf.AppBar("Form input widgets"),
		form,
	)

	shared.Render(w, r, page)
}

// handleFormValidation demonstrates client-side and server-side form validation.
func HandleFormValidation(w http.ResponseWriter, r *http.Request) {
	clientForm := wf.Form().WithName("clientForm").Add(
		"Client-side validations are done in the browser and do not perform a round-trip to the backend. ",
		"Validations are always repeated on the backend to prevent circumventing the browser checks. ",
		"Standard error popups are generated by the browser.",
		wf.SpacerBreak(),

		wf.Field().
			AddLeft("").
			AddRight(
				wf.Link("?name1=Harry Potter&ssn1=444-00-4444").Add(wf.Icon("create"), "Quick fill")),
		wf.Field().
			AddLeft("Name (required)").
			AddRight(
				wf.InputText("name1", "").
					WithRequired(true).
					RedrawIfChanged(r, "name1")),
		wf.Field().
			AddLeft("SSN (not required)").
			AddRight(
				wf.InputText("ssn1", "").
					WithPlaceholder("000-00-0000").
					WithPattern(`^[0-9]{3}-[0-9]{2}-[0-9]{4}$`).
					RedrawIfChanged(r, "ssn1")),
		wf.ButtonText("").Add("Cancel").WithHrefBack(),
		wf.ButtonFilled("").Add("Submit"),
	)

	serverForm := wf.Form().WithName("serverForm").Add(
		"Server-side validations are done in the backend and require a round-trip to the server. ",
		"The entire form is re-rendered on each request. ",
		"Error messages are generated at the backend and may be customized.",
		wf.SpacerBreak(),

		wf.Field().
			AddLeft("").
			AddRight(
				wf.Link("?name2=Tom Riddle&ssn2=777-00-7777").Add(wf.Icon("create"), "Quick fill")),
		wf.Field().
			AddLeft("Name (required)").
			AddRight(
				wf.InputText("name2", "").
					WithPredicate(func(value string) (bool, string) {
						return value != "", "A name is required"
					}).
					RedrawIfChanged(r, "name2")),
		wf.Field().
			AddLeft("SSN (not required)").
			AddRight(
				wf.InputText("ssn2", "").
					WithPlaceholder("000-00-0000").
					WithPredicate(func(value string) (bool, string) {
						if value == "" {
							return true, ""
						}
						match, _ := regexp.MatchString(`^[0-9]{3}-[0-9]{2}-[0-9]{4}$`, value)
						return match, "Enter a valid social security number (hyphens required)"
					}).
					RedrawIfChanged(r, "ssn2")),

		wf.ButtonText("").Add("Cancel").WithHrefBack(),
		wf.ButtonFilled("").Add("Submit"),
	)

	snackbar := wf.Snackbar().RedrawIf(true)
	page := wf.Page().Add(
		wf.AppBar("Form validation"),

		wf.HeadlineMedium("Client-side validation"),
		clientForm,

		wf.HeadlineMedium("Server-side validation"),
		serverForm,

		wf.Debugger(),
		snackbar,
	)

	if clientForm.ReadyToCommit(r) {
		values := clientForm.Values(r)
		snackbar.Add(values.Get("name1"), " submitted OK")
		clientForm.Reset(r)
	}
	if serverForm.ReadyToCommit(r) {
		values := serverForm.Values(r)
		snackbar.Add(values.Get("name2"), " submitted OK")
		serverForm.Reset(r)
	}

	shared.Render(w, r, page)
}
