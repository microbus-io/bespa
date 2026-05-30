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
	"net/url"
	"strconv"

	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&ButtonWidget{})      // Ensure interface
var _ = InputWidget(&ButtonWidget{}) // Ensure interface

// ButtonWidget renders a button.
type ButtonWidget struct {
	*widget.InputWidgetBase[*ButtonWidget]
	style      string
	href       string
	target     string
	back       bool
	children   []Widget
	width      string
	predicates Predicates
	errMsg     string
}

// button creates a new widget that renders a button.
func (f FormFactory) button(style string, name string) *ButtonWidget {
	x := &ButtonWidget{
		style: style,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// ButtonFilled creates a new high-emphasis Material button — use for the
// primary action of a form or screen. name is the field name posted on
// click; an empty name means the press isn't reported back via Pressed.
func (f FormFactory) ButtonFilled(name string) *ButtonWidget {
	return f.button("Filled", name)
}

// ButtonOutlined creates a new Material outlined button — medium emphasis,
// for actions paired with a primary filled button.
func (f FormFactory) ButtonOutlined(name string) *ButtonWidget {
	return f.button("Outlined", name)
}

// ButtonText creates a new Material text button — low emphasis, suitable
// for in-line and toolbar actions.
func (f FormFactory) ButtonText(name string) *ButtonWidget {
	return f.button("Text", name)
}

// ButtonElevated creates a new Material elevated button — like ButtonText
// but with a slight shadow, used when the button needs to stand out from
// a busy background.
func (f FormFactory) ButtonElevated(name string) *ButtonWidget {
	return f.button("Elevated", name)
}

// ButtonTonal creates a new Material tonal button — a softer alternative
// to ButtonFilled, used for secondary actions that still need visual weight.
func (f FormFactory) ButtonTonal(name string) *ButtonWidget {
	return f.button("Tonal", name)
}

// WithHref turns the button into a link to href instead of a submit
// button. Accepts the full action-URL grammar (`?key=`, `^?…`, `/path`,
// etc.). Resets any prior WithHrefBack().
func (wgt *ButtonWidget) WithHref(href string) *ButtonWidget {
	wgt.back = false
	wgt.href = href
	return wgt
}

// WithHrefBack turns this into a "back" button. It follows the `_back`
// state variable when set; otherwise falls back to a browser-history
// back step if the referrer is the same host. The button auto-hides
// when there's no history to go back to. Set `_back=0` to force-disable.
func (wgt *ButtonWidget) WithHrefBack() *ButtonWidget {
	wgt.back = true
	return wgt
}

// WithTarget sets the target of the button's link.
func (wgt *ButtonWidget) WithTarget(target string) *ButtonWidget {
	wgt.target = target
	return wgt
}

// WithWidth scales the button to the given width.
// Pass any CSS length, e.g. "120px", "50%" or "calc(100% - 1em)". Empty clears it.
func (wgt *ButtonWidget) WithWidth(css string) *ButtonWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithPredicate attaches a custom validator that runs only when this
// button is pressed. The button must be named for press detection to
// work. Use this for cross-field validations that don't belong to any
// individual input — e.g. "either email or phone must be filled".
func (wgt *ButtonWidget) WithPredicate(predicate func(value string) (bool, string)) *ButtonWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

// Add adds nested widgets.
func (wgt *ButtonWidget) Add(children ...any) *ButtonWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *ButtonWidget) Children() []Widget {
	return wgt.children
}

// Value returns "1" if this named button was just pressed, "" otherwise.
// Disabled buttons always return "". Prefer Pressed for boolean checks.
func (wgt *ButtonWidget) Value(r *http.Request) string {
	if wgt.Disabled() {
		return ""
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) && state.Get(wgt.Name()) != "" {
		return "1"
	}
	return ""
}

// Pressed reports whether this button was just clicked. Use this to
// distinguish between multiple buttons on the same form. Only named
// buttons can be detected.
func (wgt *ButtonWidget) Pressed(r *http.Request) bool {
	return wgt.Value(r) != ""
}

// Valid validates the field's value against all validators.
func (wgt *ButtonWidget) Valid(r *http.Request) bool {
	if wgt.Pressed(r) {
		if wgt.Disabled() {
			return false
		}
		// Predicates
		value := wgt.Value(r)
		if ok, errMsg := wgt.predicates.Validate(value); !ok {
			wgt.errMsg = errMsg
			return false
		}
	}
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *ButtonWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.Value(r) != ""
}

// Draw renders the widget's HTML.
func (wgt *ButtonWidget) Draw(w io.Writer, r *http.Request) (err error) {
	state := factory.StateOf(r)
	href := wgt.href
	if wgt.back {
		href = state.Get("_back")
		if href == "" { // Default is go back one, if on same domain
			refHost := ""
			if ref, _ := url.Parse(r.Header.Get("Referer")); ref != nil {
				refHost = ref.Host
				if ref.Port() == "" {
					switch ref.Scheme {
					case "https":
						refHost += ":443"
					case "http":
						refHost += ":80"
					}
				}
			}
			host := r.Header.Get("X-Forwarded-Host")
			if host == "" {
				host = r.Host
			}
			if refHost == host {
				href = "-1"
			}
		}
		if href == "0" { // Disable back
			href = ""
		}
	}
	var detectionScript *widget.BytesWidget
	backButtonID := ""
	if backCount, err := strconv.Atoi(href); err == nil && backCount < 0 {
		href = "javascript:history.go(" + href + ")"
		// Detect if there's enough history and hide the button if not
		backButtonID = "back" + widget.RandomAlphaNumID(8)
		detectionScript = factory.HTMLUnsafe("<script>if (history.length<", strconv.Itoa(1-backCount), ") {document.getElementById('"+backButtonID+"').style.display='none';}</script>")
	}
	// Create a hidden link for button with an action
	btnType := "submit"
	aTag := Tag("a").Hide(true)
	if href != "" && !wgt.Disabled() {
		target := wgt.target
		if target == "" {
			target = state.Get("_target")
		}
		aTag = Tag("a").
			Attr("href", href).
			Attr("target", target).
			Hide(true)
		btnType = "button"
	}

	// The hidden checkbox is used to display validation errors
	invalid := !wgt.Valid(r)
	inputTag := Tag("")
	errTag := Tag("")
	if invalid && wgt.errMsg != "" {
		randomID := widget.RandomAlphaNumID(8)
		inputTag = Tag("input").
			Attr("type", "checkbox").
			Attr("tabindex", "-1").
			Attr("value", "1").
			Attr("id", randomID)
		errTag = customValidityScript(randomID, wgt.errMsg)
	}

	buttonTag := Tag("button").
		Attr("type", btnType).
		Style(wgt.width).
		Class(wgt.style).
		Attr("name", wgt.Name()).
		Attr("id", backButtonID).
		AttrIf(wgt.Disabled(), "disabled", "1").
		AttrIf(!wgt.Disabled(), "value", "1").
		AttrIf(!wgt.Disabled(), "onclick", "button_click(event)").
		AttrIf(!wgt.Disabled(), "tabindex", "0").
		Add(wgt.children, detectionScript, aTag)
	if len(wgt.children) == 1 {
		if _, ok := wgt.children[0].(*basic.IconWidget); ok {
			buttonTag.Class("SingleIcon")
		}
	}
	return Tag("span").
		Class("SubmitButton").
		ClassIf(invalid, "Invalid").
		Style(wgt.width).
		Attr("data-id", wgt.ID()).
		Attr("type", btnType).
		Add(inputTag, buttonTag, errTag).
		When(wgt.Shown(r) && (!wgt.back || href != "")).
		Draw(w, r)
}
