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

package table

import (
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&PageSizerWidget{}) // Ensure interface

// PageSizerWidget renders a table page sizer.
type PageSizerWidget struct {
	*widget.WidgetBase[*PageSizerWidget]
	tableName string
	options   []int
}

// PageSizer creates a new widget that renders a dropdown letting the user
// pick how many rows the table shows per page. Defaults to options
// 10/25/50/100; override with WithOptions. The widget hides itself when
// the total row count fits in the smallest option. Bind to a non-default
// table with ForTable.
func (f TableFactory) PageSizer() *PageSizerWidget {
	x := &PageSizerWidget{
		tableName: "table",
		options:   []int{10, 25, 50, 100},
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// ForTable binds this page sizer to the named table. Defaults to "table".
// Set this when the page has multiple tables.
func (wgt *PageSizerWidget) ForTable(name string) *PageSizerWidget {
	wgt.tableName = name
	return wgt
}

// WithOptions overrides the page-size choices offered in the dropdown.
// Default is 10, 25, 50, 100. The table's current default page size is
// always added to the list (and sorted), so the user never sees a missing
// "currently selected" option.
func (wgt *PageSizerWidget) WithOptions(rowsPerPage ...int) *PageSizerWidget {
	wgt.options = rowsPerPage
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *PageSizerWidget) Draw(w io.Writer, r *http.Request) (err error) {
	state := factory.StateOf(r)
	defaultPageRows := 25
	v := state.Get("_" + wgt.tableName + "_rows")
	if v != "" {
		defaultPageRows, _ = strconv.Atoi(v)
		if defaultPageRows <= 0 {
			defaultPageRows = 25
		}
	}

	// Dropdown of options
	current := defaultPageRows
	if i, err := strconv.Atoi(factory.StateOf(r).Get(wgt.tableName + "_rows")); err == nil && i >= 1 {
		current = i
	}
	options := []int{current}
	options = append(options, wgt.options...)
	sort.Ints(options)

	value := strconv.Itoa(current)
	if current == defaultPageRows {
		value = ""
	}

	dropdown := factory.Dropdown(wgt.tableName+"_rows", value).WithAutoSubmit(true)
	for i := 0; i < len(options); i++ {
		if i == 0 || options[i] != options[i-1] {
			value := strconv.Itoa(options[i])
			if options[i] == defaultPageRows {
				value = ""
			}
			dropdown.AddOption(value, strconv.Itoa(options[i]))
		}
	}

	totalRows := tblTotalRows(r, wgt.tableName)

	return Tag("div").
		Class("PageSizer").
		Attr("data-id", wgt.ID()).
		Add(dropdown, "rows").
		When(wgt.Shown(r) && (totalRows < 0 || totalRows > options[0])).
		Draw(w, r)
}

func (wgt *PageSizerWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) ||
		state.Changed(wgt.tableName+"_q")
}
