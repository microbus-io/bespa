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

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&SorterWidget{}) // Ensure interface

// SorterWidget renders a sorter.
type SorterWidget struct {
	*widget.WidgetBase[*SorterWidget]
	tableName string
	sortKey   string
	children  []Widget
}

// Sorter creates a new widget that renders a clickable sort trigger,
// typically placed inside a Column. sortKey is the identifier passed to
// your store (via Table.SortOrder); label is the visible text. Clicking
// cycles ascending → descending → off, updating `<table>_sort` in state.
// Bind to a non-default table with ForTable.
func (f TableFactory) Sorter(sortKey string, label any) *SorterWidget {
	x := &SorterWidget{
		tableName: "table",
		sortKey:   sortKey,
		children:  Many(label),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// ForTable binds this sorter to the named table. Defaults to "table".
// Set this when the page has multiple tables.
func (wgt *SorterWidget) ForTable(name string) *SorterWidget {
	wgt.tableName = name
	return wgt
}

// Add adds nested widgets.
func (wgt *SorterWidget) Add(children ...any) *SorterWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *SorterWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *SorterWidget) Draw(w io.Writer, r *http.Request) (err error) {
	defaultSortOrder := factory.StateOf(r).Get("_" + wgt.tableName + "_sort")
	sortOrder := factory.StateOf(r).Get(wgt.tableName + "_sort")
	if sortOrder == "" {
		sortOrder = defaultSortOrder
	}
	nextSortOrder := wgt.sortKey
	icon := ""
	if sortOrder == wgt.sortKey {
		icon = "arrow_upward"
		nextSortOrder = "-" + wgt.sortKey
	} else if sortOrder == "-"+wgt.sortKey {
		icon = "arrow_downward"
		nextSortOrder = ""
		if sortOrder == defaultSortOrder {
			nextSortOrder = wgt.sortKey
		}
	}
	anchorTag := Tag("a").
		Class("Sorter").
		Attr("data-id", wgt.ID()).
		Attr("tabindex", "0").
		Attr("href", "?"+wgt.tableName+"_sort="+nextSortOrder).
		Add(wgt.children)
	if icon != "" {
		anchorTag.Add(" ", factory.Icon(icon))
	}
	return anchorTag.
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wb *SorterWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wb.WidgetBase.Drawn(r) ||
		state.Changed(wb.tableName+"_sort")
}
