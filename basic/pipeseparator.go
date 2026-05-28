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

var _ = Widget(&PipeSeparatorWidget{}) // Ensure interface

// PipeSeparatorWidget renders a vertical pipe separator.
type PipeSeparatorWidget struct {
	*widget.WidgetBase[*PipeSeparatorWidget]
}

// PipeSeparator creates a new widget that renders a vertical pipe ("|")
// between inline items, e.g. in a footer or breadcrumb.
func (f BasicFactory) PipeSeparator() *PipeSeparatorWidget {
	x := &PipeSeparatorWidget{}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Draw renders the widget's HTML.
func (wgt *PipeSeparatorWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Class("PipeSeparator").
		When(wgt.Shown(r)).
		Draw(w, r)
}
