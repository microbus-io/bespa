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
	"strconv"

	"github.com/microbus-io/bespa/website/shared"
)

// HandleNesting demonstrates the nesting of a page within another page.
func HandleNesting(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Nesting pages"),

		wf.Markdown(
			"A BESPA page is a self-contained unit: its own AppBar, its own ",
			"state, its own redraw boundary. The framework lets you embed one ",
			"page inside another so you can compose larger UIs out of those ",
			"units. The nested page doesn't know — or care — whether it's the ",
			"top of the document or sitting inside a modal frame somewhere.",
		),
		wf.HeadlineMedium("Four modalities"),
		wf.Markdown(
			"The same nested page can appear in any of these containers. Pick ",
			"a button to embed `/basics/nested` using that ",
			"modality — and notice how the embedded page keeps its own counter ",
			"state independent of this outer page:",
		),
		wf.Toolbar().AddLeft(
			wf.ButtonTonal("").
				Add(wf.Icon("web asset"), " Modal").
				WithHref("?modal=1").
				WithDisabled(wf.StateOf(r).Get("modal") == "1").
				RedrawIfChanged(r, "modal"),
			wf.ButtonTonal("").
				Add(wf.Icon("dock to right"), " Side panel").
				WithHref("?panel=1").
				WithDisabled(wf.StateOf(r).Get("panel") == "1").
				RedrawIfChanged(r, "panel"),
			wf.ButtonTonal("").
				Add(wf.Icon("view stream"), " Inline").
				WithHref("?inpage=1").
				WithDisabled(wf.StateOf(r).Get("inpage") == "1").
				RedrawIfChanged(r, "inpage"),
			wf.ButtonTonal("").
				Add(wf.Icon("open in new"), " Full page").
				WithHref("nested"),
		),

		wf.Modal("modal").Add(wf.EmbedHandler(mux.ServeHTTP, r, "GET", "/basics/nested?_back=^?modal=", nil)),
		wf.SidePanel("panel").Add(wf.EmbedHandler(mux.ServeHTTP, r, "GET", "/basics/nested?_back=^?panel=", nil)),

		wf.GroupingFrame("Inline embed").Add(wf.EmbedHandler(mux.ServeHTTP, r, "GET", "/basics/nested?_back=^?inpage=", nil)).
			HideIfEmpty(r, "inpage").
			RedrawIfChanged(r, "inpage"),
		wf.HeadlineMedium("How it works"),
		wf.Markdown(
			"Each container widget — `Modal`, `SidePanel`, or a plain ",
			"`GroupingFrame` — embeds the nested handler with `EmbedHandler`. ",
			"The handler is invoked against an in-memory response recorder, ",
			"the resulting HTML is sliced out of its `<body>` tag, and dropped ",
			"into the container. The nested page renders exactly the same ",
			"widgets it would render on its own URL.",
		),
		wf.Markdown(
			"Two things make the nesting clean:",
			"\n\n",
			"- **Isolated state.** A state variable named `x` in the embedded ",
			"page doesn't collide with one of the same name in the outer page. ",
			"Click the embedded \"Increment\" link a few times in different ",
			"modalities and watch each counter run independently.\n",
			"- **Return-link plumbing.** The `?_back=^?modal=` parameter passed ",
			"to `EmbedHandler` tells the nested page how to dismiss itself: ",
			"when its back button runs `RedirectBack`, the response sends a ",
			"`^?modal=` redirect — \"go up one page, clear modal\" — which ",
			"makes the outer page redraw the modal away.",
		),
		wf.HeadlineMedium("When to use which"),
		wf.Markdown(
			"- `Modal` — for a focused interaction that should pause whatever ",
			"the user was doing. Centered, dimmed background, ",
			"dismiss-via-cancel.\n",
			"- `SidePanel` — for an inspector or contextual surface that sits ",
			"beside the main content. Doesn't dim the page; the user can still ",
			"see and reference what's underneath.\n",
			"- `GroupingFrame` (or any container) — for inline composition. ",
			"The embedded page is part of the document flow.\n",
			"- **Full-page navigation** (no embed) — when the user really is ",
			"leaving for a different screen. The embedded page can also be ",
			"reached this way; its handler doesn't change.",
		),

		wf.Debugger(),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Embedded pages: overview](/basics/embedded-pages) ",
			"— the three embedding mechanisms compared in one table.",
			"\n\n",
			"[Targeting frames](/basics/frames) — how to direct a ",
			"click at a specific named embedded frame rather than the current ",
			"page.",
			"\n\n",
			"[The action-URL pattern](/basics/action-url-pattern) ",
			"— the `^` prefix used in the return-link plumbing.",
			"\n\n",
			"[Build → Modals & side panels](/build/modals) — the ",
			"consumer-side recipes for the most common embeds.",
		),
	)

	shared.Render(w, r, page)
}

// HandleNested is the page rendered inside the nesting demo's modals/panels/inline embed.
func HandleNested(w http.ResponseWriter, r *http.Request) {
	x, _ := strconv.Atoi(wf.StateOf(r).Get("x"))
	if x == 0 {
		x = 1
	}
	back := wf.StateOf(r).Get("_back")

	page := wf.Page().Add(
		wf.AppBar("Nested page"),
		wf.Markdown(
			"This page has no knowledge of what container it is nested inside ",
			"— if any. Its state, navigation, and redraws are entirely ",
			"self-contained.",
		),
		wf.TitleLarge("Page ", x).RedrawIfChanged(r, "x"),
		wf.Link("?x="+strconv.Itoa(x+1)).Add("Increment").RedrawIfChanged(r, "x"),

		wf.Block(
			wf.SpacerParagraph(),
			wf.Markdown(
				"A back arrow shows in the AppBar so the user can return to ",
				"the container. Additional close controls — link, button — ",
				"can also use `WithHrefBack`:",
			),
			wf.Toolbar().AddLeft(
				wf.Link("").Add("Close (link)").WithHrefBack(),
				wf.ButtonText("").Add("Close (button)").WithHrefBack(),
			),
		).HideIf(back == "0"),
	)

	shared.Render(w, r, page)
}
