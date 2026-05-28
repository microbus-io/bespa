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

	"github.com/microbus-io/bespa/form"
	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&QuickSearchWidget{}) // Ensure interface

// QuickSearchWidget renders a quick search input box for a table.
type QuickSearchWidget struct {
	*widget.WidgetBase[*QuickSearchWidget]
	tableName string
	input     *form.InputTextWidget
}

// QuickSearch creates a new widget that renders an auto-submitting search
// input. Typing into it sets `<table>_q` in state; read it back via
// Table.Query when fetching rows. Bind to a non-default table with
// ForTable.
func (f TableFactory) QuickSearch() *QuickSearchWidget {
	x := &QuickSearchWidget{
		tableName: "table",
		input: factory.InputText("", "").
			WithLength(0, 32).
			WithWidth(16).
			WithPlaceholder("Search"),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// ForTable binds this quick search to the named table. Defaults to
// "table". Set this when the page has multiple tables.
func (wgt *QuickSearchWidget) ForTable(name string) *QuickSearchWidget {
	wgt.tableName = name
	return wgt
}

// WithLength sets the min and max lengths, in characters, allowed by the field.
// A negative value indicates unbound.
// A field with a minimum length greater than 0 is automatically assumed as required.
func (wgt *QuickSearchWidget) WithLength(minChars int, maxChars int) *QuickSearchWidget {
	wgt.input.WithLength(minChars, maxChars)
	return wgt
}

// WithWidth sets the visual width of the field, in characters.
// By default, the field stretches to the full available width (100%).
func (wgt *QuickSearchWidget) WithWidth(chars int) *QuickSearchWidget {
	wgt.input.WithWidth(chars)
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *QuickSearchWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return wgt.input.WithName(wgt.tableName+"_q").
		WithID(wgt.ID()).
		WithAutoSubmit(true).
		HideIf(!wgt.Shown(r)).
		Draw(w, r)
}
