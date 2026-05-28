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

package build

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

const redrawCorrect = `// HeadlineMedium depends on state.name. Redraw it whenever name changes.
wf.HeadlineMedium("Hello, ", state.Get("name"), "!").
    HideIfEmpty(r, "name").
    RedrawIfChanged(r, "name"),

// InputText posts back to the same state.name. DO NOT redraw it —
// re-rendering the <input> mid-typing loses the cursor focus.
wf.Form().Add(
    wf.InputText("name", "").WithAutoSubmit(true),
),
`

const redrawHrefDepends = `// A button whose href encodes the next count value depends on state.count.
// Without RedrawIfChanged, the button keeps pointing at a stale URL.
wf.ButtonFilled("").Add("+").
    WithHref("?count=" + strconv.Itoa(count+1)).
    RedrawIfChanged(r, "count"),
`

const redrawNever = `// Pure presentation: text never changes. Don't add RedrawIfChanged.
wf.AppBar("Counter"),
wf.HeadlineLarge("Welcome"),
"Choose a number using the buttons below.",
`

// HandleRedraw covers when to opt a widget into RedrawIfChanged.
func HandleRedraw(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("When to redraw"),

		wf.Markdown(
			"Every widget on a BESPA page is drawn on the initial request. ",
			"On subsequent state changes the framework only re-renders widgets ",
			"that have opted in by calling `RedrawIfChanged` ",
			"(or `RedrawIf`, `RedrawIfEq`, or a sibling). ",
			"Picking the right set is the single most useful skill in BESPA.",
		),
		wf.HeadlineMedium("The rule"),
		wf.Markdown(
			"A widget needs `RedrawIfChanged(r, \"x\")` if and only if its ",
			"rendered output depends on the state variable `x`. Output here ",
			"means the bytes the widget writes — not just visible text, but ",
			"href, value, class, hidden/shown state, and child content.",
		),
		wf.HeadlineMedium("Output depends on state"),
		wf.Markdown(
			"The classic case: a widget displays a state value, hides itself ",
			"based on a state value, or computes its own href from a state value.",
		),
		wf.CodeBlock(redrawCorrect).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Both `HideIfEmpty` and the inline `state.Get` reads in the ",
			"children list make this widget depend on `name`. The ",
			"`RedrawIfChanged` makes the redraw happen.",
		),
		wf.HeadlineMedium("The href trap"),
		wf.Markdown(
			"A button whose href encodes the next state — like a counter that ",
			"increments — needs `RedrawIfChanged` just as much as the heading ",
			"that displays the count:",
		),
		wf.CodeBlock(redrawHrefDepends).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Without it, clicking the button once works, but clicking it a ",
			"second time would jump to the same value because the href was ",
			"never recomputed.",
		),
		wf.HeadlineMedium("Don't redraw inputs"),
		wf.Markdown(
			"An `InputText` posts its current value back via the form. If you ",
			"redraw the input itself when its state variable changes, the ",
			"browser replaces the live `<input>` element and the cursor jumps ",
			"to position 0 mid-typing. Wrap the input in a `Form` with ",
			"`WithAutoSubmit(true)` and let the dependent widgets (heading, ",
			"derived values) redraw — the input itself stays mounted.",
		),
		wf.HeadlineMedium("Output doesn't depend on state"),
		wf.Markdown(
			"If nothing about a widget's rendered output references state, ",
			"leave `RedrawIfChanged` off. Static prose, headings, and chrome ",
			"like an `AppBar` usually fall here:",
		),
		wf.CodeBlock(redrawNever).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Adding `RedrawIfChanged` on these wastes bytes — the framework ",
			"dutifully re-emits identical HTML and the client swaps it in for ",
			"nothing.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Basics → Incremental updates](/basics/incremental) ",
			"— the underlying mechanism that makes `RedrawIfChanged` work.",
		),
		wf.Markdown(
			"[Showcase → Data table](/showcase/states) ",
			"— a worked example with many widgets depending on the same state.",
		),
	)
	shared.Render(w, r, page)
}
