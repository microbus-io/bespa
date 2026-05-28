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
	"fmt"
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&TextStyleWidget{}) // Ensure interface

// TextStyleWidget alters the appearance of text.
type TextStyleWidget struct {
	*widget.WidgetBase[*TextStyleWidget]
	size     string
	color    string
	weight   string
	slant    string
	mono     string
	apply    bool
	children []Widget
}

// TextStyle creates a new widget that wraps children and lets you chain
// styling helpers (WithBold, WithItalic, WithColorPrimary, …). Each axis
// — color, weight, slant, size, monospace — is independent, but later
// calls on the same axis overwrite earlier ones (the last WithColor*
// wins, etc.). Color names come from the Material theme.
func (f BasicFactory) TextStyle(children ...any) *TextStyleWidget {
	x := &TextStyleWidget{
		children: Many(children),
		apply:    true,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// TextLightweight creates a new widget that alters the weight of the font to lightweight.
func (f BasicFactory) TextLightweight(children ...any) *TextStyleWidget {
	x := f.TextStyle(children...)
	x.WithLightweight()
	return x
}

// TextBold creates a new widget that alters the weight of the font to bold.
func (f BasicFactory) TextBold(children ...any) *TextStyleWidget {
	x := f.TextStyle(children...)
	x.WithBold()
	return x
}

// ApplyIf gates all styling on a condition: when false, the wrapper renders
// its children with no styling applied at all (size, color, weight, etc.
// are dropped). Useful for "bold the active row" patterns.
func (wgt *TextStyleWidget) ApplyIf(condition bool) *TextStyleWidget {
	wgt.apply = condition
	return wgt
}

// WithSize sets the size of the font.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
func (wgt *TextStyleWidget) WithSize(size float32, unit string) *TextStyleWidget {
	if size > 0 {
		wgt.size = fmt.Sprintf("font-size: %.2f%s; line-height:1", size, unit)
	} else {
		wgt.size = ""
	}
	return wgt
}

// WithSizeMultiplier is shorthand for setting the size with em unit.
func (wgt *TextStyleWidget) WithSizeMultiplier(multiplier float32) *TextStyleWidget {
	wgt.WithSize(multiplier, "em")
	return wgt
}

// WithColorPrimary changes the color of the text to the primary color.
func (wgt *TextStyleWidget) WithColorPrimary() *TextStyleWidget {
	wgt.color = "TextColorPrimary"
	return wgt
}

// WithColorOnPrimary sets text and background to read legibly when placed
// over the primary color — i.e. the "on-primary" pairing in M3.
func (wgt *TextStyleWidget) WithColorOnPrimary() *TextStyleWidget {
	wgt.color = "TextColorOnPrimary"
	return wgt
}

// WithColorSecondary changes the color of the text to the secondary color.
func (wgt *TextStyleWidget) WithColorSecondary() *TextStyleWidget {
	wgt.color = "TextColorSecondary"
	return wgt
}

// WithColorOnSecondary is the legible pairing for content rendered on top
// of the secondary color.
func (wgt *TextStyleWidget) WithColorOnSecondary() *TextStyleWidget {
	wgt.color = "TextColorOnSecondary"
	return wgt
}

// WithColorTertiary changes the color of the text to the tertiary color.
func (wgt *TextStyleWidget) WithColorTertiary() *TextStyleWidget {
	wgt.color = "TextColorTertiary"
	return wgt
}

// WithColorOnTertiary is the legible pairing for content rendered on top
// of the tertiary color.
func (wgt *TextStyleWidget) WithColorOnTertiary() *TextStyleWidget {
	wgt.color = "TextColorOnTertiary"
	return wgt
}

// WithColorDefault changes the color of the text to the default color.
func (wgt *TextStyleWidget) WithColorDefault() *TextStyleWidget {
	wgt.color = "TextColorDefault"
	return wgt
}

// WithColorInverse swaps to the inverse surface colors — used for content
// that should stand out against the normal page background.
func (wgt *TextStyleWidget) WithColorInverse() *TextStyleWidget {
	wgt.color = "TextColorInverse"
	return wgt
}

// WithColorError changes the color of the text to the error color.
func (wgt *TextStyleWidget) WithColorError() *TextStyleWidget {
	wgt.color = "TextColorError"
	return wgt
}

// WithColorOnError is the legible pairing for content rendered on top of
// the error color.
func (wgt *TextStyleWidget) WithColorOnError() *TextStyleWidget {
	wgt.color = "TextColorOnError"
	return wgt
}

// WithColorOK changes the text color to the OK / success color (typically green).
func (wgt *TextStyleWidget) WithColorOK() *TextStyleWidget {
	wgt.color = "TextColorOK"
	return wgt
}

// WithColorOnOK is the legible pairing for content rendered on top of the
// OK / success color.
func (wgt *TextStyleWidget) WithColorOnOK() *TextStyleWidget {
	wgt.color = "TextColorOnOK"
	return wgt
}

// WithColorDeemphasized changes the color of the text to the deemphasized color.
func (wgt *TextStyleWidget) WithColorDeemphasized() *TextStyleWidget {
	wgt.color = "TextColorDeemphasized"
	return wgt
}

// WithColorDisabled changes the color of the text to the disabled text color.
func (wgt *TextStyleWidget) WithColorDisabled() *TextStyleWidget {
	wgt.color = "TextColorDisabled"
	return wgt
}

// WithBold changes the weight of the font to bold.
func (wgt *TextStyleWidget) WithBold() *TextStyleWidget {
	wgt.weight = "TextBold"
	return wgt
}

// WithLightweight changes the weight of the font to lightweight.
func (wgt *TextStyleWidget) WithLightweight() *TextStyleWidget {
	wgt.weight = "TextLightweight"
	return wgt
}

// WithItalic changes the slant of the font to italic.
func (wgt *TextStyleWidget) WithItalic() *TextStyleWidget {
	wgt.slant = "TextItalic"
	return wgt
}

// WithMonospace changes the font to a monospace font.
func (wgt *TextStyleWidget) WithMonospace() *TextStyleWidget {
	wgt.mono = "TextMonospace"
	return wgt
}

// Add adds nested widgets.
func (wgt *TextStyleWidget) Add(children ...any) *TextStyleWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *TextStyleWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *TextStyleWidget) Draw(w io.Writer, r *http.Request) (err error) {
	color := wgt.color
	weight := wgt.weight
	slant := wgt.slant
	size := wgt.size
	mono := wgt.mono
	if !wgt.apply {
		color = ""
		weight = ""
		slant = ""
		size = ""
		mono = ""
	}
	return Tag("span").
		Class(color, weight, slant, mono).
		Attr("data-id", wgt.ID()).
		Style(size).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
