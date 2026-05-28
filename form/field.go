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

var _ = Widget(&FieldWidget{}) // Ensure interface

// FieldWidget renders a field of a form.
type FieldWidget struct {
	*widget.WidgetBase[*FieldWidget]
	label []Widget
	body  []Widget
}

// Field creates a new widget that lays out one form row with a label on
// the left and the input(s) on the right. Use AddLeft for the label
// content and AddRight for the input. On narrow viewports the label
// stacks above the input.
func (f FormFactory) Field() *FieldWidget {
	x := &FieldWidget{}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// AddRight appends widgets to the field's body — typically the input
// control(s) and any inline hints.
func (wgt *FieldWidget) AddRight(children ...any) *FieldWidget {
	wgt.body = Many(wgt.body, children)
	return wgt
}

// AddLeft appends widgets to the field's label area.
func (wgt *FieldWidget) AddLeft(children ...any) *FieldWidget {
	wgt.label = Many(wgt.label, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *FieldWidget) Children() []Widget {
	allWidgets := []Widget{}
	allWidgets = append(allWidgets, wgt.label...)
	allWidgets = append(allWidgets, wgt.body...)
	return allWidgets
}

// Draw renders the widget's HTML.
func (wgt *FieldWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Class("Field", "Block").
		Attr("data-id", wgt.ID()).
		Add(
			factory.FieldLabel(wgt.label),
			Tag("div").Add(wgt.body)).
		When(wgt.Shown(r)).
		Draw(w, r)
}
