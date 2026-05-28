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

var _ = Widget(&MessageBarWidget{}) // Ensure interface

// MessageBarWidget renders a message bar.
type MessageBarWidget struct {
	*widget.WidgetBase[*MessageBarWidget]
	style    string
	children []Widget
}

// messageBar creates a new widget that renders a message bar.
func (f BasicFactory) messageBar(style string, children ...any) *MessageBarWidget {
	x := &MessageBarWidget{
		style:    style,
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// MessageBar creates a new widget that renders message bar in the primary color.
func (f BasicFactory) MessageBar(children ...any) *MessageBarWidget {
	return f.messageBar("Primary", children...)
}

// MessageBarError creates a new widget that renders message bar in the alert color (typically red).
func (f BasicFactory) MessageBarError(children ...any) *MessageBarWidget {
	return f.messageBar("Error", children...)
}

// Add adds nested widgets.
func (wgt *MessageBarWidget) Add(children ...any) *MessageBarWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *MessageBarWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *MessageBarWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Class("MessageBar", "Block", wgt.style).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
