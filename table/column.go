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
	"strings"

	"github.com/microbus-io/bespa/widget"
	"github.com/microbus-io/errors"
)

var _ = Widget(&ColumnWidget{}) // Ensure interface

// ColumnWidget renders a table column header.
// Columns must be nested under a table.
type ColumnWidget struct {
	*widget.WidgetBase[*ColumnWidget]
	width      int
	align      string
	children   []Widget
	parent     *TableWidget
	visibility string
}

// Column creates a new widget for a single column header.
// Add it to a Table via Table.Add; the header cell's content is populated
// via Column.Add. Defaults: visible at all viewport widths, width 100
// (an arbitrary unit weighted against other columns' widths).
func (f TableFactory) Column() *ColumnWidget {
	x := &ColumnWidget{
		visibility: "Narrow Wide Expanded",
		width:      100,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Col is a shorthand factory equivalent to
// Column().WithAlignment(alignment).WithWidth(width).WithVisibility(visibility).
// Use it when you want a one-line column declaration.
func (f TableFactory) Col(visibility string, width int, alignment string) *ColumnWidget {
	return f.Column().
		WithAlignment(alignment).
		WithWidth(width).
		WithVisibility(visibility)
}

// WithAlignment sets the text alignment of the column.
// Valid values are "left", "right", "center" and the empty string "" for the default behavior.
func (wgt *ColumnWidget) WithAlignment(alignment string) *ColumnWidget {
	if alignment == "left" || alignment == "right" || alignment == "center" || alignment == "" {
		wgt.align = alignment
	}
	return wgt
}

// WithWidth sets the width of the column in relation to the total widths of all the columns.
// The default width is 100, so for example, setting a width of 200 sets the column to be twice as wide.
//
// An alternative approach that yields good results is setting the width of the columns to the expected
// number of characters of their content. Width of all columns must be explicitly set for this approach
// to work correctly.
func (wgt *ColumnWidget) WithWidth(width int) *ColumnWidget {
	if width > 0 {
		wgt.width = width
	}
	return wgt
}

// WithVisibility sets the visibility of the column based on the total width available to the table.
// The specification is a string that contains a letter for each of the following 3 cases:
// "n" for narrow (under 600px); "w" for wide (600-1199px); and "x" for expanded (1200px or more).
// By default a column is visible in all 3 situations, i.e. "nwx".
//
// In narrow spaces it typically makes sense to
// hide certain columns outright;
// merge two columns into a third by hiding the former and showing the latter;
// relocate the action menu from the right-most column to the left-most if there's a chance of it
// going off screen.
//
// In expanded spaces it typically makes sense to
// reveal less-important contextual columns.
func (wgt *ColumnWidget) WithVisibility(nwx string) *ColumnWidget {
	wgt.visibility = ""
	if strings.Contains(nwx, "n") {
		wgt.visibility += "Narrow "
	}
	if strings.Contains(nwx, "w") {
		wgt.visibility += "Wide "
	}
	if strings.Contains(nwx, "x") {
		wgt.visibility += "Expanded "
	}
	wgt.visibility = strings.TrimSpace(wgt.visibility)
	return wgt
}

// firstOfWidth returns a list of classes that indicate when this column is the first visible for its width.
func (wgt *ColumnWidget) firstOfWidth() string {
	var narrow, wide, expanded bool
	for i := range wgt.parent.cols {
		col := wgt.parent.cols[i].(*ColumnWidget)
		if col == wgt {
			break
		}
		narrow = narrow || strings.Contains(col.visibility, "Narrow")
		wide = wide || strings.Contains(col.visibility, "Wide")
		expanded = expanded || strings.Contains(col.visibility, "Expanded")
	}
	first := ""
	if !narrow && strings.Contains(wgt.visibility, "Narrow") {
		first += "FirstNarrow "
	}
	if !wide && strings.Contains(wgt.visibility, "Wide") {
		first += "FirstWide "
	}
	if !expanded && strings.Contains(wgt.visibility, "Expanded") {
		first += "FirstExpanded "
	}
	return strings.TrimSpace(first)
}

// lastOfWidth returns a list of classes that indicate when this column is the last visible for its width.
func (wgt *ColumnWidget) lastOfWidth() string {
	var narrow, wide, expanded bool
	for i := len(wgt.parent.cols) - 1; i >= 0; i-- {
		col := wgt.parent.cols[i].(*ColumnWidget)
		if col == wgt {
			break
		}
		narrow = narrow || strings.Contains(col.visibility, "Narrow")
		wide = wide || strings.Contains(col.visibility, "Wide")
		expanded = expanded || strings.Contains(col.visibility, "Expanded")
	}
	last := ""
	if !narrow && strings.Contains(wgt.visibility, "Narrow") {
		last += "LastNarrow "
	}
	if !wide && strings.Contains(wgt.visibility, "Wide") {
		last += "LastWide "
	}
	if !expanded && strings.Contains(wgt.visibility, "Expanded") {
		last += "LastExpanded "
	}
	return strings.TrimSpace(last)
}

// Add adds nested widgets.
func (wgt *ColumnWidget) Add(children ...any) *ColumnWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *ColumnWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *ColumnWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if wgt.parent == nil {
		return errors.New("column must be nested inside a table")
	}
	shown := ""
	if !wgt.Shown(r) {
		shown = "display:none"
	}
	return Tag("th").
		Attr("data-id", wgt.ID()).
		Attr("scope", "col").
		Attr("align", wgt.align).
		Style(shown).
		Class(wgt.visibility).
		Class(wgt.firstOfWidth()).
		Class(wgt.lastOfWidth()).
		Add(wgt.children).
		Draw(w, r)
}
