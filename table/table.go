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
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&TableWidget{}) // Ensure interface

// TableWidget renders a table.
type TableWidget struct {
	*widget.WidgetBase[*TableWidget]
	cols         []Widget
	rows         []Widget
	name         string
	border       string
	noHeader     bool
	emptyMessage string
	vAlign       string
}

// Table creates a new widget that renders a data table with header, sort,
// pagination, and filter hooks. Add Column and Row children, then pair with
// the companion widgets (Sorter, Paginator, PageSizer, QuickSearch) — they
// bind to this table by name (default "table").
func (f TableFactory) Table() *TableWidget {
	x := &TableWidget{
		name:         "table",
		emptyMessage: "The table is empty",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Grid creates a new widget that renders a borderless, header-less table —
// suitable for laying out tiles or summary rows where headings would add
// noise.
func (f TableFactory) Grid() *TableWidget {
	x := &TableWidget{
		name:         "table",
		border:       "NoBorder",
		noHeader:     true,
		emptyMessage: "The table is empty",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithName sets the table's name. Companion widgets (Sorter, Paginator,
// PageSizer, QuickSearch, QuickSearchUnderliner) bind to the table by this
// name via ForTable, and the state variables they read/write are scoped
// under it (`<name>_sort`, `<name>_page`, `<name>_rows`, `<name>_q`).
// Use distinct names when multiple tables share a page. Default is "table".
func (wgt *TableWidget) WithName(name string) *TableWidget {
	wgt.name = name
	return wgt
}

// WithEmptyMessage sets the message rendered as a placeholder row when the
// table has no data rows. Default is "The table is empty".
func (wgt *TableWidget) WithEmptyMessage(emptyMsg string) *TableWidget {
	wgt.emptyMessage = emptyMsg
	return wgt
}

// WithVerticalAlign sets the vertical alignment of cells within rows.
// Accepts "middle" (default), "top", or "bottom"; other values are ignored.
// Override per row with Row.WithVerticalAlign.
func (wgt *TableWidget) WithVerticalAlign(verticalAlign string) *TableWidget {
	wgt.vAlign = strings.ToLower(verticalAlign)
	return wgt
}

// WithDefaultSortOrder seeds the initial sort order used when the user
// hasn't clicked a sorter yet. The value is a sort key, optionally prefixed
// with "-" for descending. Empty means no default order. Call this each
// request — it writes to the table's reserved default state variable.
func (wgt *TableWidget) WithDefaultSortOrder(r *http.Request, sortOrder string) *TableWidget {
	state := factory.StateOf(r)
	state.Set("_"+wgt.name+"_sort", sortOrder)
	return wgt
}

// WithBorder controls whether cell borders are drawn. Default is true
// (borders shown).
func (wgt *TableWidget) WithBorder(border bool) *TableWidget {
	if !border {
		wgt.border = "NoBorder"
	} else {
		wgt.border = ""
	}
	return wgt
}

// WithHeader controls whether the column-header row is rendered. Default
// is true.
func (wgt *TableWidget) WithHeader(header bool) *TableWidget {
	wgt.noHeader = !header
	return wgt
}

// SortOrder returns the currently active sort key, falling back to the
// default set via WithDefaultSortOrder. A leading "-" indicates descending.
// Empty means unsorted. Pass this into your store query.
func (wgt *TableWidget) SortOrder(r *http.Request) string {
	state := factory.StateOf(r)
	v := state.Get(wgt.name + "_sort")
	if v == "" {
		v = state.Get("_" + wgt.name + "_sort") // Default sort order
	}
	return v
}

// WithTotalRows tells the table how many rows the underlying data source
// has in total — this drives the paginator's page count and the empty-state.
// Pass a negative value when the total is unknown (the paginator falls back
// to a "Next…" affordance). Call this after querying your store and before
// Draw.
func (wgt *TableWidget) WithTotalRows(r *http.Request, totalNumRows int) *TableWidget {
	state := factory.StateOf(r)
	if totalNumRows >= 0 {
		state.Set("_"+wgt.name+"_total", strconv.Itoa(totalNumRows))
	} else {
		state.Del("_" + wgt.name + "_total")
	}
	return wgt
}

// tblTotalRows returns the total number of rows.
// If the total number of rows is unknown, -1 is returned.
func tblTotalRows(r *http.Request, tableName string) int {
	state := factory.StateOf(r)
	v := state.Get("_" + tableName + "_total")
	if v == "" {
		return -1
	}
	if v == "0" {
		return 0
	}
	totalRows, _ := strconv.Atoi(v)
	return totalRows
}

// tblNumPages returns the total number of pages, given the total number of rows and the page size.
// If the total number of rows is unknown, -1 is returned.
func tblNumPages(r *http.Request, tableName string) int {
	totalRows := tblTotalRows(r, tableName)
	if totalRows < 0 {
		return -1
	}
	if totalRows == 0 {
		return 0
	}
	pageRows := tblPageRows(r, tableName)
	if totalRows%pageRows == 0 {
		return totalRows / pageRows
	}
	return totalRows/pageRows + 1
}

// DisplayRange returns the half-open row range [fromRow, toRow) the current
// page should show. Feed it directly into your data store query. When the
// total is unknown, toRow may extend past the real row count — just take
// whatever you get back. Changes to query/sort/page-size automatically
// reset the page to 1.
func (wgt *TableWidget) DisplayRange(r *http.Request) (fromRow int, toRow int) {
	state := factory.StateOf(r)
	itemsPerPage := tblPageRows(r, wgt.name)
	totalRows := tblTotalRows(r, wgt.name)
	if state.Changed(wgt.name+"_q") || state.Changed(wgt.name+"_sort") || state.Changed(wgt.name+"_rows") {
		state.Del(wgt.name + "_page")
	}
	if i, err := strconv.Atoi(state.Get(wgt.name + "_page")); err == nil && i >= 1 {
		if totalRows >= 0 && i*itemsPerPage > totalRows {
			return (i - 1) * itemsPerPage, totalRows
		}
		return (i - 1) * itemsPerPage, i * itemsPerPage
	}
	if totalRows >= 0 && itemsPerPage > totalRows {
		return 0, totalRows
	}
	return 0, itemsPerPage
}

// WithDefaultPageRows sets the initial page size. Without a paired
// PageSizer this becomes the fixed page size. Default if never set is 25.
// Non-positive values are ignored.
func (wgt *TableWidget) WithDefaultPageRows(r *http.Request, numRowsPerPage int) *TableWidget {
	if numRowsPerPage >= 1 {
		state := factory.StateOf(r)
		state.Set("_"+wgt.name+"_rows", strconv.Itoa(numRowsPerPage))
	}
	return wgt
}

// PageRows returns the maximum number of rows to show in the table.
// A page sizer widget can be used to allow the user to customize the page size.
func tblPageRows(r *http.Request, tableName string) int {
	state := factory.StateOf(r)
	v := state.Get(tableName + "_rows")
	if v == "" {
		v = state.Get("_" + tableName + "_rows")
	}
	if v != "" {
		i, err := strconv.Atoi(v)
		if err == nil && i >= 1 {
			return i
		}
	}
	return 25
}

// Query returns the current quick-search text entered through the paired
// QuickSearch widget (or empty if none). Apply it as a filter when
// fetching rows from your data source.
func (wgt *TableWidget) Query(r *http.Request) string {
	state := factory.StateOf(r)
	return state.Get(wgt.name + "_q")
}

/*
Add adds nested widgets to the table.
Column widgets are added as column headers.
Row widgets are added as rows.
Remaining elements are composed into a single row.

Example:

	// Add columns
	tbl.Add(
		wgt.Column(3, "Name"),
		wgt.Column(1, "DOB"),
	)
	// Add a row
	tbl.Add(
		wgt.Row("Albert Einstein", "3/14/1879"),
	)
	// Add a row
	tbl.Add(
		"Steven Hawking", "1/8/1942",
	)
*/
func (wgt *TableWidget) Add(children ...any) *TableWidget {
	others := []any{}
	for _, c := range children {
		switch v := c.(type) {
		case *ColumnWidget:
			wgt.cols = append(wgt.cols, v)
		case *RowWidget:
			wgt.rows = append(wgt.rows, v)
		default:
			others = append(others, v)
		}
	}
	if len(others) > 0 {
		wgt.rows = append(wgt.rows, factory.Row().Add(others...))
	}
	for _, c := range wgt.cols {
		c.(*ColumnWidget).parent = wgt
	}
	for _, r := range wgt.rows {
		r.(*RowWidget).parent = wgt
	}
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *TableWidget) Children() []Widget {
	return Many(wgt.cols, wgt.rows)
}

// Draw renders the widget's HTML.
func (wgt *TableWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)

	var emptyRow *widget.TagWidget
	if len(wgt.rows) == 0 {
		emptyRow = Tag("tr").
			Class("EmptyRow").
			Add(
				Tag("td").
					Attr("colspan", strconv.Itoa(len(wgt.cols))).
					Class("Narrow", "Wide", "Expanded").
					Class("FirstNarrow", "FirstWide", "FirstExpanded").
					Class("LastNarrow", "LastWide", "LastExpanded").
					Add(wgt.emptyMessage),
			)
	}

	// Calculate column percentile widths for each column for each width category (narrow, wide, expanded)
	sumNarrow := 0
	sumWide := 0
	sumExpanded := 0
	style := ""
	for _, col := range wgt.cols {
		c := col.(*ColumnWidget)
		if strings.Contains(c.visibility, "Narrow") {
			sumNarrow += c.width
		}
		if strings.Contains(c.visibility, "Wide") {
			sumWide += c.width
		}
		if strings.Contains(c.visibility, "Expanded") {
			sumExpanded += c.width
		}
	}
	for i, col := range wgt.cols {
		c := col.(*ColumnWidget)
		if strings.Contains(c.visibility, "Narrow") {
			style += fmt.Sprintf(".DataTable.Width_600#%s > TABLE > TBODY > TR >*:nth-child(%d) {width:%.2f%%}\n",
				randomID, i+1, float32(c.width)*100.0/float32(sumNarrow))
		}
		if strings.Contains(c.visibility, "Wide") {
			style += fmt.Sprintf(".DataTable.Width600_1200#%s > TABLE > TBODY > TR >*:nth-child(%d) {width:%.2f%%}\n",
				randomID, i+1, float32(c.width)*100.0/float32(sumWide))
		}
		if strings.Contains(c.visibility, "Expanded") {
			style += fmt.Sprintf(".DataTable.Width1200_#%s > TABLE > TBODY > TR > *:nth-child(%d) {width:%.2f%%}\n",
				randomID, i+1, float32(c.width)*100.0/float32(sumExpanded))
		}
	}
	styleTag := Tag("style").Add(HTMLUnsafe(style))

	header := Tag("tr").Add(wgt.cols)
	if wgt.noHeader {
		header = Tag("")
	}
	return Tag("div").
		Class("DataTable", "Block", wgt.border).
		ClassIf(wgt.vAlign == "top", "VAlignTop").
		ClassIf(wgt.vAlign == "bottom", "VAlignBottom").
		ClassIf(wgt.vAlign == "middle", "VAlignMiddle").
		Attr("id", randomID).
		Attr("data-id", wgt.ID()).
		Attr("data-observe-width", "600,1200").
		Add(
			styleTag,
			Tag("table").
				Attr("cellspacing", "0").
				Attr("cellpadding", "0").
				Add(
					// Header row
					header,
					// Data rows
					wgt.rows,
					// Empty row
					emptyRow,
				),
			// Reduce flickering
			Tag("script").Add(HTMLUnsafe("table_init('", randomID, "')")),
		).
		When(wgt.Shown(r) && len(wgt.cols) > 0).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *TableWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) ||
		state.Changed(wgt.name+"_sort") ||
		state.Changed(wgt.name+"_page") ||
		state.Changed(wgt.name+"_rows") ||
		state.Changed(wgt.name+"_q")
}
