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

var _ = Widget(&ToggleWidget{})      // Ensure interface
var _ = InputWidget(&ToggleWidget{}) // Ensure interface

// ToggleWidget renders a switch toggle.
type ToggleWidget struct {
	*widget.InputWidgetBase[*ToggleWidget]
	checked bool
}

// Toggle creates a new widget that renders a Material switch.
// Functionally identical to Checkbox (posts "1"/"") but visually
// represented as a sliding switch — prefer it for on/off settings rather
// than agree-to-terms cases.
func (f FormFactory) Toggle(name string, checked bool) *ToggleWidget {
	x := &ToggleWidget{
		checked: checked,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// Value returns "1" when the toggle is on, "" when off. Prefer Checked
// for boolean reads.
func (wgt *ToggleWidget) Value(r *http.Request) string {
	value := ""
	if wgt.checked {
		value = "1"
	}
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) {
		if state.Get(wgt.Name()) == "1" {
			value = "1"
		} else {
			value = ""
		}
	} else if wgt.Submitted(r) {
		value = ""
	}
	return value
}

// Checked reports whether the toggle is currently on.
func (wgt *ToggleWidget) Checked(r *http.Request) bool {
	return wgt.Value(r) != ""
}

// Valid validates the field's value against all validators.
func (wgt *ToggleWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		return false
	}
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *ToggleWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.checked != wgt.Checked(r)
}

// Draw renders the widget's HTML.
func (wgt *ToggleWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)
	checked := wgt.checked
	if wgt.Submitted(r) {
		checked = (value != "")
	}

	inputTag := Tag("input").
		Attr("type", "checkbox").
		Attr("tabindex", "-1").
		Attr("value", "1").
		AttrIf(checked, "checked", "1").
		ClassIf(invalid, "Invalid").
		AttrIf(wgt.Disabled(), "disabled", "1")
	if !wgt.Disabled() {
		inputTag.
			Attr("name", wgt.Name()).
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("oninput", "input_input(event)").
			Attr("oninvalid", "input_invalid(event)").
			AttrIf(wgt.Required(), "required", "1")
	}
	thumbTag := Tag("span").
		Class("Thumb").
		Add(factory.Icon("clear"), factory.Icon("check"))
	toggleTag := Tag("span").
		Attr("data-id", wgt.ID()).
		Class("Toggle").
		Add(inputTag, thumbTag).
		AttrIf(wgt.Disabled(), "disabled", "1").
		AttrIf(!wgt.Disabled(), "tabindex", "0").
		AttrIf(!wgt.Disabled(), "onkeydown", "toggle_keydown(event)").
		AttrIf(!wgt.Disabled(), "onclick", "toggle_click(event)").
		ClassIf(invalid, "Invalid").
		ClassIf(checked, "Selected")
	return toggleTag.
		When(wgt.Shown(r)).
		Draw(w, r)
}
