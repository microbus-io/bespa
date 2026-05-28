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

var _ = Widget(&RowWidget{}) // Ensure interface

// RowWidget renders a table row.
// Rows must be nested under a table.
type RowWidget struct {
	*widget.WidgetBase[*RowWidget]
	lineHeight string
	href       string
	target     string
	color      string
	cells      []Widget
	parent     *TableWidget
	vAlign     string
}

// Row creates a new widget for one table row. Add cells in declared
// column order via Row.Add; missing cells are rendered as empty.
// Add the row to a Table via Table.Add.
func (f TableFactory) Row() *RowWidget {
	x := &RowWidget{
		lineHeight: "0",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithAction makes the row clickable, navigating to href.
// Accepts the full action-URL grammar (`?key=`, `^?…`, `/path`, etc.).
func (wgt *RowWidget) WithAction(href string) *RowWidget {
	wgt.href = href
	return wgt
}

// WithTarget sets the HTML target for the row's action. Defaults to the
// page's `_target` state variable when unset.
func (wgt *RowWidget) WithTarget(target string) *RowWidget {
	wgt.target = target
	return wgt
}

// WithTextColorDisabled greys out the row's text to signal it's inactive.
// The row itself stays interactive — combine with WithAction("") to also
// drop the click target.
func (wgt *RowWidget) WithTextColorDisabled() *RowWidget {
	wgt.color = "TextColorDisabled"
	return wgt
}

// WithVerticalAlign sets the vertical alignment of the row to "middle" (default), "top" or "bottom".
func (wgt *RowWidget) WithVerticalAlign(verticalAlign string) *RowWidget {
	wgt.vAlign = strings.ToLower(verticalAlign)
	return wgt
}

// Add adds nested widgets.
func (wgt *RowWidget) Add(cells ...any) *RowWidget {
	wgt.cells = Many(wgt.cells, cells)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *RowWidget) Children() []Widget {
	return wgt.cells
}

// Draw renders the widget's HTML.
func (wgt *RowWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if wgt.parent == nil {
		return errors.New("row must be nested inside a table")
	}
	trTag := Tag("tr").
		Class("MinRowHt"+wgt.lineHeight, wgt.color).
		ClassIf(wgt.vAlign == "top", "VAlignTop").
		ClassIf(wgt.vAlign == "bottom", "VAlignBottom").
		ClassIf(wgt.vAlign == "middle", "VAlignMiddle").
		Attr("data-id", wgt.ID())
	if wgt.href != "" {
		trTag.
			Attr("data-href", "1").
			Attr("onclick", "row_click(event)").
			Attr("onkeydown", "row_keydown(event)").
			Attr("tabindex", "0")
	}
	state := factory.StateOf(r)
	target := wgt.target
	if target == "" {
		target = state.Get("_target")
	}
	for i := range wgt.parent.cols {
		col := wgt.parent.cols[i].(*ColumnWidget)
		var cell Widget
		if i < len(wgt.cells) {
			cell = wgt.cells[i]
		} else {
			cell = Text("")
		}
		tdTag := Tag("td").
			Attr("align", col.align).
			Class(col.visibility).
			Class(col.firstOfWidth()).
			Class(col.lastOfWidth()).
			Hide(!col.Shown(r)).
			Add(cell)
		if i == 0 && wgt.href != "" {
			tdTag.Add(Tag("a").
				Attr("href", wgt.href).
				Hide(true).
				Attr("target", target))
		}
		trTag.Add(tdTag)
	}
	return trTag.
		When(wgt.Shown(r)).
		Draw(w, r)
}
