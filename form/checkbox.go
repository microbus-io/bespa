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

var _ = Widget(&CheckboxWidget{})      // Ensure interface
var _ = InputWidget(&CheckboxWidget{}) // Ensure interface

// CheckboxWidget renders a Checkbox.
type CheckboxWidget struct {
	*widget.InputWidgetBase[*CheckboxWidget]
	checked    bool
	children   []Widget
	predicates Predicates
	errMsg     string
}

// Checkbox creates a new widget that renders a Material checkbox.
// name is the state variable; checked is the initial state. The posted
// value is "1" when checked and "" when not — read it via Checked.
// Pair with WithRequired to force the user to tick it (e.g. T&Cs).
func (f FormFactory) Checkbox(name string, checked bool) *CheckboxWidget {
	x := &CheckboxWidget{
		checked: checked,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// Add adds nested widgets.
func (wgt *CheckboxWidget) Add(children ...any) *CheckboxWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *CheckboxWidget) Children() []Widget {
	return wgt.children
}

// WithPredicate adds a custom validator. value is "1" when checked and
// "" otherwise; predicates run on both states.
func (wgt *CheckboxWidget) WithPredicate(predicate func(value string) (valid bool, errMsg string)) *CheckboxWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

// Value returns "1" when the checkbox is currently checked, "" otherwise.
// Prefer Checked for boolean reads.
func (wgt *CheckboxWidget) Value(r *http.Request) string {
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
		} else { // Could be "0"
			value = ""
		}
	} else if wgt.Submitted(r) {
		value = ""
	}
	return value
}

// Checked reports whether the checkbox is currently checked, taking into
// account the user's input and the initial state.
func (wgt *CheckboxWidget) Checked(r *http.Request) bool {
	return wgt.Value(r) != ""
}

// Valid validates the field's value against all validators.
func (wgt *CheckboxWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value != "1" {
		return false
	}
	// Predicates
	if ok, errMsg := wgt.predicates.Validate(value); !ok {
		wgt.errMsg = errMsg
		return false
	}
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *CheckboxWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.checked != wgt.Checked(r)
}

// Draw renders the widget's HTML.
func (wgt *CheckboxWidget) Draw(w io.Writer, r *http.Request) (err error) {
	invalid := !wgt.Valid(r)
	checked := wgt.Checked(r)
	randomID := widget.RandomAlphaNumID(8)
	inputTag := Tag("input").
		Attr("id", randomID).
		Attr("type", "checkbox").
		Attr("value", "1").
		ClassIf(invalid, "Invalid").
		AttrIf(checked, "checked", "1").
		AttrIf(wgt.Disabled(), "disabled", "1")
	if !wgt.Disabled() {
		inputTag.Attr("name", wgt.Name()).
			Attr("tabindex", "0").
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("oninput", "input_input(event)").
			Attr("oninput", "input_initBackgroundIconColor('"+randomID+"')").
			Attr("oninvalid", "input_invalid(event)").
			AttrIf(wgt.Required(), "required", "1")
	}
	labelTag := Tag("")
	if len(wgt.children) > 0 {
		labelTag = Tag("label").
			Attr("for", randomID).
			Add(wgt.children)
	}
	errTag := Tag("")
	if invalid && wgt.errMsg != "" {
		errTag = customValidityScript(randomID, wgt.errMsg)
	}
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Class("Checkbox").
		Add(
			inputTag,
			labelTag,
			errTag,
			Tag("script").Add(factory.HTMLUnsafe("input_initBackgroundIconColor('", randomID, "')")),
		).
		When(wgt.Shown(r)).
		Draw(w, r)
}
