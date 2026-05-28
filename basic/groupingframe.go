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

var _ = Widget(&GroupingFrameWidget{}) // Ensure interface

// GroupingFrameWidget renders a border around a grouping of widgets.
type GroupingFrameWidget struct {
	*widget.WidgetBase[*GroupingFrameWidget]
	children []Widget
	title    string
}

// GroupingFrame creates a new widget that renders a border around a grouping of widgets.
// The title is not shown if the grouping frame is nested right inside a tab switcher.
func (f BasicFactory) GroupingFrame(title string) *GroupingFrameWidget {
	x := &GroupingFrameWidget{
		title: title,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *GroupingFrameWidget) Add(children ...any) *GroupingFrameWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *GroupingFrameWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *GroupingFrameWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Class("GroupingFrame", "Block").
		Add(
			Tag("span").Add(wgt.title).Hide(wgt.title == ""),
			Tag("div").Add(wgt.children)).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
