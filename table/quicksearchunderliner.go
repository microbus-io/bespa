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

var _ = Widget(&QuickSearchUnderlinerWidget{}) // Ensure interface

// QuickSearchUnderlinerWidget highlights its content based on the search box associated with a table.
type QuickSearchUnderlinerWidget struct {
	*widget.WidgetBase[*QuickSearchUnderlinerWidget]
	tableName string
	content   string
	prefix    bool
}

// QuickSearchUnderliner creates a new widget that renders content with any
// matches of the paired QuickSearch query underlined. Use it inside row
// cells to highlight the matched substring. Bind to a non-default table
// with ForTable.
func (f TableFactory) QuickSearchUnderliner(content string) *QuickSearchUnderlinerWidget {
	x := &QuickSearchUnderlinerWidget{
		tableName: "table",
		content:   content,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// ForTable binds this underliner to the named table. Defaults to "table".
// Set this when the page has multiple tables.
func (wgt *QuickSearchUnderlinerWidget) ForTable(name string) *QuickSearchUnderlinerWidget {
	wgt.tableName = name
	return wgt
}

// WithPrefixOnly restricts matching to word prefixes. When false (default),
// search terms match anywhere within a word.
func (wgt *QuickSearchUnderlinerWidget) WithPrefixOnly(prefix bool) *QuickSearchUnderlinerWidget {
	wgt.prefix = prefix
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *QuickSearchUnderlinerWidget) Draw(w io.Writer, r *http.Request) (err error) {
	state := factory.StateOf(r)
	q := state.Get(wgt.tableName + "_q")
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(
			factory.Underliner(wgt.content, q).WithPrefixOnly(wgt.prefix),
		).
		Draw(w, r)
}
