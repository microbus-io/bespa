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

// HandleToolbar demonstrates the toolbar widget.
func HandleNavigation(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Navigation widgets"),
		`Navigation widgets enable the user to traverse a hierarchical menu structure.
		Navigation targets comprise of a label and an optional icon and may be aligned to
		the top or bottom, or left or right, of the navigation widget.`,

		// Navigation strip
		wf.HeadlineMedium("Navigation strip"),
		`The navigation strip is displayed as a horizontal menu strip with
		sideways animation enabling traversal of a hierarchical menu structure.`,
		wf.SpacerParagraph(),
		wf.NavStrip().AddTop(
			wf.NavTarget("Home", "?clicked=Home").WithIcon("home"),
			wf.NavTarget("Inventory", "?clicked=Inventory").WithIcon("inventory"),
			wf.NavTarget("Analytics", "?clicked=Analytics").WithIcon("analytics"),
			wf.NavTarget("Options", "").WithIcon("settings").WithSubMenu(
				wf.NavStrip().AddTop(
					wf.NavTargetBack("Back"),
					wf.NavTarget("Language", "?clicked=Language").WithIcon("language"),
					wf.NavTarget("Colors", "?clicked=Colors").WithIcon("palette"),
					wf.NavTarget("More options", "").WithIcon("folder open").WithSubMenu(wf.NavStrip().AddTop(
						wf.NavTargetBack("Back"),
						wf.NavTarget("One", "?clicked=One").WithIcon("looks one"),
						wf.NavTarget("Two", "?clicked=Two").WithIcon("looks two"),
					)),
				),
			),
		).AddBottom(
			wf.NavTarget("Profile", "?clicked=Profile").WithIcon("person"),
		),

		// Navigation drawer
		wf.HeadlineMedium("Navigation drawer"),
		`The navigation drawer is displayed as a vertical menu.
		Sideways animation enables navigation of a hierarchical menu structure.`,
		wf.SpacerParagraph(),
		wf.NavDrawer().AddTop(
			wf.NavTarget("Home", "?clicked=Home").WithIcon("home"),
			wf.NavTarget("Inventory", "?clicked=Inventory").WithIcon("inventory"),
			wf.NavTarget("Analytics", "?clicked=Analytics").WithIcon("analytics"),
			wf.NavTarget("Options", "").WithIcon("settings").WithSubMenu(
				wf.NavDrawer().AddTop(
					wf.NavTargetBack("Back"),
					wf.NavTarget("Language", "?clicked=Language").WithIcon("language"),
					wf.NavTarget("Colors", "?clicked=Colors").WithIcon("palette"),
					wf.NavTarget("More options", "").WithIcon("folder open").WithSubMenu(wf.NavDrawer().AddTop(
						wf.NavTargetBack("Back"),
						wf.NavTarget("One", "?clicked=One").WithIcon("looks one"),
						wf.NavTarget("Two", "?clicked=Two").WithIcon("looks two"),
					)),
				),
			),
		).AddBottom(
			wf.NavTarget("Profile", "?clicked=Profile").WithIcon("person"),
		),

		// Navigation rail
		wf.HeadlineMedium("Navigation rail"),
		`The navigation rail is displayed as a vertical rail of icons.
		A drawer opens up to enable navigation of a hierarchical menu structure.`,
		wf.SpacerParagraph(),
		wf.NavRail().AddTop(
			wf.NavTarget("Home", "?clicked=Home").WithIcon("home"),
			wf.NavTarget("Inventory", "?clicked=Inventory").WithIcon("inventory"),
			wf.NavTarget("Analytics", "?clicked=Analytics").WithIcon("analytics"),
			wf.NavTarget("Options", "").WithIcon("settings").WithSubMenu(
				wf.NavDrawer().AddTop(
					wf.NavTarget("Language", "?clicked=Language").WithIcon("language"),
					wf.NavTarget("Colors", "?clicked=Colors").WithIcon("palette"),
					wf.NavTarget("More options", "").WithIcon("folder open").WithSubMenu(wf.NavDrawer().AddTop(
						wf.NavTargetBack("Back"),
						wf.NavTarget("One", "?clicked=One").WithIcon("looks one"),
						wf.NavTarget("Two", "?clicked=Two").WithIcon("looks two"),
					)),
				),
			),
		).AddBottom(
			wf.NavTarget("Profile", "?clicked=Profile").WithIcon("person"),
		),

		wf.Snackbar().Add(
			wf.StateOf(r).Get("clicked"), " clicked",
		).HideIfEmpty(r, "clicked").RedrawIfChanged(r, "clicked"),

		wf.Debugger(),
	)
	shared.Render(w, r, page)
}
