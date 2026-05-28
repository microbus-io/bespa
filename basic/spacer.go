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

var _ = Widget(&SpacerWidget{}) // Ensure interface

// SpacerWidget renders a vertical space.
type SpacerWidget struct {
	*widget.WidgetBase[*SpacerWidget]
	space float32
}

// Spacer creates a new widget that inserts vertical whitespace measured in
// 16-pixel "lines". Lines can be fractional or negative (negative pulls
// surrounding content closer together).
func (f BasicFactory) Spacer(lines float32) *SpacerWidget {
	x := &SpacerWidget{
		space: lines,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// SpacerNewLine creates a spacer with 0 lines of vertical space.
func (f BasicFactory) SpacerNewLine() *SpacerWidget {
	return f.Spacer(0)
}

// SpacerBreak creates a spacer with 0.5 lines of vertical space.
func (f BasicFactory) SpacerBreak() *SpacerWidget {
	return f.Spacer(0.5)
}

// SpacerParagraph creates a spacer with 1 lines of vertical space.
func (f BasicFactory) SpacerParagraph() *SpacerWidget {
	return f.Spacer(1)
}

// Draw renders the widget's HTML.
func (wgt *SpacerWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Style(fmt.Sprintf("height:%.2fpx;", 16*wgt.space)).
		When(wgt.Shown(r)).
		Draw(w, r)
}
