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

var _ = Widget(&DebuggerWidget{}) // Ensure interface

// DebuggerWidget renders a debugger tool.
type DebuggerWidget struct {
	*widget.WidgetBase[*DebuggerWidget]
}

// Debugger creates a new widget that renders a floating in-page debug
// panel. The panel exposes the current state, request details, and
// partial-redraw activity — useful during development. Drop it from
// production pages.
func (f BasicFactory) Debugger() *DebuggerWidget {
	x := &DebuggerWidget{}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Draw renders the widget's HTML.
func (wgt *DebuggerWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)
	return Tag("span").
		Attr("id", randomID).
		Attr("data-id", wgt.ID()).
		Add(
			Tag("div").
				Class("Debugger").
				Attr("onclick", "debugger_click(event)"),
			Tag("script").
				Add(HTMLUnsafe("debugger_init('", randomID, "')"))).
		When(wgt.Shown(r)).
		Draw(w, r)
}
