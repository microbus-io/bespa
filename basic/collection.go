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

var _ = Widget(&CollectionWidget{}) // Ensure interface

// CollectionWidget renders a collection of widgets.
type CollectionWidget struct {
	*widget.WidgetBase[*CollectionWidget]
	children []Widget
	block    string
	tag      string
	width    string
	maxWidth string
}

// Collection creates a new widget that wraps a group of widgets in an inline
// <span>. Use it to bundle siblings under a single redraw boundary or apply
// shared width/visibility. For a block-level wrapper with vertical rhythm,
// use Block.
func (f BasicFactory) Collection(children ...any) *CollectionWidget {
	x := &CollectionWidget{
		children: Many(children),
		tag:      "span",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Block creates a new widget that wraps children in a block-level <div>
// with the standard bottom margin used for vertical rhythm between sections.
func (f BasicFactory) Block(children ...any) *CollectionWidget {
	x := f.Collection(children...)
	x.block = "Block"
	x.tag = "div"
	return x
}

// WithWidth sets an explicit width. Empty clears it.
// Pass any CSS length, e.g. "400px", "100%" or "calc(100vw - 2em)".
func (wgt *CollectionWidget) WithWidth(css string) *CollectionWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithMaxWidth caps the width while allowing the content to be narrower.
// Empty clears it.
// Pass any CSS length, e.g. "400px", "80%" or "calc(100vw - 2em)".
func (wgt *CollectionWidget) WithMaxWidth(css string) *CollectionWidget {
	if css != "" {
		wgt.maxWidth = "max-width:" + css
	} else {
		wgt.maxWidth = ""
	}
	return wgt
}

// Add adds nested widgets.
func (wgt *CollectionWidget) Add(children ...any) *CollectionWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *CollectionWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *CollectionWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag(wgt.tag).
		Attr("data-id", wgt.ID()).
		Style(wgt.width, wgt.maxWidth).
		Class(wgt.block).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
