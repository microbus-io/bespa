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

// HandleIncremental demonstrates how widgets are dynamically redrawn in response to state changes.
func HandleIncremental(w http.ResponseWriter, r *http.Request) {
	state := wf.StateOf(r)

	page := wf.Page().Add(
		wf.AppBar("Incremental updates"),

		wf.Markdown(
			"Every BESPA page carries a small bag of key-value pairs called ",
			"*state*. When a state value changes, the framework re-renders ",
			"only the widgets that asked to be notified, sends those HTML ",
			"fragments back to the browser, and the client swaps them into the ",
			"DOM. No JS diffing, no full reload, no flash.",
		),
		wf.HeadlineMedium("A live state viewer"),
		wf.Markdown(
			"The grouping frame below displays the current value of the state ",
			"variable named `x`. It opts in to redraws on changes to `x` by ",
			"calling `RedrawIfChanged`:",
		),
		wf.GroupingFrame(`State variable "x" viewer`).
			Add(`x = "`, state.Get("x"), `"`).
			RedrawIfChanged(r, "x"),
		wf.SpacerParagraph(),

		wf.Markdown(
			"State variables move when you click a link whose href starts with ",
			"`?`. Try it:",
		),
		wf.Link("?x=foo").Add(wf.Code("?x=foo")),
		wf.PipeSeparator(),
		wf.Link("?x=bar").Add(wf.Code("?x=bar")),
		wf.PipeSeparator(),
		wf.Link("?x=").Add(wf.Code("?x= (clear)")),
		wf.SpacerParagraph(),

		wf.Markdown(
			"What happened on each click: the browser posted the state delta ",
			"back to the server over a small `fetch` request, the server ",
			"re-rendered only the widgets opted in for `x` (just the viewer ",
			"above), and the client swapped that single fragment into the DOM. ",
			"The rest of the page — including the AppBar, headings, and this ",
			"paragraph — stayed exactly as it was.",
		),
		wf.HeadlineMedium("Show, hide, and redraw"),
		wf.Markdown(
			"A widget can also choose to disappear entirely based on state. ",
			"The secret widget below combines `HideIfNotEq` (visibility) with ",
			"`RedrawIfChanged` (re-render trigger) — both keyed on the variable ",
			"`secret`:",
		),
		wf.Link("?secret=show").Add(wf.Code("?secret=show")),
		wf.PipeSeparator(),
		wf.Link("?secret=hide").Add(wf.Code("?secret=hide")),
		wf.SpacerBreak(),
		wf.GroupingFrame("Secret").
			Add("It is revealed!").
			HideIfNotEq(r, "secret", "show").
			RedrawIfChanged(r, "secret"),
		wf.SpacerBreak(),
		wf.Markdown(
			"When `HideIf` hides a widget, the framework still emits an empty ",
			"placeholder element in its place — so a later state change can ",
			"swap real content back in without breaking the DOM shape. The ",
			"hide isn't `display:none` CSS; it's the widget disappearing ",
			"entirely.",
		),
		wf.HeadlineMedium("Compared to a full reload"),
		wf.Markdown(
			"For contrast, the following link uses a regular URL — no `?` ",
			"prefix — which triggers a full page reload with `x` set in the ",
			"initial query string:",
		),
		wf.Link("incremental?x=baz&_back="+state.Get("_back")).Add(wf.Code("incremental?x=baz")),
		wf.SpacerBreak(),
		wf.Markdown(
			"Both routes get you to the same end state — `x = baz` — but the ",
			"full reload re-renders every widget on the page and the browser ",
			"repaints the whole document. The incremental path touches only ",
			"the viewer.",
		),
		wf.HeadlineMedium("The wire format"),
		wf.Markdown(
			"Underneath, partial redraws are a small protocol you can inspect ",
			"in the browser's network panel — or hit with `curl` for debugging.",
			"\n\n",
			"**Request.** The client posts to the page's URL with the header ",
			"`Bespa-Fetch: 1` and a form body that contains:\n",
			"\n",
			"- the new value of every state variable the page tracks, including ",
			"those that haven't changed (so the server sees a complete state ",
			"snapshot), and\n",
			"- a `_changed` field whose value is a comma-separated list of the ",
			"keys that actually moved.\n",
			"\n",
			"Example body: `name=Ada&theme=dark&_changed=name`. The server uses ",
			"`_changed` to decide which widgets to re-render; `RedrawIfChanged(r, \"name\")` ",
			"is true when `\"name\"` appears in the `_changed` list.",
			"\n\n",
			"**Response.** Plain HTML — a series of fragments concatenated. Each ",
			"fragment is the root element of one redrawn widget, carrying a ",
			"`data-id` attribute that matches the element it replaces in the ",
			"live DOM. A typical response body:\n",
			"\n",
			"```\n",
			"<span data-id=\"abc\" class=\"GroupingFrame\">…new contents…</span>\n",
			"<button data-id=\"def\" class=\"ButtonFilled\">…</button>\n",
			"```\n",
			"\n",
			"The client walks the response, matches each `data-id` against the ",
			"current DOM, and swaps. Anything not in the response stays put.",
			"\n\n",
			"**Why `Bespa-Fetch: 1`.** The header tells the server this is a ",
			"partial-redraw request, not a full page load. Redirect responses ",
			"also use a different envelope when the header is set — the server ",
			"writes `Location: …` as a plain-text response body instead of an ",
			"HTTP 3xx, so the client can follow it without the browser losing ",
			"the in-flight `fetch` context.",
		),
		wf.HeadlineMedium("Live updates and push"),
		wf.Markdown(
			"The protocol above is one-way: the client initiates, the server ",
			"responds. BESPA does not maintain persistent connections — ",
			"there's no built-in websocket or SSE layer, and no managed ",
			"subscription engine.",
			"\n\n",
			"That's deliberate. Persistent connections require the server to ",
			"track who's connected, what they're subscribed to, and how to ",
			"fan messages out — state that lives across requests. BESPA is ",
			"stateless throughout, which is what keeps it small and ",
			"horizontally trivial to scale.",
			"\n\n",
			"You have two options when you need live-ish behavior:",
			"\n\n",
			"- **Polling.** The `Progress` widget already polls on an interval ",
			"— set `WithRefreshURL(...)` and `WithRefreshInterval(...)` and ",
			"it asks the server for an updated value every N seconds. The ",
			"same pattern works for any \"refresh until done\" use case: a ",
			"job status, a deployment indicator, an unread-count badge. The ",
			"server answers each poll statelessly; nothing is held open.\n",
			"- **External push.** When polling isn't tight enough, layer an ",
			"SSE or websocket connection on top. A small client-side widget ",
			"opens the connection to your own endpoint and, on each received ",
			"event, fires a BESPA action — set a state variable (`?ping=1`) ",
			"to trigger a redraw, navigate to a path, open a modal, or call ",
			"`page_fetch` directly with the URL to refresh. The transport ",
			"stays outside the framework; the action it triggers is plain ",
			"BESPA. You own the connection management; BESPA stays stateless.",
		),
		wf.HeadlineMedium("Watch it happen"),
		wf.Markdown(
			"The debugger pinned to the bottom-right corner shows the current ",
			"state, the URLs that have been fetched, and the HTML fragments ",
			"that came back. Open it and repeat the clicks above to see the ",
			"protocol live.",
		),
		wf.Debugger(),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[The action-URL pattern](/basics/action-url-pattern) ",
			"— the full grammar of link / form URLs, of which `?` is one ",
			"prefix.",
			"\n\n",
			"[Build → When to redraw](/build/redraw) — picking which ",
			"widgets opt in to which state variables.",
			"\n\n",
			"[Nesting pages](/basics/nesting) — how the same ",
			"protocol drives modals, side panels, and embedded sub-pages.",
		),
	)

	shared.Render(w, r, page)
}
