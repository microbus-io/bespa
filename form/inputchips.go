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
	"strconv"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&InputChipsWidget{})      // Ensure interface
var _ = InputWidget(&InputChipsWidget{}) // Ensure interface

// InputChipsWidget renders a chips input.
type InputChipsWidget struct {
	*widget.InputWidgetBase[*InputChipsWidget]
	dataURL     string
	value       string
	titles      map[string]string
	dedup       bool
	predicates  Predicates
	errMsg      string
	maxItems    int
	maxLen      int
	minLen      int
	placeholder string
	width       string
}

/*
InputChips creates a new widget that renders a chips input.
The data URL is an endpoint that responds with an array of options in response to a query
in the form:

	GET dataURL?q=harry

The response is expected to JSON object with a single JSON array field "options".

	{
		"options": [
			{"title":"Harry Potter", "value":"12345678", "desc":"harry.potter@hogwarts.edu"},
			{"title":"Dirty Harry",  "value":"90909090", "desc":"dirty.harry@sfpd.gov"},
			...
		]
	}

The first 8 options that are not already selected will be displayed in the dropdown.
It is therefore recommended to return more than 8 results.
*/
func (f FormFactory) InputChips(name string, dataURL string) *InputChipsWidget {
	x := &InputChipsWidget{
		dataURL: dataURL,
		titles:  map[string]string{},
		minLen:  -1, // Unbound
		maxLen:  -1, // Unbound
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// AddChip seeds an initial chip. value is what gets posted; title is the
// visible label (defaults to value when empty). Call once per pre-selected
// item before Draw.
func (wgt *InputChipsWidget) AddChip(value string, title string) *InputChipsWidget {
	if wgt.value != "" {
		wgt.value += "\n"
	}
	wgt.value += value
	if title != "" {
		wgt.titles[value] = title
	} else {
		wgt.titles[value] = value
	}
	return wgt
}

// WithDedup drops duplicate values from the posted result. Default is
// false (duplicates allowed). Comparison is by chip value, not title.
func (wgt *InputChipsWidget) WithDedup(dedup bool) *InputChipsWidget {
	wgt.dedup = dedup
	return wgt
}

// WithPredicate checks the field's value against a predicate function.
// Predicates are run against empty fields as well.
func (wgt *InputChipsWidget) WithPredicate(predicate func(value string) (bool, string)) *InputChipsWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

// WithMaxItems caps the number of chips the user can add. Submitting
// more fails validation. 0 (default) means no limit.
func (wgt *InputChipsWidget) WithMaxItems(maxItems int) *InputChipsWidget {
	wgt.maxItems = maxItems
	return wgt
}

// WithLength bounds the length of each individual chip's value, in
// characters. Pass a negative value for either bound to leave it
// unbounded. A non-zero minimum does not imply Required.
func (wgt *InputChipsWidget) WithLength(minChars int, maxChars int) *InputChipsWidget {
	wgt.minLen = minChars
	wgt.maxLen = maxChars
	return wgt
}

// WithPlaceholder sets the placeholder text of the field.
func (wgt *InputChipsWidget) WithPlaceholder(placeholder string) *InputChipsWidget {
	wgt.placeholder = placeholder
	return wgt
}

// WithWidth sets the visible width. Pass any CSS length, e.g. "16ch",
// "200px" or "100%". Empty lets it fill the container — the default.
func (wgt *InputChipsWidget) WithWidth(css string) *InputChipsWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// Value returns the chips' values concatenated with "\n" separators.
// Split on "\n" to recover the individual chip values. Empty when no
// chips are present.
func (wgt *InputChipsWidget) Value(r *http.Request) string {
	dedup := func(v string) string {
		var vv []string
		m := map[string]bool{}
		for _, i := range strings.Split(v, "\n") {
			if !m[i] {
				vv = append(vv, i)
				m[i] = true
			}
		}
		return strings.Join(vv, "\n")
	}
	value := wgt.value
	if wgt.dedup {
		value = dedup(value)
	}
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) { // || wgt.Submitted(r)
		value = state.Get(wgt.Name())
		if state.Has(wgt.Name() + "_title") {
			title := state.Get(wgt.Name() + "_title")
			titles := strings.Split(title, "\n")
			for i, v := range strings.Split(value, "\n") {
				if i < len(titles) { // Should always be true
					wgt.titles[v] = titles[i]
				}
			}
		}
		if wgt.dedup {
			value = dedup(value)
		}
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *InputChipsWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		wgt.errMsg = "A value is required"
		return false
	}
	if wgt.maxItems > 0 && len(strings.Split(value, "\n")) > wgt.maxItems {
		wgt.errMsg = "Too many items"
		return false
	}
	// Length
	for _, v := range strings.Split(value, "\n") {
		if wgt.maxLen >= 0 && len([]rune(v)) > wgt.maxLen {
			return false
		}
		if wgt.minLen >= 0 && len([]rune(v)) < wgt.minLen {
			return false
		}
	}
	// Predicates
	if ok, errMsg := wgt.predicates.Validate(value); !ok {
		wgt.errMsg = errMsg
		return false
	}
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *InputChipsWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *InputChipsWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)
	title := ""

	chipTags := []widget.Widget{}
	if value != "" {
		for _, v := range strings.Split(value, "\n") {
			removeTag := Tag("span").
				Class("Remove").
				Add(factory.Icon("close"))
			titleTag := Tag("span").Add(wgt.titles[v])
			chipTag := Tag("span").
				Class("Chip").
				Attr("data-value", v).
				Add(titleTag, removeTag)
			if !wgt.Disabled() {
				chipTag.
					Attr("tabindex", "0").
					Attr("onkeydown", "inputchips_keydownRemove(event)")
				removeTag.
					Attr("onclick", "inputchips_remove(event)")
			}
			chipTags = append(chipTags, chipTag)
			if title != "" {
				title += "\n"
			}
			title += wgt.titles[v]
		}
	}
	blankChipTag := Tag("")
	if !wgt.Disabled() {
		blankChipTag = Tag("span").
			Class("BlankChip").
			Attr("tabindex", "0").
			Attr("onkeydown", "inputchips_keydownRemove(event)").
			Add(
				Tag("span").
					Add("Blank"),
				Tag("span").
					Class("Remove").
					Attr("onclick", "inputchips_remove(event)").
					Add(factory.Icon("close")),
			)
	}

	inputTag := Tag("input").
		Attr("type", "text").
		Attr("placeholder", wgt.placeholder).
		AttrIf(wgt.Disabled(), "disabled", "1").
		ClassIf(wgt.maxItems > 0 && len(chipTags) >= wgt.maxItems, "Saturated")
	popupTag := Tag("")
	hiddenValueTag := Tag("")
	hiddenTitleTag := Tag("")
	errTag := Tag("")
	if !wgt.Disabled() {
		randomID := widget.RandomAlphaNumID(8)
		inputTag.
			Attr("id", randomID).
			Attr("oninput", "inputchips_input(event)").
			Attr("onkeydown", "inputchips_keydown(event)").
			Attr("tabindex", "0").
			AttrIf(wgt.minLen >= 0, "minlength", strconv.Itoa(wgt.minLen)).
			AttrIf(wgt.maxLen >= 0, "maxlength", strconv.Itoa(wgt.maxLen))

		popupTag = Tag("ul").
			Attr("onclick", "inputchips_click(event)").
			Attr("onmouseenter", "inputchips_mouseenter(event)").
			Attr("onmouseleave", "inputchips_mouseleave(event)")

		hiddenValueTag = Tag("input").
			Attr("type", "hidden").
			Attr("name", wgt.Name()).
			Attr("oninput", "input_input(event)").
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("value", value)
		hiddenTitleTag = Tag("input").
			Attr("type", "hidden").
			Attr("name", wgt.Name()+"_title").
			Attr("value", title)

		if invalid && wgt.errMsg != "" {
			errTag = customValidityScript(randomID, wgt.errMsg)
		}
	}

	return Tag("div").
		Class("InputChips").
		Attr("data-id", wgt.ID()).
		Attr("data-url", wgt.dataURL).
		AttrIf(wgt.maxItems > 0, "data-maxitems", strconv.Itoa(wgt.maxItems)).
		Attr("onfocusout", "inputchips_focusout(event)").
		Attr("onfocusin", "inputchips_focusin(event)").
		ClassIf(wgt.Disabled(), "Disabled").
		ClassIf(invalid, "Invalid").
		Style(wgt.width).
		Add(hiddenValueTag, hiddenTitleTag, blankChipTag, chipTags, inputTag, popupTag, errTag).
		When(wgt.Shown(r)).
		Draw(w, r)
}
