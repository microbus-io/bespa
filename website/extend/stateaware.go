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

const stateAwareDataID = `// factory is the package-local composite factory; see Packaging as a library.
func (g *GreetingWidget) Draw(w io.Writer, r *http.Request) error {
    return factory.Tag("span").
        Class("Greeting").
        Attr("data-id", g.ID()).       // ← required for partial-redraw swap
        Add("Hello, ", g.name, "!").
        When(g.Shown(r)).              // ← respects HideIf*
        Draw(w, r)
}
`

const stateAwareShownDrawn = `// Default behavior — inherited from WidgetBase:
//   Drawn(r) is true on full draws, and on partial draws only if RedrawIf
//   has been called.
//   Shown(r) is true unless HideIf has been called.

// Custom Shown — modal is visible only when its state variable is non-empty:
func (m *ModalWidget) Shown(r *http.Request) bool {
    return m.WidgetBase.Shown(r) && factory.StateOf(r).Get(m.name) != ""
}

// Custom Drawn — modal needs to re-render whenever its state variable moves:
func (m *ModalWidget) Drawn(r *http.Request) bool {
    return m.WidgetBase.Drawn(r) || factory.StateOf(r).Get(m.name) != ""
}
`

const stateAwarePlaceholder = `// When the widget is hidden, the page still needs an anchor with the
// same data-id so a future redraw can swap content back in. Tag.When(false)
// renders <span class="Empty" data-id="…"></span> for that purpose:

return factory.Tag("div").
    Class("MyWidget").
    Attr("data-id", w.ID()).
    Add(w.children).
    When(w.Shown(r)).      // false → render the placeholder span
    Draw(out, r)
`

// HandleStateAware covers Drawn / Shown / RedrawIfChanged in custom widgets.
func HandleStateAware(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("State-aware widgets"),

		wf.Markdown(
			"Two methods on every widget control how it interacts with the ",
			"partial-redraw protocol: `Drawn(r)` says whether the widget should ",
			"be re-rendered on a partial-redraw pass; `Shown(r)` says whether it ",
			"should be visible at all. The defaults from `WidgetBase` cover most ",
			"cases, and `RedrawIfChanged` / `HideIf` are the front-end controls ",
			"callers reach for. This page is about the rare cases where you need ",
			"to override them in a custom widget.",
		),
		wf.HeadlineMedium("The data-id contract"),
		wf.Markdown(
			"For the partial-redraw protocol to work, the widget's root ",
			"rendered element MUST carry a `data-id` attribute set to ",
			"`widget.ID()`. The client uses this to find and replace the element ",
			"when the server returns a redrawn fragment:",
		),
		wf.CodeBlock(stateAwareDataID).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Forget the `data-id` and partial redraws silently no-op — the ",
			"widget renders fine on a full page load but stops updating on state ",
			"changes. This is the single most common bug in a hand-rolled widget.",
		),
		wf.HeadlineMedium("Drawn and Shown defaults"),
		wf.Markdown(
			"`WidgetBase` already implements both. `Drawn` returns true on full ",
			"page draws, and on partial-redraw passes only if `RedrawIfChanged` ",
			"(or a sibling) was called by the consumer. `Shown` returns true ",
			"unless `HideIf` was called. Most widgets never need to override ",
			"either.",
		),
		wf.HeadlineMedium("When to override"),
		wf.Markdown(
			"Override these methods when the widget has its OWN state contract — ",
			"a name it owns and binds to a state variable, like a modal or a ",
			"tab switcher. The modal is the canonical example:",
		),
		wf.CodeBlock(stateAwareShownDrawn).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Both overrides AND with the WidgetBase result — that lets callers ",
			"still use `HideIf` or `RedrawIfChanged` externally if they want ",
			"extra control.",
		),
		wf.HeadlineMedium("Hidden-state placeholders"),
		wf.Markdown(
			"A widget that's currently hidden still needs to leave a marker in ",
			"the DOM so the partial-redraw mechanism can put it back later. ",
			"`Tag.When(false)` already does this for you — when `When` is false ",
			"but the tag has a `data-id`, it emits an invisible ",
			"`<span class=\"Empty\" data-id=\"…\"></span>` instead of nothing:",
		),
		wf.CodeBlock(stateAwarePlaceholder).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This is why every well-behaved widget pipes `Shown(r)` through ",
			"`When` on its outer `Tag` call — the framework handles the ",
			"show/hide round-trip without breaking the DOM swap.",
		),
		wf.HeadlineMedium("Children traversal"),
		wf.Markdown(
			"The framework's traversal pass calls `Children()` on every widget ",
			"to assign stable IDs and collect form inputs. If your widget's ",
			"`Draw` builds a sub-tree internally, return that same sub-tree from ",
			"`Children()` — see Extend → Composing existing widgets for ",
			"the pattern.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Basics → Incremental updates](/basics/incremental) ",
			"— the underlying protocol Drawn / Shown plug into.",
			"\n\n",
			"[Build → When to redraw](/build/redraw) — the ",
			"consumer-side view of the same machinery.",
			"\n\n",
			"Read `basic/modal.go` for a real widget that overrides both methods.",
		),
	)
	shared.Render(w, r, page)
}
