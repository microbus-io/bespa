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
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&CodeWidget{}) // Ensure interface

// CodeWidget renders code.
type CodeWidget struct {
	*widget.WidgetBase[*CodeWidget]
	block    bool
	children []Widget
	language string
}

// Code creates a new widget that renders an inline <code> snippet.
// For multi-line blocks use PlainCodeBlock, or the code package's CodeBlock
// for syntax highlighting.
func (f BasicFactory) Code(code string) *CodeWidget {
	x := &CodeWidget{
		children: Many(code),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// PlainCodeBlock creates a new widget that renders a block of code with no
// syntax highlighting — just <pre><code> wrapping. For Chroma-driven syntax
// highlighting, use the code package's CodeBlock instead.
func (f BasicFactory) PlainCodeBlock(code ...any) *CodeWidget {
	x := &CodeWidget{
		block:    true,
		children: Many(code),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithLanguage tags the code with a `language-<lang>` CSS class.
// It does not perform syntax highlighting on its own; use the code package
// for that.
func (wgt *CodeWidget) WithLanguage(language string) *CodeWidget {
	wgt.language = "language-" + strings.ToLower(language)
	return wgt
}

// Add adds nested widgets.
func (wgt *CodeWidget) Add(children ...any) *CodeWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *CodeWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *CodeWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if wgt.block {
		return Tag("pre").
			Attr("data-id", wgt.ID()).
			Class("Code", "Block").
			Add(
				Tag("code").Class(wgt.language).Add(wgt.children),
			).
			When(wgt.Shown(r) && len(wgt.children) > 0).
			Draw(w, r)
	}
	return Tag("code").
		Attr("data-id", wgt.ID()).
		Class(wgt.language).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
