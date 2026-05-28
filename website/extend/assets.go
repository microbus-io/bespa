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

const assetsFactory = `package mywidget

import (
    "embed"
    "github.com/microbus-io/bespa/widget"
)

// Globbing by extension keeps .go source files (and READMEs, .txt, etc.)
// out of the binary. List every extension the package actually ships.
//go:embed *.css *.js *.woff2
var bundle embed.FS

func init() {
    widget.AssetRegistry.RegisterFS(bundle)
}
`

const assetsExplicit = `// For a single asset with a non-default key, register it explicitly:
//go:embed mywidget.css
var mywidgetCSS string

func init() {
    widget.AssetRegistry.RegisterStyle("mywidget", mywidgetCSS)
}
`

const assetsIsolated = `//go:embed chart-big-lib.js
var chartLib string

func init() {
    // Served at /bespa/chart-big-lib.js with its own <script> tag and
    // its own cache key — does NOT join the aggregated /bespa/script.js bundle.
    widget.AssetRegistry.RegisterIsolatedScript("chart-big-lib", chartLib)
}
`

const assetsHandler = `widget.AssetRegistry.RegisterHandler("maps/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // r.URL.Path is "/bespa/maps/<spec>.json"; serve whatever the spec asks for.
    spec := strings.TrimPrefix(r.URL.Path, "/bespa/maps/")
    spec = strings.TrimSuffix(spec, ".json")
    // ... compute the JSON for spec, write it ...
}))
`

const assetsRegisterFile = `//go:embed flags.png
var flags []byte

func init() {
    widget.AssetRegistry.RegisterFile("flags.png", flags)
}
// Served at /bespa/flags.png with content type sniffed from the bytes.
`

// HandleAssets covers asset registration and the /bespa/ namespace.
func HandleAssets(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Assets & CSS"),

		wf.Markdown(
			"A widget's Go code lives in the source tree; its CSS, JavaScript, ",
			"fonts, images, and any other static assets need to ship to the ",
			"browser somehow. BESPA does this through a single global registry ",
			"that every page emits links to in its `<head>`.",
		),
		wf.HeadlineMedium("The pattern"),
		wf.Markdown(
			"Each widget package embeds its assets at compile time with ",
			"`//go:embed`, hands the resulting `embed.FS` to `RegisterFS`, and ",
			"the framework picks up everything it knows how to serve:",
		),
		wf.CodeBlock(assetsFactory).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"That's the whole thing. `RegisterFS` walks the FS and routes each ",
			"file by extension — `.css` joins `/bespa/style.css`, `.js` joins ",
			"`/bespa/script.js` (or is isolated if it carries a sourcemap ",
			"reference), `.woff2` / `.png` / `.jpg` / `.svg` are served as ",
			"static files at `/bespa/<name>`. Unknown extensions are ignored.",
			"\n\n",
			"List every extension the package actually ships in the directive. ",
			"`//go:embed *` would also work but would bundle the package's `.go` ",
			"source files into the binary — wasted bytes, since `RegisterFS` ",
			"ignores them. Each pattern must match at least one file or the ",
			"build fails, so drop extensions that the package doesn't yet have.",
		),
		wf.HeadlineMedium("Explicit registration"),
		wf.Markdown(
			"When you need a non-default key, asset content built at init time, ",
			"or just one file in the package, register directly:",
		),
		wf.CodeBlock(assetsExplicit).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"`RegisterStyle` / `RegisterScript` take a key — the filename stem ",
			"in the `RegisterFS` path — and the asset content. Use them when ",
			"the auto-derived key isn't what you want.",
		),
		wf.HeadlineMedium("Large libraries: isolated scripts"),
		wf.Markdown(
			"For anything sizeable enough that you don't want to inflate the ",
			"main bundle (Quill, ECharts, Chroma — anything in the hundreds of ",
			"kilobytes or more), register it as an isolated script:",
		),
		wf.CodeBlock(assetsIsolated).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"It gets its own `<script>` tag, its own cache entry, and its own ",
			"cache-buster query string. Pages that don't use the widget don't ",
			"pay for it (apart from the empty `<script>` tag, which the browser ",
			"short-circuits on second load).",
		),
		wf.HeadlineMedium("Static binary assets"),
		wf.Markdown(
			"PNG, WOFF2, JSON, anything else — use `RegisterFile`:",
		),
		wf.CodeBlock(assetsRegisterFile).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The Content-Type is detected from the leading bytes via ",
			"`http.DetectContentType`. A 30-day cache header is set ",
			"automatically — the cache-buster query string on the URL handles ",
			"invalidation across deploys.",
		),
		wf.HeadlineMedium("Dynamic assets: handlers"),
		wf.Markdown(
			"Sometimes the asset is too big to embed wholesale and too varied ",
			"to pre-compute. Register an `http.Handler` that owns a sub-path:",
		),
		wf.CodeBlock(assetsHandler).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Any request under `/bespa/maps/` routes to this handler. ",
			"`chart/maps` does exactly this to serve 200+ country / subdivision ",
			"GeoJSON slices out of two embedded master files.",
		),
		wf.HeadlineMedium("URL conventions"),
		wf.Markdown(
			"Everything the registry serves lives under `/bespa/`. Grouped ",
			"assets use sub-paths: `/bespa/maps/usa.json`, ",
			"`/bespa/icons/heart.svg`. When you mount the registry on your mux, ",
			"use the trailing-slash pattern so all those paths route to it:",
		),
		wf.CodeBlock(`mux.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)`).WithLanguage("go"),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Packaging as a library](/extend/packaging) — the full file ",
			"layout for a publishable widget package.",
			"\n\n",
			"Read `widget/assets.go` for the registry implementation and ",
			"`chart/maps/maps.go` for the dynamic-handler pattern.",
		),
	)
	shared.Render(w, r, page)
}
