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

package showcase

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

// HandleHome renders the showcase overview as a responsive 1/2/4-column deck of
// outlined cards, ordered from simpler concepts to more complex. Each card has
// a leading icon (matching the side-menu icon) and three modality launchers:
// open in the same window, in a modal dialog, or in a side panel.
func HandleHome(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Web widgets showcase"),

		wf.Modal("modal").Add(wf.EmbedHandler(mux.ServeHTTP, r, "GET", wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil)),
		wf.SidePanel("panel").Add(wf.EmbedHandler(mux.ServeHTTP, r, "GET", wf.StateOf(r).Get("panel")+"?_back=^?panel=", nil)),

		wf.Debugger(),
	)

	type screen struct {
		heading string
		icon    string
		desc    string
		path    string
	}

	// Ordered simpler → more complex.
	screens := []screen{
		{
			heading: "Text formatting",
			icon:    "text format",
			desc:    "The Material typography scale and inline text-style helpers — colors, weights, sizes — across the framework's font stack.",
			path:    "/showcase/text-formatting",
		},
		{
			heading: "Code blocks",
			icon:    "code blocks",
			desc:    "Server-side syntax highlighting via Chroma, with token classes mapped to Material design colors so the same markup recolors on theme change.",
			path:    "/showcase/code",
		},
		{
			heading: "Progress widget",
			icon:    "pending",
			desc:    "Determinate and indeterminate progress indicators, plus the finite/infinite progress-status pattern for long-running operations.",
			path:    "/showcase/progress",
		},
		{
			heading: "Toolbar widget",
			icon:    "construction",
			desc:    "An action toolbar with grouped buttons and separators. Buttons that share a parent are automatically gathered into a toolbar without explicit configuration.",
			path:    "/showcase/toolbar",
		},
		{
			heading: "Gallery widget",
			icon:    "collections",
			desc:    "A horizontally-scrolling gallery of cards. Useful for featured-content rows or any list where horizontal swipe beats vertical scroll.",
			path:    "/showcase/gallery",
		},
		{
			heading: "Deck of cards widget",
			icon:    "dashboard",
			desc:    "A responsive grid of cards that re-flows between 1, 2, and more columns depending on the viewport width. Good for dashboards or any uniform collection.",
			path:    "/showcase/deck",
		},
		{
			heading: "Tab switcher widget",
			icon:    "tab",
			desc:    "Tabbed content with smooth switching between sections. Each tab body is rendered server-side and revealed via the same state-driven redraw that powers the rest of the framework.",
			path:    "/showcase/tab-switcher",
		},
		{
			heading: "Navigation widgets",
			icon:    "explore",
			desc:    "The three flavors of side navigation — a compact strip, a labeled rail, and a full drawer — all driven by the same widget tree and adapting to the viewport.",
			path:    "/showcase/navigation",
		},
		{
			heading: "Form input widgets",
			icon:    "edit",
			desc:    "Every form input widget on one page — text, email, password, OTP, dates, ranges, dropdowns, file uploads, color pickers, checkboxes, radios, star/sentiment ratings, filter and input chips, toggle switches, and the rich-text editor.",
			path:    "/showcase/form-input",
		},
		{
			heading: "Form validation",
			icon:    "checklist",
			desc:    "Client-side and server-side form validation working together: required fields, pattern matching, length limits, and custom predicate functions, with inline error messages that survive a server round-trip.",
			path:    "/showcase/form-validation",
		},
		{
			heading: "Data table",
			icon:    "table view",
			desc:    "A sortable, filterable, paginated table over a small US-states dataset. Quick-search, column sort, and page-size controls all use the same incremental-redraw plumbing the rest of the framework uses.",
			path:    "/showcase/states",
		},
		{
			heading: "CRUD",
			icon:    "edit note",
			desc:    "Standard create / read / update / delete flow over a per-session in-memory directory of contacts. Demonstrates form-driven persistence with validation, duplicate detection, and a record cap.",
			path:    "/showcase/dir-list",
		},
		{
			heading: "Charts",
			icon:    "bar chart",
			desc:    "Apache ECharts wrapped as bespa widgets. Bar, pie, column, and a US-states map all themed against the Material design tokens, with live theme switching.",
			path:    "/showcase/charts",
		},
		{
			heading: "Mermaid",
			icon:    "account tree",
			desc:    "Mermaid diagrams (flowcharts, sequence, class, state, gantt) rendered client-side, themed against the Material design tokens, with optional zoom and pan interaction.",
			path:    "/showcase/mermaid",
		},
	}

	deck := wf.Deck(1, 2, 4)
	for _, scrn := range screens {
		deck.Add(
			wf.CardOutlined().Add(
				wf.TitleLarge(
					wf.Icon(scrn.icon),
					" ",
					scrn.heading,
				),
				wf.Spacer(0.25),
				scrn.desc,
				wf.SpacerBreak(),
				wf.Link(scrn.path).
					Add(wf.Icon("open_in_new").WithAltText("Open in window")),
				wf.PipeSeparator(),
				wf.Link("?modal="+scrn.path+"&panel=").
					Add(wf.Icon("web_asset").WithAltText("Open in modal")),
				wf.PipeSeparator(),
				wf.Link("?panel="+scrn.path+"&modal=").
					Add(wf.Icon("dock_to_right").WithAltText("Open in side panel")),
			),
		)
	}
	page.Add(deck)

	shared.Render(w, r, page)
}
