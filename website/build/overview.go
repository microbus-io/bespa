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

// HandleOverview is the landing page for the "Build apps" track. It frames
// BESPA as a library and lists the practical-use topics that will be written
// to cover how to wire up real-world UI.
func HandleOverview(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Build"),

		"This track is for developers using BESPA as a library to build an ",
		"application. It focuses on the techniques that come up over and over: ",
		"which widgets to redraw on which state change, how to wire forms and ",
		"validation, how to compose modals and side panels, and the patterns ",
		"that keep code readable as a page grows.",
		wf.HeadlineMedium("Topics"),
		wf.Deck(1, 2, 3).Add(
			topicCard("autorenew", "When to redraw", "/build/redraw",
				"The single most useful skill in BESPA: deciding which widgets opt into "+
					"RedrawIfChanged for which state variables, and which ones never need to. "+
					"Includes the cursor-loss trap on InputText."),
			topicCard("hub", "State patterns", "/build/state",
				"Naming conventions for state variables, scoping them to a page or a "+
					"sub-page, return-link plumbing, and what to put in a session vs. the URL."),
			topicCard("edit", "Forms & validation", "/build/forms",
				"Form widgets, auto-submit, predicate validators, server-side guarding, "+
					"and the round-trip that surfaces inline errors next to the right input."),
			topicCard("table view", "Data tables", "/build/tables",
				"Building sortable, filterable, paginated tables. Quick-search, server-side "+
					"filtering, and connecting a table to a backing data source."),
			topicCard("web asset", "Modals & side panels", "/build/modals",
				"Embedding a page inside a modal or a side panel using the named-frame "+
					"pattern, and the link-action conventions that drive them."),
			topicCard("bolt", "Live data", "/build/live-data",
				"Why BESPA doesn't ship server-push, and the small custom widget that "+
					"plugs SSE or WebSockets into the action-URL flow when you need it."),
			topicCard("alt route", "Navigation patterns", "/build/navigation",
				"Nav rail, drawer, strip, and the back-link plumbing. When to nest a "+
					"sub-section under its own NavTarget vs. flatten."),
			topicCard("palette", "Theming", "/build/theming",
				"Picking dark / light / system mode, choosing a Material key-color palette, "+
					"persisting per-user preferences, and how the resolved tokens reach the browser."),
			topicCard("rocket launch", "Deployment", "/build/deployment",
				"Single-binary packaging, the CSP posture and why it needs unsafe-inline, "+
					"automatic cache-busting, reverse-proxy headers, and compression."),
		),
		wf.SpacerParagraph(),
	)
	shared.Render(w, r, page)
}
