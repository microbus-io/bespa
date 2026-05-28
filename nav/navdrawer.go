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
	"fmt"
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&NavDrawerWidget{}) // Ensure interface

// NavDrawerWidget renders a navigation drawer.
type NavDrawerWidget struct {
	*widget.WidgetBase[*NavDrawerWidget]
	topChildren    []Widget
	bottomChildren []Widget
	style          string
	height         string
}

// NavDrawer creates a new widget that renders a vertical sliding panel
// of nav targets — the side drawer slot of a MainMenu on narrow screens.
// Populate via AddTop and AddBottom; bottom items stick to the bottom
// only when WithHeight is set (otherwise the drawer fits its content
// and there's no spare space to push them down).
func (f NavFactory) NavDrawer() *NavDrawerWidget {
	x := &NavDrawerWidget{
		style: "Drawer",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// AddTop appends entries to the top section of the drawer. NavTargets,
// section titles, and Rule separators are the typical children.
func (wgt *NavDrawerWidget) AddTop(topChildren ...any) *NavDrawerWidget {
	wgt.topChildren = Many(wgt.topChildren, topChildren)
	return wgt
}

// AddBottom appends entries pinned to the bottom of the drawer.
// Pinning requires WithHeight to be set; without it the items just
// follow the top section.
func (wgt *NavDrawerWidget) AddBottom(bottomChildren ...any) *NavDrawerWidget {
	wgt.bottomChildren = Many(wgt.bottomChildren, bottomChildren)
	return wgt
}

// WithHeight pins the drawer to a fixed height. Default is auto (fits
// content). Set this to 100% (or similar) so AddBottom items can stick
// to the bottom; without it, the drawer has no spare space and bottom
// items just trail the top section.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
func (wgt *NavDrawerWidget) WithHeight(height float32, unit string) *NavDrawerWidget {
	if height > 0 {
		wgt.height = fmt.Sprintf("height:%f%s", height, unit)
	} else {
		wgt.height = ""
	}
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *NavDrawerWidget) Children() []Widget {
	return Many(wgt.topChildren, wgt.bottomChildren)
}

// Draw renders the widget's HTML.
func (wgt *NavDrawerWidget) Draw(w io.Writer, r *http.Request) (err error) {
	sliderTag := Tag("div").Class("Slider")
	mainPanel := Tag("div").
		Class("Panel").
		Add(
			wgt.topChildren,
			Tag("div").Class("AutoMargin").Hide(len(wgt.bottomChildren) == 0),
			wgt.bottomChildren,
		)
	sliderTag.Add(mainPanel)
	for _, c := range wgt.Children() {
		if tgt, ok := c.(*NavTargetWidget); ok && tgt.submenu != nil {
			childPanel := Tag("div").
				Class("Panel").
				Attr("id", tgt.submenuID).
				Add(tgt.submenu)
			sliderTag.Add(childPanel)
		}
	}
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Class("Nav"+wgt.style).
		Attr("onclick", "navdrawer_click(event)").
		Attr("onkeydown", "navdrawer_keydown(event)").
		Style(wgt.height).
		Add(sliderTag).
		When(wgt.Shown(r) && len(wgt.topChildren)+len(wgt.bottomChildren) > 0).
		Draw(w, r)
}
