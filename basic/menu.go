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

package basic

import (
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&MenuWidget{}) // Ensure interface

// MenuWidget renders a popup menu.
type MenuWidget struct {
	*widget.WidgetBase[*MenuWidget]
	title   Widget
	actions []Widget
}

// Menu creates a new widget that renders a popup menu opened by hovering or
// clicking the title. Items are added via Add and must be Link widgets so
// each action carries its own href.
func (f BasicFactory) Menu(title any) *MenuWidget {
	x := &MenuWidget{
		title: Any(title),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// MenuEllipsisH creates a new widget that renders a popup menu with a horizontal ellipsis menu for its title.
func (f BasicFactory) MenuEllipsisH() *MenuWidget {
	x := &MenuWidget{
		title: factory.Icon("more_horiz"),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// MenuEllipsisV creates a new widget that renders a popup menu with a vertical ellipsis menu for its title.
func (f BasicFactory) MenuEllipsisV() *MenuWidget {
	x := &MenuWidget{
		title: factory.Icon("more_vert"),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add appends menu items. Each item must be a Link — its href determines
// what the menu action does.
func (wgt *MenuWidget) Add(actions ...*LinkWidget) *MenuWidget {
	for _, a := range actions {
		wgt.actions = append(wgt.actions, a)
	}
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *MenuWidget) Children() []Widget {
	return Many(wgt.title, wgt.actions)
}

// Draw renders the widget's HTML.
func (wgt *MenuWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(
			Tag("div").
				Class("Menu").
				Attr("onmouseenter", "menu_mouseenter(event)").
				Attr("onclick", "menu_click(event)").
				Add(wgt.title),
			Tag("div").
				Add(wgt.actions)).
		When(wgt.Shown(r) && len(wgt.actions) > 0).
		Draw(w, r)
}
