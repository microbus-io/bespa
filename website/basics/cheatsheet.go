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

// The cheatsheet is a dense one-page BESPA reference. The repository's
// AGENTS.md file at the repo root is a paired copy — when you change one,
// change the other. Keep this dense: an LLM should be able to ingest it as
// a single fetch and write correct BESPA code from it.

const cheatMinimumCode = `package main

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
}`

const cheatFactoryCode = `var wf = struct {
    bespa.DefaultFactory
    chart.ChartFactory       // optional: Apache ECharts
    code.CodeFactory         // optional: Chroma syntax highlighting
    richedit.RichEditFactory // optional: Quill rich text editor
    myorg.MyOrgFactory       // your own widget library
}{}`

const cheatPageSkeletonCode = `func handle(w http.ResponseWriter, r *http.Request) {
    state := wf.StateOf(r)
    page := wf.Page().Add(
        wf.AppBar("Title"),
        // ... widgets ...
    )
    page.Draw(w, r)
}`

const cheatStateAccessorsCode = `state := wf.StateOf(r)
state.Get("x")                    // value or ""
state.Has("x")                    // bool
state.Changed(r, "x")             // changed since last request
state.HasChanges(r, "a", "b")     // any of the keys changed`

const cheatRedrawCode = `wf.HeadlineMedium("Hello, ", state.Get("name"), "!").
    RedrawIfChanged(r, "name")  // redraw when name changes`

const cheatFormsCode = `form := wf.Form().Add(
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
}`

const cheatModalCode = `wf.Modal("modal").Add(
    wf.EmbedHandler(mux.ServeHTTP, r, "GET",
        wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
),
// Anywhere on the page:
wf.Link("?modal=/orders/new").Add("New order")`

const cheatTablesCode = `tbl := wf.Table().
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
wf.PageSizer()`

const cheatCustomWidgetCode = `type GreetingWidget struct {
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
}`

const cheatPackageLayoutTree = `mywidgets/
├── doc.go         // package doc with binary-footprint note
├── factory.go     // factory type + go:embed + RegisterFS
├── mywidget.go    // one widget per file
├── mywidget.css   // optional sibling
└── mywidget.js    // optional sibling`

const cheatFactoryFileCode = `package mywidgets

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
}`

const cheatAssetRegistryCode = `//go:embed *.css *.js
var bundle embed.FS
widget.AssetRegistry.RegisterFS(bundle)            // auto-route by extension

widget.AssetRegistry.RegisterStyle("key", css)     // explicit
widget.AssetRegistry.RegisterScript("key", js)
widget.AssetRegistry.RegisterIsolatedScript("lib", js) // own <script> tag
widget.AssetRegistry.RegisterFile("name.png", b)
widget.AssetRegistry.RegisterHandler("prefix/", h)     // dynamic`

