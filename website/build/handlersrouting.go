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

const handlersMinimum = `package main

import (
    "net/http"

    "github.com/microbus-io/bespa"
    "github.com/microbus-io/bespa/widget"
)

var wf = bespa.DefaultFactory{}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)
    mux.HandleFunc("/", handleHome)
    http.ListenAndServe(":8080", mux)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    wf.Page().Add("Hello, world!").Draw(w, r)
}
`

const handlersChi = `// chi
r := chi.NewRouter()
r.Handle("/bespa/*", widget.AssetRegistry)
r.Get("/orders/{id}", handleOrder)

// gorilla/mux
m := mux.NewRouter()
m.PathPrefix("/bespa/").Handler(widget.AssetRegistry)
m.HandleFunc("/orders/{id}", handleOrder)

// stdlib (Go 1.22+) — method + pattern + path variables
mux := http.NewServeMux()
mux.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)
mux.HandleFunc("GET /orders/{id}", handleOrder)
`

const handlersMiddleware = `// Wrap a handler:
mux.Handle("/admin/", requireRole("admin", adminHandler))

// Wrap the whole mux:
http.ListenAndServe(":8080", logging(compress(mux)))
`

const handlersEmbedTrap = `// mux is your *http.ServeMux (or whatever router you wired above).
wf.Modal("modal").Add(
    wf.EmbedHandler(mux.ServeHTTP, r, "GET",
        wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
)
`

const handlersCSP = `Content-Security-Policy:
  default-src 'self';
  script-src  'self' 'unsafe-inline';
  style-src   'self' 'unsafe-inline';
  font-src    'self' data:;
  img-src     'self' data:;
`

// HandleHandlersRouting covers wiring BESPA handlers into a mux, middleware
// composition, the EmbedHandler integration point, and CSP / CORS notes.
func HandleHandlersRouting(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Handlers & routing"),

		wf.Markdown(
			"A BESPA page is an `http.Handler`. The framework has no router of ",
			"its own — you register handlers against whatever mux you already ",
			"use, and `Page().Draw(w, r)` writes HTML to the response. This ",
			"page is the wiring contract.",
		),
		wf.HeadlineMedium("The minimum"),
		wf.CodeBlock(handlersMinimum).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Two things are required:",
			"\n\n",
			"- **`/bespa/` must route to `AssetRegistry.ServeHTTP`.** The ",
			"trailing slash matters — it's a subtree pattern. Everything the ",
			"framework serves (CSS, JS, fonts, dynamic GeoJSON, anything a ",
			"widget registered) lives under this prefix. Without it, your ",
			"page's `<link>` and `<script>` tags 404.\n",
			"- **The rest is yours.** Page handlers can sit at any path. ",
			"Path parameters work — Go 1.22+ stdlib supports ",
			"`mux.HandleFunc(\"GET /orders/{id}\", ...)` and you read ",
			"`r.PathValue(\"id\")` inside the handler.",
		),
		wf.HeadlineMedium("With another router"),
		wf.Markdown(
			"Anything that speaks `http.Handler` works:",
		),
		wf.CodeBlock(handlersChi).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The path-parameter accessor changes per router; the BESPA-side ",
			"code does not.",
		),
		wf.HeadlineMedium("Middleware"),
		wf.Markdown(
			"Standard `net/http` middleware composes the way it always does. ",
			"Wrap individual handlers, or wrap the whole mux. BESPA is unaware:",
		),
		wf.CodeBlock(handlersMiddleware).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Compression middleware (brotli, gzip, deflate) is the most common ",
			"outer wrapper. `EmbedHandler` notices the `Content-Encoding` ",
			"header on the response from your handler and decompresses ",
			"transparently before parsing — so a mux wrapped in compression ",
			"can still be passed to `EmbedHandler` directly. You don't have to ",
			"hand it the inner unwrapped mux to avoid double-encoding.",
		),
		wf.HeadlineMedium("EmbedHandler — the integration point"),
		wf.Markdown(
			"`EmbedHandler` is how one page renders another inside itself. It ",
			"takes a handler function (typically `mux.ServeHTTP`), the current ",
			"request, an HTTP method, a path, and a body. It invokes the ",
			"handler against an in-memory recorder, extracts the `<body>` ",
			"content, and returns it as a widget you can drop anywhere in your ",
			"page tree:",
		),
		wf.CodeBlock(handlersEmbedTrap).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Three things to know:",
			"\n\n",
			"- The embedded handler runs in-process. No network round-trip; ",
			"it's a function call.\n",
			"- The embedded page renders with its own state (the framework ",
			"isolates the query string).\n",
			"- The `_back=^?modal=` suffix is how the embedded page knows how ",
			"to dismiss itself — see [Nesting pages](/basics/nesting).",
		),
		wf.HeadlineMedium("`mux` in code samples"),
		wf.Markdown(
			"Throughout these docs, `mux` is *your* `*http.ServeMux` (or ",
			"whatever router you used) — the one that knows how to resolve a ",
			"path back to a handler. You're free to name it anything; the ",
			"pages use `mux` because that's what `EmbedHandler` typically ",
			"receives.",
		),
		wf.HeadlineMedium("CSP"),
		wf.Markdown(
			"BESPA emits a small amount of inline `<script>` (the client ",
			"bootstrap) and inline `<style>` (the resolved theme tokens). A ",
			"strict Content-Security-Policy needs `'unsafe-inline'` for both ",
			"— or a nonce/hash plumbing layer, which the framework does not ",
			"currently provide:",
		),
		wf.CodeBlock(handlersCSP).WithLanguage("plaintext"),
		wf.SpacerBreak(),
		wf.Markdown(
			"If your deployment requires nonce-based CSP, the inline ",
			"emissions live in `widget/page.go` — wrap or fork.",
		),
		wf.HeadlineMedium("CORS"),
		wf.Markdown(
			"BESPA renders HTML, not an API. CORS only matters if you embed a ",
			"BESPA page in a cross-origin iframe or accept cross-origin POSTs. ",
			"Add the headers at the middleware layer; nothing in BESPA opts ",
			"in or out.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Basics → Nesting pages](/basics/nesting) — what `EmbedHandler` ",
			"actually does under the hood.",
			"\n\n",
			"[Modals & side panels](/build/modals) — the most common reason to ",
			"embed.",
			"\n\n",
			"`website/main.go` in the bespa repo — the canonical wiring this ",
			"site uses: a logging + brotli wrapper around `*http.ServeMux`, ",
			"`AssetRegistry` at `/bespa/`, and sub-packages each calling ",
			"`Init(mux.ServeMux)` to register on the inner mux.",
		),
	)
	shared.Render(w, r, page)
}
