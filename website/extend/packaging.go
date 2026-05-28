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

package extend

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

const packagingLayout = `myorg-widgets/
├── doc.go              ← package doc with import path, license, footprint
├── factory.go          ← MyOrgFactory + asset registration in init()
├── mywidget.go         ← MyWidgetWidget + builder methods + Draw
├── mywidget.css        ← optional sibling CSS
├── mywidget.js         ← optional sibling JS
└── anotherwidget.go    ← one widget per file
`

const packagingFactory = `package myorg

import (
    "embed"

    "github.com/microbus-io/bespa/basic"
    "github.com/microbus-io/bespa/widget"
)

// MyOrgFactory aggregates the widget constructors of this package.
type MyOrgFactory struct{}

// factory is a private composite the package uses internally to construct
// nested framework widgets.
var factory = struct {
    widget.WidgetFactory
    basic.BasicFactory
    MyOrgFactory
}{}

// Convenience aliases — let the Draw methods in this package write
// Tag(…) and Text(…) instead of factory.Tag(…) and factory.Text(…).
// Not the public entry point; consumers compose MyOrgFactory into their
// own factory struct (see "How consumers compose" below).
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

// Bundle every static asset the package ships. Globbing by extension
// keeps .go source files (and READMEs, .txt, etc.) out of the binary.
// Add or remove extensions to match what the package contains — each
// pattern must match at least one file.
//go:embed *.css *.js
var bundle embed.FS

func init() {
    widget.AssetRegistry.RegisterFS(bundle)
}
`

const packagingConsumer = `import (
    "github.com/microbus-io/bespa"
    "github.com/myorg/myorg-widgets"
)

// Compose MyOrgFactory alongside the framework defaults:
var wf = struct {
    bespa.DefaultFactory
    myorg.MyOrgFactory
}{}

// Now wf.MyWidget(...) is available with full type-safety on the chain.
`

// HandlePackaging covers structuring a widget library.
func HandlePackaging(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Packaging as a library"),

		wf.Markdown(
			"A widget package — whether it's a small in-app helpers file or a ",
			"separately-versioned library — follows a consistent shape so ",
			"consumers can mix it into their factory in two lines. This page is ",
			"the contract.",
		),
		wf.HeadlineMedium("File layout"),
		wf.Markdown(
			"One widget per `.go` file, named after the widget. Optional CSS and ",
			"JS siblings carrying the same base name. A `doc.go` for the package ",
			"doc and a `factory.go` for the factory type and asset registration:",
		),
		wf.CodeBlock(packagingLayout).WithLanguage("plaintext"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This is the layout every package under the framework root uses — ",
			"see `basic/`, `form/`, or any of the optional packages (`chart/`, ",
			"`code/`, `richedit/`).",
		),
		wf.HeadlineMedium("The factory file"),
		wf.Markdown("Four responsibilities, all in one short file:"),
		wf.CodeBlock(packagingFactory).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"- `MyOrgFactory` is the empty type whose methods are the ",
			"constructors for your widgets. Consumers embed it.\n",
			"- The package-local `factory` is the composite your widgets reach ",
			"for internally — to call `factory.Tag`, `factory.Text`, etc.\n",
			"- **Convenience aliases** — package-level `Tag`, `Text`, `Widget`, ",
			"etc. that point at the corresponding entries on the composite ",
			"factory. They let your widgets' Draw methods write `Tag(\"div\")` ",
			"instead of `factory.Tag(\"div\")`. The full set every framework ",
			"package keeps in sync is: `Any`, `Many`, `Text`, `HTML`, ",
			"`HTMLUnsafe`, `Bytes`, `Tag` (functions) and `Widget`, ",
			"`InputWidget`, `BytesWidget` (types). They are exported only ",
			"because Go requires it for cross-file visibility within the ",
			"package — they are not the intended consumer entry point, which is ",
			"the factory composition shown below.\n",
			"- The `init()` call registers every CSS / JS asset the package ",
			"owns. The framework's asset registry handles bundling and delivery.",
		),
		wf.HeadlineMedium("How consumers compose"),
		wf.Markdown(
			"From the calling side, your widgets fall in next to the framework ",
			"ones:",
		),
		wf.CodeBlock(packagingConsumer).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Because `MyOrgFactory` is at the shallowest depth in the composite ",
			"struct, Go's selector resolution picks its methods over any ",
			"same-named method on a more deeply embedded factory. That's how ",
			"you can override a built-in widget by providing a same-named ",
			"constructor.",
		),
		wf.HeadlineMedium("Name collisions"),
		wf.Markdown(
			"If two embedded factories both define a method with the same name ",
			"— say `Button` — Go's selector resolution falls back to the ",
			"*shallowest* match. When two libraries collide at the same depth, ",
			"`go build` reports an ambiguous-selector error. Two ways out:",
			"\n\n",
			"- **Win by depth.** Put the factory you want to override at a ",
			"shallower position in the consumer's composite. The factory ",
			"highest in the struct wins, and BESPA itself doesn't fight you — ",
			"every standard widget is overrideable this way.\n",
			"- **Rename at the package level.** A library that needs to coexist ",
			"with the framework's `Button` can publish its constructor as ",
			"`PrimaryButton`, `IconButton`, or anything namespaced. Keep the ",
			"struct type unambiguous (e.g. `MyOrgButtonWidget`) even when the ",
			"factory method shares a name.",
			"\n\n",
			"For a published library, document the methods it exports. Two ",
			"libraries can be composed without surprises as long as their ",
			"surface-area docs are honest about what names they claim.",
		),
		wf.HeadlineMedium("Versioning"),
		wf.Markdown(
			"Semantic Versioning, the usual rules:",
			"\n\n",
			"- Pin your widget library against a single BESPA major. When the ",
			"framework bumps a major version (renaming `WidgetBase` to take a ",
			"new type parameter, changing the partial-redraw wire format, ",
			"etc.), publish a corresponding major of your library.\n",
			"- Patch and minor versions of your library should keep working ",
			"against any patch / minor of the BESPA major they target. ",
			"Resist adding `WithFoo` methods that only exist on framework ",
			"versions newer than your declared `go.mod` requirement.\n",
			"- Beyond v1, the import path includes the major: ",
			"`github.com/myorg/widgets/v2`. Go's module system enforces this; ",
			"consumers can run v1 and v2 of your library side-by-side under ",
			"different import paths if they really need to.\n",
			"- For pre-1.0 (`v0.x`) releases, expect breakage on every minor ",
			"version and say so in the package doc. Most consumers will pin ",
			"to an exact tag during this phase.",
		),
		wf.HeadlineMedium("Naming the binary cost"),
		wf.Markdown(
			"If your package bundles anything weighty — a JS library, a font, ",
			"large GeoJSON — say so in the package doc:",
		),
		wf.CodeBlock(`// Package mywidget embeds Foo.js (~800 KB, MIT). Importing this
// package adds ~1 MB to your binary. See ATTRIBUTIONS.md.
package mywidget`).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Consumers should know what they're taking on. The pattern in this ",
			"codebase is that anything over ~100 KB earns its own opt-in package ",
			"rather than living in `basic/`.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Assets & CSS](/extend/assets) — the asset-registration API ",
			"the `init()` calls into.",
			"\n\n",
			"[Material theming](/extend/theming) — how widget CSS should ",
			"reference design tokens.",
			"\n\n",
			"For a worked example, read the entirety of `chart/` — it's a ",
			"single-widget package with a JS dependency, CSS, a dynamic-asset ",
			"handler, and a clean factory composition.",
		),
	)
	shared.Render(w, r, page)
}
