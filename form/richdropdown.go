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

var _ = Widget(&RichDropdownWidget{})      // Ensure interface
var _ = InputWidget(&RichDropdownWidget{}) // Ensure interface

// RichDropdownWidget renders a rich dropdown selector.
type RichDropdownWidget struct {
	*widget.InputWidgetBase[*RichDropdownWidget]
	value       string
	optionsKeys []string
	options     map[string][]Widget
}

// RichDropdown creates a new widget that renders a dropdown whose options
// can contain arbitrary widgets (icons, avatars, multi-line descriptions),
// not just text. Use AddOption to attach a value and its rendered body.
// For a plain text dropdown, use Dropdown instead.
func (f FormFactory) RichDropdown(name string, value string) *RichDropdownWidget {
	x := &RichDropdownWidget{
		value:   value,
		options: map[string][]Widget{},
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// AddOption appends a choice. value is what gets posted when selected;
// body is the widget tree shown for this option (both in the closed
// dropdown and in the open list). Options render in insertion order.
func (wgt *RichDropdownWidget) AddOption(value string, body ...any) *RichDropdownWidget {
	wgt.optionsKeys = append(wgt.optionsKeys, value) // Keep the order of insertion
	wgt.options[value] = Many(body)
	return wgt
}

// Value returns the value of the field.
func (wgt *RichDropdownWidget) Value(r *http.Request) string {
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
func (wgt *RichDropdownWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		return false
	}
	for v := range wgt.options {
		if v == value {
			return true
		}
	}
	return false
}

// Changed indicates if the value of the field changed.
func (wgt *RichDropdownWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *RichDropdownWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)

	// Hidden input to hold the value
	hiddenTag := Tag("")
	if !wgt.Disabled() {
		hiddenTag = Tag("input").
			Attr("type", "hidden").
			Attr("name", wgt.Name()).
			Attr("value", value)
	}

	// Current value
	randomID := widget.RandomAlphaNumID(8)
	divTag := Tag("div").
		Class("RichDropdown").
		Attr("id", randomID).
		ClassIf(invalid, "Invalid").
		AttrIf(wgt.Disabled(), "disabled", "1").
		AttrIf(!wgt.Disabled(), "onclick", "richdropdown_click(event)").
		AttrIf(!wgt.Disabled(), "onmouseleave", "richdropdown_mouseleave(event)").
		AttrIf(!wgt.Disabled(), "tabindex", "0")
	if x, ok := wgt.options[value]; ok {
		// Show the initial value in the drop down
		divTag.Add(Tag("div").Add(x))
	} else {
		// The initial value was not found, show the first option hidden to give the popup the correct height
		if len(wgt.optionsKeys) > 0 {
			divTag.Add(Tag("div").
				Style("visibility:hidden").
				Add(wgt.options[wgt.optionsKeys[0]]))
		} else {
			divTag.Add(HTMLUnsafe("&#8203;")) // Zero-width space
		}
	}

	// All options
	if !wgt.Disabled() {
		ulTag := Tag("ul")
		for _, v := range wgt.optionsKeys {
			ulTag.Add(
				Tag("li").Attr("value", v).
					Attr("onclick", "richdropdown_optionClick(event)").
					Add(wgt.options[v]))
		}
		divTag.Add(ulTag)
	}

	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(
			hiddenTag,
			divTag,
			Tag("script").Add(factory.HTMLUnsafe("input_initBackgroundIconColor('", randomID, "')")),
		).
		When(wgt.Shown(r)).
		Draw(w, r)
}
