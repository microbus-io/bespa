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

var _ = Widget(&CopyToClipboardWidget{}) // Ensure interface

// CopyToClipboardWidget renders a small icon button that writes a string to
// the system clipboard when clicked.
type CopyToClipboardWidget struct {
	*widget.WidgetBase[*CopyToClipboardWidget]
	text    string
	altText string
}

// CopyToClipboard creates a new widget that renders a "content copy" icon.
// Clicking it writes the given text to the system clipboard and briefly
// swaps the icon to a check mark as feedback. The text is embedded in the
// rendered HTML, so prefer this widget for short-to-medium strings rather
// than large documents.
func (f BasicFactory) CopyToClipboard(text string) *CopyToClipboardWidget {
	x := &CopyToClipboardWidget{
		text:    text,
		altText: "Copy to clipboard",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithAltText overrides the default "Copy to clipboard" label used for the
// hover tooltip and screen-reader announcement.
func (wgt *CopyToClipboardWidget) WithAltText(altText string) *CopyToClipboardWidget {
	wgt.altText = altText
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *CopyToClipboardWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("span").
		Class("CopyToClipboard").
		Attr("data-id", wgt.ID()).
		Attr("data-copy-text", wgt.text).
		Attr("role", "button").
		Attr("tabindex", "0").
		Attr("title", wgt.altText).
		Attr("aria-label", wgt.altText).
		Attr("onclick", "copytoclipboard_click(event)").
		Attr("onkeydown", "copytoclipboard_keydown(event)").
		Add(factory.Icon("content_copy")).
		When(wgt.Shown(r)).
		Draw(w, r)
}
