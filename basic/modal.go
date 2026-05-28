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
	"fmt"
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&ModalWidget{}) // Ensure interface

// ModalWidget renders a modal window.
type ModalWidget struct {
	*widget.WidgetBase[*ModalWidget]
	width     string
	minHeight string
	children  []Widget
	name      string
}

// Modal creates a new widget that renders a modal window.
// The modal is bound to the named state variable: it opens when the variable
// is non-empty and closes when it is cleared. The typical pattern is to set
// the variable to the URL of an embedded handler and close with `^?name=`
// from inside the modal's content.
func (f BasicFactory) Modal(name string) *ModalWidget {
	x := &ModalWidget{
		name: name,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithWidth sets the width of the modal window.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
// The default is 826px which fits 800px content internally.
// In any case the modal will not span more than 90% of the width of the viewport.
func (wgt *ModalWidget) WithWidth(width float32, unit string) *ModalWidget {
	if width > 0 {
		wgt.width = fmt.Sprintf("width:%f%s", width, unit)
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithMinHeight sets the minimum height of the modal window.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
// The default is 240px.
// Set to 0 to adjust the height to the content of the modal.
// In any case the modal will not span more than 90% of the height of the viewport.
func (wgt *ModalWidget) WithMinHeight(minHeight float32, unit string) *ModalWidget {
	if minHeight <= 0 {
		wgt.minHeight = ""
		return wgt
	}
	wgt.minHeight = fmt.Sprintf("min-height:%f%s", minHeight, unit)
	return wgt
}

// Add adds nested widgets.
func (wgt *ModalWidget) Add(children ...any) *ModalWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Embed is shorthand for modal.Add(factory.Embed(fetcher)).
func (wgt *ModalWidget) Embed(fetcher func() (res *http.Response, err error)) *ModalWidget {
	return wgt.Add(factory.Embed(fetcher))
}

// Children are the widgets nested under this widget.
func (wgt *ModalWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *ModalWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)
	if !wgt.Shown(r) || len(wgt.children) == 0 {
		return Tag("span").
			Attr("data-id", wgt.ID()).
			Attr("id", randomID).
			Add(Tag("script").Add(HTMLUnsafe("modal_close('", randomID, "')"))).
			Hide(true).
			Draw(w, r)
	}
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Attr("id", randomID).
		Attr("role", "dialog").
		Attr("aria-modal", "true").
		Class("Modal").
		Add(
			Tag("div").
				Style(wgt.width, wgt.minHeight).
				Add(wgt.children),
			Tag("script").Add(HTMLUnsafe("modal_open('", randomID, "')")),
		).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *ModalWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) || state.Changed(wgt.name)
}

// Shown returns whether the widget is shown or hidden.
// A widget that is hidden is not rendered.
func (wgt *ModalWidget) Shown(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Shown(r) && state.Get(wgt.name) != ""
}
