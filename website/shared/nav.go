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

package shared

import (
	"github.com/microbus-io/bespa/nav"
	"github.com/microbus-io/bespa/widget"
)

// navMenu creates the top navigation menu of the website.
func navMenu() *nav.MainMenuWidget {
	rail := wf.NavRail().AddTop(
		wf.NavTarget("Home", "/").WithIcon("home"),
		wf.NavTarget("Get started", "/start").WithIcon("rocket launch"),
		showcase(false),
		basics(false),
		build(false),
		extend(false),
		wf.NavTarget("Git it", "https://github.com/microbus-io/bespa").WithIcon("download"),
		wf.NavTarget("Find us", "/contact").WithIcon("phone"),
	).AddBottom(
		wf.NavTarget("Profile", "/profile").WithIcon("person").WithBadge("!"),
	)

	vertical := wf.NavDrawer().AddTop(
		wf.NavTarget("Home", "/").WithIcon("home"),
		wf.NavTarget("Get started", "/start").WithIcon("rocket launch"),
		showcase(true),
		basics(true),
		build(true),
		extend(true),
		wf.NavTarget("Git it", "https://github.com/microbus-io/bespa").WithIcon("download"),
		wf.NavTarget("Find us", "/contact").WithIcon("phone"),
	).AddBottom(
		wf.NavTarget("Profile", "/profile").WithIcon("person").WithBadge("!"),
	)

	horizontal := wf.NavStrip().AddTop(
		wf.NavTarget("BESPA", "/").WithIcon("home"),
	)

	return wf.MainMenu().
		WithRail(rail).
		WithVertical(vertical).
		WithHorizontal(horizontal)
}

// showcase returns the nav targets of the showcase section of the menu.
func showcase(backToMainMenu bool) widget.Widget {
	targets := wf.NavDrawer()
	if backToMainMenu {
		targets.AddTop(wf.NavTargetBack("Main menu"), wf.Rule())
	}
	targets.AddTop(
		wf.NavTarget("Overview", "/showcase/overview").WithIcon("play arrow"),
		wf.NavTarget("Text formatting", "/showcase/text-formatting").WithIcon("text format"),
		wf.NavTarget("Toolbar", "/showcase/toolbar").WithIcon("construction"),
		wf.NavTarget("Gallery", "/showcase/gallery").WithIcon("collections"),
		wf.NavTarget("Deck of cards", "/showcase/deck").WithIcon("dashboard"),
		wf.NavTarget("Tab switcher", "/showcase/tab-switcher").WithIcon("tab"),
		wf.NavTarget("Navigation", "/showcase/navigation").WithIcon("explore"),
		wf.NavTarget("Progress", "/showcase/progress").WithIcon("pending"),
		wf.NavTarget("Form input", "/showcase/form-input").WithIcon("edit"),
		wf.NavTarget("Form validation", "/showcase/form-validation").WithIcon("checklist"),
		wf.NavTarget("Data table", "/showcase/states").WithIcon("table view"),
		wf.NavTarget("CRUD", "/showcase/dir-list").WithIcon("edit note"),
		wf.NavTarget("Code blocks", "/showcase/code").WithIcon("code blocks"),
		wf.NavTarget("Charts", "/showcase/charts").WithIcon("bar chart"),
		wf.NavTarget("Mermaid", "/showcase/mermaid").WithIcon("account tree"),
	)
	return wf.NavTarget("Showcase", "/showcase").WithIcon("play arrow").WithSubMenu(targets)
}

// basics returns the nav targets of the Basics section — how BESPA's
// rendering / state / partial-redraw protocol works internally.
func basics(backToMainMenu bool) widget.Widget {
	targets := wf.NavDrawer()
	if backToMainMenu {
		targets.AddTop(wf.NavTargetBack("Main menu"), wf.Rule())
	}
	targets.AddTop(
		wf.NavTarget("Overview", "/basics/overview").WithIcon("school"),
		wf.NavTarget("Declarative views", "/basics/declarative-views").WithIcon("dataset"),
		wf.NavTarget("Action-URL pattern", "/basics/action-url-pattern").WithIcon("link"),
		wf.NavTarget("Incremental updates", "/basics/incremental").WithIcon("autorenew"),
		wf.NavTarget("Embedded pages", "/basics/embedded-pages").WithIcon("dashboard customize"),
		wf.NavTarget("Nesting pages", "/basics/nesting").WithIcon("view quilt"),
		wf.NavTarget("Targeting frames", "/basics/frames").WithIcon("space dashboard"),
		wf.NavTarget("Cheat sheet", "/basics/cheatsheet").WithIcon("summarize"),
	)
	return wf.NavTarget("Basics", "/basics/overview").WithIcon("school").WithSubMenu(targets)
}

// build returns the nav targets of the Build section — practical techniques
// for using BESPA as a library: when to redraw, state, forms, tables, modals.
func build(backToMainMenu bool) widget.Widget {
	targets := wf.NavDrawer()
	if backToMainMenu {
		targets.AddTop(wf.NavTargetBack("Main menu"), wf.Rule())
	}
	targets.AddTop(
		wf.NavTarget("Overview", "/build/overview").WithIcon("build"),
		wf.NavTarget("Handlers & routing", "/build/handlers-and-routing").WithIcon("schema"),
		wf.NavTarget("State patterns", "/build/state").WithIcon("hub"),
		wf.NavTarget("When to redraw", "/build/redraw").WithIcon("autorenew"),
		wf.NavTarget("Forms & validation", "/build/forms").WithIcon("edit"),
		wf.NavTarget("Data tables", "/build/tables").WithIcon("table view"),
		wf.NavTarget("Modals & side panels", "/build/modals").WithIcon("web asset"),
		wf.NavTarget("Navigation patterns", "/build/navigation").WithIcon("alt route"),
		wf.NavTarget("Errors", "/build/errors").WithIcon("error"),
		wf.NavTarget("Sessions & auth", "/build/sessions-and-auth").WithIcon("lock"),
		wf.NavTarget("Live data", "/build/live-data").WithIcon("bolt"),
		wf.NavTarget("Theming", "/build/theming").WithIcon("palette"),
		wf.NavTarget("Deployment", "/build/deployment").WithIcon("rocket launch"),
	)
	return wf.NavTarget("Build", "/build/overview").WithIcon("build").WithSubMenu(targets)
}

// extend returns the nav targets of the Extend section — building custom
// widgets and packaging reusable widget libraries.
func extend(backToMainMenu bool) widget.Widget {
	targets := wf.NavDrawer()
	if backToMainMenu {
		targets.AddTop(wf.NavTargetBack("Main menu"), wf.Rule())
	}
	targets.AddTop(
		wf.NavTarget("Overview", "/extend/overview").WithIcon("extension"),
		wf.NavTarget("Widget anatomy", "/extend/anatomy").WithIcon("widgets"),
		wf.NavTarget("State-aware widgets", "/extend/state-aware").WithIcon("autorenew"),
		wf.NavTarget("Composing", "/extend/composing").WithIcon("dashboard"),
		wf.NavTarget("Assets & CSS", "/extend/assets").WithIcon("style"),
		wf.NavTarget("Material theming", "/extend/theming").WithIcon("palette"),
		wf.NavTarget("Custom form inputs", "/extend/form-input-widgets").WithIcon("edit"),
		wf.NavTarget("Packaging", "/extend/packaging").WithIcon("inventory 2"),
	)
	return wf.NavTarget("Extend", "/extend/overview").WithIcon("extension").WithSubMenu(targets)
}
