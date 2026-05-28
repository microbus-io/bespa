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

package nav

import (
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&MainMenuWidget{})               // Ensure interface
var _ = widget.NavAreaMarker(&MainMenuWidget{}) // Ensure interface

// MainMenuWidget renders a navigation menu.
type MainMenuWidget struct {
	*widget.WidgetBase[*MainMenuWidget]
	rail       *NavRailWidget
	vertical   *NavDrawerWidget
	horizontal *NavStripWidget
}

// MainMenu creates a new widget that renders the page's primary navigation,
// switching between three viewports automatically: a NavRail on wide
// screens, a NavStrip across the top on narrow screens, and a NavDrawer
// that slides in when the user opens the menu. Any of the three sections
// you don't set defaults to a single Home target. Implements
// NavAreaMarker so the parent Page renders it in the navigation slot —
// MainMenu is not drawn during partial redraws.
func (f NavFactory) MainMenu() *MainMenuWidget {
	x := &MainMenuWidget{}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// NavAreaMarker is a framework hook telling the enclosing page to place
// this widget in its navigation area. App code should not call this.
func (wgt *MainMenuWidget) NavAreaMarker() {
}

// WithRail sets the rail shown on wide (desktop) screens. Defaults to
// a single Home target when not set.
func (wgt *MainMenuWidget) WithRail(rail *NavRailWidget) *MainMenuWidget {
	wgt.rail = rail
	return wgt
}

// WithVertical sets the sliding drawer shown when the user opens the
// menu on narrow (mobile) screens. Defaults to a single Home target.
func (wgt *MainMenuWidget) WithVertical(vertical *NavDrawerWidget) *MainMenuWidget {
	wgt.vertical = vertical
	return wgt
}

// WithHorizontal sets the persistent top strip shown alongside the
// hamburger toggle on narrow (mobile) screens. Defaults to a single
// Home target.
func (wgt *MainMenuWidget) WithHorizontal(horizontal *NavStripWidget) *MainMenuWidget {
	wgt.horizontal = horizontal
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *MainMenuWidget) Children() []Widget {
	return Many(wgt.rail, wgt.vertical, wgt.horizontal)
}

// Draw renders the widget's HTML.
func (wgt *MainMenuWidget) Draw(w io.Writer, r *http.Request) (err error) {
	// Optimize: do not render for nested pages
	if r.Header.Get("Bespa-Fetch") == "1" {
		return nil
	}

	vertical := wgt.vertical
	if vertical == nil {
		vertical = factory.NavDrawer()
		vertical.AddTop(factory.NavTarget("Home", "/").WithIcon("home"))
	}
	horizontal := wgt.horizontal
	if horizontal == nil {
		horizontal = factory.NavStrip()
		horizontal.AddTop(factory.NavTarget("Home", "/").WithIcon("home"))
	}
	rail := wgt.rail
	if rail == nil {
		rail = factory.NavRail()
		rail.AddTop(factory.NavTarget("Home", "/").WithIcon("home"))
	}

	return Tag("div").
		Class("MainMenu").
		Attr("onmouseleave", "mainmenu_mouseleave(event)").
		Attr("data-id", wgt.ID()).
		Add(
			Tag("div").
				Class("HorizontalSection").
				Add(
					Tag("a").
						Class("MenuToggle").
						Attr("tabindex", "0").
						Attr("href", "").
						Attr("onclick", "mainmenu_reveal(event)").
						Add(factory.Icon("menu")),
					horizontal,
				),
			Tag("div").
				Class("Backdrop"),
			Tag("div").
				Class("VerticalSection").
				Add(
					Tag("a").
						Class("MenuToggle").
						Attr("tabindex", "0").
						Attr("href", "").
						Attr("onclick", "mainmenu_conceal(event)").
						Add(factory.Icon("menu open")),
					vertical,
				),
			Tag("div").
				Class("RailSection").
				Attr("onmouseenter", "mainmenu_railMousenter(event)").
				Add(rail),
		).
		When(wgt.Shown(r)).
		Draw(w, r)
}
