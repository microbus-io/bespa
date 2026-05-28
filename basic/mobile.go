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

var _ = Widget(&MobileWidget{}) // Ensure interface

// MobileWidget renders its content to be visible or hidden depending on the screen size.
type MobileWidget struct {
	*widget.WidgetBase[*MobileWidget]
	mobile   string
	children []Widget
}

// MobileOnly creates a new widget whose children are shown only on
// narrow (mobile-sized) viewports — useful for hamburger triggers or
// compact summaries that desktop layouts shouldn't display. Visibility is
// CSS-driven; the markup still ships to every client.
func (f BasicFactory) MobileOnly(children ...any) *MobileWidget {
	x := &MobileWidget{
		mobile:   "Only",
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// MobileHide creates a new widget whose children are hidden on narrow
// (mobile-sized) viewports — the mirror of MobileOnly. Visibility is
// CSS-driven; the markup still ships to every client.
func (f BasicFactory) MobileHide(children ...any) *MobileWidget {
	x := &MobileWidget{
		mobile:   "Hide",
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *MobileWidget) Add(children ...any) *MobileWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *MobileWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *MobileWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("span").
		Class("Mobile"+wgt.mobile).
		Attr("data-id", wgt.ID()).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
