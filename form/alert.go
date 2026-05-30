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
	"net/url"

	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&AlertWidget{}) // Ensure interface

// AlertWidget renders an alert modal.
type AlertWidget struct {
	*widget.WidgetBase[*AlertWidget]
	title   Widget
	form    *FormWidget
	icon    Widget
	name    string
	modal   *basic.ModalWidget
	buttons []*ButtonWidget
}

// Alert creates a new widget that renders a confirmation modal — a
// primary-coloured icon on the left next to a heading and action
// buttons. Like Modal, the alert opens when the named state variable is
// non-empty; the typical pattern is to set it to an ID (e.g. ?delete=42)
// from a list-row action. Pass icon="" for an icon-less alert. For an
// error-styled (red) icon, use AlertError.
func (f FormFactory) Alert(name string, icon string, title string) *AlertWidget {
	var iconWidget widget.Widget
	if icon != "" {
		iconWidget = factory.TextStyle(factory.Icon(icon)).WithColorPrimary()
	}
	x := &AlertWidget{
		name:  name,
		title: Any(title),
		icon:  iconWidget,
		form:  factory.Form().WithName("alert"),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	x.createModal()
	return x
}

// AlertError is the error-styled variant of Alert — the icon is rendered
// in the theme's error color, suitable for destructive confirmations.
func (f FormFactory) AlertError(name string, icon string, title string) *AlertWidget {
	var iconWidget widget.Widget
	if icon != "" {
		iconWidget = factory.TextStyle(factory.Icon(icon)).WithColorError()
	}
	x := &AlertWidget{
		name:  name,
		title: Any(title),
		icon:  iconWidget,
		form:  factory.Form().WithName("alert"),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	x.createModal()
	return x
}

func (wgt *AlertWidget) createModal() {
	var body widget.Widget
	if wgt.icon != nil {
		body = factory.Splitter(1, 4).WithWrap(false).
			AddLeft(
				factory.Block(
					factory.SpacerParagraph(),
					factory.TextAlignCenter(wgt.icon),
				),
			).
			AddRight(
				factory.SpacerBreak(),
				factory.HeadlineSmall(wgt.title),
				wgt.form,
				factory.SpacerNewLine(),
			)
	} else {
		body = factory.Collection(
			factory.SpacerBreak(),
			factory.HeadlineSmall(wgt.title),
			wgt.form,
			factory.SpacerNewLine(),
		)
	}
	wgt.modal = factory.Modal(wgt.name).
		WithWidth("450px").
		WithMinHeight("1px").
		Add(body)
}

// WithWidth sets the width of the alert window.
// Pass any CSS length, e.g. "450px", "90%" or "calc(100vw - 2em)". Empty clears it.
// The default is 450px.
func (wgt *AlertWidget) WithWidth(css string) *AlertWidget {
	wgt.modal.WithWidth(css)
	return wgt
}

// Add appends widgets to the alert's body. Buttons get special handling:
// clicking any button clears the alert's state variable (closing the
// modal), and a named button additionally writes its own state variable
// to the closed alert's prior value — so the handler can dispatch on
// which button was pressed and still see the row ID. Example: alert
// "delete" set to "42"; clicking button "confirm" results in
// `?delete=&confirm=42`.
func (wgt *AlertWidget) Add(children ...any) *AlertWidget {
	for _, w := range children {
		if btn, ok := w.(*ButtonWidget); ok {
			wgt.buttons = append(wgt.buttons, btn)
		}
	}
	wgt.form.Add(children...)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *AlertWidget) Children() []Widget {
	return Many(wgt.modal)
}

// Draw renders the widget's HTML.
func (wgt *AlertWidget) Draw(w io.Writer, r *http.Request) (err error) {
	state := factory.StateOf(r)
	for _, btn := range wgt.buttons {
		if btn.href != "" {
			continue
		}
		if btn.Name() == "" {
			btn.WithHref("?" + wgt.name + "=")
		} else {
			v := state.Get(wgt.name)
			btn.WithHref("?" + wgt.name + "=&" + btn.Name() + "=" + url.QueryEscape(v))
		}
	}
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Class("Alert").
		Add(wgt.modal).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *AlertWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) || state.Changed(wgt.name)
}
