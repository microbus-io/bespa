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

var _ = Widget(&InfoLinkWidget{}) // Ensure interface

// InfoLinkWidget renders an info or help link.
type InfoLinkWidget struct {
	*widget.WidgetBase[*InfoLinkWidget]
	icon   Widget
	href   string
	target string
}

// InfoLink creates a new widget that renders an info icon as a clickable
// link. The href accepts the full action-URL grammar (`?key=`, `^?…`,
// `/path`, etc.). For an in-page popover instead of a navigation, use
// InfoBubble.
func (f BasicFactory) InfoLink(href string) *InfoLinkWidget {
	x := &InfoLinkWidget{
		icon: factory.Icon("info"),
		href: href,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// HelpLink is identical to InfoLink but uses a help (?) icon — use it for
// "what is this?" links rather than informational ones.
func (f BasicFactory) HelpLink(href string) *InfoLinkWidget {
	x := &InfoLinkWidget{
		icon: factory.Icon("help_outline"),
		href: href,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithTarget sets the target of the link.
func (wgt *InfoLinkWidget) WithTarget(target string) *InfoLinkWidget {
	wgt.target = target
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *InfoLinkWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("span").
		Class("InfoLink").
		Attr("data-id", wgt.ID()).
		Add(factory.Link(wgt.href).Add(wgt.icon).WithTarget(wgt.target)).
		When(wgt.Shown(r)).
		Draw(w, r)
}
