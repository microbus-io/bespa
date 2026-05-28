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
	"errors"
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&TabBodyWidget{}) // Ensure interface

// TabBodyWidget renders a tab.
type TabBodyWidget struct {
	*widget.WidgetBase[*TabBodyWidget]
	key      string
	children []Widget
	parent   *TabSwitcherWidget
}

// TabBody creates a new widget for the content of one tab. key identifies
// the tab and must match a TabLabel's key within the same TabSwitcher.
// TabBody must be added directly to a TabSwitcher; rendering raises an
// error otherwise.
func (f BasicFactory) TabBody(key string) *TabBodyWidget {
	x := &TabBodyWidget{
		key: key,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets to the body.
func (wgt *TabBodyWidget) Add(children ...any) *TabBodyWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *TabBodyWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *TabBodyWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if wgt.parent == nil {
		return errors.New("tab body must be contained in a tab switcher")
	}
	nested := wgt.children
	if !wgt.Shown(r) {
		nested = nil
	}
	selected := wgt.parent.currentTabKey(r) == wgt.key
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Attr("data-name", wgt.parent.name).
		Attr("data-tab", wgt.key).
		Attr("role", "tabpanel").
		Attr("id", "tabbody_"+wgt.parent.name+"_"+wgt.key).
		Attr("aria-labelledby", "tablabel_"+wgt.parent.name+"_"+wgt.key).
		AttrIf(!selected, "hidden", "true").
		Class("TabBody", "Block").
		ClassIf(selected, "selected").
		Add(nested).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}

// Shown returns whether the widget is shown or hidden.
// A widget that is hidden is not rendered.
func (wgt *TabBodyWidget) Shown(r *http.Request) bool {
	return wgt.WidgetBase.Shown(r) && (!wgt.parent.dynamic || wgt.parent.currentTabKey(r) == wgt.key)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *TabBodyWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) ||
		(wgt.parent.dynamic && state.Changed(wgt.parent.name))
}
