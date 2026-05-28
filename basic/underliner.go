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
	"html"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&UnderlinerWidget{}) // Ensure interface

// UnderlinerWidget underlines search terms in a text.
type UnderlinerWidget struct {
	*widget.WidgetBase[*UnderlinerWidget]
	text   string
	terms  []string
	prefix bool
}

// Underliner creates a new widget that renders text with any matches of
// the given search terms wrapped in <u>. terms is whitespace-split into
// individual words; matching is case-insensitive. Useful for highlighting
// the query in QuickSearch results — see also QuickSearchUnderliner.
func (f BasicFactory) Underliner(text string, terms string) *UnderlinerWidget {
	word := ""
	words := []string{}
	for _, r := range terms {
		if unicode.IsSpace(r) {
			if word != "" {
				words = append(words, word)
			}
			word = ""
		} else {
			word += string(r)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	x := &UnderlinerWidget{
		text:  text,
		terms: words,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithPrefixOnly restricts matching to word prefixes. When false (default),
// a term matches anywhere within a word. Word boundaries are runs of
// whitespace.
func (wgt *UnderlinerWidget) WithPrefixOnly(prefix bool) *UnderlinerWidget {
	wgt.prefix = prefix
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *UnderlinerWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if len(wgt.terms) == 0 {
		return Tag("span").
			Attr("data-id", wgt.ID()).
			Add(wgt.text).
			Draw(w, r)
	}
	word := ""
	words := []string{}
	for _, r := range wgt.text {
		if unicode.IsSpace(r) {
			if word != "" {
				words = append(words, word)
			}
			word = ""
		} else {
			word += string(r)
		}
	}
	if word != "" {
		words = append(words, word)
	}
	var sb strings.Builder
	for i, word := range words {
		if i > 0 {
			sb.WriteString(" ")
		}
		wordLower := strings.ToLower(word)
		printed := false
		for _, term := range wgt.terms {
			term := strings.ToLower(term)
			p := strings.Index(wordLower, term)
			if wgt.prefix && p == 0 || !wgt.prefix && p >= 0 {
				sb.WriteString(html.EscapeString(word[:p]))
				sb.WriteString("<u>")
				sb.WriteString(html.EscapeString(word[p : p+len(term)]))
				sb.WriteString("</u>")
				sb.WriteString(html.EscapeString(word[p+len(term):]))
				printed = true
				break
			}
		}
		if !printed {
			sb.WriteString(html.EscapeString(word))
		}
	}
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(factory.HTMLUnsafe(sb.String())).
		Draw(w, r)
}
