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

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&RatingWidget{}) // Ensure interface

// RatingWidget renders a rating input widget.
type RatingWidget struct {
	*widget.InputWidgetBase[*RatingWidget]
	value     string
	sentiment bool
}

// RatingStars creates a new widget that renders a 1-to-5 star rating
// input. valueStars is the initial selection (clamped to 0..5); pass 0
// for "no rating yet". Posted as a decimal integer string.
func (f FormFactory) RatingStars(name string, valueStars int) *RatingWidget {
	if valueStars > 5 {
		valueStars = 5
	}
	if valueStars < 0 {
		valueStars = 0
	}
	x := &RatingWidget{
		value:     strconv.Itoa(valueStars),
		sentiment: false,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// RatingSentiment creates a new widget that renders a 1-to-5 sentiment
// rating (very-dissatisfied → very-satisfied face icons). Same posting
// semantics as RatingStars. Pass 0 for "no rating yet".
func (f FormFactory) RatingSentiment(name string, valueStars int) *RatingWidget {
	if valueStars > 5 {
		valueStars = 5
	}
	if valueStars < 0 {
		valueStars = 0
	}
	x := &RatingWidget{
		value:     strconv.Itoa(valueStars),
		sentiment: true,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// Value returns the value of the field.
func (wgt *RatingWidget) Value(r *http.Request) string {
	value := wgt.value
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) {
		value = state.Get(wgt.Name())
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *RatingWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		return false
	}
	if value != "" {
		v, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		if v < 1 || v > 5 {
			return false
		}
	}
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *RatingWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *RatingWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if value == "0" {
		value = ""
	}

	tag := Tag("span").
		Class("Rating").
		Attr("data-id", wgt.ID()).
		AttrIf(wgt.Disabled(), "disabled", "1").
		ClassIf(invalid && !wgt.Disabled(), "Invalid")
	// Sentiment
	sentimentIcons := [5]string{
		"sentiment very dissatisfied",
		"sentiment dissatisfied",
		"sentiment neutral",
		"sentiment satisfied",
		"sentiment very satisfied",
	}
	for i := 0; i < 5; i++ {
		full := ""
		if valueInt == i+1 {
			full = "Full"
		}
		if !wgt.sentiment && valueInt >= i+1 {
			full = "Full"
		}
		starTag := Tag("span").
			Class("Star", full).
			Attr("value", strconv.Itoa(i+1))
		if wgt.sentiment {
			starTag.Add(factory.Icon(sentimentIcons[i]))
		} else {
			starTag.Add(factory.Icon("star"))
		}
		if !wgt.Disabled() {
			starTag.Attr("tabindex", "0").
				Attr("onclick", "starrating_click(event)").
				Attr("onkeydown", "starrating_keydown(event)")
		}
		tag.Add(starTag)
		if wgt.sentiment {
			tag.Attr("data-style", "sentiment")
		}
	}
	if !wgt.Disabled() {
		tag.Add(Tag("input").
			Attr("type", "hidden").
			Attr("name", wgt.Name()).
			Attr("value", value).
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("oninput", "input_input(event)")).
			Attr("oninvalid", "input_invalid(event)")
	}
	return tag.
		When(wgt.Shown(r)).
		Draw(w, r)
}
