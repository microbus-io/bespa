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

package form

import (
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&InputHiddenWidget{})      // Ensure interface
var _ = InputWidget(&InputHiddenWidget{}) // Ensure interface

// InputHiddenWidget renders a hidden input field.
type InputHiddenWidget struct {
	*widget.InputWidgetBase[*InputHiddenWidget]
	value string
}

// InputHidden creates a new widget that renders a hidden form field —
// posted with the form but invisible to the user. Use this to carry an
// ID or context value through a submission. The value is fixed at
// construction; user input cannot change it.
func (f FormFactory) InputHidden(name string, value string) *InputHiddenWidget {
	x := &InputHiddenWidget{
		value: value,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// Value returns the value of the field.
func (wgt *InputHiddenWidget) Value(r *http.Request) string {
	return wgt.value
}

// Valid validates the field's value against all validators.
func (wgt *InputHiddenWidget) Valid(r *http.Request) bool {
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *InputHiddenWidget) Changed(r *http.Request) bool {
	return false
}

// Draw renders the widget's HTML.
func (wgt *InputHiddenWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("input").
		Attr("type", "hidden").
		Attr("name", wgt.Name()).
		Attr("value", wgt.value).
		Attr("data-id", wgt.ID()).
		When(wgt.Shown(r)).
		Draw(w, r)
}
