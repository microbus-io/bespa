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

var _ = Widget(&SnackbarWidget{}) // Ensure interface

// SnackbarWidget renders a transient toast notification.
type SnackbarWidget struct {
	*widget.WidgetBase[*SnackbarWidget]
	children []Widget
}

// Snackbar creates a new widget that renders a Material snackbar: a brief
// message that slides in, auto-dismisses, and pauses dismissal while
// hovered. Pair it with RedrawIfChanged on the state variable that drives
// the message so it appears in response to actions like saves or errors.
func (f BasicFactory) Snackbar() *SnackbarWidget {
	x := &SnackbarWidget{}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *SnackbarWidget) Add(children ...any) *SnackbarWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *SnackbarWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *SnackbarWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)
	return Tag("div").
		Class("Snackbar").
		Attr("data-id", wgt.ID()).
		Attr("id", randomID).
		Attr("role", "status").
		Attr("aria-live", "polite").
		Attr("aria-atomic", "true").
		Attr("onmouseenter", "snackbar_mouseenter(event)").
		Attr("onmouseleave", "snackbar_mouseleave(event)").
		Add(
			wgt.children,
			Tag("script").Add(HTMLUnsafe("snackbar_init('", randomID, "')")),
		).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
