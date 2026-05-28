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

const modalsAtPageLevel = `// mux is your *http.ServeMux (see Handlers & routing).
wf.Page().Add(
    wf.AppBar("Dashboard"),

    // The modal slot. Empty state.modal means hidden.
    wf.Modal("modal").Add(
        wf.EmbedHandler(mux.ServeHTTP, r, "GET",
            wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
    ),

    // Anywhere on the page: open a path inside the modal.
    wf.Link("?modal=/orders/new").Add("New order"),
    wf.ButtonFilled("").
        Add(wf.Icon("settings"), " Settings").
        WithHref("?modal=/settings"),
)
`

const modalsPanel = `wf.Page().Add(
    wf.SidePanel("panel").Add(
        wf.EmbedHandler(mux.ServeHTTP, r, "GET",
            wf.StateOf(r).Get("panel")+"?_back=^?panel=", nil),
    ),

    wf.Link("?panel=/orders/123/details").Add("Show details"),
)
`

const modalsClose = `// Inside the embedded page — close the parent modal:
wf.Link("^?modal=").Add("Cancel"),

// Or a button helper that reads state._back:
wf.ButtonText("").Add("Back").WithHrefBack(),
`

// HandleModals covers modals and side panels.
func HandleModals(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Modals & side panels"),

		wf.Markdown(
			"A modal and a side panel are the same idea in two shapes: an ",
			"overlay that displays another BESPA page, controlled by a single ",
			"state variable. When the variable is empty, the overlay is hidden; ",
			"when it points at a path, that path's handler is invoked and its ",
			"output rendered inside the overlay.",
		),
		wf.HeadlineMedium("Anatomy of a modal"),
		wf.Markdown(
			"Place a `Modal` widget at the top of your page, named by the ",
			"state key it watches. Its child is the embedded page:",
		),
		wf.CodeBlock(modalsAtPageLevel).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Three things to notice:",
			"\n\n",
			"- `wf.Modal(\"modal\")` — the argument names the state variable. ",
			"Use a different name if you have several modals or want to namespace ",
			"by page.\n",
			"- `EmbedHandler` calls back into your mux to render whatever page ",
			"the state value points at. The `?_back=^?modal=` suffix tells the ",
			"embedded page how to dismiss itself.\n",
			"- Any link or button can open the modal by writing the path into ",
			"state — no special widget needed.",
		),
		wf.HeadlineMedium("Side panels"),
		wf.Markdown(
			"A `SidePanel` is the same pattern with a different shape: a ",
			"sliding pane attached to one edge of the viewport. Use it for ",
			"inspectors, contextual help, or any secondary surface that should ",
			"sit next to the main content rather than over it:",
		),
		wf.CodeBlock(modalsPanel).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"You can have both on the same page — each watches its own state ",
			"key. To prevent one from showing while the other opens, write ",
			"`?modal=/foo&panel=` to open the modal while clearing the panel.",
		),
		wf.HeadlineMedium("Closing"),
		wf.Markdown(
			"From inside the embedded page, two idioms close the modal: setting ",
			"the state var to empty (via `^?modal=` — the `^` means \"act on the ",
			"parent page\"), or using a back button that reads the `_back` value ",
			"the modal installed when it opened:",
		),
		wf.CodeBlock(modalsClose).WithLanguage("go"),
		wf.HeadlineMedium("What goes inside"),
		wf.Markdown(
			"The embedded page is a full BESPA page — its own `AppBar`, its own ",
			"form, its own state. The framework keeps the parent and child ",
			"state isolated; a state variable named `name` in the modal doesn't ",
			"collide with one of the same name in the parent.",
		),
		wf.HeadlineMedium("Focus, keyboard, and accessibility"),
		wf.Markdown(
			"The `Modal` widget renders with `role=\"dialog\"` and ",
			"`aria-modal=\"true\"`, so screen readers announce it as a modal ",
			"context. The framework's client runtime moves keyboard focus to ",
			"the first focusable element inside the modal when it opens, and ",
			"`Esc` dismisses the modal by clearing the controlling state ",
			"variable (the same effect as a link to `^?modal=`).",
			"\n\n",
			"`SidePanel` does the same with `role=\"complementary\"` — it's a ",
			"secondary surface, not a blocking dialog, so the page underneath ",
			"remains interactive.",
		),
		wf.HeadlineMedium("Returning data to the parent"),
		wf.Markdown(
			"When a modal completes a task (a form posted, an item selected) ",
			"and you want the result reflected on the parent page, write the ",
			"result into a state variable on the parent at the same time you ",
			"close the modal. The `^?key=value` action-URL prefix does both ",
			"in one move:",
			"\n\n",
			"```go\n",
			"// In the embedded handler, after the user picked an item:\n",
			"wf.Redirect(w, r, \"^?modal=&picked=\"+id)\n",
			"```\n",
			"\n",
			"On the parent, any widget that called `RedrawIfChanged(r, \"picked\")` ",
			"re-renders with the new value. The modal disappears (because its ",
			"state variable went empty) and the parent surfaces the result — ",
			"one round trip.",
		),
		wf.HeadlineMedium("Dirty-state confirmation"),
		wf.Markdown(
			"For modal forms that the user might dismiss by mistake, gate the ",
			"close on whether anything was edited. Track \"dirty\" as a state ",
			"variable the form's inputs flip, and override the dismiss action:",
			"\n\n",
			"```go\n",
			"if wf.StateOf(r).Get(\"dirty\") == \"1\" {\n",
			"    closeHref = \"?confirm=close\"   // open a confirm sub-modal\n",
			"} else {\n",
			"    closeHref = \"^?modal=\"         // close immediately\n",
			"}\n",
			"```\n",
			"\n",
			"The framework doesn't impose this — it's an app-level pattern ",
			"because what counts as \"dirty\" depends on your form.",
		),
		wf.HeadlineMedium("Nested modals"),
		wf.Markdown(
			"A modal *can* open another modal — the embedded page can have ",
			"its own `Modal(\"sub\")` slot watching a different state key. ",
			"It works, but it's almost never the right design: two stacked ",
			"overlays is usually a sign that one of the steps should be a ",
			"full page instead. The convention in this codebase is to keep ",
			"modal nesting to at most one level.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Basics → Embedded pages: overview](/basics/embedded-pages) ",
			"— how modals compare to inline embeds and named frames.",
			"\n\n",
			"[Basics → Targeting frames](/basics/frames) ",
			"— how the `^` and `~` path prefixes work and how to target a ",
			"specific named frame.",
			"\n\n",
			"[Showcase → Overview](/showcase/overview) — every demo can be ",
			"opened in a modal or panel.",
		),
	)
	shared.Render(w, r, page)
}
