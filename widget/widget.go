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

package widget

import (
	"io"
	"net/http"
)

// Widget is a UI component that is rendered on a page.
// WidgetBase is used to implement this interface.
type Widget interface {
	// ID is a unique identifier within the scope of the page.
	ID() string

	// SetID sets a unique identifier within the scope of the page.
	SetID(id string)

	// Children are the widgets nested under this widget, or nil if none.
	Children() []Widget

	// Draw renders the widget to the writer, given an HTTP request.
	Draw(w io.Writer, r *http.Request) (err error)

	// Drawn indicates whether this widget needs to be drawn in either a full or partial page rendering.
	Drawn(r *http.Request) bool

	// Shown indicates whether this widget is shown or hidden.
	Shown(r *http.Request) bool
}
