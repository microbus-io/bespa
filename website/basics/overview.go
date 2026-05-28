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

// HandleOverview is the landing page for the Fundamentals track.
// It explains the technical model — server-side rendering + incremental
// updates + a state-driven redraw pass — and links to the deep-dive pages.
func HandleOverview(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Basics"),

		"This track walks through what BESPA is doing on every request and ",
		"every interaction. Read this if you want a precise model of incremental ",
		"updates, page nesting, and how a click in the browser turns into a ",
		"targeted DOM swap.",
		wf.HeadlineMedium("Topics"),
		wf.Deck(1, 2, 3).Add(
			topicCard("autorenew", "Incremental updates", "/basics/incremental",
				"State variables drive partial-DOM redraws. The smallest possible patch "+
					"goes over the wire and replaces only the affected widgets."),
			topicCard("link", "The action-URL pattern", "/basics/action-url-pattern",
				"The four URL shapes the framework interprets: ? for state, ^ for the "+
					"parent page, ~ for the top page, and everything else for normal navigation."),
			topicCard("dashboard customize", "Embedded pages: overview", "/basics/embedded-pages",
				"The three mechanisms for putting one page inside another — inline embeds, "+
					"overlays, named frames — side by side, with a table you can use to pick."),
			topicCard("view quilt", "Nesting pages", "/basics/nesting",
				"Embed one BESPA page inside another to compose dashboards, split "+
					"panes, modal flows, or any structure where a sub-page has its own state."),
			topicCard("space dashboard", "Targeting frames", "/basics/frames",
				"Direct a link or form submission at a specific named page in the "+
					"tree using target names — _parent, _top, or a custom frame name."),
		),
		wf.SpacerParagraph(),
	)
	shared.Render(w, r, page)
}
