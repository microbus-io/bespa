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

package main

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

// --- Snippets shown on the Start page --------------------------------------

const startInstall = `go get github.com/microbus-io/bespa
`

// startHelloWorld is the same runnable program shown on the landing page and
// linked from the README. Keeping it identical here means a developer copying
// from any of the three locations gets a working server.
const startHelloWorld = `package main

import (
    "net/http"

    "github.com/microbus-io/bespa"
    "github.com/microbus-io/bespa/widget"
)

var wf = bespa.DefaultFactory{}

func handleHome(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    page := wf.Page().Add(
        wf.AppBar("Hello"),
        wf.Form().Add(
            wf.InputText("name", "").
                WithPlaceholder("Your name").
                WithAutoSubmit(true),
        ),
        wf.HeadlineMedium("Hello, ", state.Get("name"), "!").
            HideIfEmpty(r, "name").
            RedrawIfChanged(r, "name"),
    )
    page.Draw(w, r)
}

func main() {
    http.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)
    http.HandleFunc("/", handleHome)
    http.ListenAndServe(":8080", nil)
}
`

const startCounter = `func handleCounter(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    count, _ := strconv.Atoi(state.Get("count"))
    wf.Page().Add(
        wf.AppBar("Counter"),
        wf.HeadlineLarge(count).
            RedrawIfChanged(r, "count"),
        wf.Toolbar().AddLeft(
            wf.ButtonFilled("").Add("-").
                WithHref("?count=" + strconv.Itoa(count-1)).
                RedrawIfChanged(r, "count"),
            wf.ButtonFilled("").Add("+").
                WithHref("?count=" + strconv.Itoa(count+1)).
                RedrawIfChanged(r, "count"),
        ),
    ).Draw(w, r)
}
`

const startFactory = `import (
    "github.com/microbus-io/bespa"
    "github.com/microbus-io/bespa/chart"
    "github.com/microbus-io/bespa/code"
)

// Compose the default widgets with optional packages and your own.
var wf = struct {
    bespa.DefaultFactory
    chart.ChartFactory
    code.CodeFactory
    // myorg.MyOrgFactory   // your own widget library
}{}
`

// HandleStart serves the Getting started page — the friendliest entry point
// for a developer (or agent) who has just heard of BESPA and wants to see a
// runnable example, then know where to look next.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Get started"),

		wf.Markdown(
			"BESPA is a Go library — no toolchain, no bundler, no separate ",
			"frontend repo. Add one `go get` to your project, write a handler ",
			"that returns a tree of widgets, and serve it. This page walks ",
			"through the smallest end-to-end example.",
		),
		wf.HeadlineMedium("1. Install"),
		wf.Markdown("Add the module to your project:"),
		wf.CodeBlock(startInstall).WithLanguage("bash"),
		wf.HeadlineMedium("2. Hello, world!"),
		wf.Markdown(
			"Put this in `main.go`. It's a complete program — a name input ",
			"and a greeting that updates as you type:",
		),
		wf.CodeBlock(startHelloWorld).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.ButtonFilled("").
			Add(wf.Icon("play arrow"), " Try it").
			WithHref("?modal=/try/hello"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Then run it and open the page:",
		),
		wf.CodeBlock("go run .\n").WithLanguage("bash"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Visit `http://localhost:8080`. Type your name. Notice the greeting ",
			"updates without a full page reload, the URL changes to ",
			"`?name=…`, and the cursor stays in the input. That's the entire ",
			"BESPA loop: state in the URL, server re-renders affected ",
			"widgets, the client swaps fragments in place.",
		),
		wf.Markdown(
			"Two lines in `main` matter:",
			"\n\n",
			"- `http.HandleFunc(\"/bespa/\", widget.AssetRegistry.ServeHTTP)` — ",
			"required. The framework serves its CSS, JavaScript, fonts, and ",
			"any widget assets out of this prefix. Without it your `<link>` ",
			"and `<script>` tags 404.\n",
			"- `var wf = bespa.DefaultFactory{}` — the bag of widgets you ",
			"already have. `wf.Page`, `wf.Form`, `wf.InputText`, ",
			"`wf.HeadlineMedium`, etc., all come from this one value.",
		),
		wf.HeadlineMedium("3. State drives the page"),
		wf.Markdown(
			"A second handler shows the core idiom: every link is a state ",
			"change, every widget that depends on the state opts in to redraw, ",
			"and nothing else moves. A counter:",
		),
		wf.CodeBlock(startCounter).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.ButtonFilled("").
			Add(wf.Icon("play arrow"), " Try it").
			WithHref("?modal=/try/counter"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Click `+` or `−` and only the heading and the two buttons ",
			"(whose `href`s change) re-render. The toolbar, AppBar, and the ",
			"rest of the page stay put. See ",
			"[Basics → Incremental updates](/basics/incremental) for the ",
			"underlying protocol.",
		),
		wf.HeadlineMedium("4. Add more widgets"),
		wf.Markdown(
			"`bespa.DefaultFactory` gives you everything in `basic/`, `form/`, ",
			"`table/`, and `nav/` — the standard universe. Heavier widgets ",
			"are opt-in packages: `chart` (Apache ECharts), `code` (Chroma ",
			"syntax highlighting), `richedit` (Quill 2). Compose whichever ",
			"you need into a single factory:",
		),
		wf.CodeBlock(startFactory).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Go's embedded-struct method resolution does the rest — ",
			"`wf.Chart(...)` and `wf.CodeBlock(...)` show up alongside ",
			"`wf.Page(...)` with no conflict. See [Extend → Packaging](/extend/packaging) ",
			"to write your own.",
		),
		wf.HeadlineMedium("Where to go next"),
		wf.Deck(1, 2, 4).Add(
			startCard("Basics", "school", "/basics/overview",
				"How the framework works internally — state, incremental updates, nesting pages, frames."),
			startCard("Build apps", "build", "/build/overview",
				"Practical recipes — handlers and routing, forms, tables, modals, navigation, theming."),
			startCard("Extend", "extension", "/extend/overview",
				"Writing your own widgets and packaging widget libraries."),
			startCard("Showcase", "play arrow", "/showcase/overview",
				"Every widget the framework ships with, live and copyable."),
		),
		wf.SpacerParagraph(),
	)

	// Hero example modal slot, so the "Try it" buttons can open inline.
	page.Add(
		wf.Modal("modal").Add(
			wf.EmbedHandler(mux.ServeHTTP, r, "GET",
				wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
		),
	)

	shared.Render(w, r, page)
}

func startCard(heading, icon, path, desc string) any {
	return wf.CardOutlined().WithHref(path).Add(
		wf.TitleLarge(wf.Icon(icon), " ", heading),
		wf.Spacer(0.25),
		desc,
	)
}
