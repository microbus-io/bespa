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

var _ = Widget(&DropdownWidget{})      // Ensure interface
var _ = InputWidget(&DropdownWidget{}) // Ensure interface

// DropdownWidget renders a dropdown selector.
type DropdownWidget struct {
	*widget.InputWidgetBase[*DropdownWidget]
	value       string
	optionsKeys []string
	options     map[string]option
	errMsg      string
	predicates  Predicates
	stretch     string
}

type option struct {
	Caption  string
	Disabled bool
}

// Dropdown creates a new widget that renders a native <select>. Add
// choices via AddOption (or AddDisabledOption); insertion order is
// preserved. A blank choice is prepended automatically unless the field
// is required or already includes one. Submitted values that aren't in
// the options list fail validation. For a searchable dropdown, see
// RichDropdown.
func (f FormFactory) Dropdown(name, value string) *DropdownWidget {
	x := &DropdownWidget{
		value:   value,
		options: map[string]option{},
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// AddOption appends a selectable option. value is what gets posted;
// caption is shown to the user. Options appear in insertion order.
func (wgt *DropdownWidget) AddOption(value string, caption string) *DropdownWidget {
	wgt.optionsKeys = append(wgt.optionsKeys, value) // Keep the order of insertion
	wgt.options[value] = option{
		Caption: caption,
	}
	return wgt
}

// AddDisabledOption appends a non-selectable option — useful for section
// headers or otherwise-grey-out choices. Submitting a disabled option
// fails validation.
func (wgt *DropdownWidget) AddDisabledOption(value string, caption string) *DropdownWidget {
	wgt.optionsKeys = append(wgt.optionsKeys, value) // Keep the order of insertion
	wgt.options[value] = option{
		Caption:  caption,
		Disabled: true,
	}
	return wgt
}

// WithStretch makes the dropdown fill the available width. By default
// it's sized to its longest option label.
func (wgt *DropdownWidget) WithStretch(stretch bool) *DropdownWidget {
	if stretch {
		wgt.stretch = "width:100%"
	} else {
		wgt.stretch = ""
	}
	return wgt
}

// WithPredicate adds a custom validator. Runs after the option-list
// check, so it always receives a value that exists in the options.
func (wgt *DropdownWidget) WithPredicate(predicate func(value string) (bool, string)) *DropdownWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

// Value returns the value of the field.
func (wgt *DropdownWidget) Value(r *http.Request) string {
	value := wgt.value
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) { // || wgt.Submitted(r)
		value = state.Get(wgt.Name())
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *DropdownWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	// Required
	if value == "" && wgt.Required() {
		wgt.errMsg = "A value is required"
		return false
	}
	if value == "" && !wgt.Required() {
		return true
	}
	// Valid values
	option, found := wgt.options[value]
	if !found || option.Disabled {
		wgt.errMsg = "Invalid value"
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
func (wgt *DropdownWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *DropdownWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)

	randomID := widget.RandomAlphaNumID(8)
	tagSelect := Tag("select").
		Attr("id", randomID).
		AttrIf(wgt.Disabled(), "disabled", "1").
		Attr("value", value).
		ClassIf(invalid, "Invalid").
		Style(wgt.stretch)
	if !wgt.Disabled() {
		tagSelect.
			Attr("name", wgt.Name()).
			Attr("tabindex", "0").
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("oninput", "input_input(event)").
			Attr("oninvalid", "input_invalid(event)").
			AttrIf(wgt.Required(), "required", "1")
	}
	if _, ok := wgt.options[""]; !ok && !wgt.Required() && (!wgt.Disabled() || value == "") {
		// Add an empty option
		tagSelect.Add(Tag("option").
			Attr("value", "").
			AttrIf(value == "", "selected", "1").
			Add(""))
	}
	for _, option := range wgt.optionsKeys {
		selected := option == value
		if wgt.Disabled() && !selected {
			// Only draw the selected option when disabled
			continue
		}
		tagSelect.Add(Tag("option").
			Attr("value", option).
			AttrIf(selected, "selected", "1").
			AttrIf(wgt.options[option].Disabled, "disabled", "1").
			Add(wgt.options[option].Caption))
	}
	errTag := Tag("")
	if invalid && wgt.errMsg != "" {
		errTag = customValidityScript(randomID, wgt.errMsg)
	}
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(
			tagSelect,
			Tag("script").Add(factory.HTMLUnsafe("input_initBackgroundIconColor('", randomID, "')")),
			errTag,
		).
		When(wgt.Shown(r)).
		Draw(w, r)
}
