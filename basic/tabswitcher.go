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
	"github.com/microbus-io/errors"
)

var _ = Widget(&TabSwitcherWidget{}) // Ensure interface

// TabSwitcherWidget renders a tab switcher.
type TabSwitcherWidget struct {
	*widget.WidgetBase[*TabSwitcherWidget]
	toolbar       *ToolbarWidget
	children      []Widget
	name          string
	initialTabKey string
	line          bool
	dynamic       bool
	labels        []*TabLabelWidget
	bodies        []*TabBodyWidget
}

// TabSwitcher creates a new widget that groups a row of TabLabels with
// their corresponding TabBodies (matched by key). The active tab is held
// in a state variable whose name defaults to "tab" — use WithName to
// disambiguate when there are multiple switchers on a page.
func (f BasicFactory) TabSwitcher() *TabSwitcherWidget {
	x := &TabSwitcherWidget{
		name:    "tab",
		toolbar: factory.Toolbar(),
		line:    true,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithName sets the state variable that holds the active tab key.
// Multiple tab switchers on the same page must use distinct names so they
// don't fight over the same key. Default is "tab".
func (wgt *TabSwitcherWidget) WithName(name string) *TabSwitcherWidget {
	wgt.name = name
	return wgt
}

// WithSelected sets the tab to show on first render. Overridden once the
// user clicks a tab and the state variable is set. Defaults to the first
// label/body added.
func (wgt *TabSwitcherWidget) WithSelected(tabKey string) *TabSwitcherWidget {
	wgt.initialTabKey = tabKey
	return wgt
}

// WithLine indicates if to draw a line below the tab labels.
// The default behavior is to draw a line.
func (wgt *TabSwitcherWidget) WithLine(line bool) *TabSwitcherWidget {
	wgt.line = line
	return wgt
}

// WithDynamic switches the tab switcher to server-side rendering: only the
// active tab's body is sent, and switching tabs triggers a partial redraw.
// Default is false — all tabs are rendered up-front and switched on the
// client. Turn this on when individual tabs are expensive to render or
// contain data that must reflect a tab change.
func (wgt *TabSwitcherWidget) WithDynamic(dynamic bool) *TabSwitcherWidget {
	wgt.dynamic = dynamic
	return wgt
}

// AddLeft places TabLabels into the left side of the tab strip and queues
// their TabBodies. Non-tab children are silently discarded.
func (wgt *TabSwitcherWidget) AddLeft(children ...any) *TabSwitcherWidget {
	for _, c := range children {
		if lbl, ok := c.(*TabLabelWidget); ok {
			lbl.parent = wgt // Link the label to this switcher
			wgt.labels = append(wgt.labels, lbl)
			wgt.toolbar.AddLeft(c)
		}
		if bod, ok := c.(*TabBodyWidget); ok {
			bod.parent = wgt // Link the body to this switcher
			wgt.bodies = append(wgt.bodies, bod)
			wgt.children = Many(wgt.children, bod)
		}
	}
	return wgt
}

// AddRight places TabLabels into the right side of the tab strip and queues
// their TabBodies. Non-tab children are silently discarded.
func (wgt *TabSwitcherWidget) AddRight(children ...any) *TabSwitcherWidget {
	for _, c := range children {
		if lbl, ok := c.(*TabLabelWidget); ok {
			lbl.parent = wgt // Link the label to this switcher
			wgt.labels = append(wgt.labels, lbl)
			wgt.toolbar.AddRight(c)
		}
		if bod, ok := c.(*TabBodyWidget); ok {
			bod.parent = wgt // Link the body to this switcher
			wgt.bodies = append(wgt.bodies, bod)
			wgt.children = Many(wgt.children, bod)
		}
	}
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *TabSwitcherWidget) Children() []Widget {
	return Many(wgt.toolbar, wgt.children)
}

// Draw renders the widget's HTML.
func (wgt *TabSwitcherWidget) Draw(w io.Writer, r *http.Request) (err error) {
	err = Tag("div").
		Attr("data-id", wgt.ID()).
		Attr("data-name", wgt.name).
		Class("TabSwitcher", "Block").
		ClassIf(!wgt.line, "NoLine").
		Add(
			Tag("div").
				Class("TabLabels", "Block").
				Attr("role", "tablist").
				Add(wgt.toolbar).
				When(len(wgt.toolbar.Children()) > 0),
			Tag("div").Class("TabBodies", "Block").Add(wgt.children).When(len(wgt.children) > 0),
		).
		When(wgt.Shown(r) && len(wgt.labels)+len(wgt.bodies) > 0).
		Draw(w, r)
	return errors.Trace(err)
}

// currentTabKey returns the currently selected tab key.
// It is called by the nested tab labels and tab bodies to detect if they are the current tab.
func (wgt *TabSwitcherWidget) currentTabKey(r *http.Request) string {
	state := factory.StateOf(r)
	currentTabKey := wgt.initialTabKey
	if currentTabKey == "" {
		if len(wgt.labels) > 0 {
			currentTabKey = wgt.labels[0].key
		} else if len(wgt.bodies) > 0 {
			currentTabKey = wgt.bodies[0].key
		}
	}
	if state.Has(wgt.name) {
		currentTabKey = state.Get(wgt.name)
	}
	return currentTabKey
}
