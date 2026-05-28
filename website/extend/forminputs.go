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

package extend

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

const formInputsInterface = `// widget/inputwidget.go
type InputWidget interface {
    Name() string                  // name of the state variable
    Value(r *http.Request) string  // current value (initial OR posted)
    Valid(r *http.Request) bool    // validation result
    Changed(r *http.Request) bool  // posted value differs from initial
    SetFormName(formName string)   // wired up by the parent Form
}
`

const formInputsSkeleton = `type ColorWheelWidget struct {
    *widget.InputWidgetBase[*ColorWheelWidget]
    value string
}

var _ = widget.Widget(&ColorWheelWidget{})      // interface assertion
var _ = widget.InputWidget(&ColorWheelWidget{}) // interface assertion

func (f MyFactory) ColorWheel(name string) *ColorWheelWidget {
    x := &ColorWheelWidget{}
    x.InputWidgetBase = widget.NewInputWidgetBase(x)
    x.WithName(name)
    return x
}

func (x *ColorWheelWidget) WithValue(v string) *ColorWheelWidget {
    x.value = v
    return x
}

func (x *ColorWheelWidget) Value(r *http.Request) string {
    if x.Disabled() {
        return x.value
    }
    state := factory.StateOf(r)
    if state.Has(x.Name()) {
        return state.Get(x.Name())
    }
    return x.value
}

func (x *ColorWheelWidget) Valid(r *http.Request) bool {
    if x.Disabled() || !x.Submitted(r) {
        return true
    }
    v := x.Value(r)
    if v == "" && x.Required() {
        return false
    }
    return colorRegex.MatchString(v)
}

func (x *ColorWheelWidget) Changed(r *http.Request) bool {
    if x.Disabled() || !x.Submitted(r) {
        return false
    }
    return x.value != x.Value(r)
}
`

const formInputsDraw = `func (x *ColorWheelWidget) Draw(w io.Writer, r *http.Request) error {
    // The hidden input is what the form submission picks up. Its "name"
    // attribute MUST equal x.Name(). The visible UI — the color wheel —
    // is whatever you want; sync its value into the hidden input via JS.
    _, err := fmt.Fprintf(w,
        "<span data-id=%q class=\"ColorWheel\">"+
        "<input type=\"hidden\" name=%q value=%q>"+
        "<canvas data-target=%q></canvas>"+
        "<script>colorwheel_init(%q)</script>"+
        "</span>",
        x.ID(), x.Name(), x.Value(r), x.ID(), x.ID())
    return err
}
`

const formInputsJS = `// Companion JS, registered via widget.AssetRegistry.RegisterScript.
// Whenever the visible widget changes, write the new value into the hidden
// input so the form submission picks it up.

function colorwheel_init(id) {
    const root = document.getElementById(id);
    const hidden = root.querySelector('input[type=hidden]');
    const canvas = root.querySelector('canvas');
    canvas.addEventListener('colorpicked', (e) => {
        hidden.value = e.detail.color;
    });
}
`

