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

const stateRead = `state := wf.StateOf(r)
name := state.Get("name")             // empty string if absent
hasName := state.Has("name")
nameChanged := state.Changed("name")  // true during a partial redraw if "name" moved
`

const stateLink = `// Set state.foo to "bar" — the page redraws with the new value.
wf.Link("?foo=bar").Add("Set foo"),

// Clear state.foo by setting it to empty.
wf.Link("?foo=").Add("Clear foo"),

// Multiple variables in one click.
wf.Link("?foo=bar&panel=").Add("Set foo, close panel"),
`

const stateSession = `// website/shared/session.go
type Session struct {
    ID       string
    Theme    string
    Palette  string
    // ... per-session-but-not-URL data
}

func SessionOf(w http.ResponseWriter, r *http.Request) *Session {
    // Issues an HttpOnly SessionID cookie and stores Session in memory.
    ...
}
`

// HandleState covers state patterns.
func HandleState(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("State patterns"),

		wf.Markdown(
			"BESPA's state model is intentionally tiny: every state variable is ",
			"a key in the URL's query string. The browser holds the current ",
			"state, the server reads it on each request, and `RedrawIfChanged` ",
			"decides what re-renders when a value moves.",
		),
		wf.HeadlineMedium("Reading state"),
		wf.Markdown("Every handler starts the same way:"),
		wf.CodeBlock(stateRead).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"`state.Get` returns the empty string for missing keys — this is ",
			"the same shape as `url.Values`, which is what `StateOf` is internally.",
		),
		wf.HeadlineMedium("Writing state"),
		wf.Markdown(
			"Links and form submissions targeting `?key=value` merge into the ",
			"page state and trigger an incremental redraw:",
		),
		wf.CodeBlock(stateLink).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Empty values delete the key. `&` between params lets a single ",
			"click move several variables at once — useful for mutually-exclusive ",
			"UI like \"open modal, close panel.\"",
		),
		wf.HeadlineMedium("Reserved keys"),
		wf.Markdown(
			"Names starting with `_` are reserved by the framework. The ones ",
			"you'll see most:",
			"\n\n",
			"- `_changed` — comma-separated list of state vars that moved on this ",
			"request. `state.HasChanges` checks for it; `state.Changed` queries it.\n",
			"- `_back` — return URL used by `RedirectBack` and the `HrefBack`-style ",
			"button helpers. Conventionally set by the calling page so the destination ",
			"knows where to send the user.\n",
			"- `_target` — names a nested page; see Basics → Targeting frames.\n",
			"- `_submit` — set automatically on form submissions to name the form.",
		),
		wf.HeadlineMedium("Naming conventions"),
		wf.Markdown(
			"State is a flat key/value store, but a couple of conventions keep ",
			"it readable:",
			"\n\n",
			"- Prefix per-table state with the table name — `table_q` for ",
			"quick-search, `table_sort` for the current sort column. This lets ",
			"two tables coexist on one page.\n",
			"- Modal/panel state is conventionally a single key (e.g. `modal`) ",
			"whose value is the path of the embedded page.\n",
			"- Use form-field names that match your data shape — they post back ",
			"through `form.Values(r)` into the same struct field names.",
		),
		wf.HeadlineMedium("URL state vs. session state"),
		wf.Markdown(
			"State that should survive a refresh or a shared-link paste goes in ",
			"the URL. State that's per-user-but-not-shareable goes in the session. ",
			"The showcase uses an in-memory session keyed by an HttpOnly cookie:",
		),
		wf.CodeBlock(stateSession).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"For a real app, swap the in-memory map for whatever persistence ",
			"layer you already have — the API surface (one method per accessor) ",
			"keeps the page code unchanged. See ",
			"[Sessions & auth](/build/sessions-and-auth) for the integration ",
			"pattern with real session libraries.",
		),
		wf.HeadlineMedium("Query string vs. form body"),
		wf.Markdown(
			"`wf.StateOf(r)` reads from both the URL query string and the ",
			"posted form body. When the same key appears in both — typical on ",
			"a form post that has the page's existing query string in the URL ",
			"— **the form body wins.** That's what you want: the user's most ",
			"recent input takes precedence over the URL that put them on the ",
			"page.",
			"\n\n",
			"State written via `state.Set(...)` lives only in the current ",
			"request's view of state; the response either echoes it back as ",
			"new query parameters (on a redirect) or folds it into the next ",
			"redraw fragment's links. The framework handles round-tripping; ",
			"you don't write cookies or persistent stores from `Set`.",
		),
		wf.HeadlineMedium("What survives a navigation"),
		wf.Markdown(
			"Each handler sees only the state encoded in the URL it was called ",
			"with, plus whatever the form body added. **Nothing carries over ",
			"implicitly between page navigations.** If you want a state ",
			"variable to ride along to the next page, put it in the link:",
			"\n\n",
			"```go\n",
			"wf.Link(\"/orders?filter=\"+state.Get(\"filter\")).Add(\"All orders\")\n",
			"```\n",
			"\n",
			"For state that should follow the user across every page — theme, ",
			"locale, login — use a session (above). For state that's just ",
			"\"return to where I was\", use the `_back` convention: pass ",
			"`?_back=<url>` when linking *into* a page, and call ",
			"`wf.RedirectBack(w, r)` to return.",
		),
		wf.HeadlineMedium("Isolated state in nested pages"),
		wf.Markdown(
			"When a page is embedded inside another (modal, side panel, named ",
			"frame), each page has its **own state namespace** — the framework ",
			"isolates them. A state variable named `q` in the modal's URL does ",
			"*not* collide with `q` on the parent page; they're independent.",
			"\n\n",
			"This is what makes `EmbedHandler` safe to point at any handler: ",
			"the embedded page renders as if it were on its own URL, with its ",
			"own query string. State written by the embedded page stays in the ",
			"embedded page's URL until the page asks to write back to the ",
			"parent — that's what the `^?key=` action-URL prefix is for. See ",
			"[Basics → Nesting pages](/basics/nesting) for the full model.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[When to redraw](/build/redraw) — how to react to state moves.",
			"\n\n",
			"[Basics → Incremental updates](/basics/incremental) ",
			"— the underlying protocol.",
		),
	)
	shared.Render(w, r, page)
}
