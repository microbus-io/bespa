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
)

// llmsTxt is the content of /llms.txt — a curated, machine-readable index of
// the site's docs aimed at AI agents ingesting the framework as context. The
// convention (introduced by Jeremy Howard / Answer.AI, late 2024) is a single
// markdown file at /llms.txt with H1 + blockquote summary, then H2-grouped
// "- [Title](url): one-line description" entries.
//
// Keep this short enough that an agent can fetch it once and route subsequent
// requests intelligently. Skip helper / fixture routes (frame1, nested,
// receiver, etc.) that have no value out of context.
const llmsTxt = `# BESPA

> Backend Single Page Application — a Go library for building interactive web
> apps where all HTML is rendered server-side from a tree of Go widget structs,
> state lives in the URL query string, and a tiny client-side runtime swaps in
> updated fragments as state changes. No npm, no bundler, no separate frontend
> repo. Material Design 3 ships in the box.

## Getting started
- [Get started](https://bespa.io/start): Install, hello-world, factory composition, where to go next
- [README](https://github.com/microbus-io/bespa/blob/main/README.md): Repo landing page with the same hello-world example
- [GitHub](https://github.com/microbus-io/bespa): Source, releases, issues

## Basics — how the framework works internally
- [Overview](https://bespa.io/basics/overview): Index of the basics track
- [Cheat sheet](https://bespa.io/basics/cheatsheet): One-page reference — factory composition, page skeleton, action-URL grammar, redraw rule, state accessors, form / modal / table / custom-widget skeletons. Designed for LLM ingestion in a single fetch.
- [Declarative views](https://bespa.io/basics/declarative-views): The SwiftUI / Compose / Flutter / Blazor / Vaadin / Wt family — BESPA is the Go cell
- [Incremental updates](https://bespa.io/basics/incremental): How state changes drive partial-redraw fragments to the browser
- [Action-URL pattern](https://bespa.io/basics/action-url-pattern): The ?, ^, ~, and plain-path prefix grammar interpreted by every link and form
- [Embedded pages: overview](https://bespa.io/basics/embedded-pages): Side-by-side comparison of inline embeds, overlays (modal / side panel), and named frames — when to pick which
- [Nesting pages](https://bespa.io/basics/nesting): Embedding a BESPA page inside another (modal, side panel, inline)
- [Targeting frames](https://bespa.io/basics/frames): Named embedded frames; the _top / _parent / custom target keywords

## Build apps — using BESPA as a library
- [Overview](https://bespa.io/build/overview): Index of the build-apps track
- [Handlers & routing](https://bespa.io/build/handlers-and-routing): Mounting BESPA into stdlib / chi / gorilla / microbus mux; middleware, CSP, CORS; the EmbedHandler integration point
- [Sessions & auth](https://bespa.io/build/sessions-and-auth): Integrating a session library (scs / gorilla / your own); auth middleware; login forms; CSRF; OAuth round-trip with _back; per-user preferences
- [Errors](https://bespa.io/build/errors): 404 not-found pages, 500 panic-recovery middleware, form/persist errors, errors inside an embedded handler, logging vs surfacing
- [When to redraw](https://bespa.io/build/redraw): Which widgets opt in to RedrawIfChanged and which never should; the cursor-loss trap on inputs
- [State patterns](https://bespa.io/build/state): Naming, scoping, _back plumbing, URL vs session boundaries
- [Forms & validation](https://bespa.io/build/forms): Validators, predicates, after-submit redirects, modal-form-and-close, multiple submit buttons, persist errors, file uploads
- [Data tables](https://bespa.io/build/tables): Sortable / paginated / quick-search tables with server-side backing stores
- [Modals & side panels](https://bespa.io/build/modals): Open a path inside an overlay using a single state variable
- [Live data](https://bespa.io/build/live-data): Why BESPA stays stateless and ships no server-push mechanism; the small custom-widget pattern (SSE or WebSocket → hidden-link click → normal action-URL flow) that plugs in real-time when you need it
- [Navigation patterns](https://bespa.io/build/navigation): Rails, drawers, strips, sub-menus, back-button plumbing
- [Theming](https://bespa.io/build/theming): Light / dark / system appearance, key-color palettes, persisting the user's choice
- [Deployment](https://bespa.io/build/deployment): Single-binary packaging via go:embed; why BESPA requires 'unsafe-inline' for script-src and style-src (per-instance widget init scripts, inline theme tokens) and why nonces aren't yet plumbed; automatic ?id= cache-busting tied to process start; X-Forwarded-Host/Proto/Prefix; compression as middleware

## Extend — writing custom widgets and packaging libraries
- [Overview](https://bespa.io/extend/overview): Index of the extend track
- [Widget anatomy](https://bespa.io/extend/anatomy): Implementing the Widget interface and the generic WidgetBase[T] embed
- [Composing existing widgets](https://bespa.io/extend/composing): Wrapping a Form + InputText into a self-contained search box
- [Assets & CSS](https://bespa.io/extend/assets): //go:embed and RegisterFS, isolated scripts for large JS libraries, dynamic asset handlers
- [State-aware widgets](https://bespa.io/extend/state-aware): Drawn / Shown / RedrawIfChanged interactions, the data-id contract, hidden-state placeholders
- [Custom form inputs](https://bespa.io/extend/form-input-widgets): The InputWidget contract (Name, Value, Valid, Changed, SetFormName); the hidden-input rule that lets canvas / custom-element widgets participate in form submission; the JS bridge pattern; richedit (Quill) as a worked external-package example
- [Packaging as a library](https://bespa.io/extend/packaging): Factory type, convenience aliases, composing into the consumer's factory
- [Material theming](https://bespa.io/extend/theming): Referencing --md-sys-* tokens in widget CSS; resolving them at draw time for canvas-rendered widgets

## Showcase — every widget live with code
- [Overview](https://bespa.io/showcase/overview): Index of all live widget demos
- [Text formatting](https://bespa.io/showcase/text-formatting): Headlines, body, labels, typography scale
- [Code blocks](https://bespa.io/showcase/code): Server-side syntax highlighting via Chroma
- [Progress](https://bespa.io/showcase/progress): Linear / circular indicators
- [Toolbar](https://bespa.io/showcase/toolbar): Action rows above and below content
- [Gallery](https://bespa.io/showcase/gallery): Responsive image grid
- [Deck of cards](https://bespa.io/showcase/deck): Responsive card layout
- [Tab switcher](https://bespa.io/showcase/tab-switcher): Multi-pane content with state-driven active tab
- [Navigation](https://bespa.io/showcase/navigation): Rails, drawers, strips, main menu
- [Form input](https://bespa.io/showcase/form-input): Every input widget type
- [Form validation](https://bespa.io/showcase/form-validation): Built-in validators and custom predicates
- [Data table](https://bespa.io/showcase/states): Sort, filter, paginate over US states
- [CRUD](https://bespa.io/showcase/dir-list): Full create / read / update / delete flow
- [Charts](https://bespa.io/showcase/charts): Apache ECharts wrapped as BESPA widgets
- [Mermaid](https://bespa.io/showcase/mermaid): Mermaid diagrams (flowchart, sequence, class, state, gantt) with optional zoom/pan

## Optional
- [Attributions](https://github.com/microbus-io/bespa/blob/main/ATTRIBUTIONS.md): Third-party license notices for ECharts, Mermaid, Quill, Chroma, Roboto, Natural Earth, and others
- [License](https://github.com/microbus-io/bespa/blob/main/LICENSE): Apache 2.0
`

// HandleLLMs serves /llms.txt — the LLM-readable site index. Served as
// text/plain (the convention) so agents can ingest it without HTML parsing.
func HandleLLMs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write([]byte(llmsTxt))
}
