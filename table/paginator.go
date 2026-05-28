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

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&PaginatorWidget{}) // Ensure interface

// PaginatorWidget renders a paginator.
type PaginatorWidget struct {
	*widget.WidgetBase[*PaginatorWidget]
	tableName string
}

// Paginator creates a new widget that renders page-number flippers for a
// table. The paginator hides itself when there's only one page; when the
// total row count is unknown (see Table.WithTotalRows) it falls back to a
// "Next…" affordance. Bind to a non-default table with ForTable.
func (f TableFactory) Paginator() *PaginatorWidget {
	x := &PaginatorWidget{
		tableName: "table",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// ForTable binds this paginator to the named table. Defaults to "table".
// Set this when the page has multiple tables.
func (wgt *PaginatorWidget) ForTable(name string) *PaginatorWidget {
	wgt.tableName = name
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *PaginatorWidget) Draw(w io.Writer, r *http.Request) (err error) {
	// Pages start from 1
	state := factory.StateOf(r)
	page := 1
	if state.Changed(wgt.tableName+"_q") || state.Changed(wgt.tableName+"_sort") || state.Changed(wgt.tableName+"_rows") {
		state.Del(wgt.tableName + "_page")
	}
	if i, err := strconv.Atoi(state.Get(wgt.tableName + "_page")); err == nil && i >= 1 {
		page = i
	}

	// Prepare the flippers
	hrefOf := func(index int) string {
		if index == 1 {
			return fmt.Sprintf("?%s_page=", wgt.tableName)
		} else {
			return fmt.Sprintf("?%s_page=%d", wgt.tableName, index)
		}
	}
	toTheSide := 2
	flippers := factory.Collection("Page ")
	// [1] ...
	if page > toTheSide+1 {
		flippers.Add(factory.Link(hrefOf(1)).Add(1).WithDisabled(page == 1))
		if page > toTheSide+2 {
			flippers.Add(HTMLUnsafe(`<span class="Ellipsis">`), factory.Icon("more_horiz"), HTMLUnsafe("</span>"))
		}
	}
	// [n-2] [n-1] [n]
	for i := page - toTheSide; i <= page; i++ {
		if i >= 1 {
			flippers.Add(factory.Link(hrefOf(i)).Add(i).WithDisabled(page == i))
		}
	}
	numPages := tblNumPages(r, wgt.tableName)
	if numPages >= 0 {
		// [n+1] [n+2]
		for i := page + 1; i <= page+toTheSide && i <= numPages; i++ {
			flippers.Add(factory.Link(hrefOf(i)).Add(i).WithDisabled(page == i))
		}
		// ... [max]
		if page+toTheSide < numPages {
			if page+toTheSide+1 < numPages {
				flippers.Add(HTMLUnsafe(`<span class="Ellipsis">`), factory.Icon("more_horiz"), HTMLUnsafe("</span>"))
			}
			flippers.Add(factory.Link(hrefOf(numPages)).Add(numPages).WithDisabled(page == numPages))
		}
	} else {
		// [next] ...
		flippers.Add(factory.Link(hrefOf(page + 1)).Add("Next"))
		flippers.Add(HTMLUnsafe(`<span class="Ellipsis">`), factory.Icon("more_horiz"), HTMLUnsafe("</span>"))
	}

	return Tag("div").
		Class("Paginator").
		Attr("data-id", wgt.ID()).
		Add(flippers.Children()).
		When(wgt.Shown(r) && numPages != 0 && numPages != 1).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *PaginatorWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) ||
		state.Changed(wgt.tableName+"_page") ||
		state.Changed(wgt.tableName+"_rows") ||
		state.Changed(wgt.tableName+"_sort") ||
		state.Changed(wgt.tableName+"_q")
}
