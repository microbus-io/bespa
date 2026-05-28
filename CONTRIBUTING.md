# Contributing to BESPA

This file orients human contributors and AI agents to the BESPA codebase
itself — what's where, the conventions in play, and what to think about
when adding a feature. Skim it before editing.

It's a companion to [AGENTS.md](./AGENTS.md), which gives the user-facing
model of the framework. If you're modifying any package in this repo,
read AGENTS.md first for context, then come back here for the
contributor-specific conventions.

## What BESPA is

BESPA = **Backend Single Page Application**. A Go library for building
interactive web apps where:

- All HTML is rendered on the server from a tree of Go widgets.
- State lives in the URL query string. Links and form submissions update it.
- A tiny client-side runtime (`widget/page.js`, ~5 KB) handles partial-redraw
  requests: a state change posts the form back, the server re-renders only
  the affected widgets, and the client swaps the changed fragments into the
  DOM.
- Material Design 3 tokens, typography scale, and components ship in the box.

No npm, no bundler, no separate frontend repo. The framework is a Go package
you import with `go get`.

## Repository layout

| Path | Purpose |
| --- | --- |
| `widget/` | Core: `Widget` interface, `WidgetBase[T]`, `PageWidget`, `State`, `AssetRegistry`, the `<script>` client runtime. |
| `basic/` | Primitives — cards, headings, icons, modal, side panel, gallery, deck, tab switcher, code, link, etc. |
| `form/` | Input widgets — text, dropdowns, checkbox, radio, button, file, rating, chips, etc. |
| `table/` | Data table with sort / filter / page. |
| `nav/` | Navigation widgets — drawer, rail, strip, main menu. |
| `chart/` | Apache ECharts wrapper (Apache 2.0). Optional opt-in import. |
| `chart/maps/` | Country/state subdivision GeoJSON served via a dynamic handler. |
| `code/` | Server-side syntax highlighting via Chroma (MIT). Optional opt-in import. |
| `richedit/` | Rich-text editor — Quill 2 (BSD-3) + quill-mention (MIT). |
| `css/`, `hct/` | Material color tokens, typography scale, HCT color math. |
| `website/` | The example site that lives at `bespa.io`. Showcase + learn pages. |
| `ATTRIBUTIONS.md` | Third-party license notices for everything bundled. |

## Three learning tracks

The Learn section of the website is organized into three parallel tracks. If
you are writing documentation, sample apps, or designing new widgets, keep
this taxonomy in mind — content should fit into one of the three:

### 1. Fundamentals — `/learn/fundamentals/*`

Deep technical material on how BESPA works internally. The audience is a
developer who wants a precise model of incremental rendering, the
state-change protocol, page nesting, and named targets. Existing pages:

- `/learn/fundamentals/overview`
- `/learn/fundamentals/incremental` — state-driven partial redraws
- `/learn/fundamentals/nesting` — embedding pages inside pages
- `/learn/fundamentals/frames` — targeting named frames

### 2. Build apps — `/learn/apps/*`

Practical-use techniques for someone using BESPA as a library to build an
application. The audience is a developer with a UI to ship who wants to know
which patterns work. Topics this track should cover (most still to be
written):

- **When to redraw** — which widgets opt into `RedrawIfChanged` for which
  state variables, and which should never redraw. The single most useful
  skill in BESPA.
- **State patterns** — naming, scoping, `_back` plumbing, what to keep in a
  session vs. the URL.
- **Forms & validation** — form widgets, auto-submit, predicates, server
  guards.
- **Data tables** — sorting, filtering, pagination, backing data sources.
- **Modals & side panels** — embedding pages inside modal/panel frames.
- **Navigation patterns** — when to nest under a `NavTarget` vs. flatten.

### 3. Build widgets — `/learn/widgets/*`

Extending the framework with custom widgets and reusable widget libraries.
The audience is a developer who needs something that doesn't exist yet or
wants to publish a library. Topics this track should cover:

- **Widget anatomy** — the `Widget` interface, embedding `WidgetBase[T]`, the
  typed-builder pattern with generics.
- **Composing existing widgets** — wrapping a Form + InputText pair into a
  reusable search box, etc.
- **Assets & CSS** — `AssetRegistry`, the bespa-namespaced URL convention,
  dynamic handlers (see `chart/maps/` for a worked example).
- **State-aware widgets** — implementing `Drawn` / `Shown` / `RedrawIfChanged`
  so a widget plays correctly with the incremental-redraw protocol.
- **Packaging as a library** — factory type, function/type aliases, embedding.
- **Material theming** — resolving `--md-sys-*` tokens at draw time; live
  retheming on theme switch (see `chart/chart.js` for the pattern).

## Conventions in this codebase

- **One package per widget family.** Anything heavy (Chroma, ECharts, Quill)
  goes in its own sub-package so it's opt-in. The default `bespa.DefaultFactory`
  pulls only the small/universal packages.
- **Builder methods return the typed pointer.** `WidgetBase[T]` is generic so
  chained `WithFoo()` calls keep their concrete type. Don't break this — it
  is what lets callers `.WithA().WithB().WithC()` without re-asserting.
- **Material Design 3 sentence case** for headings — `"Build apps"`, not
  `"Build Apps"`. Reflects M3 guidance.
- **Server-rendered HTML is the truth.** JS should never construct UI. The
  client runtime only swaps server-rendered fragments and forwards state
  changes back. Resist adding client-side state.
- **Bespa-namespaced asset URLs** — everything served by `AssetRegistry` lives
  at `/bespa/<key>`. Grouped assets use sub-paths: `/bespa/maps/usa.json`.
- **Showcase pages double as tests.** Every widget should be exercised on
  some page under `website/showcase/`. If you add a widget, add a showcase.

## When asked to add a feature

1. **Place it in the right package** — primitive → `basic/`, form input →
   `form/`, etc. If it pulls a third-party dependency, give it its own
   package.
2. **Add it to the showcase** under `website/showcase/`.
3. **Update `ATTRIBUTIONS.md`** if it bundles a third-party library or
   dataset, with the full license text.
4. **Update `chart/factory.go`-style doc comments** explaining what binary
   weight the user is taking on by importing it.
5. **Don't write client-side state.** If a feature needs interactivity, it
   should still flow through the state → server → redraw cycle.
