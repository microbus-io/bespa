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

var _ = Widget(&AlignWidget{}) // Ensure interface

// AlignWidget renders a collection of widgets aligned to the left, right or center.
type AlignWidget struct {
	*widget.WidgetBase[*AlignWidget]
	direction string
	children  []Widget
}

// align creates a new widget that renders a collection of widgets aligned to the left, right or center.
func (f BasicFactory) align(direction string, children ...any) *AlignWidget {
	x := &AlignWidget{
		direction: direction,
		children:  Many(children...),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// AlignLeft creates a new widget that renders a collection of widgets aligned to the left.
func (f BasicFactory) AlignLeft(children ...any) *AlignWidget {
	return f.align("Left", children)
}

// AlignRight creates a new widget that renders a collection of widgets aligned to the right.
func (f BasicFactory) AlignRight(children ...any) *AlignWidget {
	return f.align("Right", children)
}

// AlignCenter creates a new widget that renders a collection of widgets center-aligned.
func (f BasicFactory) AlignCenter(children ...any) *AlignWidget {
	return f.align("Center", children)
}

// Add adds nested widgets.
func (wgt *AlignWidget) Add(children ...any) *AlignWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *AlignWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *AlignWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Class("Align", wgt.direction).
		Attr("data-id", wgt.ID()).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
