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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/microbus-io/bespa/website/resources"
	"github.com/microbus-io/bespa/website/shared"
)

// --- Code samples shown on the landing page ---------------------------------
// Kept as top-level consts so they stay alongside HandleRoot but don't clutter
// the page-construction code.

// tinyHello is the smallest possible BESPA handler: a Page with a single
// string child. Featured right under the hero pitch.
const tinyHello = `func(w http.ResponseWriter, r *http.Request) {
	wf.Page().Add("Hello, World!").Draw(w, r)
}
`

// helloWorld is the full runnable program that backs the Try-it button below
// the hero example. It mirrors HandleTryHello byte-for-byte (modulo the
// outer shared.Render wrapper) so what visitors copy is what they run.
const helloWorld = `package main

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

// counter is the snippet used in the "State drives the page" section. It
// shows the URL-state-changes-drive-partial-redraws idiom: every widget that
// depends on the "count" state variable opts into RedrawIfChanged, so a click
// re-renders just the heading and the two buttons (whose hrefs have moved).
const counter = `state := wf.StateOf(r)
count, _ := strconv.Atoi(state.Get("count"))

page := wf.Page().Add(
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
)
`

// HandleRoot serves the landing page at bespa.io. It introduces the framework
// to developers, demonstrates the core idiom with code blocks, and links to
// live demos that open in modal frames.
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Modal/panel embedding pattern — same as the showcase home so demo
	// links open inline instead of full-page navigation.
	page := wf.Page().Add(
		wf.AppBar("BESPA"),

		wf.Modal("modal").Add(
			wf.EmbedHandler(mux.ServeHTTP, r, "GET",
				wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
		),

		// --- Hero ---------------------------------------------------------
		wf.HeadlineLarge("Build single-page apps with just backend Go."),
		wf.Markdown(
			"Declarative typed widgets in the shape of [SwiftUI, Jetpack ",
			"Compose, or Flutter](/basics/declarative-views) — except the ",
			"renderer is your `net/http` server. BESPA composes a tree of Go ",
			"structs, ships HTML to the browser, and quietly swaps in updated ",
			"fragments as state changes. No JavaScript framework. No build ",
			"step. No bundler. Material Design 3 is the default.",
		),
		wf.CodeBlock(tinyHello).WithLanguage("go"),
		wf.SpacerParagraph(),

		// --- Value props deck ---------------------------------------------
		wf.HeadlineMedium("Why BESPA?"),
		wf.Deck(1, 2, 4).Add(
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("dns"), " Server-side rendering"),
				wf.Spacer(0.25),
				"Every byte of HTML is produced by your Go code. ",
				"The first request returns a complete, indexable, accessible page — not a JS bootstrap shell.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("autorenew"), " Incremental redraws"),
				wf.Spacer(0.25),
				"A click on `?x=foo` posts the state back, the server re-renders only the affected widgets, ",
				"and a ~5 KB client swaps the changed fragments into the DOM. No page flash, no full reload.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("verified"), " Strongly typed widgets"),
				wf.Spacer(0.25),
				"Every widget is a Go struct with chained builder methods. ",
				"Refactor with confidence; the compiler catches what would be a runtime stack trace in JS.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("palette"), " Material Design 3"),
				wf.Spacer(0.25),
				"Color tokens, typography scale, elevation, and components ship in the box. ",
				"Light and dark themes work without you writing a line of CSS.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("bolt"), " Tiny client runtime"),
				wf.Spacer(0.25),
				"The total JavaScript needed to run a BESPA app is in the kilobytes, ",
				"and it ships with the framework. You don't write or maintain any of it.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("library add"), " Drop-in library"),
				wf.Spacer(0.25),
				"BESPA is a Go package. Embed it in any `net/http` mux, in a Microbus.io service, ",
				"or behind whatever HTTP middleware you already use. No new infrastructure.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("rocket launch"), " Fast"),
				wf.Spacer(0.25),
				"Server-rendered HTML reaches the browser in milliseconds — no JS bundle to parse, ",
				"no client-side hydration. State changes ship as minimal HTML patches and Go's ",
				"`net/http` does the rest at native speed.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("extension"), " Extensible"),
				wf.Spacer(0.25),
				"Need a widget that doesn't exist yet? Define a Go struct, write a `Draw` method, ",
				"register its CSS and JS with the asset registry, and it composes alongside ",
				"everything that ships in the box — including your own published libraries.",
			),
			wf.CardOutlined().Add(
				wf.TitleLarge(wf.Icon("accessibility new"), " ARIA compliant"),
				wf.Spacer(0.25),
				"Real semantic HTML, proper ARIA roles on modals, tabs, snackbars, and nav, ",
				"keyboard focus management, and screen-reader-friendly labels on icons. ",
				"Your app gets accessibility for free from the framework.",
			),
		),
		wf.SpacerParagraph(),

		// --- Hello, world! ------------------------------------------------
		wf.HeadlineMedium("Hello, world!"),
		"The smallest BESPA app — a name input and a greeting that updates as you type:",
		wf.SpacerBreak(),
		wf.CodeBlock(helloWorld).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.ButtonFilled("").
			Add(wf.Icon("play arrow"), " Try it").
			WithHref("?modal=/try/hello"),
		wf.SpacerBreak(),
		"That's it. No package.json, no webpack config, no separate frontend repo. ",
		"`go run .` and you have a working interactive app at ",
		wf.Code("http://localhost:8080"), ".",
		wf.SpacerParagraph(),

		// --- The core idiom -----------------------------------------------
		wf.HeadlineMedium("State drives the page"),
		"BESPA's only state-management primitive is the URL's query string. ",
		"Links and form submissions update state; widgets that depend on a state variable ",
		"declare it with ", wf.Code("RedrawIfChanged"),
		" and re-render whenever it moves. A counter looks like this:",
		wf.SpacerBreak(),
		wf.CodeBlock(counter).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.ButtonFilled("").
			Add(wf.Icon("play arrow"), " Try it").
			WithHref("?modal=/try/counter"),
		wf.SpacerBreak(),
		"Click + or − and only the heading and the buttons re-render. ",
		"See ", wf.Link("?modal=/basics/incremental").Add("incremental updates"),
		" for a walkthrough of the underlying mechanism.",
		wf.SpacerParagraph(),

		// --- Live demos ---------------------------------------------------
		wf.HeadlineMedium("See it in action"),
		"Each demo below opens in a modal frame — itself another BESPA widget. ",
		"Or jump to the ", wf.Link("/showcase/overview").Add("full showcase"), ".",
		wf.SpacerBreak(),
		wf.Deck(1, 2, 4).Add(
			demoCard("Forms", "edit", "/showcase/form-input",
				"Every input widget — text, dates, ranges, dropdowns, ratings, the rich-text editor — on one page."),
			demoCard("Validation", "checklist", "/showcase/form-validation",
				"Client-side and server-side validation working together with inline error messages."),
			demoCard("Data table", "table view", "/showcase/states",
				"Sorting, filtering, paging, all driven by the same incremental-redraw plumbing."),
			demoCard("CRUD", "edit note", "/showcase/dir-list",
				"Create / read / update / delete flow over a per-session in-memory directory."),
			demoCard("Charts", "bar chart", "/showcase/charts",
				"Apache ECharts wrapped as bespa widgets, themed against Material design tokens."),
			demoCard("Mermaid", "account tree", "/showcase/mermaid",
				"Mermaid diagrams themed against Material tokens, with optional zoom and pan."),
			demoCard("Code blocks", "code blocks", "/showcase/code",
				"Server-side syntax highlighting via Chroma, retheming live with the rest of the page."),
		),
		wf.SpacerParagraph(),

		// --- Learn more ---------------------------------------------------
		wf.HeadlineMedium("Learn more"),
		wf.Deck(1, 2, 4).Add(
			learnCard("Basics", "school", "/basics/overview",
				"How state changes flow from the browser to the server and back as targeted DOM swaps."),
			learnCard("Build", "build", "/build/overview",
				"Practical techniques for using BESPA as a library — when to redraw, forms, tables, modals."),
			learnCard("Extend", "extension", "/extend/overview",
				"Adding your own widgets to the framework and packaging them as reusable libraries."),
			learnCard("Showcase", "play arrow", "/showcase/overview",
				"Every widget the framework ships with, with code-ready examples."),
		),
		wf.SpacerParagraph(),
	)
	shared.Render(w, r, page)
}

// demoCard builds an outlined card linking to a showcase via the page-level modal.
func demoCard(heading, icon, path, desc string) any {
	return wf.CardOutlined().Add(
		wf.TitleLarge(wf.Icon(icon), " ", heading),
		wf.Spacer(0.25),
		desc,
		wf.SpacerBreak(),
		wf.Link("?modal="+path).
			Add(wf.Icon("web_asset").WithAltText("Open demo in modal")),
		wf.PipeSeparator(),
		wf.Link(path).
			Add(wf.Icon("open_in_new").WithAltText("Open demo in window")),
	)
}

// learnCard builds an outlined card linking to a learn page (always full-window).
func learnCard(heading, icon, path, desc string) any {
	return wf.CardOutlined().WithHref(path).Add(
		wf.TitleLarge(wf.Icon(icon), " ", heading),
		wf.Spacer(0.25),
		desc,
	)
}

// HandleTryHello serves a live, runnable copy of the Hello, world! example
// shown on the landing page. The Go source displayed there reduces to
// exactly this handler, so what visitors see in the modal IS what the code
// produces.
func HandleTryHello(w http.ResponseWriter, r *http.Request) {
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
	shared.Render(w, r, page)
}

// HandleTryCounter serves a runnable copy of the counter example. The handler
// matches the source displayed in the "State drives the page" section — each
// + / − click updates the count state variable and redraws only the heading
// and the recomputed button links.
func HandleTryCounter(w http.ResponseWriter, r *http.Request) {
	state := wf.StateOf(r)
	count, _ := strconv.Atoi(state.Get("count"))
	page := wf.Page().Add(
		wf.AppBar("Counter"),
		wf.HeadlineLarge(count).
			RedrawIfChanged(r, "count"),
		wf.Toolbar().AddLeft(
			wf.ButtonFilled("").Add("−").
				WithHref("?count="+strconv.Itoa(count-1)).
				RedrawIfChanged(r, "count"),
			wf.ButtonFilled("").Add("+").
				WithHref("?count="+strconv.Itoa(count+1)).
				RedrawIfChanged(r, "count"),
		),
	)
	shared.Render(w, r, page)
}

// HandleImages serves an image from the resources directory.
func HandleImages(w http.ResponseWriter, r *http.Request) {
	p := strings.LastIndex(r.URL.Path, "/images/")
	if p >= 0 {
		filePath := strings.ReplaceAll(r.URL.Path[p+1:], "/", string(filepath.Separator))
		b, err := resources.Bundle.ReadFile(filePath)
		if err == nil {
			w.Header().Set("Content-Type", http.DetectContentType(b))
			w.Write(b)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
