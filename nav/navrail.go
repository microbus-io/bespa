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

var _ = Widget(&NavRailWidget{}) // Ensure interface

// NavRailWidget renders a navigation rail.
type NavRailWidget struct {
	*widget.WidgetBase[*NavRailWidget]
	topChildren    []Widget
	bottomChildren []Widget
	height         string
}

// NavRail creates a new widget that renders a narrow vertical strip of
// icon-only nav targets — the desktop slot of a MainMenu. Populate via
// AddTop and AddBottom; bottom items pin to the bottom only when
// WithHeight is set.
func (f NavFactory) NavRail() *NavRailWidget {
	x := &NavRailWidget{}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithHeight pins the rail to a fixed height. Default is auto (fits
// content). Required if you want AddBottom items to stick to the bottom.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
func (wgt *NavRailWidget) WithHeight(height float32, units string) *NavRailWidget {
	if height > 0 {
		wgt.height = fmt.Sprintf("height:%.2f%s", height, units)
	} else {
		wgt.height = ""
	}
	return wgt
}

// AddTop appends entries to the top of the rail. NavTargets, section
// titles, and Rule separators are the typical children.
func (wgt *NavRailWidget) AddTop(topChildren ...any) *NavRailWidget {
	wgt.topChildren = Many(wgt.topChildren, topChildren)
	return wgt
}

// AddBottom appends entries pinned to the bottom of the rail. Pinning
// requires WithHeight to be set; without it the items follow the top
// section instead.
func (wgt *NavRailWidget) AddBottom(bottomChildren ...any) *NavRailWidget {
	wgt.bottomChildren = Many(wgt.bottomChildren, bottomChildren)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *NavRailWidget) Children() []Widget {
	return Many(wgt.topChildren, wgt.bottomChildren)
}

// Draw renders the widget's HTML.
func (wgt *NavRailWidget) Draw(w io.Writer, r *http.Request) (err error) {
	sliderTag := Tag("div").
		Class("Slider").
		Style(fmt.Sprintf("z-index:%d;", -1))
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
		Class("NavRail").
		Attr("onclick", "navrail_click(event)").
		Attr("onkeydown", "navrail_keydown(event)").
		Attr("onmouseover", "navrail_mouseover(event)").
		Attr("onmouseleave", "navrail_mouseleave(event)").
		Add(
			sliderTag,
			Tag("div").Class("Panel").Add(
				wgt.topChildren,
				Tag("div").Class("AutoMargin").Hide(len(wgt.topChildren) == 0 || len(wgt.bottomChildren) == 0),
				wgt.bottomChildren,
			),
		).
		When(wgt.Shown(r) && len(wgt.topChildren)+len(wgt.bottomChildren) > 0).
		Draw(w, r)
}
