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

package form

import (
	"io"
	"net/http"
	"strconv"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&InputFileWidget{})      // Ensure interface
var _ = InputWidget(&InputFileWidget{}) // Ensure interface

// InputFileWidget renders a file selector.
type InputFileWidget struct {
	*widget.InputWidgetBase[*InputFileWidget]
	value    string
	accept   string
	receiver string
	maxSize  string
}

/*
InputFile creates a new widget that renders a file selector.
A non-empty filename indicates the current value of the widget.

The receiver is a web endpoint that the client can communicate with to upload the file.
The endpoint must support the following API:

	POST path?name=filename.txt

Accepts the body of the first or only data chunk. It must return a JSON struct with a "key"
property to identify the upload sequence: {"key":"123456"}

	POST path?key=123456

Accepts additional chunks to the upload sequence represented by the key.

	GET path?key=123456

Returns the consolidated chunks in the body of the message, and the file name in the
Content-Disposition header.
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
*/
func (f FormFactory) InputFile(name string, fileName string, receiverURL string) *InputFileWidget {
	x := &InputFileWidget{
		value:    fileName,
		receiver: receiverURL,
		maxSize:  strconv.Itoa(64 * 1024 * 1024),
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

/*
WithAccept determines what file types are accepted by the widget.
Types is a comma-separated list of any of the following:
(1) A valid case-insensitive filename extension, starting with a period (".") character;
(2) A valid MIME type string, with no extensions;
(3) The string audio/* meaning "any audio file";
(4) The string video/* meaning "any video file";
(5) The string image/* meaning "any image file".

See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/accept .
*/
func (wgt *InputFileWidget) WithAccept(types string) *InputFileWidget {
	wgt.accept = types
	return wgt
}

// WithMaxSize caps the upload size in bytes. Default is 64 MiB. Pass 0
// (or a negative value) to remove the limit. Enforced client-side only —
// always re-validate at the receiver endpoint.
func (wgt *InputFileWidget) WithMaxSize(sizeBytes int) *InputFileWidget {
	if sizeBytes > 0 {
		wgt.maxSize = strconv.Itoa(sizeBytes)
	} else {
		wgt.maxSize = ""
	}
	return wgt
}

// Value returns the value of the field.
func (wgt *InputFileWidget) Value(r *http.Request) string {
	value := wgt.value
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) { // || wgt.Submitted(r)
		value = state.Get(wgt.Name())
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *InputFileWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		return false
	}
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *InputFileWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *InputFileWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	// invalid := !wgt.Valid(r)

	dropTargetTag := Tag("div").
		Class("DropZone").
		Add("Drop a file or click to browse")
	if !wgt.Disabled() {
		dropTargetTag.
			Attr("ondrop", "inputfile_drop(event)").
			Attr("ondragover", "inputfile_dragover(event)").
			Attr("ondragenter", "inputfile_dragenter(event)").
			Attr("ondragleave", "inputfile_dragleave(event)")
	}
	inputFileTag := Tag("input").
		Attr("type", "file").
		Attr("accept", wgt.accept).
		Attr("onchange", "inputfile_fileSelected(event)")
	if wgt.Disabled() {
		inputFileTag = nil
	}
	inputHiddenTag := Tag("input").
		Attr("type", "hidden").
		Attr("name", wgt.Name()).
		AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
		Attr("oninput", "input_input(event)").
		Attr("oninvalid", "input_invalid(event)").
		AttrIf(value != "", "value", "0")
	if wgt.Disabled() {
		inputHiddenTag = nil
	}

	progressTag := Tag("progress").
		Attr("max", "100")

	statusTag := Tag("div").
		Class("FileName").
		Add(
			Tag("div").Add(value),
			factory.Icon("cancel"),
		)
	if wgt.Disabled() {
		statusTag = Tag("div").
			Class("FileName").
			Add(value)
	}

	state := "Upload"
	if value != "" {
		state = "Uploaded"
	}
	return Tag("div").
		Class("InputFile", state).
		Attr("data-id", wgt.ID()).
		Attr("data-receiver", wgt.receiver).
		Attr("data-max-size", wgt.maxSize).
		Attr("tabindex", "0").
		Attr("onclick", "inputfile_click(event)").
		Attr("onkeydown", "inputfile_keydown(event)").
		AttrIf(wgt.Disabled(), "disabled", "1").
		Add(dropTargetTag, progressTag, statusTag, inputFileTag, inputHiddenTag).
		When(wgt.Shown(r)).
		Draw(w, r)
}
