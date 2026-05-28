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

var _ = Widget(&IfWidget{}) // Ensure interface

// IfWidget renders nested widgets based on a condition.
type IfWidget struct {
	*widget.WidgetBase[*IfWidget]
	children []Widget
	cond     bool
}

// If creates a new widget that picks between Then/Else branches at build
// time based on cond. Note this is evaluated when the page is constructed,
// not on each redraw — for visibility that reacts to state changes use
// HideIf / HideIfEq / HideIfEmpty instead.
func (f BasicFactory) If(cond bool) *IfWidget {
	x := &IfWidget{
		cond: cond,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Then adds nested widgets to be rendered when the condition is true.
func (wgt *IfWidget) Then(children ...any) *IfWidget {
	if wgt.cond {
		wgt.children = Many(wgt.children, children)
	}
	return wgt
}

// Else adds nested widgets to be rendered when the condition is false.
func (wgt *IfWidget) Else(children ...any) *IfWidget {
	if !wgt.cond {
		wgt.children = Many(wgt.children, children)
	}
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *IfWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *IfWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
