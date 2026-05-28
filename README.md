# BESPA

**Backend Single Page Application** — a Go library for building interactive web apps where all HTML is rendered on the server, state lives in the URL, and a tiny client-side runtime swaps in updated fragments as state changes.

No npm. No bundler. No separate frontend repo. Add one `go get` to your project and ship.

```go
package main

import (
    "net/http"

    "github.com/microbus-io/bespa"
    "github.com/microbus-io/bespa/widget"
)

var wf = bespa.DefaultFactory{}

func handleHome(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    wf.Page().Add(
        wf.AppBar("Hello"),
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

func main() {
    http.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)
    http.HandleFunc("/", handleHome)
    http.ListenAndServe(":8080", nil)
}
```

Type your name in the input — the heading updates as you type. No client-side JavaScript framework, no JSON API, no JSX. The widget that depends on the `name` state variable opts in with `RedrawIfChanged`; everything else stays put.

## Why BESPA

- **Server-side rendering.** Every byte of HTML is produced by your Go code. First request returns a complete, indexable, accessible page.
- **Incremental redraws.** State changes post back as small `?key=value` requests; the server re-renders only the affected widgets and a ~5 KB client swaps the fragments into the DOM. No page flash, no full reload.
- **Strongly typed widgets.** Every widget is a Go struct with chained builder methods. The compiler catches what would be a runtime stack trace in JS.
- **Material Design 3 in the box.** Color tokens, typography scale, elevation, and components ship with the framework. Light/dark and palette switches work without writing CSS.
- **Tiny client runtime.** Total JavaScript needed to run a BESPA app is in the kilobytes, ships with the framework, and you don't write or maintain any of it.
- **Drop-in.** BESPA is a Go package. Embed it in any `net/http` mux, in a Microbus.io service, or behind whatever HTTP middleware you already use.
- **Extensible.** Define a Go struct, write a `Draw` method, register CSS / JS with the asset registry, and your widget composes alongside everything that ships in the box — including published widget libraries.

## How it works

1. **Page tree.** Build a tree of widgets with chained builder methods. `wf.Page().Add(...)` is the root.
2. **State in the URL.** Every page carries a small bag of key-value pairs called *state*, derived from the URL query string. Links like `?q=hello` mutate state; form submissions write state too.
3. **Partial redraws.** When state moves, the browser posts the delta to the server. The server re-renders only the widgets that asked to be notified (`widget.RedrawIfChanged(r, "q")`) and returns those HTML fragments. The client swaps them in by `data-id`. The cursor stays in your input. The scroll stays where it was.
4. **Nesting.** A page can embed another page (modal, side panel, inline frame). Each embedded page has its own state, its own redraw boundary, and can be navigated independently — see [bespa.io/basics/nesting](https://bespa.io/basics/nesting).

## Quickstart

```sh
go get github.com/microbus-io/bespa
```

Then create a `main.go` like the example at the top, and:

```sh
go run .
```

Visit `http://localhost:8080`. Three things to remember:

- Register the asset handler at `/bespa/` — the framework serves its CSS and client JavaScript out of that namespace.
- Pages render to anything implementing `http.ResponseWriter`. You can put them behind any middleware.
- A `bespa.DefaultFactory` value gives you every standard widget. Mix in optional packages (`chart`, `code`, `richedit`) by embedding their factories alongside.

## Repository layout

| Path | Purpose |
| --- | --- |
| `widget/` | Core: `Widget` interface, `WidgetBase[T]`, `PageWidget`, `State`, `AssetRegistry`, the `<script>` client runtime. |
| `basic/` | Primitives — cards, headings, icons, modals, gallery, deck, tab switcher, etc. |
| `form/` | Input widgets — text, dropdowns, checkbox, radio, button, file, rating, chips. |
| `table/` | Data table with sort / filter / page. |
| `nav/` | Navigation widgets — drawer, rail, strip, main menu. |
| `chart/` | Apache ECharts wrapper. Opt-in (large dependency). |
| `chart/maps/` | Country / state-subdivision GeoJSON served via a dynamic handler. |
| `code/` | Server-side syntax highlighting via Chroma. Opt-in. |
| `richedit/` | Rich-text editor — Quill 2 + quill-mention. Opt-in. |
| `css/`, `hct/` | Material color tokens, typography scale, HCT color math. |
| `website/` | The example site that lives at [bespa.io](https://bespa.io). Showcase + learn pages. |

## Documentation

- [AGENTS.md](./AGENTS.md) — one-page cheat sheet for writing BESPA code. The dense version for AI coding agents and quick lookups.
- [CONTRIBUTING.md](./CONTRIBUTING.md) — for working on the framework itself: repo layout, package conventions, what to think about when adding a feature.
- [bespa.io](https://bespa.io) — landing page, hero examples, links to everything below.
- [bespa.io/basics/cheatsheet](https://bespa.io/basics/cheatsheet) — same content as AGENTS.md, served on the site.
- [bespa.io/showcase/overview](https://bespa.io/showcase/overview) — every widget the framework ships with, with code-ready examples.
- [bespa.io/basics](https://bespa.io/basics) — how BESPA works internally: state, incremental updates, nesting pages, frames.
- [bespa.io/build](https://bespa.io/build) — practical recipes for building apps: when to redraw, forms, tables, modals, navigation, theming.
- [bespa.io/extend](https://bespa.io/extend) — writing custom widgets and packaging widget libraries.

## Third-party licenses

See [ATTRIBUTIONS.md](./ATTRIBUTIONS.md) for the full text of every bundled dataset and dependency license — Apache ECharts, Quill, Chroma, Material Symbols, Roboto, Natural Earth, and others.

## License

Apache 2.0. See the file headers and `LICENSE` (or the SPDX identifier in each Go file) for details.
