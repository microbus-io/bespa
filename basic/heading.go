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

var _ = Widget(&HeadingWidget{}) // Ensure interface

// HeadingWidget renders a heading.
type HeadingWidget struct {
	*widget.WidgetBase[*HeadingWidget]
	level    string
	children []Widget
}

// heading creates a new widget that renders a heading.
func (f BasicFactory) heading(level string, children ...any) *HeadingWidget {
	x := &HeadingWidget{
		level:    level,
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// HeadlineLarge creates a new widget that renders a H1 heading.
func (f BasicFactory) HeadlineLarge(children ...any) *HeadingWidget {
	return f.heading("1", children...)
}

// HeadlineMedium creates a new widget that renders a H2 heading.
func (f BasicFactory) HeadlineMedium(children ...any) *HeadingWidget {
	return f.heading("2", children...)
}

// HeadlineSmall creates a new widget that renders a H3 heading.
func (f BasicFactory) HeadlineSmall(children ...any) *HeadingWidget {
	return f.heading("3", children...)
}

// TitleLarge creates a new widget that renders a H4 heading.
func (f BasicFactory) TitleLarge(children ...any) *HeadingWidget {
	return f.heading("4", children...)
}

// TitleMedium creates a new widget that renders a H5 heading.
func (f BasicFactory) TitleMedium(children ...any) *HeadingWidget {
	return f.heading("5", children...)
}

// TitleSmall creates a new widget that renders a H6 heading.
func (f BasicFactory) TitleSmall(children ...any) *HeadingWidget {
	return f.heading("6", children...)
}

// Add adds nested widgets.
func (wgt *HeadingWidget) Add(children ...any) *HeadingWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *HeadingWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *HeadingWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("h"+wgt.level).
		Class("Heading").
		Attr("data-id", wgt.ID()).
		Add(wgt.children).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
