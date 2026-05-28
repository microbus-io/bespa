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

const dvSwiftUIView = `// SwiftUI
struct GreetingView: View {
    @State var name: String = ""
    var body: some View {
        VStack {
            TextField("Your name", text: $name)
            Text("Hello, \(name)!")
                .font(.headline)
        }
    }
}
`

const dvBespaWidget = `// BESPA
func handleHome(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    wf.Page().Add(
        wf.Form().Add(
            wf.InputText("name", "").
                WithPlaceholder("Your name").
                WithAutoSubmit(true),
        ),
        wf.HeadlineMedium("Hello, ", state.Get("name"), "!").
            HideIfEmpty(r, "name").
            RedrawIfChanged(r, "name"),
    ).Draw(w, r)
}
`

const dvGeneric = `// SwiftUI — opaque return type, fluent chain stays on the concrete type
Text("Hi").font(.title).foregroundColor(.blue)

// BESPA — generic WidgetBase[T] returns the concrete pointer, so chains
// stay typed all the way through.
wf.HeadlineMedium("Hi").
    WithColorPrimary().
    RedrawIfChanged(r, "name")
`

// HandleDeclarativeViews positions BESPA in the declarative typed view-tree
// family (SwiftUI, Jetpack Compose, Flutter). The audience is a developer
// coming from one of those frameworks who wants confirmation that their
// mental model transfers.
func HandleDeclarativeViews(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Declarative views"),

		wf.Markdown(
			"If you've written **SwiftUI**, **Jetpack Compose**, or **Flutter**, ",
			"BESPA's programming model will feel familiar. All four belong to ",
			"the same family — *declarative typed view-trees* — where:",
			"\n\n",
			"- Views/widgets are **typed value structures**, not strings or templates.\n",
			"- You **compose** them in trees with **chained builder methods**.\n",
			"- **State changes drive rebuilds**; you don't imperatively mutate the UI.\n",
			"- The **framework owns the diff** and decides what to repaint.",
			"\n\n",
			"BESPA is the same shape, applied to the server. The renderer is ",
			"your `net/http` handler instead of Apple's graphics pipeline, and ",
			"the diff travels over HTTP as HTML fragments instead of being ",
			"applied locally — but the *programming model* is the same.",
		),
		wf.HeadlineMedium("Side-by-side"),
		wf.Markdown(
			"A greeting that updates as the user types — in SwiftUI:",
		),
		wf.CodeBlock(dvSwiftUIView).WithLanguage("swift"),
		wf.SpacerBreak(),
		wf.Markdown("And in BESPA:"),
		wf.CodeBlock(dvBespaWidget).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Different language, different runtime, same shape: a function ",
			"that returns a tree of typed view nodes, where each node opts in ",
			"to the state it cares about. The framework figures out the rest.",
		),
		wf.HeadlineMedium("The mental-model map"),
		wf.Markdown(
			"| SwiftUI / Compose / Flutter | BESPA |\n",
			"|---|---|\n",
			"| `struct ContentView: View` / `@Composable fun X()` / `Widget build()` | A Go handler that returns a widget tree |\n",
			"| `var body: some View { ... }` | `wf.Page().Add(...)` |\n",
			"| `Text(\"Hi\").font(.title)` | `wf.HeadlineMedium(\"Hi\")` |\n",
			"| `VStack { ... }`, `HStack { ... }`, `Row { ... }` | `wf.Toolbar().AddLeft(...)`, `wf.Splitter(...).AddLeft(...)` |\n",
			"| `@State var x: String` triggers rebuild | URL state `?x=...` triggers rebuild |\n",
			"| `@Binding`, `@ObservedObject`, `@EnvironmentObject` | `wf.StateOf(r).Get(\"x\")` everywhere |\n",
			"| The framework diffs the view tree and repaints | The framework diffs which widgets opted in, swaps fragments by `data-id` |\n",
			"| Hot reload via Xcode preview | Hot reload via `go run .` |\n",
			"| Strongly typed: refactor a `.font()` rename and the compiler catches it | Strongly typed: refactor a `WithColorPrimary()` and the compiler catches it |\n",
		),
		wf.HeadlineMedium("The generic-self-typing trick"),
		wf.Markdown(
			"SwiftUI's `some View` keeps `.font(.title)` returning the ",
			"concrete view type, so the chain doesn't lose its identity. ",
			"BESPA does the same thing with Go generics: every widget embeds ",
			"`*WidgetBase[T]` where `T` is the widget's own pointer type, so ",
			"the `WithFoo` methods on the base return `*T`, not `*WidgetBase`. ",
			"The chain stays typed all the way to the end:",
		),
		wf.CodeBlock(dvGeneric).WithLanguage("go"),
		wf.HeadlineMedium("Where it diverges"),
		wf.Markdown(
			"The model is the same; the *runtime* isn't. Things to recalibrate:",
			"\n\n",
			"- **Where the rebuild happens.** SwiftUI rebuilds locally on ",
			"`@State` change. BESPA rebuilds on the server after a round-trip ",
			"— a `fetch` posts the state delta, the server re-renders the ",
			"affected widgets, and the client swaps fragments by `data-id`.\n",
			"- **Where state lives.** SwiftUI: a typed reactive graph rooted ",
			"in the view (`@State`, `@Binding`, `@EnvironmentObject`). BESPA: ",
			"the URL query string. Less expressive, but persistent, ",
			"bookmarkable, and trivially debuggable.\n",
			"- **Children syntax.** SwiftUI's `@ViewBuilder` lets you write ",
			"`VStack { Text(\"a\"); Text(\"b\") }` without a separator. Go ",
			"variadic functions get the same effect with `.Add(a, b, c)`.\n",
			"- **No declarative animations** (yet). SwiftUI's `.animation()` ",
			"and `withAnimation { }` have no direct analog — animations in ",
			"BESPA are CSS transitions on the rendered HTML.\n",
			"- **No previews / instant rebuild loop.** Xcode preview, ",
			"Compose preview, and Flutter hot-reload are tighter feedback ",
			"loops than `go run .`. The Go side is still fast (sub-second ",
			"compiles for a typical app) — just not interactive.",
		),
		wf.HeadlineMedium("What transfers immediately"),
		wf.Markdown(
			"If you came from one of those frameworks, the following habits ",
			"work on day one:",
			"\n\n",
			"- **Compose, don't inherit.** Wrap existing widgets to make new ",
			"ones — see [Composing existing widgets](/extend/composing).\n",
			"- **Push state up, pass it down.** Keep state at the page level ",
			"and let widgets read it via `wf.StateOf(r).Get(...)`. Same idea ",
			"as lifting `@State` to a parent view.\n",
			"- **Small focused widgets.** Lots of small widget functions ",
			"beats one big page handler. Same as SwiftUI's `body` discipline.\n",
			"- **The framework owns identity.** Don't try to hold onto ",
			"references to live widgets; treat each render as a fresh tree. ",
			"Same as SwiftUI's value-type View.",
		),
		wf.HeadlineMedium("Related frameworks"),
		wf.Markdown(
			"BESPA isn't a new idea — it applies a well-established programming ",
			"model to Go. Other frameworks in the same architectural family:",
			"\n\n",
			"- **[Blazor Server](https://dotnet.microsoft.com/apps/aspnet/web-apps/blazor)** ",
			"(C#) — the closest match. Typed Razor components, server-side ",
			"render, automatic DOM sync over SignalR. If you've used it, BESPA ",
			"will feel immediately familiar. The main difference is BESPA ",
			"emits all-Go (no Razor template syntax) and patches over plain HTTP.\n",
			"- **[Vaadin Flow](https://vaadin.com/flow)** (Java) — typed Java ",
			"component trees, server-rendered HTML, DOM kept in sync via ",
			"WebSocket. Mature, enterprise-grade. The Java analog of BESPA.\n",
			"- **[Wt](https://www.webtoolkit.eu/wt)** (C++) — pre-dates most of ",
			"this list (2005). Pure-code widget trees in C++ that render as ",
			"HTML, with reactive updates. *Architecturally* identical to BESPA; ",
			"C++ ergonomics held wider adoption back.\n",
			"- **[Phoenix LiveView](https://www.phoenixframework.org/)** ",
			"(Elixir) — adjacent cousin. Server-authoritative, surgical DOM ",
			"patches over WebSocket. Closer to template-first than ",
			"view-tree-first, but the *philosophy* (server owns rendering, ",
			"client is a thin patch applier) is the same.",
		),

		wf.Markdown(
			"**Adjacent but different** — these solve the same end-goal ",
			"(server-rendered HTML with partial updates) without the typed ",
			"view-tree:",
			"\n\n",
			"- **[HTMX](https://htmx.org/)** — HTML-with-attributes rather than ",
			"widgets-that-emit-HTML. Backend-agnostic, smaller surface area, ",
			"zero opinion on your component model.\n",
			"- **[Hotwire / Turbo](https://hotwired.dev/)** (Rails) — partial ",
			"updates via `<turbo-frame>` and `<turbo-stream>`. Template-first.\n",
			"- **[Datastar](https://data-star.dev/)** — reactive signals + SSE ",
			"patches. New (2024), template-first, backend-agnostic.",
		),

		wf.Markdown(
			"**Client-side cousins** — same view-tree shape, but they render in ",
			"the browser rather than the server:",
			"\n\n",
			"- **SwiftUI / Jetpack Compose / Flutter** — the family BESPA ",
			"borrows from. Native rendering, not HTML.\n",
			"- **[Lit](https://lit.dev/)** — typed Web Components with ",
			"template literals. Same chained-builder feel, client-side.\n",
			"- **[Solid.js](https://www.solidjs.com/)**, **[Mithril](https://mithril.js.org/)**, ",
			"**hyperapp** — JSX or builder-style view trees, client-rendered.",
		),

		wf.Markdown(
			"The four-cell positioning: typed-view-tree × server-rendered. ",
			"Blazor Server is the C# cell, Vaadin Flow the Java cell, Wt the ",
			"C++ cell, and BESPA the Go cell. Of those four, the only one ",
			"that emits HTML in pure handler-returns-tree form (no template ",
			"language alongside) is BESPA.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Basics → Incremental updates](/basics/incremental) — the ",
			"diff/repaint protocol, the part of BESPA that doesn't have a ",
			"SwiftUI analog.",
			"\n\n",
			"[Extend → Widget anatomy](/extend/anatomy) — the ",
			"`WidgetBase[T]` generic and how it gives chained methods their ",
			"return type.",
			"\n\n",
			"[Get started](/start) — install and run the equivalent of a ",
			"\"hello, world\" view function.",
		),
	)
	shared.Render(w, r, page)
}
