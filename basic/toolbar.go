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

var _ = Widget(&ToolbarWidget{}) // Ensure interface

// ToolbarWidget renders a toolbar.
type ToolbarWidget struct {
	*widget.WidgetBase[*ToolbarWidget]
	leftChildren  []Widget
	rightChildren []Widget
	align         string
	wrap          string
}

// Toolbar creates a new widget that renders a horizontal toolbar with two
// item groups — Add{Left,Right} populate them, and the right group is
// flush-right. Defaults: vertically centered, wraps onto multiple rows
// on narrow viewports.
func (f BasicFactory) Toolbar() *ToolbarWidget {
	x := &ToolbarWidget{
		align: "AlignCenter",
		wrap:  "Wrap",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithWrap controls whether the toolbar wraps onto multiple rows on narrow
// viewports. When wrapped, the left group appears above the right group.
// Default is true.
func (wgt *ToolbarWidget) WithWrap(wrap bool) *ToolbarWidget {
	if !wrap {
		wgt.wrap = ""
	} else {
		wgt.wrap = "Wrap"
	}
	return wgt
}

// AddLeft adds nested widgets aligned to the left.
func (wgt *ToolbarWidget) AddLeft(leftChildren ...any) *ToolbarWidget {
	wgt.leftChildren = Many(wgt.leftChildren, leftChildren)
	return wgt
}

// AddRight adds nested widgets aligned to the right.
func (wgt *ToolbarWidget) AddRight(rightChildren ...any) *ToolbarWidget {
	wgt.rightChildren = Many(wgt.rightChildren, rightChildren)
	return wgt
}

// WithAlignCenter centers the toolbar widgets vertically.
// This is the default behavior.
func (wgt *ToolbarWidget) WithAlignCenter() *ToolbarWidget {
	wgt.align = "AlignCenter"
	return wgt
}

// WithAlignBottom aligns the toolbar widgets to the bottom of the toolbar.
func (wgt *ToolbarWidget) WithAlignBottom() *ToolbarWidget {
	wgt.align = "AlignBottom"
	return wgt
}

// WithAlignTop aligns the toolbar widgets to the top.
func (wgt *ToolbarWidget) WithAlignTop() *ToolbarWidget {
	wgt.align = "AlignTop"
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *ToolbarWidget) Children() []Widget {
	return Many(wgt.leftChildren, wgt.rightChildren)
}

// Draw renders the widget's HTML.
func (wgt *ToolbarWidget) Draw(w io.Writer, r *http.Request) (err error) {
	left := Tag("div")
	for _, c := range wgt.leftChildren {
		left.Add(Tag("div").Add(c))
	}
	right := Tag("div")
	for _, c := range wgt.rightChildren {
		right.Add(Tag("div").Add(c))
	}
	return Tag("div").
		Class("Toolbar", "Block", wgt.align, wgt.wrap).
		Attr("data-id", wgt.ID()).
		Add(
			left.Hide(len(wgt.leftChildren) == 0),
			right,
		).
		When(wgt.Shown(r) && len(wgt.leftChildren)+len(wgt.rightChildren) > 0).
		Draw(w, r)
}
