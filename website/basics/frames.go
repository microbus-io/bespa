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

// HandleFrames demonstrates targeting a named embedded frame.
func HandleFrames(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Targeting frames"),

		wf.Markdown(
			"When a page has more than one embedded sub-page, you need a way ",
			"to say which one a click should land in. BESPA borrows the idea — ",
			"and the syntax — from classic HTML frames: each embedded page can ",
			"be named, and a link or form can specify a `target` naming where ",
			"its response should go.",
		),
		wf.HeadlineMedium("Two named frames"),
		wf.Markdown(
			"The two grouping frames below contain embeds named `one` and ",
			"`two`. The links above them target one or the other by name — the ",
			"clicked path renders into the matching frame, leaving the rest of ",
			"the page untouched:",
		),
		wf.Link("/basics/frame1").WithTarget("one").Add("Load page A into frame ", wf.Code("one")),
		wf.PipeSeparator(),
		wf.Link("/basics/nested").WithTarget("two").Add("Load page B into frame ", wf.Code("two")),
		wf.SpacerBreak(),
		wf.Splitter(1, 1).AddLeft(
			wf.GroupingFrame("Frame one").Add(
				wf.EmbedHandler(HandleFrameEmpty, r, "GET", "/basics/frameempty", nil).WithName("one"),
			),
		).AddRight(
			wf.GroupingFrame("Frame two").Add(
				wf.EmbedHandler(HandleFrameEmpty, r, "GET", "/basics/frameempty", nil).WithName("two"),
			),
		),
		wf.SpacerParagraph(),

		wf.Markdown(
			"Click into frame `one` first. The loaded page contains its own ",
			"link that targets frame `two` — cross-frame navigation from ",
			"inside a frame. The page in frame `two` has a link targeting ",
			"`_top` that replaces this whole document — escaping the frame ",
			"structure entirely.",
		),
		wf.HeadlineMedium("Target keywords"),
		wf.Markdown(
			"Three special target names are reserved, plus any custom name you ",
			"assign with `EmbedHandler(...).WithName(\"…\")`:",
			"\n\n",
			"- `_top` — replace the outermost page in the browser tab. The ",
			"reader leaves the current frame structure and lands in a fresh ",
			"top-level document.\n",
			"- `_parent` — replace the immediate parent page in the embedding ",
			"hierarchy. Acts on the page that contains me.\n",
			"- `_blank` — open in a new browser tab. The current document is ",
			"unaffected.\n",
			"- Any custom name (e.g. `\"one\"`, `\"two\"`) — finds the embed ",
			"widget with the matching `WithName` anywhere in the page tree and ",
			"renders into it.",
		),
		wf.HeadlineMedium("Default targets"),
		wf.Markdown(
			"A page can declare a default target for all its links and forms ",
			"via `Page.WithTarget(name)`. This is what makes the action-URL ",
			"pattern's `^` prefix interesting: an embedded page that uses `^` ",
			"walks up to its parent, picks up the parent's default target, and ",
			"lands the response there. Without a default target on the parent, ",
			"`^/some/path` behaves the same as `/some/path`.",
		),
		wf.HeadlineMedium("When to use frames"),
		wf.Markdown(
			"Frames shine when a single screen has multiple independent ",
			"surfaces that update on their own schedule — a chat thread next ",
			"to a contact list, a code editor next to a preview pane, an ",
			"inspector beside a main canvas. Each surface is its own BESPA ",
			"page, with its own state and its own redraw boundary, and the ",
			"user can navigate inside one without disturbing the others.",
		),
		wf.Markdown(
			"For single-purpose embeds (one modal at a time, one inline ",
			"editor), see [Nesting pages](/basics/nesting) — ",
			"usually simpler.",
		),

		wf.Debugger(),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Embedded pages: overview](/basics/embedded-pages) ",
			"— the three embedding mechanisms side by side, with a table for ",
			"picking the right one.",
			"\n\n",
			"[Nesting pages](/basics/nesting) — the unnamed-embed ",
			"case, for when there's only one sub-page.",
			"\n\n",
			"[The action-URL pattern](/basics/action-url-pattern) ",
			"— the `^` / `~` prefixes that interact with frame targets.",
		),
	)
	shared.Render(w, r, page)
}

// HandleFrameEmpty renders a placeholder page used as the initial content of
// an unloaded named frame.
func HandleFrameEmpty(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.TextStyle().WithColorDeemphasized().Add("(empty — load a page using one of the links above)"),
	)
	shared.Render(w, r, page)
}

// HandleFrameOne is the demo page loaded into frame "one". It contains a link
// that targets frame "two" — cross-frame navigation.
func HandleFrameOne(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.TitleLarge("Page A"),
		wf.Markdown(
			"This page is loaded into frame `one`. Its own links can target ",
			"other frames:",
		),
		wf.Link("/basics/frame2").WithTarget("two").
			Add("Load page B into frame ", wf.Code("two")),
	)
	shared.Render(w, r, page)
}

// HandleFrameTwo is the demo page loaded into frame "two". It contains a link
// that targets _top — escaping the frame structure.
func HandleFrameTwo(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.TitleLarge("Page B"),
		wf.Markdown(
			"This page is loaded into frame `two`. The link below targets ",
			"`_top` — it replaces the whole document, not just this frame:",
		),
		wf.Link("/basics/frames").WithTarget("_top").
			Add("Reload the frames demo at the top level"),
	)
	shared.Render(w, r, page)
}
