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
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&NavTargetWidget{}) // Ensure interface

// NavTargetWidget renders a navigation target.
type NavTargetWidget struct {
	*widget.WidgetBase[*NavTargetWidget]
	href      string
	label     string
	icon      string
	badge     string
	back      string
	submenu   Widget
	submenuID string
	tabulated string
}

// NavTarget creates a new widget for one entry in a nav rail, drawer,
// or strip. label is the visible text; href is the destination (accepts
// the full action-URL grammar). The target auto-highlights as "selected"
// when the current request path equals or is nested under its href path.
// Add an icon with WithIcon, a notification badge with WithBadge, or a
// nested sub-menu with WithSubMenu.
func (f NavFactory) NavTarget(label string, href string) *NavTargetWidget {
	x := &NavTargetWidget{
		label: label,
		href:  href,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// NavTargetBack creates a new "back" entry for use inside a sub-menu —
// clicking it slides the user back up to the parent panel rather than
// navigating away. Pass an empty label to use the default "Back".
func (f NavFactory) NavTargetBack(label string) *NavTargetWidget {
	if label == "" {
		label = "Back"
	}
	x := &NavTargetWidget{
		label: label,
		icon:  "arrow back",
		back:  "1",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithIcon sets the leading icon. Accepts any Material Symbol name
// (e.g. "home", "settings"). When omitted, an invisible placeholder is
// rendered so labels still align with their siblings.
func (wgt *NavTargetWidget) WithIcon(icon string) *NavTargetWidget {
	wgt.icon = icon
	return wgt
}

// WithTabulation indents the target — used to render a flat list as a
// visual sub-hierarchy without nesting it under a sub-menu.
func (wgt *NavTargetWidget) WithTabulation(tab bool) *NavTargetWidget {
	if tab {
		wgt.tabulated = "Tabulated"
	} else {
		wgt.tabulated = ""
	}
	return wgt
}

// WithBadge attaches a notification badge to the target, typically a
// small count like "3" or "99+". Pass "." for a dotless indicator (no
// text — just a colored dot). Empty hides the badge.
func (wgt *NavTargetWidget) WithBadge(badge string) *NavTargetWidget {
	wgt.badge = badge
	return wgt
}

// WithSubMenu attaches a nested panel that slides in when the target is
// activated. Add a NavTargetBack at the top of the sub-menu so users can
// return to the parent.
func (wgt *NavTargetWidget) WithSubMenu(submenu Widget) *NavTargetWidget {
	wgt.submenu = submenu
	wgt.submenuID = widget.RandomAlphaNumID(8)
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *NavTargetWidget) Draw(w io.Writer, r *http.Request) (err error) {
	linkTag := Tag("a").
		Attr("data-id", wgt.ID()).
		Class("NavTarget").
		Class(wgt.tabulated).
		Attr("tabindex", "0").
		Attr("href", wgt.href).
		Attr("data-back", wgt.back)
	prefix := r.Header.Get("X-Forwarded-Prefix")
	path := r.Header.Get("X-Forwarded-Path")
	if path == "" {
		path = r.URL.Path
	}
	hrefPath := wgt.href
	if p := strings.Index(hrefPath, "?"); p >= 0 {
		hrefPath = hrefPath[:p]
	}
	if hrefPath != "" && (prefix+path == hrefPath || strings.HasPrefix(prefix+path, hrefPath+"/")) {
		linkTag.Class("Selected")
	}
	emptyIcon := ""
	icon := wgt.icon
	if icon == "" {
		icon = "crop square"
		emptyIcon = "TgtNoIcon"
	}
	microBadge := ""
	badge := wgt.badge
	if badge == "." {
		microBadge = "MicroBadge"
		badge = ""
	}
	linkTag.Add(
		Tag("div").
			Class("TgtIcon", emptyIcon).
			Add(
				factory.Icon(icon),
				Tag("span").Class("MiniBadge", microBadge).Add(badge).Hide(wgt.badge == ""),
			),
	)
	linkTag.Add(
		Tag("div").Class("TgtLabel").Add(wgt.label).Hide(wgt.label == ""),
		Tag("div").Class("Badge").Add(badge),
		Tag("div").Class("TgtNext").Add(factory.Icon("arrow forward")),
	)
	if wgt.submenu != nil {
		linkTag.
			Attr("data-next", wgt.submenuID).
			Attr("aria-haspopup", "menu").
			Attr("aria-expanded", "false").
			Attr("aria-controls", wgt.submenuID)
	}
	return linkTag.
		When(wgt.Shown(r)).
		Draw(w, r)
}
