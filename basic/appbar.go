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
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&AppBarWidget{})                 // Ensure interface
var _ = widget.PageTitleMarker(&AppBarWidget{}) // Ensure interface

// AppBarWidget renders the top application bar.
type AppBarWidget struct {
	*widget.WidgetBase[*AppBarWidget]
	leftChildren   []Widget
	rightChildren  []Widget
	bottomChildren []Widget
	help           Widget
	backArrow      *LinkWidget
	h1             *HeadingWidget
}

// AppBar creates a new widget that renders the top application bar.
// The title children are rendered as a HeadlineLarge; a back-arrow link is
// included by default and follows the `_back` state variable.
func (f BasicFactory) AppBar(titleChildren ...any) *AppBarWidget {
	x := &AppBarWidget{
		h1:        factory.HeadlineLarge(titleChildren),
		backArrow: factory.Link("").WithHrefBack().Add(factory.Icon("arrow_back")),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithHelpBubble sets help for the app bar.
func (wgt *AppBarWidget) WithHelpBubble(helpChildren ...any) *AppBarWidget {
	wgt.help = factory.HelpBubble(helpChildren...)
	return wgt
}

// WithHelpLink sets help for the app bar.
func (wgt *AppBarWidget) WithHelpLink(href string) *AppBarWidget {
	wgt.help = factory.HelpLink(href)
	return wgt
}

// WithBackLink sets an explicit href for the back arrow.
// Without it, the back arrow falls back to the `_back` state variable.
// Accepts the full action-URL grammar (`?key=`, `^?…`, `/path`, etc.).
func (wgt *AppBarWidget) WithBackLink(href string) *AppBarWidget {
	wgt.backArrow.WithHref(href)
	return wgt
}

// AddLeft adds nested widgets aligned to the left.
func (wgt *AppBarWidget) AddLeft(children ...any) *AppBarWidget {
	wgt.leftChildren = Many(wgt.leftChildren, children)
	return wgt
}

// AddRight adds nested widgets aligned to the right.
func (wgt *AppBarWidget) AddRight(children ...any) *AppBarWidget {
	wgt.rightChildren = Many(wgt.rightChildren, children)
	return wgt
}

// AddBottom adds nested widgets in the bottom section of the app bar.
func (wgt *AppBarWidget) AddBottom(children ...any) *AppBarWidget {
	wgt.bottomChildren = Many(wgt.bottomChildren, children)
	return wgt
}

// PageTitleMarker is a framework hook that lets the enclosing page use the
// app bar's title as its <title>. App code should not need to call this.
func (wgt *AppBarWidget) PageTitleMarker() {
}

// PageTitle returns the rendered text of the app bar's title heading, used
// by the enclosing page as its <title>.
func (wgt *AppBarWidget) PageTitle() string {
	type Stringer interface {
		String() string
	}
	var sb strings.Builder
	for _, c := range wgt.h1.children {
		if str, ok := c.(Stringer); ok {
			sb.WriteString(str.String())
		}
		if icon, ok := c.(*IconWidget); ok {
			if icon.symbol == "chevron_right" {
				sb.WriteString("\u203a")
			}
		}
		sb.WriteString(" ")
	}
	return sb.String()
}

// Children are the widgets nested under this widget.
func (wgt *AppBarWidget) Children() []Widget {
	return Many(wgt.backArrow, wgt.h1, wgt.help, wgt.leftChildren, wgt.rightChildren, wgt.bottomChildren)
}

// Draw renders the widget's HTML.
func (wgt *AppBarWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Class("AppBar").
		Attr("data-id", wgt.ID()).
		Add(
			factory.Toolbar().
				AddLeft(
					HTMLUnsafe(Tag("div").Class("AppBarTitle").Add(wgt.backArrow, wgt.h1, wgt.help).String(r)),
					wgt.leftChildren,
				).
				AddRight(wgt.rightChildren),
			wgt.bottomChildren,
		).
		When(wgt.Shown(r)).
		Draw(w, r)
}
