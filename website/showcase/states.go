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
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/microbus-io/bespa/website/resources"
	"github.com/microbus-io/bespa/website/shared"
	"github.com/microbus-io/bespa/website/storage"
	"github.com/microbus-io/bespa/widget"
)

// HandleStates renders an interactive table listing stats about the 50 US States and
// the District of Columbia.
// The table supports sorting, pagination, page sizing, and quick searching.
func HandleStates(w http.ResponseWriter, r *http.Request) {
	pageSizer := wf.PageSizer()
	paginator1 := wf.Paginator()
	paginator2 := wf.Paginator()
	tbl := wf.Table().
		WithDefaultPageRows(r, 10).
		WithDefaultSortOrder(r, "abbrev").
		Add(
			wf.Col("nwx", 6, "left").Add(wf.Sorter("abbrev", "Code")),
			wf.Col("x", 15, "left").Add(wf.Sorter("name", "Name")),
			wf.Col("wx", 10, "right").Add(wf.Sorter("pop", "Population")),
			wf.Col("n", 10, "right").Add(wf.Sorter("pop", "Pop")),
			wf.Col("wx", 10, "right").Add(wf.Sorter("land", "Land (sq mi)")),
			wf.Col("n", 10, "right").Add(wf.Sorter("land", "Land (sq mi)")),
			wf.Col("wx", 8, "right").Add(wf.Sorter("density", "Pop density")),
			wf.Col("n", 8, "right").Add(wf.Sorter("density", "Den\u200bsity")),
			wf.Col("wx", 10, "right").Add(wf.Sorter("gdp", "GDP")),
			wf.Col("n", 10, "right").Add(wf.Sorter("gdp", "GDP")),
			wf.Col("wx", 8, "right").Add(wf.Sorter("gdpcapita", "GDP/capita")),
			wf.Col("n", 8, "right").Add(wf.Sorter("gdpcapita", "GDP per capita")),
		)
	// Query and populate the table, but only if the table is re/drawn
	if tbl.Drawn(r) {
		// Query the data source
		rowFrom, rowTo := tbl.DisplayRange(r)
		states, totalRecords, _ := storage.USQuery(tbl.Query(r), tbl.SortOrder(r), rowFrom, rowTo)
		tbl.WithTotalRows(r, totalRecords) // So paginators know the last page number
		// Populate the visible rows
		for _, state := range states {
			tbl.Add(wf.Row().Add(
				wf.QuickSearchUnderliner(state.Abbrev),                          // Abbrev
				wf.QuickSearchUnderliner(state.Name),                            // Name
				state.Population,                                                // Population (wide)
				wf.Collection((state.Population+500000)/1000000, "M"),           // Population (narrow)
				wf.Float(state.Land, 0),                                         // Land (wide)
				wf.Collection(wf.Float(state.Land/1000.0, 1), "K"),              // Land (narrow)
				wf.Float(state.PopulationDensity, 2),                            // Density (wide)
				wf.Float(state.PopulationDensity, 0),                            // Density (narrow)
				wf.Collection("$", state.GDP, "M"),                              // GDP (wide)
				wf.Collection("$", state.GDP/1000.0, "B"),                       // GDP (narrow)
				wf.Collection("$", wf.Float(state.GDPPerCapita, 0)),             // GDP/capita (wide)
				wf.Collection("$", wf.Float(state.GDPPerCapita/1000.0, 0), "K"), // GDP/capita (narrow)
			).WithAction("state?st=" + state.Abbrev))
		}
	}

	page := wf.Page().Add(
		wf.AppBar("50 States"),
		`This screen brings together multiple widgets to create an interactive user experience.
		A table widget is supported by sorters located in each of the column titles,
		two paginators, a quick search filter and a page size selector.
		The widgets interact with each other via the page's state variables
		which can be observed by opening the debugger in the lower right corner.
		Resize the window to see how columns adjust dynamically to different widths.`,

		wf.Toolbar().
			AddLeft(wf.QuickSearch(),
				wf.SpacerParagraph()).
			AddRight(wf.AlignRight(paginator1)),
		tbl,
		wf.Toolbar().
			AddLeft(pageSizer).
			AddRight(wf.AlignRight(paginator2)),
		wf.Debugger(),
	)

	shared.Render(w, r, page)
}

// HandleState display stats about a single state.
func HandleState(w http.ResponseWriter, r *http.Request) {
	// Query the data source
	abbrev := wf.StateOf(r).Get("st")
	unitedState, ok, _ := storage.USLookup(abbrev)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var flagImage widget.Widget
	b, _ := resources.Bundle.ReadFile(filepath.Join("images", "stateflags", strings.ToLower(abbrev)+".webp"))
	if b != nil {
		flagImage = wf.Image("/images/stateflags/" + strings.ToLower(abbrev) + ".webp")
	}

	page := wf.Page().Add(
		wf.AppBar(unitedState.Name, " (", unitedState.Abbrev, ")"),
		wf.Splitter(0, 0, 1).AddRight(
			wf.Field().AddLeft("Population").AddRight(unitedState.Population),
			wf.Field().AddLeft("Land area (sq mi)").AddRight(unitedState.Land),
			wf.Field().AddLeft("Population density").AddRight(unitedState.PopulationDensity),
			wf.Field().AddLeft("GDP").AddRight("$", unitedState.GDP, ",000,000"),
			wf.Field().AddLeft("GDP/capita").AddRight("$", wf.Float(unitedState.GDPPerCapita, 0)),
		).AddLeft(
			flagImage,
			wf.SpacerParagraph(),
		),
	)

	shared.Render(w, r, page)
}

// HandleQueryStates queries the list of states in order to support a chips input widget.
func HandleQueryStates(w http.ResponseWriter, r *http.Request) {
	type option struct {
		Title string `json:"title"`
		Value string `json:"value"`
		Desc  string `json:"desc"`
	}
	result := struct {
		Options []*option `json:"options"`
	}{}
	abbrev := wf.StateOf(r).Get("q")
	unitedStates, _, _ := storage.USQuery(abbrev, "", 0, 8)
	for _, us := range unitedStates {
		result.Options = append(result.Options, &option{
			Title: us.Name,
			Value: us.Abbrev,
			Desc:  us.Abbrev,
		})
	}
	json.NewEncoder(w).Encode(result)
}
