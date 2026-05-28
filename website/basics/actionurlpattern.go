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

const actionURLQuestion = `// Set state.q to "hello" — partial redraw on the current page.
wf.Link("?q=hello").Add("Search hello"),

// Set several keys at once.
wf.Link("?q=hello&page=1").Add("Reset to page 1"),

// Clear a key by setting it to empty.
wf.Link("?q=").Add("Clear search"),
`

const actionURLCaret = `// Inside a modal — close it by clearing the parent's "modal" state.
wf.Link("^?modal=").Add("Cancel"),

// Inside a nested page — submit to the parent's form handler.
// (Only different from "/orders" when the parent has a target frame; see below.)
wf.Form().WithAction("^/orders").Add(
    // ... fields ...
)
`

const actionURLTilde = `// From any depth — switch to dark mode on the top page.
wf.Link("~?theme=dark").Add("Dark mode"),

// Navigate the top page somewhere new, ignoring all the frames you're in.
wf.Link("~/").Add("Home"),
`

const actionURLPath = `// Navigate the current page to a new URL — full handler invocation.
wf.Link("/orders").Add("All orders"),

// Relative paths resolve against the page's data-location attribute.
wf.Link("123").Add("Open order 123"),

// External link — leaves the site.
wf.Link("https://github.com/microbus-io/bespa").Add("Source"),
`

const actionURLCombos = `// Bookmark for a parent page from a child embed — "close me and go home":
wf.Link("^/").Add("Home"),

// Set a state variable on the top page from any nesting depth:
wf.Link("~?palette=violet").Add("Violet"),

// Cascading: ^^? would delegate twice up (grandparent), but most apps
// don't need this — the modal-and-panel pattern stays 2 levels deep.
`

// HandleActionURLPattern documents the special URL prefixes the framework
// interprets in link hrefs and form actions: ?, ^, ~, and bare paths.
func HandleActionURLPattern(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("The action-URL pattern"),

		wf.Markdown(
			"Every link's `href` and every form's `action` in a BESPA page is ",
			"interpreted by the client-side router. Most of the time the URL ",
			"means what URLs always mean — navigate the browser to that path. ",
			"But four prefix characters give the URL a special meaning: `?`, ",
			"`^`, `~`, and the absence of any of them. Together they're a tiny ",
			"grammar that drives the framework's incremental updates, modal ",
			"dismissal, and cross-frame navigation.",
		),
		wf.HeadlineMedium("? — apply to state, redraw current page"),
		wf.Markdown(
			"A URL beginning with `?` doesn't navigate. It merges the query ",
			"string into the current page's state and triggers a partial ",
			"redraw. The widgets that opted into `RedrawIfChanged` for a ",
			"changed key re-render; everything else stays put. The cursor stays ",
			"in your input, the scroll stays where it was.",
		),
		wf.CodeBlock(actionURLQuestion).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Empty values delete the key. Multiple keys in one click move them ",
			"all atomically — useful for mutually-exclusive UI ",
			"(`?modal=/foo&panel=` opens the modal while clearing the panel).",
		),
		wf.HeadlineMedium("^ — delegate to the parent page"),
		wf.Markdown(
			"When pages are nested (modals, side panels, embedded pages), each ",
			"has its own state and its own redraw boundary. The `^` prefix says ",
			"\"strip me off, then re-process the rest of this URL as if I had ",
			"been emitted by my parent.\"",
		),
		wf.CodeBlock(actionURLCaret).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This is how an embedded page closes itself: the parent owns the ",
			"state variable that controls the modal's visibility, so clearing ",
			"that variable has to happen at the parent's level. `^?modal=` ",
			"reads as \"go up one, set modal to empty.\"",
		),
		wf.Markdown(
			"The form-action example is less obvious — for a regular path like ",
			"`/orders`, you might expect `^/orders` to behave the same as plain ",
			"`/orders`. It does, unless the parent has a target frame set (via ",
			"`WithTarget` on its Page, or a `target=\"…\"` attribute on the ",
			"link/form). In that case, the parent's target is what decides ",
			"which frame receives the response. Submitting with `^/orders` ",
			"walks up to the parent first, picks up that target, and the ",
			"response lands in the named frame — possibly an entirely different ",
			"frame than the modal that originated the click. Without the `^`, ",
			"the response would just replace the modal's content. See ",
			"[Targeting frames](/basics/frames) for the ",
			"frame-resolution machinery.",
		),
		wf.HeadlineMedium("~ — delegate to the top page"),
		wf.Markdown(
			"Same idea, but jumps all the way to the topmost (outermost) page ",
			"regardless of how deeply you're nested. Useful for global controls ",
			"— theme switches, language pickers, sign-out links — that should ",
			"act on the whole window, not on the modal that contains them.",
		),
		wf.CodeBlock(actionURLTilde).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"If there's no nesting, `~` behaves identically to `^` — both ",
			"resolve to the current page, which is already the top.",
		),
		wf.HeadlineMedium("No prefix — navigate"),
		wf.Markdown(
			"Any URL without one of the special prefixes is a real navigation. ",
			"The browser fetches the URL and replaces the page (or, if the link ",
			"has a `target` attribute, replaces the named frame inside the ",
			"current page):",
		),
		wf.CodeBlock(actionURLPath).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Relative URLs resolve against the value of the page's ",
			"`data-location` attribute, which the framework sets to the path ",
			"the page was served at (without trailing slash for non-root paths, ",
			"`/` for the root).",
		),
		wf.HeadlineMedium("Combinations"),
		wf.Markdown(
			"The prefixes stack with everything else. `^` and `~` just shift ",
			"the context, so any URL valid at the new context is valid after ",
			"the prefix:",
		),
		wf.CodeBlock(actionURLCombos).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The shapes you'll actually write most often are `?key=val` (state ",
			"change), `^?key=` (close-modal idiom), and `/some/path` (normal ",
			"navigation). The others are escape hatches.",
		),
		wf.HeadlineMedium("Server-side too"),
		wf.Markdown(
			"The same grammar applies on the server. A handler can emit a ",
			"redirect to `\"^/orders\"` to mean \"close me and navigate the ",
			"parent to /orders\", or to `\"~?theme=dark\"` to set state on the ",
			"top page. The redirect helpers understand the prefixes; you don't ",
			"have to do anything special.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Incremental updates](/basics/incremental) — what ",
			"happens after a `?` URL applies to state.",
			"\n\n",
			"[Nesting pages](/basics/nesting) — what `^` and `~` ",
			"actually walk through.",
			"\n\n",
			"[Targeting frames](/basics/frames) — the related ",
			"`target=\"_parent\"` / `target=\"_top\"` / named-frame mechanism ",
			"that picks WHICH page receives a regular navigation.",
		),
	)
	shared.Render(w, r, page)
}
