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

var _ = Widget(&TextAlignWidget{}) // Ensure interface

// TextAlignWidget aligns text content via CSS text-align. For aligning
// block-level children (cards, buttons, etc.) horizontally as a group,
// use AlignLeft/Center/Right instead.
type TextAlignWidget struct {
	*widget.WidgetBase[*TextAlignWidget]
	direction string
	children  []Widget
}

// textAlign creates a new widget that renders a collection of widgets aligned to the left, right or center.
func (f BasicFactory) textAlign(direction string, children ...any) *TextAlignWidget {
	x := &TextAlignWidget{
		direction: direction,
		children:  Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// TextAlignLeft creates a new widget that renders a collection of widgets aligned to the left.
func (f BasicFactory) TextAlignLeft(children ...any) *TextAlignWidget {
	return f.textAlign("Left", children...)
}

// TextAlignRight creates a new widget that renders a collection of widgets aligned to the right.
func (f BasicFactory) TextAlignRight(children ...any) *TextAlignWidget {
	return f.textAlign("Right", children...)
}

// TextAlignCenter creates a new widget that renders a collection of widgets center-aligned.
func (f BasicFactory) TextAlignCenter(children ...any) *TextAlignWidget {
	return f.textAlign("Center", children...)
}

// Add adds nested widgets.
func (wgt *TextAlignWidget) Add(children ...any) *TextAlignWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *TextAlignWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *TextAlignWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Class("TextAlign", wgt.direction).
		Attr("data-id", wgt.ID()).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
