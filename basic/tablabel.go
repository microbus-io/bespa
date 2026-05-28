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
	"net/url"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&TabLabelWidget{}) // Ensure interface

// TabLabelWidget renders a tab.
type TabLabelWidget struct {
	*widget.WidgetBase[*TabLabelWidget]
	key      string
	children []Widget
	href     string
	parent   *TabSwitcherWidget
}

// TabLabel creates a new widget that renders one clickable tab. key
// identifies the tab and must match a sibling TabBody's key within the
// same TabSwitcher. TabLabel must be added directly to a TabSwitcher.
// By default clicking the tab sets the switcher's state variable to key;
// override the target with WithHref.
func (f BasicFactory) TabLabel(key string) *TabLabelWidget {
	x := &TabLabelWidget{
		key: key,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithHref overrides the tab's default action with an explicit href.
// Use this to point a tab at another page instead of just selecting itself
// in the current switcher. Accepts the full action-URL grammar.
func (wgt *TabLabelWidget) WithHref(href string) *TabLabelWidget {
	wgt.href = href
	return wgt
}

// Add adds nested widgets to the label.
func (wgt *TabLabelWidget) Add(children ...any) *TabLabelWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *TabLabelWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *TabLabelWidget) Draw(w io.Writer, r *http.Request) (err error) {
	href := wgt.href
	if href == "" {
		href = "?" + wgt.parent.name + "=" + url.QueryEscape(wgt.key)
	}
	singleIcon := false
	if len(wgt.children) == 1 {
		_, singleIcon = wgt.children[0].(*IconWidget)
	}
	selected := wgt.parent.currentTabKey(r) == wgt.key
	return Tag("span").
		Class("TabLabel").
		Attr("data-id", wgt.ID()).
		ClassIf(selected, "selected").
		ClassIf(singleIcon, "SingleIcon").
		Attr("tabindex", "0").
		Attr("role", "tab").
		Attr("aria-selected", boolAttr(selected)).
		Attr("aria-controls", "tabbody_"+wgt.parent.name+"_"+wgt.key).
		Attr("id", "tablabel_"+wgt.parent.name+"_"+wgt.key).
		Attr("data-name", wgt.parent.name).
		Attr("data-tab", wgt.key).
		AttrIf(wgt.parent.dynamic, "data-dynamic", "1").
		Attr("onclick", "tabswitcher_click(event)").
		Attr("onkeydown", "tabswitcher_keydown(event)").
		Add(
			Tag("div").Add(wgt.children),
			Tag("u"),
			Tag("a").Attr("href", href).Hide(true)).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}

// boolAttr formats a boolean as the "true"/"false" string expected by ARIA attributes.
func boolAttr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *TabLabelWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) ||
		(wgt.parent.dynamic && state.Changed(wgt.parent.name))
}
