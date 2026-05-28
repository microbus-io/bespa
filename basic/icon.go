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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&IconWidget{}) // Ensure interface

// IconWidget renders an icon, either from the Material Symbols font or
// from inline SVG markup.
type IconWidget struct {
	*widget.WidgetBase[*IconWidget]
	altText string
	symbol  string
	fill    string
	svg     string
	size    string
}

/*
Icon creates a new widget that renders an icon.
The spec can be either a material symbol name or valid <svg></svg> markup.

Material symbol names are listed at https://fonts.google.com/icons?icon.set=Material+Symbols.
Material symbols are a font and are therefore sized according to the current font size.
See https://developers.google.com/fonts/docs/material_symbols for the developer's guide.

The SVG should not specify its own width or height so that
it is rendered using the height of the current font size (i.e. 1em).
The SVG should use the "currentColor" for fill and stroke colors to best fit
with the color palette of the page.
*/
func (f BasicFactory) Icon(spec string) *IconWidget {
	x := &IconWidget{}
	spec = strings.TrimSpace(spec)
	if strings.HasPrefix(spec, "<svg ") && strings.HasSuffix(spec, "</svg>") {
		x.svg = spec
	} else {
		if spec == "" {
			spec = "\u00a0" // non-breaking space
		}
		spec = strings.ReplaceAll(strings.ToLower(spec), " ", "_")
		x.symbol = spec
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithAltText labels the icon for screen readers and as the hover tooltip.
// Without alt text the icon is marked decorative (`aria-hidden`), which is
// usually right when it sits inside an already-labelled control.
func (wgt *IconWidget) WithAltText(altText string) *IconWidget {
	wgt.altText = altText
	return wgt
}

// WithFill switches a Material Symbol from its default outlined style to
// its filled variant. Has no effect on SVG icons.
func (wgt *IconWidget) WithFill(fill bool) *IconWidget {
	if fill {
		wgt.fill = "font-variation-settings:'FILL' 1"
	} else {
		wgt.fill = ""
	}
	return wgt
}

// WithSizeMultiplier scales the icon relative to the surrounding font size
// (1.0 = match, 2.0 = double). Values <= 0 are ignored.
func (wgt *IconWidget) WithSizeMultiplier(size float32) *IconWidget {
	if size > 0 {
		wgt.size = fmt.Sprintf("font-size:%.2fem", size)
	} else {
		wgt.size = ""
	}
	return wgt
}

// Draw renders the widget's HTML.
//
// Accessibility: when alt text is provided, it is emitted as BOTH `title`
// (browser hover tooltip) and `aria-label` (screen-reader name). When no
// alt text is provided the icon is treated as decorative and hidden from
// assistive tech via `aria-hidden="true"` — appropriate when the icon
// sits inside an already-named control like a labelled button.
func (wgt *IconWidget) Draw(w io.Writer, r *http.Request) (err error) {
	hasLabel := wgt.altText != ""
	if wgt.symbol != "" {
		return Tag("i").
			Class("Icon", "material-symbols-outlined").
			Attr("data-id", wgt.ID()).
			AttrIf(hasLabel, "title", wgt.altText).
			AttrIf(hasLabel, "aria-label", wgt.altText).
			AttrIf(hasLabel, "role", "img").
			AttrIf(!hasLabel, "aria-hidden", "true").
			Style(wgt.fill, wgt.size).
			Add(wgt.symbol).
			When(wgt.Shown(r)).
			Draw(w, r)
	}
	if wgt.svg != "" {
		svg := strings.TrimSpace(wgt.svg)
		if !strings.HasPrefix(svg, "<svg ") || !strings.HasSuffix(svg, "</svg>") {
			return errors.New("invalid SVG markup")
		}
		return Tag("span").
			Attr("data-id", wgt.ID()).
			Class("Icon").
			AttrIf(hasLabel, "title", wgt.altText).
			AttrIf(hasLabel, "aria-label", wgt.altText).
			AttrIf(hasLabel, "role", "img").
			AttrIf(!hasLabel, "aria-hidden", "true").
			Style(wgt.size).
			Add(HTMLUnsafe(`<svg class="SVG" `), HTMLUnsafe(svg[5:])).
			When(wgt.Shown(r)).
			Draw(w, r)
	}
	return Tag("span").
		Attr("data-id", wgt.ID()).
		When(wgt.Shown(r)).
		Draw(w, r)
}
