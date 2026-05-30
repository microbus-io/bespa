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

var _ = Widget(&SidePanelWidget{}) // Ensure interface

// SidePanelWidget renders a floating side panel.
type SidePanelWidget struct {
	*widget.WidgetBase[*SidePanelWidget]
	width    string
	children []Widget
	name     string
}

// SidePanel creates a new widget that renders a floating side panel.
// Like Modal, the panel is bound to the named state variable: it opens
// when the variable is non-empty and closes when cleared. The typical
// pattern is to set the variable to the URL of an embedded handler and
// close with `^?name=` from inside the panel's content.
func (f BasicFactory) SidePanel(name string) *SidePanelWidget {
	x := &SidePanelWidget{
		name: name,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	x.WithWidth("400px")
	return x
}

// WithWidth sets the width of the side panel. Default is 400px.
// Pass any CSS length, e.g. "400px", "50%" or "calc(100vh - 50px)". Empty clears it.
// The panel is capped at 90% of the viewport width regardless.
func (wgt *SidePanelWidget) WithWidth(css string) *SidePanelWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// Add adds nested widgets.
func (wgt *SidePanelWidget) Add(children ...any) *SidePanelWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *SidePanelWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *SidePanelWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Class("SidePanel").
		Style(wgt.width).
		Attr("onmousedown", "sidepanel_mousedown(event)").
		Add(
			Tag("div").Add(Tag("div").Add(wgt.children)),
			Tag("script").Add(HTMLUnsafe("sidepanel_open('", wgt.ID(), "')")),
		).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *SidePanelWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) || state.Changed(wgt.name)
}

// Shown returns whether the widget is shown or hidden.
// A widget that is hidden is not rendered.
func (wgt *SidePanelWidget) Shown(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Shown(r) && state.Get(wgt.name) != ""
}