// HandleFormInputs covers writing custom form-input widgets that participate
// in BESPA forms: the InputWidget contract, the hidden-input rule, and the
// JS bridge for widgets backed by a third-party library.
func HandleFormInputs(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Custom form inputs"),

		wf.Markdown(
			"Most apps never need to write one — `DefaultFactory` ships text, ",
			"number, date, dropdown, chips, file, toggle, range, color, rating, ",
			"and timezone inputs. You write one when you want a UI BESPA ",
			"doesn't have: a canvas-based editor, a map picker, an integration ",
			"with a third-party JS widget like Quill or Monaco. The framework ",
			"is set up for this — the contract is small.",
		),

		wf.HeadlineMedium("The contract"),
		wf.Markdown(
			"An input widget is a `Widget` that also satisfies the ",
			"`InputWidget` interface. Five methods, three of which the base ",
			"type implements for you:",
		),
		wf.CodeBlock(formInputsInterface).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"`Name()` and `SetFormName()` come from `*InputWidgetBase[T]`. ",
			"You implement `Value`, `Valid`, and `Changed`. The parent `Form` ",
			"walks its children at submit time, picks out everything ",
			"satisfying `InputWidget`, calls `Valid` on each, and only fires ",
			"the post-submit hook if every input returns `true`.",
		),

		wf.HeadlineMedium("Skeleton"),
		wf.Markdown(
			"A custom `ColorWheel` widget, factory constructor, and the three ",
			"interface methods:",
		),
		wf.CodeBlock(formInputsSkeleton).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Three conventions to follow:",
			"\n\n",
			"- The two `var _ = ...(&MyWidget{})` lines assert at compile time ",
			"that the widget satisfies both interfaces. Drop one and you'll ",
			"get a runtime miss instead of a build error.\n",
			"- The constructor calls `widget.NewInputWidgetBase(x)` (not the ",
			"plain `NewWidgetBase`) so the input-specific bookkeeping ",
			"(`Name`, `Disabled`, `Required`, `AutoSubmit`, `Submitted`) is ",
			"wired up.\n",
			"- `Value(r)` returns the initial value when the form hasn't been ",
			"submitted and the posted value when it has. The ",
			"`state.Has(Name())` check distinguishes the two — `Has` ",
			"recognizes a key that was POSTed, even if empty.",
		),

		wf.HeadlineMedium("Draw: the hidden-input rule"),
		wf.Markdown(
			"BESPA's client runtime gathers form values by walking ",
			"`<input>`, `<textarea>`, and `<select>` elements inside the ",
			"`<form>`. That's the *only* mechanism. If your widget renders ",
			"its UI as a canvas, a `<div>` tree, or a custom element, you ",
			"need a hidden `<input>` (or `<textarea>`) carrying the value ",
			"and using `Name()` as its `name` attribute. Otherwise the form ",
			"submit sees nothing.",
		),
		wf.CodeBlock(formInputsDraw).WithLanguage("go"),

		wf.HeadlineMedium("JS bridge"),
		wf.Markdown(
			"For widgets backed by a JS library, the bridge is: on every ",
			"value change in the JS widget, write the new value into the ",
			"hidden input. The form-submit handler reads from the DOM, so as ",
			"long as the hidden input is current at submit time, BESPA sees ",
			"the right value.",
		),
		wf.CodeBlock(formInputsJS).WithLanguage("js"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Register the companion JS via `widget.AssetRegistry.RegisterScript` ",
			"— see [Extend → Assets & CSS](/extend/assets).",
		),

		wf.HeadlineMedium("Validation"),
		wf.Markdown(
			"`Valid(r)` runs at submit time. The convention used across the ",
			"built-in inputs is to short-circuit when disabled or not yet ",
			"submitted, then check `Required` / length / type rules, then ",
			"delegate to a registered predicate chain so callers can attach ",
			"their own validators via `WithValidator(...)`.",
			"\n\n",
			"Store the error message on the widget so `Draw` can render it ",
			"alongside the field. The built-in `Field` wrapper widget already ",
			"reads `Valid` and surfaces an error label — pair your input with ",
			"`Field` and you get the standard error layout for free.",
		),

		wf.HeadlineMedium("A worked example: richedit"),
		wf.Markdown(
			"The `richedit` package ([github.com/microbus-io/bespa/richedit]",
			"(https://pkg.go.dev/github.com/microbus-io/bespa/richedit)) is ",
			"a separate package — not part of `DefaultFactory` — that wraps ",
			"the [Quill](https://quilljs.com/) editor as a BESPA input. It's ",
			"~300 lines of Go plus the embedded JS/CSS, and a good reference ",
			"for the full shape:",
			"\n\n",
			"- Embeds `*widget.InputWidgetBase[*InputRichTextWidget]`.\n",
			"- Implements `Value`, `Valid`, `Changed`, `Draw` exactly as ",
			"described above.\n",
			"- Renders a hidden `<textarea>` whose `name` matches `Name()` — ",
			"the form sees this textarea at submit time. Quill manages the ",
			"visible UI and writes back into the textarea.\n",
			"- Adds its own validation: max length, min length, predicate ",
			"chain, plus a script-tag / event-handler sanitizer (since the ",
			"value is HTML).\n",
			"- Exposes a small factory (`richedit.RichEditFactory`) that the ",
			"app composes alongside `bespa.DefaultFactory` — see ",
			"[Packaging](/extend/packaging).",
			"\n\n",
			"Read it end-to-end before writing your own — most of the ",
			"questions you'll have are already answered there.",
		),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Widget anatomy](/extend/anatomy) — the underlying ",
			"`Widget` interface and `WidgetBase`. `InputWidgetBase` is the ",
			"input-specific subclass.\n",
			"\n",
			"[Assets & CSS](/extend/assets) — how to register the ",
			"companion JS and CSS.\n",
			"\n",
			"[Packaging as a library](/extend/packaging) — how to ",
			"ship your input as a separate Go package that consumers compose ",
			"into their factory.\n",
			"\n",
			"[Build → Forms & validation](/build/forms) — the ",
			"consumer-side view: how a form built with your custom input ",
			"behaves at submit time.",
		),
	)
	shared.Render(w, r, page)
}
