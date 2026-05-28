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

package basics

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

// HandleEmbeddedPages is the comparison page that orients the reader across
// the framework's three ways of putting one page inside another: inline
// embeds, overlays (modal / side panel), and named frames. The deep-dive
// pages then explain each in detail.
func HandleEmbeddedPages(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Embedded pages: overview"),

		wf.Markdown(
			"BESPA gives you three ways to render one page inside another. ",
			"They share the same mechanism — `EmbedHandler` invokes a handler ",
			"against an in-memory response recorder, slices the resulting ",
			"`<body>`, and drops it into a container widget. What differs is ",
			"the *container* and the *state contract* around it.",
			"\n\n",
			"Use this page to pick the right one. Each row links to its ",
			"deep-dive.",
		),

		wf.HeadlineMedium("At a glance"),
		wf.Markdown(
			"| Mechanism | Container | Controlled by | Visually | Reach for it when |\n",
			"|---|---|---|---|---|\n",
			"| **Inline embed** | `GroupingFrame`, any container | A boolean-ish state key (`HideIfEmpty`) | Part of document flow | A sub-page belongs in-line — a detail strip, a settings block under a header |\n",
			"| **Modal** | `Modal(\"key\")` | One state key holding the embedded path | Centered overlay, dims background | A focused interaction that should pause the rest of the page |\n",
			"| **Side panel** | `SidePanel(\"key\")` | One state key holding the embedded path | Sliding pane, doesn't dim | A contextual surface the user reads beside the main content |\n",
			"| **Named frame** | `GroupingFrame(...).WithName(\"x\")` | Links that set `WithTarget(\"x\")` | Multiple independent surfaces on one screen | A dashboard, IDE-style split, or chat-and-thread layout — each pane navigates on its own |\n",
		),

		wf.HeadlineMedium("How they relate"),
		wf.Markdown(
			"All four are the same shape under the hood:",
			"\n\n",
			"```go\n",
			"container.Add(wf.EmbedHandler(handler, r, \"GET\", path, nil))\n",
			"```\n",
			"\n",
			"What changes is the *container widget* and *how `path` is decided*:",
			"\n\n",
			"- **Inline embed** — `path` is fixed in code (or computed once). ",
			"The container shows or hides based on a state predicate.\n",
			"- **Modal / Side panel** — `path` is read from a state variable. ",
			"Empty state means the overlay isn't rendered. Setting the state ",
			"to a path opens it; clearing the state closes it.\n",
			"- **Named frame** — `path` is fixed at first paint, but the ",
			"frame is given a name via `WithName`. Subsequent links elsewhere ",
			"on the page can swap the frame's content by targeting that name.",
		),

		wf.HeadlineMedium("State isolation is the same everywhere"),
		wf.Markdown(
			"Whichever mechanism you pick, the embedded page sees an isolated ",
			"slice of state. A variable named `x` in the embed doesn't collide ",
			"with one of the same name in the parent. That's what lets you ",
			"compose pages without coordinating their state vocabularies.",
		),

		wf.HeadlineMedium("Which one am I looking at?"),
		wf.Markdown(
			"A useful disambiguation when you're reading example code:",
			"\n\n",
			"- See `wf.Modal(\"...\")` or `wf.SidePanel(\"...\")` at the top of ",
			"the page → **overlay**. State key in the constructor.\n",
			"- See `wf.GroupingFrame(...).WithName(\"...\")` → **named frame**. ",
			"Look for links with `WithTarget(\"...\")` matching that name.\n",
			"- See `wf.GroupingFrame(...)` without `WithName`, holding an ",
			"`EmbedHandler` → **inline embed**. Usually paired with ",
			"`HideIfEmpty` or a similar predicate.",
		),

		wf.HeadlineMedium("Common traps"),
		wf.Markdown(
			"- **Two modals stacked.** Allowed, almost never right — usually a ",
			"sign one of the steps should be a full page. The convention is at ",
			"most one level of modal nesting.\n",
			"- **Overlay state and frame target collision.** Don't reuse the ",
			"same name for a `Modal` state key and a frame name. They live in ",
			"different namespaces but a reader of your code won't know that.\n",
			"- **Forgetting `_back`.** Overlays install `?_back=^?modal=` (or ",
			"similar) so the embedded page can dismiss itself. Inline embeds ",
			"and named frames usually don't — the dismiss model is different.\n",
			"- **Expecting a named frame to survive a top-level reload.** ",
			"Frame contents are part of the page render. A full reload starts ",
			"the frame at its initial content.",
		),

		wf.HeadlineMedium("Deep dives"),
		wf.Markdown(
			"- [Nesting pages](/basics/nesting) — the four containers ",
			"side by side in one demo (modal, side panel, inline, full page), ",
			"with the `_back` plumbing explained.\n",
			"- [Targeting frames](/basics/frames) — named frames and the ",
			"`_top` / `_parent` / `_blank` / custom-name target keywords.\n",
			"- [Build → Modals & side panels](/build/modals) — the ",
			"consumer-side recipes: opening, closing, returning data, ",
			"dirty-state confirmation.\n",
			"- [The action-URL pattern](/basics/action-url-pattern) ",
			"— the `^` and `~` prefixes that overlays and frames lean on.",
		),
	)

	shared.Render(w, r, page)
}
