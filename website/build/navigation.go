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

const navMainMenu = `func navMenu() *nav.MainMenuWidget {
    rail := wf.NavRail().AddTop(
        wf.NavTarget("Home", "/").WithIcon("home"),
        wf.NavTarget("Orders", "/orders").WithIcon("receipt").WithBadge("3"),
        wf.NavTarget("Reports", "/reports").WithIcon("bar chart"),
    ).AddBottom(
        wf.NavTarget("Profile", "/profile").WithIcon("person"),
    )

    drawer := wf.NavDrawer().AddTop(
        // same targets, full-width labels
    )

    return wf.MainMenu().WithRail(rail).WithVertical(drawer)
}
`

const navSubmenu = `wf.NavTarget("Reports", "/reports/overview").
    WithIcon("bar chart").
    WithSubMenu(wf.NavDrawer().AddTop(
        wf.NavTargetBack("Main menu"),
        wf.Rule(),
        wf.NavTarget("Daily", "/reports/daily").WithIcon("today"),
        wf.NavTarget("Weekly", "/reports/weekly").WithIcon("date range"),
        wf.NavTarget("Custom range", "/reports/custom").WithIcon("date range"),
    ))
`

const navBack = `// Default back button — uses state._back set by the caller.
wf.ButtonText("").Add("Back").WithHrefBack(),

// Explicit destination:
wf.ButtonText("").Add("Back").WithHref("/orders"),

// Add _back to a link so the destination has a return path:
wf.Link("/orders/123?_back=" + url.QueryEscape("/orders")).Add("Open"),
`

// HandleNavigation covers navigation patterns.
func HandleNavigation(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Navigation patterns"),

		wf.Markdown(
			"Navigation in BESPA is three layered widgets: a strip across the ",
			"top, a rail or drawer down the side, and the in-page back-link ",
			"conventions that let any page return to where the user came from. ",
			"The framework adapts between rail and drawer based on viewport width.",
		),
		wf.HeadlineMedium("The main menu"),
		wf.Markdown(
			"Most apps define their menu once in a shared helper and reuse it on ",
			"every page. The pattern from this website's `shared/nav.go`:",
		),
		wf.CodeBlock(navMainMenu).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"- `NavRail` is the compact icon-and-label vertical bar shown on ",
			"wider viewports.\n",
			"- `NavDrawer` is the same content as a full-width labelled drawer, ",
			"used on narrower screens.\n",
			"- `NavStrip` is a horizontal bar — typically only used at the very ",
			"top for branding.",
		),
		wf.HeadlineMedium("Targets and decorations"),
		wf.Markdown(
			"`NavTarget(label, path)` is the link. Common chained options:",
			"\n\n",
			"- `WithIcon(name)` — Material Symbols name (space-separated words ",
			"allowed).\n",
			"- `WithBadge(text)` — small annotation chip; useful for unread ",
			"counts. Use `\".\"` for a dot-only indicator.\n",
			"- `WithSubMenu(drawer)` — turns the entry into an expandable ",
			"section with its own targets.",
		),
		wf.HeadlineMedium("Sub-menus"),
		wf.Markdown(
			"For sections with their own internal navigation, attach a ",
			"sub-drawer. Open the section to replace the main menu with the ",
			"section's targets, plus a `NavTargetBack` breadcrumb back to the ",
			"parent:",
		),
		wf.CodeBlock(navSubmenu).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This site's Learn section is exactly this pattern — try clicking ",
			"\"Learn\" in the side menu and you'll see the sub-drawer slide in.",
		),
		wf.HeadlineMedium("Back-link plumbing"),
		wf.Markdown(
			"BESPA uses the `_back` state variable to remember where a user ",
			"came from. Pages that open from a list view (\"Edit\" links, detail ",
			"views) should set it; the destination uses `WithHrefBack` or ",
			"`RedirectBack` to return cleanly:",
		),
		wf.CodeBlock(navBack).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"When you call `WithHrefBack` without an explicit back URL, the ",
			"button falls back to `history.go(-1)` if it can detect that there's ",
			"history to walk; otherwise it's hidden automatically.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Showcase → Navigation widgets](/showcase/navigation) ",
			"— all three flavors side by side.",
			"\n\n",
			"[Basics → Targeting frames](/basics/frames) ",
			"— when a click should target a nested page rather than the top.",
		),
	)
	shared.Render(w, r, page)
}
