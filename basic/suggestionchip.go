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

package basic

import (
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&SuggestionChipWidget{}) // Ensure interface

// SuggestionChipWidget renders a suggestion chip.
type SuggestionChipWidget struct {
	*widget.WidgetBase[*SuggestionChipWidget]
	href     string
	disabled bool
	target   string
	children []Widget
}

// SuggestionChip creates a new widget that renders a clickable Material
// suggestion chip. href accepts the full action-URL grammar. Group multiple
// chips inside a Gallery so they wrap nicely.
func (f BasicFactory) SuggestionChip(href string) *SuggestionChipWidget {
	x := &SuggestionChipWidget{
		href: href,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithHRef sets the URL of the suggestion chip.
func (wgt *SuggestionChipWidget) WithHref(href string) *SuggestionChipWidget {
	wgt.href = href
	return wgt
}

// WithTarget sets the target of the suggestion chip.
func (wgt *SuggestionChipWidget) WithTarget(target string) *SuggestionChipWidget {
	wgt.target = target
	return wgt
}

// WithDisabled disables the suggestion chip.
func (wgt *SuggestionChipWidget) WithDisabled(disabled bool) *SuggestionChipWidget {
	wgt.disabled = disabled
	return wgt
}

// Add adds nested widgets.
func (wgt *SuggestionChipWidget) Add(children ...any) *SuggestionChipWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *SuggestionChipWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *SuggestionChipWidget) Draw(w io.Writer, r *http.Request) (err error) {
	state := factory.StateOf(r)
	href := wgt.href
	target := wgt.target
	if target == "" {
		target = state.Get("_target")
	}
	tagName := "a"
	if wgt.disabled || href == "" {
		tagName = "span"
	}
	return Tag(tagName).
		Class("SuggestionChip").
		ClassIf(wgt.disabled, "Disabled").
		Attr("data-id", wgt.ID()).
		Attr("href", href).
		Attr("target", target).
		AttrIf(tagName == "a", "tabindex", "0").
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
