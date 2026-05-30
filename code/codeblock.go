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

package code

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/microbus-io/bespa/widget"
	"github.com/microbus-io/errors"
)

var _ = Widget(&CodeBlockWidget{}) // Ensure interface

// CodeBlockWidget renders a block of source code with syntax highlighting.
type CodeBlockWidget struct {
	*widget.WidgetBase[*CodeBlockWidget]
	code        string
	language    string
	lineNumbers bool
	maxHeight   string // CSS value driving the inner pre's max-height; empty disables.
	frame       bool
}

// CodeBlock creates a new widget that renders a syntax-highlighted code
// block. Highlighting runs server-side on every Draw; for cheap, plain
// <pre><code> output without highlighting use basic.PlainCodeBlock.
// Without WithLanguage the lexer is auto-detected from the code; if
// detection fails, the code renders unhighlighted.
func (f CodeFactory) CodeBlock(code string) *CodeBlockWidget {
	x := &CodeBlockWidget{
		code:  code,
		frame: true,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithFrame toggles the surrounding border and tinted background.
// On by default. Pass false for a bare, transparent code block — useful
// when embedding inside a Card or other container that already provides
// the frame.
func (wgt *CodeBlockWidget) WithFrame(on bool) *CodeBlockWidget {
	wgt.frame = on
	return wgt
}

// WithLanguage sets the language used to tokenize the code.
// Use any name recognized by Chroma (e.g. "go", "javascript", "python", "rust",
// "sql", "json", "yaml", "html", "css", "bash"). The full list lives at
// https://github.com/alecthomas/chroma/tree/master/lexers/embedded.
// An empty string falls back to auto-detection.
func (wgt *CodeBlockWidget) WithLanguage(name string) *CodeBlockWidget {
	wgt.language = name
	return wgt
}

// WithLineNumbers toggles line-number rendering in the left gutter.
// Off by default.
func (wgt *CodeBlockWidget) WithLineNumbers(on bool) *CodeBlockWidget {
	wgt.lineNumbers = on
	return wgt
}

// WithMaxRows caps the visible height of the code block at the given number of
// lines. Content beyond that scrolls vertically. Set to 0 to remove the cap.
func (wgt *CodeBlockWidget) WithMaxRows(rows int) *CodeBlockWidget {
	if rows <= 0 {
		wgt.maxHeight = ""
		return wgt
	}
	// One row of body-medium line-height per row, plus the inner pre's vertical
	// padding (1em top + 1em bottom).
	wgt.maxHeight = fmt.Sprintf(
		"calc(var(--md-sys-typescale-body-medium-line-height) * %d + 2em)",
		rows,
	)
	return wgt
}

// WithMaxHeight caps the visible height of the code block at the given CSS
// length, e.g. "300px", "50vh" or "calc(100vh - 50px)". Content beyond the cap
// scrolls vertically. Empty removes the cap.
func (wgt *CodeBlockWidget) WithMaxHeight(css string) *CodeBlockWidget {
	if css == "" {
		wgt.maxHeight = ""
		return wgt
	}
	wgt.maxHeight = css
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *CodeBlockWidget) Draw(w io.Writer, r *http.Request) (err error) {
	lexer := lexers.Fallback
	if wgt.language != "" {
		if l := lexers.Get(wgt.language); l != nil {
			lexer = l
		}
	} else if l := lexers.Analyse(wgt.code); l != nil {
		lexer = l
	}
	lexer = chroma.Coalesce(lexer)

	// Class-based output lets the codeblock.css stylesheet drive the look,
	// so light/dark theme switching works without re-running the highlighter.
	// We let Chroma emit its own <pre class="chroma"> wrapper rather than
	// suppressing it — Chroma's inline line-number emission is gated on the
	// same code path that produces the <pre>, so suppressing the wrapper
	// silently disables line numbers as well.
	formatter := html.New(
		html.WithClasses(true),
		html.WithLineNumbers(wgt.lineNumbers),
	)

	iterator, err := lexer.Tokenise(nil, wgt.code)
	if err != nil {
		return errors.Trace(err)
	}
	var buf bytes.Buffer
	if err := formatter.Format(&buf, styles.Fallback, iterator); err != nil {
		return errors.Trace(err)
	}

	styleAttr := ""
	if wgt.maxHeight != "" {
		styleAttr = "--codeblock-max-height:" + wgt.maxHeight
	}
	tag := Tag("div").
		Class("CodeBlock").
		Attr("data-id", wgt.ID()).
		Attr("data-language", lexer.Config().Name).
		Style(styleAttr).
		Add(factory.HTMLUnsafe(buf.String())).
		When(wgt.Shown(r))
	if !wgt.frame {
		tag = tag.Class("NoFrame")
	}
	return tag.Draw(w, r)
}
