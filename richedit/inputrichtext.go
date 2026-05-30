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

package richedit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	htmlparser "golang.org/x/net/html"

	"github.com/microbus-io/bespa/form"
	"github.com/microbus-io/bespa/widget"
	"github.com/microbus-io/errors"
)

var _ = Widget(&InputRichTextWidget{})      // Ensure interface
var _ = InputWidget(&InputRichTextWidget{}) // Ensure interface

// InputRichTextWidget renders a rich text input field backed by Quill.
type InputRichTextWidget struct {
	*widget.InputWidgetBase[*InputRichTextWidget]
	value        string
	minRows      int
	maxRows      int
	predicates   form.Predicates
	width        string
	errMsg       string
	autoFocus    bool
	toolbar      []string
	placeholder  string
	mentionFeeds []*mentionFeed
	maxLen       int
	minLen       int
}

// mentionFeed defines an autocompletion feed surfaced when the user types its marker.
type mentionFeed struct {
	Marker            string   `json:"marker"`
	MinimumCharacters int      `json:"minimumCharacters"`
	Feed              []string `json:"feed"`
}

// InputRichText creates a new widget that renders a Quill rich-text editor.
// name is the state variable; value is the initial HTML (sanitized on
// submit — <script>, javascript: URIs, and on* handlers are rejected).
// Defaults: 8 minimum rows / 12 maximum, the common formatting toolbar.
// Customize the toolbar with WithToolbar; add @-mentions via
// WithMentionFeed.
func (f RichEditFactory) InputRichText(name string, value string) *InputRichTextWidget {
	x := &InputRichTextWidget{
		value:   value,
		minRows: 8,
		maxRows: 12,
		minLen:  -1,
		maxLen:  -1,
		toolbar: []string{
			"bold", "italic", "underline", "strikethrough", "fontBackgroundColor", "fontColor", "|",
			"alignment", "numberedList", "bulletedList", "outdent", "indent", "blockQuote", "|",
			"removeFormat",
		},
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// WithWidth caps the editor's width (max-width). Pass any CSS length,
// e.g. "60ch", "800px" or "100%". Empty lets it fill the container — the
// default.
func (wgt *InputRichTextWidget) WithWidth(css string) *InputRichTextWidget {
	if css != "" {
		wgt.width = "max-width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithRows sets the editor's visible height range, in rows. The editor
// starts at minRows tall and grows up to maxRows before scrolling.
// Defaults: 8 / 12. Values below 1 are clamped to 1; maxRows is clamped
// up to minRows.
func (wgt *InputRichTextWidget) WithRows(minRows int, maxRows int) *InputRichTextWidget {
	if minRows < 1 {
		minRows = 1
	}
	if maxRows < minRows {
		maxRows = minRows
	}
	wgt.minRows = minRows
	wgt.maxRows = maxRows
	return wgt
}

// WithAutoFocus automatically focuses the cursor on the text input field.
// Auto-focus is off by default.
func (wgt *InputRichTextWidget) WithAutoFocus(autoFocus bool) *InputRichTextWidget {
	wgt.autoFocus = autoFocus
	return wgt
}

// WithPredicate adds a custom server-side validator. value is the raw
// HTML markup. Predicates run after length checks but before the
// built-in script-tag / event-handler sanitizer.
func (wgt *InputRichTextWidget) WithPredicate(predicate func(value string) (bool, string)) *InputRichTextWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

/*
WithMentionFeed adds an autocomplete feed triggered by a single-character marker.
When the user types the marker (e.g. @ or #), a suggestion panel appears next to
the caret; the chosen entry is inserted into the document. Multiple feeds may be
registered, each with its own marker.

Example:

	wgt.WithMentionFeed("@", 0, "Harry Potter", "Tom Riddle")

Implemented via the quill-mention add-on (MIT). See ATTRIBUTIONS.md.
*/
func (wgt *InputRichTextWidget) WithMentionFeed(marker string, minChars int, feed ...string) *InputRichTextWidget {
	wgt.mentionFeeds = append(wgt.mentionFeeds, &mentionFeed{
		Marker:            marker,
		MinimumCharacters: minChars,
		Feed:              feed,
	})
	return wgt
}

/*
WithToolbar sets the buttons to show in the toolbar.
Button names follow the legacy ckeditor naming for source compatibility.
Supported tokens:

	"|"                       toolbar group separator
	"bold", "italic", "underline", "strikethrough", "code", "codeBlock"
	"subscript", "superscript"
	"fontColor", "fontBackgroundColor", "fontFamily", "fontSize"
	"alignment", "alignment:left", "alignment:center", "alignment:right", "alignment:justify"
	"numberedList", "bulletedList", "todoList"
	"outdent", "indent"
	"blockQuote", "heading"
	"link", "image", "mediaEmbed"
	"removeFormat"

Tokens not in this list are silently dropped. Features without an equivalent
in Quill (find & replace, special characters menu, source view, page break,
table-properties dialogs) are not supported.
*/
func (wgt *InputRichTextWidget) WithToolbar(buttons ...string) *InputRichTextWidget {
	wgt.toolbar = buttons
	return wgt
}

// WithPlaceholder sets the placeholder text of the field.
func (wgt *InputRichTextWidget) WithPlaceholder(placeholder string) *InputRichTextWidget {
	wgt.placeholder = placeholder
	return wgt
}

// WithLength bounds the editor's value length, measured against the raw
// HTML markup (not the visible text). Pass a negative value for either
// bound to leave it unbounded. A non-zero minimum does not imply Required.
func (wgt *InputRichTextWidget) WithLength(minChars int, maxChars int) *InputRichTextWidget {
	wgt.minLen = minChars
	wgt.maxLen = maxChars
	return wgt
}

// Value returns the value of the field.
func (wgt *InputRichTextWidget) Value(r *http.Request) string {
	value := wgt.value
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) {
		value = state.Get(wgt.Name())
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *InputRichTextWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if value == "" && wgt.Required() {
		wgt.errMsg = "A value is required"
		return false
	}
	if wgt.maxLen >= 0 && len([]rune(value)) > wgt.maxLen {
		wgt.errMsg = "Value exceeds maximum length"
		return false
	}
	if wgt.minLen >= 0 && len([]rune(value)) < wgt.minLen {
		wgt.errMsg = "Value is not long enough"
		return false
	}
	if ok, errMsg := wgt.predicates.Validate(value); !ok {
		wgt.errMsg = errMsg
		return false
	}
	if value == "" {
		return true
	}
	if ok, errMsg := detectScript(value); !ok {
		wgt.errMsg = errMsg
		return false
	}
	return true
}

// detectScript checks if the value contains a script or event handler attribute.
func detectScript(value string) (ok bool, errMsg string) {
	tokenizer := htmlparser.NewTokenizer(strings.NewReader(value))
loop:
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case htmlparser.ErrorToken:
			if errors.Is(tokenizer.Err(), io.EOF) {
				break loop
			}
			return false, tokenizer.Err().Error()
		case htmlparser.StartTagToken, htmlparser.SelfClosingTagToken:
			tagName, hasAttr := tokenizer.TagName()
			if string(tagName) == "script" {
				return false, "Script tags are not allowed"
			}
			for hasAttr {
				var attrName []byte
				var attrVal []byte
				attrName, attrVal, hasAttr = tokenizer.TagAttr()
				if string(attrName) == "src" || string(attrName) == "href" {
					if strings.Contains(strings.ToLower(string(attrVal)), "javascript:") {
						return false, "Script URIs are not allowed"
					}
				}
				if strings.HasPrefix(string(attrName), "on") {
					return false, "Event handlers are not allowed"
				}
			}
		}
		if errors.Is(tokenizer.Err(), io.EOF) {
			break loop
		}
	}
	if len(tokenizer.Raw()) > 0 {
		return false, "Malformed HTML"
	}
	return true, ""
}

// Changed indicates if the value of the field changed.
func (wgt *InputRichTextWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *InputRichTextWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)
	randomID := widget.RandomAlphaNumID(8)

	// Hidden textarea retains the form value, accessibility hooks, and the bespa
	// auto-submit infrastructure. Quill writes to it via the inputrichtext.js wrapper.
	textAreaTag := Tag("textarea").
		Attr("id", randomID).
		Attr("name", wgt.Name()).
		AttrIf(wgt.Disabled(), "disabled", "1").
		Attr("rows", strconv.Itoa(wgt.maxRows)).
		Class("RichEditTextArea").
		Style(wgt.width)
	if !wgt.Disabled() {
		textAreaTag.
			Attr("tabindex", "0").
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("oninput", "input_input(event)").
			Attr("onchange", "input_change(event)").
			Attr("oninvalid", "input_invalid(event)").
			AttrIf(wgt.Required(), "required", "1").
			AttrIf(wgt.autoFocus, "autofocus", "1").
			Add(value)
	} else {
		textAreaTag.Add(value)
	}

	// Quill init options serialized as JSON for the JS wrapper.
	opts := map[string]any{
		"textareaID":  randomID,
		"toolbar":     wgt.toolbar,
		"placeholder": wgt.placeholder,
		"autoFocus":   wgt.autoFocus,
		"disabled":    wgt.Disabled(),
	}
	if len(wgt.mentionFeeds) > 0 {
		opts["mentions"] = wgt.mentionFeeds
	}
	optsJSON, err := json.Marshal(opts)
	if err != nil {
		return errors.Trace(err)
	}

	scriptTag := Tag("script").Add(factory.HTMLUnsafe(
		"richedit_init(", string(optsJSON), ");",
	))

	errTag := Tag("")
	if invalid && wgt.errMsg != "" {
		msgJSON, _ := json.Marshal(wgt.errMsg)
		encoded := string(msgJSON)
		encoded = strings.ReplaceAll(encoded, "<", "\\u003c")
		encoded = strings.ReplaceAll(encoded, ">", "\\u003e")
		errTag = Tag("script").Add(factory.HTMLUnsafe(
			`input_setCustomValidity("`, randomID, `", `, encoded, `)`,
		))
	}

	// Per-instance min/max height tracked via CSS variables on the wrapper element.
	rowHeight := "var(--md-sys-typescale-body-large-line-height)"
	containerStyle := fmt.Sprintf(
		"--richedit-min-height:calc(%s * %d);--richedit-max-height:calc(%s * %d);",
		rowHeight, wgt.minRows+1, rowHeight, wgt.maxRows+1,
	)

	return Tag("span").
		Class("RichEdit").
		Attr("id", "richedit"+randomID).
		Attr("data-id", wgt.ID()).
		Style(containerStyle).
		ClassIf(invalid, "Invalid").
		AttrIf(wgt.Disabled(), "disabled", "1").
		Add(textAreaTag, scriptTag, errTag).
		When(wgt.Shown(r)).
		Draw(w, r)
}
