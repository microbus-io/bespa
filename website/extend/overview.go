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

// HandleOverview is the landing page for the "Build widgets" track —
// extending the framework with custom widgets and reusable widget libraries.
func HandleOverview(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Extend"),

		wf.Markdown(
			"This track is for developers extending the framework — writing ",
			"new widgets, composing existing ones into higher-level building ",
			"blocks, or packaging a set of widgets as a reusable library. The ",
			"shape is the same whether the widget is a one-off helper in your ",
			"app or the start of a published library: a Go struct, a Draw ",
			"method, and (often) a sliver of CSS or JS alongside.",
		),
		wf.HeadlineMedium("Topics"),
		wf.Deck(1, 2, 3).Add(
			topicCard("widgets", "Widget anatomy", "/extend/anatomy",
				"The Widget interface, embedding WidgetBase, the typed-builder pattern "+
					"with generics, and how Children / Draw / Drawn / Shown fit together."),
			topicCard("dashboard", "Composing existing widgets", "/extend/composing",
				"Most useful widgets are built from existing ones. Wrapping a Form + "+
					"InputText pair into a self-contained search box, for example, takes "+
					"a dozen lines."),
			topicCard("style", "Assets & CSS", "/extend/assets",
				"Registering CSS, JS, and dynamic handlers with the AssetRegistry. Material "+
					"design tokens, theme-friendly class names, and the bespa-namespaced "+
					"URL convention."),
			topicCard("autorenew", "State-aware widgets", "/extend/state-aware",
				"Implementing Drawn / Shown / RedrawIfChanged so a custom widget plays "+
					"correctly with the incremental-redraw protocol. Common pitfalls and "+
					"the placeholder-span trick."),
			topicCard("edit", "Custom form inputs", "/extend/form-input-widgets",
				"The InputWidget contract, the hidden-input rule, and the JS bridge "+
					"for inputs backed by a third-party library — with richedit as a worked example."),
			topicCard("inventory 2", "Packaging as a library", "/extend/packaging",
				"Structuring a widget package: factory type, function aliases, type "+
					"aliases, asset embedding, and how the framework's default factory "+
					"composes your library in."),
			topicCard("palette", "Material theming", "/extend/theming",
				"How widget CSS should reference --md-sys-* tokens so it themes for free, "+
					"and how canvas-rendered widgets resolve those tokens at draw time and "+
					"re-render on theme change."),
		),
		wf.HeadlineMedium("Reading the source"),
		wf.Markdown(
			"Every package under the framework root is a worked example of ",
			"building a widget library: `basic/` for primitives, `form/` for ",
			"input widgets, `chart/` for an external-library wrapper, `code/` ",
			"for a server-side-processing widget. Pick whichever matches what ",
			"you want to build and read it end-to-end — a typical widget is ",
			"under 200 lines of Go.",
		),
	)
	shared.Render(w, r, page)
}
