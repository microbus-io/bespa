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

package form

import (
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&FieldLabelWidget{}) // Ensure interface

// FieldLabelWidget renders a field label.
type FieldLabelWidget struct {
	*widget.WidgetBase[*FieldLabelWidget]
	children []Widget
}

// FieldLabel creates a new widget for a styled form-label block.
// Field already wraps its left content in a FieldLabel — call this
// directly only when you're laying out a form by hand.
func (f FormFactory) FieldLabel(children ...any) *FieldLabelWidget {
	x := &FieldLabelWidget{
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *FieldLabelWidget) Add(children ...any) *FieldLabelWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *FieldLabelWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *FieldLabelWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Class("FieldLabel").
		Attr("data-id", wgt.ID()).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
