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
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&RadioWidget{})      // Ensure interface
var _ = InputWidget(&RadioWidget{}) // Ensure interface

// RadioWidget renders a radio button.
type RadioWidget struct {
	*widget.InputWidgetBase[*RadioWidget]
	value       string
	optionsKeys []string
	options     map[string]string
	appearance  string
	predicates  Predicates
	errMsg      string
}

// Radio creates a new widget that renders a group of radio buttons.
// name is the state variable; value is the option pre-selected on first
// render. Add choices via AddOption. Defaults to a vertical layout.
// A submitted value that matches none of the options fails validation.
func (f FormFactory) Radio(name string, value string) *RadioWidget {
	x := &RadioWidget{
		value:      value,
		options:    map[string]string{},
		appearance: "Vertical",
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// WithVertical stacks the radio buttons in a column. This is the default.
func (wgt *RadioWidget) WithVertical() *RadioWidget {
	wgt.appearance = "Vertical"
	return wgt
}

// WithHorizontal lays the radio buttons out in a row.
func (wgt *RadioWidget) WithHorizontal() *RadioWidget {
	wgt.appearance = "Horizontal"
	return wgt
}

// AddOption appends a choice. value is what gets posted when selected;
// caption is the label shown to the user. Options render in insertion
// order.
func (wgt *RadioWidget) AddOption(value string, caption string) *RadioWidget {
	wgt.optionsKeys = append(wgt.optionsKeys, value) // Keep the order of insertion
	wgt.options[value] = caption
	return wgt
}

// WithPredicate adds a custom validator. value is the selected option's
// value, or "" if nothing is selected.
func (wgt *RadioWidget) WithPredicate(predicate func(value string) (valid bool, errMsg string)) *RadioWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

// Value returns the value of the field.
func (wgt *RadioWidget) Value(r *http.Request) string {
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
func (wgt *RadioWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		return false
	}
	// Predicates
	if ok, errMsg := wgt.predicates.Validate(value); !ok {
		wgt.errMsg = errMsg
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
func (wgt *RadioWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *RadioWidget) Draw(w io.Writer, r *http.Request) (err error) {
	var randomIDs []string
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)
	allTag := Tag("span").
		Class("RadioStrip", wgt.appearance).
		Attr("data-id", wgt.ID())
	for _, option := range wgt.optionsKeys {
		randomID := widget.RandomAlphaNumID(8)
		randomIDs = append(randomIDs, randomID)
		radioTag := Tag("input").
			Attr("type", "radio").
			Attr("value", option).
			Attr("id", randomID).
			ClassIf(invalid, "Invalid").
			AttrIf(option == value, "checked", "1").
			AttrIf(wgt.Disabled(), "disabled", "1")
		if !wgt.Disabled() {
			radioTag.Attr("name", wgt.Name()).
				Attr("tabindex", "0").
				AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
				Attr("oninput", "input_input(event)").
				Attr("oninvalid", "input_invalid(event)").
				AttrIf(wgt.Required(), "required", "1")
		}
		var labelTag *widget.TagWidget
		label := wgt.options[option]
		if label != "" {
			id4Label := widget.RandomAlphaNumID(8)
			radioTag.Attr("id", id4Label)
			labelTag = Tag("label").
				Attr("for", id4Label).
				Add(label)
		}
		allTag.Add(Tag("span").Class("Radio").Add(radioTag, labelTag))
	}
	errTag := Tag("")
	if invalid && wgt.errMsg != "" {
		errTag = customValidityScript(strings.Join(randomIDs, " "), wgt.errMsg)
	}
	allTag.Add(errTag)
	return allTag.
		When(wgt.Shown(r)).
		Draw(w, r)
}
