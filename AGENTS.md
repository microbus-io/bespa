# BESPA agent instructions

You're working with the BESPA framework. There are two possible roles —
identify which one applies, because it affects what else you should read:

1. **Contributing to the BESPA framework itself** — adding or modifying a
   widget, fixing the runtime, editing the framework's docs, working on
   any file under `widget/`, `basic/`, `form/`, `table/`, `nav/`,
   `chart/`, `code/`, `richedit/`, `css/`, `hct/`, or `website/`.
   Read this whole file *and* [CONTRIBUTING.md](./CONTRIBUTING.md). The
   cheat sheet below gives you the user-facing model; CONTRIBUTING.md
   covers repo conventions, packaging rules, and what to think about
   when adding a feature.

2. **Building an app *with* BESPA** — using it as a library in a Go
   project that lives outside this repo. The cheat sheet below is your
   one-page reference; you can ignore CONTRIBUTING.md.

---

The cheat sheet below is also served at
[bespa.io/basics/cheatsheet](https://bespa.io/basics/cheatsheet). Both
are paired copies of the const in `website/basics/cheatsheet.go` — when
one moves, the others should move too.

## The minimum

A complete BESPA program:

```go
package main

import (
    "net/http"

    "github.com/microbus-io/bespa"
    "github.com/microbus-io/bespa/widget"
)

var wf = bespa.DefaultFactory{}

func main() {
    http.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP) // required
    http.HandleFunc("/", home)
    http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    wf.Page().Add(
        wf.AppBar("Hello"),
        wf.Form().Add(wf.InputText("name", "").WithAutoSubmit(true)),
        wf.HeadlineMedium("Hello, ", state.Get("name"), "!").
            HideIfEmpty(r, "name").
            RedrawIfChanged(r, "name"),
    ).Draw(w, r)
}
```

Two things every program needs:

- Register `/bespa/` → `widget.AssetRegistry.ServeHTTP` so CSS, JS, and asset
  files are served.
- A factory value (e.g. `bespa.DefaultFactory{}`) — the bag of widget
  constructors.

## Composing a factory

Mix in optional packages and your own widget libraries:

```go
var wf = struct {
    bespa.DefaultFactory
    chart.ChartFactory       // optional: Apache ECharts
    code.CodeFactory         // optional: Chroma syntax highlighting
    richedit.RichEditFactory // optional: Quill rich text editor
    myorg.MyOrgFactory       // your own widget library
}{}
```

The first embedded type's methods shadow same-named methods on deeper
embeds, so you can override any built-in widget by providing a same-named
constructor.

## Page skeleton

Every handler builds a page tree and draws it:

```go
func handle(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    page := wf.Page().Add(
        wf.AppBar("Title"),
        // ... widgets ...
    )
    page.Draw(w, r)
}
```

## Action-URL prefixes

Links and form actions are interpreted by the framework:

| Prefix | Meaning |
|---|---|
| `?key=value` | Apply to state, partial-redraw the current page |
| `?key=` | Clear the key from state |
| `?a=1&b=2` | Apply multiple keys atomically |
| `/path` | Full navigation to a new handler |
| `path` | Relative to the page's `data-location` |
| `^?key=` | Apply to the *parent* page (modal-close idiom) |
| `^/path` | Navigate the parent page (uses parent's target frame if set) |
| `~?key=` | Apply to the *top* page from any nesting depth |
| `~/path` | Navigate the top page |

## Reserved state keys

| Key | Purpose |
|---|---|
| `_back` | URL the page should return to (`RedirectBack` reads it) |
| `_target` | Named embedded frame to render into |
| `_submit` | Name of the submit button that fired the form |
| `_changed` | Comma-separated keys that triggered a partial redraw |

## State accessors

```go
state := wf.StateOf(r)
state.Get("x")                    // value or ""
state.Has("x")                    // bool
state.Changed(r, "x")             // changed since last request
state.HasChanges(r, "a", "b")     // any of the keys changed
```

## When to redraw

Every widget that depends on a state variable opts in:

```go
wf.HeadlineMedium("Hello, ", state.Get("name"), "!").
    RedrawIfChanged(r, "name")  // redraw when name changes
```

**Never** put `RedrawIfChanged` on an input the user is typing into — the
cursor will be lost. Put it on the downstream widget that displays the
value.

Variants:

- `RedrawIfChanged(r, "k1", "k2", ...)` — any of these change
- `RedrawIf(predicate)` — arbitrary predicate
- `HideIf(predicate)`, `HideIfEmpty(r, "k")`, `HideIfEq(r, "k", v)`,
  `HideIfNotEq(r, "k", v)`

## Forms

```go
form := wf.Form().Add(
    wf.Field().AddLeft("Name").AddRight(
        wf.InputText("name", "").WithRequired(true).WithLength(2, 40),
    ),
    wf.Field().AddLeft("Email").AddRight(
        wf.InputEmail("email", "").WithRequired(true),
    ),
    wf.ButtonFilled("save").Add("Save"),
)

if form.ReadyToCommit(r) {
    values := form.Values(r)
    if err := persist(values); err == nil {
        wf.RedirectBack(w, r)
        return
    }
}
```

Built-in validators: `WithRequired`, `WithLength(min, max)`,
`WithPattern(regex)`, `WithMin`, `WithMax`,
`WithPredicate(func(v) (bool, msg))`.

After submit:

- `wf.RedirectBack(w, r)` — return to `_back`
- `wf.Redirect(w, r, "/some/path")` — explicit destination (action-URL
  grammar applies)
- `wf.Redirect(w, r, "^?modal=")` — close a parent modal

Multiple submit buttons: read `state.Get("_submit")` to discriminate.

## Modal embed pattern

In the parent page:

```go
wf.Modal("modal").Add(
    wf.EmbedHandler(mux.ServeHTTP, r, "GET",
        wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
),
// Anywhere on the page:
wf.Link("?modal=/orders/new").Add("New order")
```

In the embedded handler: do work, then `wf.Redirect(w, r, "^?modal=")` to
close. `SidePanel` works the same way.

`EmbedHandler` decompresses `Content-Encoding` automatically, so a
compressed outer mux is fine.

## Tables

```go
tbl := wf.Table().
    WithDefaultPageRows(r, 25).
    Column("First name").Sortable(true).
    Column("Last name").Sortable(true)

rowFrom, rowTo := tbl.DisplayRange(r)
rows, total, _ := store.Query(
    tbl.Query(r),       // quick-search text
    tbl.SortOrder(r),   // sort column + direction
    rowFrom, rowTo,
)
tbl.WithTotalRows(r, total)

for _, p := range rows {
    tbl.Add(wf.Row().Add(p.FirstName, p.LastName))
}

// Companion widgets bind by table name (default "table"):
wf.Toolbar().AddLeft(wf.QuickSearch()).AddRight(wf.Paginator())
wf.PageSizer()
```

## Custom widgets

```go
type GreetingWidget struct {
    *widget.WidgetBase[*GreetingWidget]
    name string
}

func (f MyFactory) Greeting(name string) *GreetingWidget {
    g := &GreetingWidget{name: name}
    g.WidgetBase = widget.NewWidgetBase(g)
    return g
}

func (g *GreetingWidget) Draw(w io.Writer, r *http.Request) error {
    return factory.Tag("span").
        Class("Greeting").
        Attr("data-id", g.ID()).   // required for partial-redraw swap
        Add("Hello, ", g.name).
        When(g.Shown(r)).          // respects HideIf*
        Draw(w, r)
}
```

The widget's root rendered element MUST carry `data-id` = `g.ID()` or
partial redraws silently no-op.

## Widget package layout

```
mywidgets/
├── doc.go         // package doc with binary-footprint note
├── factory.go     // factory type + go:embed + RegisterFS
├── mywidget.go    // one widget per file
├── mywidget.css   // optional sibling
└── mywidget.js    // optional sibling
```

Factory file:

```go
package mywidgets

import (
    "embed"
    "github.com/microbus-io/bespa/widget"
)

type MyOrgFactory struct{}

var factory = struct {
    widget.WidgetFactory
    MyOrgFactory
}{}

// Convenience aliases — used by Draw methods inside this package.
var Any = factory.Any
var Many = factory.Many
var Text = factory.Text
var HTML = factory.HTML
var HTMLUnsafe = factory.HTMLUnsafe
var Bytes = factory.Bytes
var Tag = factory.Tag

type Widget = widget.Widget
type InputWidget = widget.InputWidget
type BytesWidget = widget.BytesWidget

//go:embed *.css *.js
var bundle embed.FS

func init() {
    widget.AssetRegistry.RegisterFS(bundle)
}
```

Glob by extension (`*.css *.js *.woff2`) — never `*` — to keep `.go` source
files out of the binary.

## Common widget surface

**Headings**: `HeadlineLarge/Medium/Small`, `TitleLarge/Medium/Small`,
`BodyLarge/Medium/Small`, `LabelLarge/Medium/Small`.

**Text**: `Code` (inline), `CodeBlock(...).WithLanguage("go")`, `Markdown`,
`Link`, `Icon`, `TextStyle().WithBold()`/`WithItalic()`/`WithColorPrimary()`.

**Layout**: `Page`, `AppBar`, `Block`, `Splitter`,
`Toolbar().AddLeft().AddRight()`, `Deck(1, 2, 4)`, `Card`, `CardOutlined`,
`Gallery`, `TabSwitcher`, `Spacer`, `SpacerBreak`, `SpacerParagraph`,
`Rule`, `PipeSeparator`.

**Form**: `Form`, `Field`, `InputText`, `InputEmail`, `InputURL`,
`InputInteger`, `InputDecimal`, `InputDate`, `InputTime`, `InputFile`,
`InputHidden`, `InputRichText`, `Checkbox`, `Toggle`, `Radio`, `Dropdown`,
`RichDropdown`, `FilterChip`, `Rating`, `ButtonFilled`, `ButtonTonal`,
`ButtonOutlined`, `ButtonText`, `ButtonIcon`.

**Nav**: `MainMenu`, `NavRail`, `NavDrawer`, `NavStrip`, `NavTarget`,
`NavTargetBack`.

**Tables**: `Table`, `Column`, `Row`, `Cell`, `QuickSearch`,
`QuickSearchUnderliner`, `Paginator`, `PageSizer`.

**Surfaces**: `Modal`, `SidePanel`, `GroupingFrame`, `Snackbar`,
`InfoBubble`, `InfoLink`.

**Other**: `Embed`, `EmbedHandler`, `Debugger`.

## Asset registry

```go
//go:embed *.css *.js
var bundle embed.FS
widget.AssetRegistry.RegisterFS(bundle)            // auto-route by extension

widget.AssetRegistry.RegisterStyle("key", css)     // explicit
widget.AssetRegistry.RegisterScript("key", js)
widget.AssetRegistry.RegisterIsolatedScript("lib", js) // own <script> tag
widget.AssetRegistry.RegisterFile("name.png", b)
widget.AssetRegistry.RegisterHandler("prefix/", h)     // dynamic
```

## Glossary

- **State** — the URL query string plus posted form body, accessed via
  `wf.StateOf(r)`. Drives every redraw.
- **Action-URL prefix** — the `?`, `^`, `~`, or no-prefix on a link `href`
  or form `action`. Tells the framework whether to apply to state, delegate
  up, or navigate.
- **Partial-redraw protocol** — request: POST with header `Bespa-Fetch: 1`
  and body `…&_changed=k1,k2`. Response: concatenated HTML fragments. See
  `/basics/incremental`.
- **Fragment** — one widget's HTML in a partial-redraw response. Root
  element carries `data-id`.
- **Redraw boundary** — a widget that called `RedrawIfChanged(r, "k")`.
  It's the unit re-rendered when `k` changes; children render too as part
  of it.
- **data-id contract** — every widget's root rendered element MUST have
  `data-id` = `widget.ID()`. Without it, partial redraws find no swap
  target and silently no-op.
- **Frame** — a named embedded sub-page (`EmbedHandler(...).WithName("x")`).
  Addressed by `target="x"` or `_target=x`.
- **Nested page** — any page embedded inside another via `Modal`,
  `SidePanel`, `GroupingFrame`, or a raw `EmbedHandler`. Has its own
  isolated state.

## Where to look for more

- **[bespa.io/basics](https://bespa.io/basics/overview)** — incremental updates, action-URL grammar, nesting, frames, declarative views
- **[bespa.io/build](https://bespa.io/build/overview)** — handlers, forms, tables, modals, navigation, theming
- **[bespa.io/extend](https://bespa.io/extend/overview)** — writing widgets, packaging, assets, theming
- **[bespa.io/showcase](https://bespa.io/showcase/overview)** — every widget live with code
- **[bespa.io/llms.txt](https://bespa.io/llms.txt)** — machine-readable site index
- **pkg.go.dev/github.com/microbus-io/bespa** — godoc
