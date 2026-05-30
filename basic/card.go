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

var _ = Widget(&CardWidget{}) // Ensure interface

// CardWidget renders a card.
type CardWidget struct {
	*widget.WidgetBase[*CardWidget]
	href      string
	target    string
	disabled  bool
	style     string
	minHeight string
	children  []Widget
}

// CardElevated creates a new widget that renders a Material elevated card
// (the default raised-surface style).
// A leading or trailing BannerImage child is rendered edge-to-edge above or
// below the card body.
func (f BasicFactory) CardElevated() *CardWidget {
	x := &CardWidget{
		style: "Elevated",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// CardOutlined creates a new widget that renders a Material outlined card.
// A leading or trailing BannerImage child is rendered edge-to-edge above or
// below the card body.
func (f BasicFactory) CardOutlined() *CardWidget {
	x := &CardWidget{
		style: "Outlined",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// CardFilled creates a new widget that renders a Material filled card.
// A leading or trailing BannerImage child is rendered edge-to-edge above or
// below the card body.
func (f BasicFactory) CardFilled() *CardWidget {
	x := &CardWidget{
		style: "Filled",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithHref makes the whole card clickable, navigating to href on click.
// Accepts the full action-URL grammar (`?key=`, `^?…`, `/path`, etc.).
func (wgt *CardWidget) WithHref(href string) *CardWidget {
	wgt.href = href
	return wgt
}

// WithTarget sets the HTML target for the card's link. Defaults to the
// page's `_target` state variable when unset.
func (wgt *CardWidget) WithTarget(target string) *CardWidget {
	wgt.target = target
	return wgt
}

// WithDisabled greys out the card and prevents the click from firing.
func (wgt *CardWidget) WithDisabled(disabled bool) *CardWidget {
	wgt.disabled = disabled
	return wgt
}

// WithMinHeight sets the card's minimum height. Default is 240px; empty
// sizes to content instead.
// Pass any CSS length, e.g. "240px", "50%" or "calc(100vh - 50px)".
func (wgt *CardWidget) WithMinHeight(css string) *CardWidget {
	if css == "" {
		wgt.minHeight = ""
		return wgt
	}
	wgt.minHeight = "min-height:" + css
	return wgt
}

// Add adds nested widgets.
func (wgt *CardWidget) Add(children ...any) *CardWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *CardWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *CardWidget) Draw(w io.Writer, r *http.Request) (err error) {
	var firstImage *BannerImageWidget
	var lastImage *BannerImageWidget
	children := wgt.children
	if len(children) > 0 {
		if img, ok := children[0].(*BannerImageWidget); ok {
			firstImage = img
			children = children[1:]
		}
	}
	if len(children) > 0 {
		n := len(children)
		if img, ok := children[n-1].(*BannerImageWidget); ok {
			lastImage = img
			children = children[:n-1]
		}
	}

	linkTag := Tag("")
	if wgt.href != "" {
		state := factory.StateOf(r)
		target := wgt.target
		if target == "" {
			target = state.Get("_target")
		}
		linkTag = Tag("a").
			Attr("href", wgt.href).
			Attr("target", target).
			Hide(true)
	}
	clickable := !wgt.disabled && wgt.href != ""
	boxTag := Tag("div").
		AttrIf(clickable, "tabindex", "0").
		AttrIf(clickable, "onclick", "card_click(event)").
		AttrIf(clickable, "onkeydown", "card_keydown(event)").
		ClassIf(wgt.disabled, "Disabled", "1").
		Style(wgt.minHeight).
		Add(firstImage, Tag("div").Add(children), lastImage)

	return Tag("div").
		Class("Card", "Block", wgt.style).
		Attr("data-id", wgt.ID()).
		Add(boxTag, linkTag).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