// HandleCheatsheet serves /basics/cheatsheet — the dense one-page reference.
// Same content lives at AGENTS.md in the repo root; agents reaching the repo
// via grep find it there, agents reaching the site via llms.txt find it here.
func HandleCheatsheet(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Cheat sheet"),
		wf.Markdown(
			"A single-page reference for writing BESPA code.",
		),

		wf.HeadlineMedium("The minimum"),
		wf.Markdown("A complete BESPA program:"),
		wf.CodeBlock(cheatMinimumCode).WithLanguage("go"),
		wf.Markdown(
			"Two things every program needs:\n\n",
			"- Register `/bespa/` → `widget.AssetRegistry.ServeHTTP` so CSS, JS, and asset files are served.\n",
			"- A factory value (e.g. `bespa.DefaultFactory{}`) — the bag of widget constructors.\n",
		),

		wf.HeadlineMedium("Composing a factory"),
		wf.Markdown("Mix in optional packages and your own widget libraries:"),
		wf.CodeBlock(cheatFactoryCode).WithLanguage("go"),
		wf.Markdown(
			"The first embedded type's methods shadow same-named methods on deeper embeds, ",
			"so you can override any built-in widget by providing a same-named constructor.",
		),

		wf.HeadlineMedium("Page skeleton"),
		wf.Markdown("Every handler builds a page tree and draws it:"),
		wf.CodeBlock(cheatPageSkeletonCode).WithLanguage("go"),

		wf.HeadlineMedium("Action-URL prefixes"),
		wf.Markdown(
			"Links and form actions are interpreted by the framework:\n\n",
			"| Prefix | Meaning |\n",
			"|---|---|\n",
			"| `?key=value` | Apply to state, partial-redraw the current page |\n",
			"| `?key=` | Clear the key from state |\n",
			"| `?a=1&b=2` | Apply multiple keys atomically |\n",
			"| `/path` | Full navigation to a new handler |\n",
			"| `path` | Relative to the page's `data-location` |\n",
			"| `^?key=` | Apply to the *parent* page (modal-close idiom) |\n",
			"| `^/path` | Navigate the parent page (uses parent's target frame if set) |\n",
			"| `~?key=` | Apply to the *top* page from any nesting depth |\n",
			"| `~/path` | Navigate the top page |\n",
		),

		wf.HeadlineMedium("Reserved state keys"),
		wf.Markdown(
			"| Key | Purpose |\n",
			"|---|---|\n",
			"| `_back` | URL the page should return to (`RedirectBack` reads it) |\n",
			"| `_target` | Named embedded frame to render into |\n",
			"| `_submit` | Name of the submit button that fired the form |\n",
			"| `_changed` | Comma-separated keys that triggered a partial redraw |\n",
		),

		wf.HeadlineMedium("State accessors"),
		wf.CodeBlock(cheatStateAccessorsCode).WithLanguage("go"),

		wf.HeadlineMedium("When to redraw"),
		wf.Markdown("Every widget that depends on a state variable opts in:"),
		wf.CodeBlock(cheatRedrawCode).WithLanguage("go"),
		wf.Markdown(
			"**Never** put `RedrawIfChanged` on an input the user is typing into — ",
			"the cursor will be lost. Put it on the downstream widget that displays the value.\n\n",
			"Variants:\n\n",
			"- `RedrawIfChanged(r, \"k1\", \"k2\", ...)` — any of these change\n",
			"- `RedrawIf(predicate)` — arbitrary predicate\n",
			"- `HideIf(predicate)`, `HideIfEmpty(r, \"k\")`, `HideIfEq(r, \"k\", v)`, `HideIfNotEq(r, \"k\", v)`\n",
		),

		wf.HeadlineMedium("Forms"),
		wf.CodeBlock(cheatFormsCode).WithLanguage("go"),
		wf.Markdown(
			"Built-in validators: `WithRequired`, `WithLength(min, max)`, `WithPattern(regex)`, ",
			"`WithMin`, `WithMax`, `WithPredicate(func(v) (bool, msg))`.\n\n",
			"After submit:\n\n",
			"- `wf.RedirectBack(w, r)` — return to `_back`\n",
			"- `wf.Redirect(w, r, \"/some/path\")` — explicit destination (action-URL grammar applies)\n",
			"- `wf.Redirect(w, r, \"^?modal=\")` — close a parent modal\n\n",
			"Multiple submit buttons: read `state.Get(\"_submit\")` to discriminate.",
		),

		wf.HeadlineMedium("Modal embed pattern"),
		wf.Markdown("In the parent page:"),
		wf.CodeBlock(cheatModalCode).WithLanguage("go"),
		wf.Markdown(
			"In the embedded handler: do work, then `wf.Redirect(w, r, \"^?modal=\")` ",
			"to close. `SidePanel` works the same way.\n\n",
			"`EmbedHandler` decompresses `Content-Encoding` automatically, so a compressed outer mux is fine.",
		),

		wf.HeadlineMedium("Tables"),
		wf.CodeBlock(cheatTablesCode).WithLanguage("go"),

		wf.HeadlineMedium("Custom widgets"),
		wf.CodeBlock(cheatCustomWidgetCode).WithLanguage("go"),
		wf.Markdown(
			"The widget's root rendered element MUST carry `data-id` = `g.ID()` or ",
			"partial redraws silently no-op.",
		),

		wf.HeadlineMedium("Widget package layout"),
		wf.CodeBlock(cheatPackageLayoutTree).WithLanguage("plaintext"),
		wf.Markdown("Factory file:"),
		wf.CodeBlock(cheatFactoryFileCode).WithLanguage("go"),
		wf.Markdown(
			"Glob by extension (`*.css *.js *.woff2`) — never `*` — to keep `.go` source ",
			"files out of the binary.",
		),

		wf.HeadlineMedium("Common widget surface"),
		wf.Markdown(
			"**Headings**: `HeadlineLarge/Medium/Small`, `TitleLarge/Medium/Small`, ",
			"`BodyLarge/Medium/Small`, `LabelLarge/Medium/Small`.\n\n",
			"**Text**: `Code` (inline), `CodeBlock(...).WithLanguage(\"go\")`, `Markdown`, ",
			"`Link`, `Icon`, `TextStyle().WithBold()`/`WithItalic()`/`WithColorPrimary()`.\n\n",
			"**Layout**: `Page`, `AppBar`, `Block`, `Splitter`, `Toolbar().AddLeft().AddRight()`, ",
			"`Deck(1, 2, 4)`, `Card`, `CardOutlined`, `Gallery`, `TabSwitcher`, `Spacer`, ",
			"`SpacerBreak`, `SpacerParagraph`, `Rule`, `PipeSeparator`.\n\n",
			"**Form**: `Form`, `Field`, `InputText`, `InputEmail`, `InputURL`, `InputInteger`, ",
			"`InputDecimal`, `InputDate`, `InputTime`, `InputFile`, `InputHidden`, `InputRichText`, ",
			"`Checkbox`, `Toggle`, `Radio`, `Dropdown`, `RichDropdown`, `FilterChip`, `Rating`, ",
			"`ButtonFilled`, `ButtonTonal`, `ButtonOutlined`, `ButtonText`, `ButtonIcon`.\n\n",
			"**Nav**: `MainMenu`, `NavRail`, `NavDrawer`, `NavStrip`, `NavTarget`, `NavTargetBack`.\n\n",
			"**Tables**: `Table`, `Column`, `Row`, `Cell`, `QuickSearch`, `QuickSearchUnderliner`, ",
			"`Paginator`, `PageSizer`.\n\n",
			"**Surfaces**: `Modal`, `SidePanel`, `GroupingFrame`, `Snackbar`, `InfoBubble`, `InfoLink`.\n\n",
			"**Other**: `Embed`, `EmbedHandler`, `Debugger`.",
		),

		wf.HeadlineMedium("Asset registry"),
		wf.CodeBlock(cheatAssetRegistryCode).WithLanguage("go"),

		wf.HeadlineMedium("Glossary"),
		wf.Markdown(
			"- **State** — the URL query string plus posted form body, accessed via ",
			"`wf.StateOf(r)`. Drives every redraw.\n",
			"- **Action-URL prefix** — the `?`, `^`, `~`, or no-prefix on a link `href` or ",
			"form `action`. Tells the framework whether to apply to state, delegate up, or navigate.\n",
			"- **Partial-redraw protocol** — request: POST with header `Bespa-Fetch: 1` and body ",
			"`…&_changed=k1,k2`. Response: concatenated HTML fragments. See `/basics/incremental`.\n",
			"- **Fragment** — one widget's HTML in a partial-redraw response. Root element carries `data-id`.\n",
			"- **Redraw boundary** — a widget that called `RedrawIfChanged(r, \"k\")`. It's the ",
			"unit re-rendered when `k` changes; children render too as part of it.\n",
			"- **data-id contract** — every widget's root rendered element MUST have `data-id` ",
			"= `widget.ID()`. Without it, partial redraws find no swap target and silently no-op.\n",
			"- **Frame** — a named embedded sub-page (`EmbedHandler(...).WithName(\"x\")`). ",
			"Addressed by `target=\"x\"` or `_target=x`.\n",
			"- **Nested page** — any page embedded inside another via `Modal`, `SidePanel`, ",
			"`GroupingFrame`, or a raw `EmbedHandler`. Has its own isolated state.\n",
		),

		wf.HeadlineMedium("Where to look for more"),
		wf.Markdown(
			"- **[bespa.io/basics](/basics/overview)** — incremental updates, action-URL grammar, nesting, frames, declarative views\n",
			"- **[bespa.io/build](/build/overview)** — handlers, forms, tables, modals, navigation, theming\n",
			"- **[bespa.io/extend](/extend/overview)** — writing widgets, packaging, assets, theming\n",
			"- **[bespa.io/showcase](/showcase/overview)** — every widget live with code\n",
			"- **[bespa.io/llms.txt](/llms.txt)** — machine-readable site index\n",
			"- **pkg.go.dev/github.com/microbus-io/bespa** — godoc\n",
		),
	)
	shared.Render(w, r, page)
}
