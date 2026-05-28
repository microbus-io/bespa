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
	"strconv"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&PrintfWidget{}) // Ensure interface

// PrintfWidget renders a collection of widgets formatted using a printf-like string.
type PrintfWidget struct {
	*widget.WidgetBase[*PrintfWidget]
	children []Widget
	format   string
}

// Printf creates a new widget that renders a collection of widgets formatted using a printf-like string.
// {n} appearing in the format string is replaced with the appropriate 0-indexed child.
// For example, "Lorem {0} ipsum {2} dolor sit amet {1}"
func (f BasicFactory) Printf(format string, children ...any) *PrintfWidget {
	x := &PrintfWidget{
		format:   format,
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *PrintfWidget) Add(children ...any) *PrintfWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *PrintfWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *PrintfWidget) Draw(w io.Writer, r *http.Request) (err error) {
	collection := factory.Collection()
	segments := strings.Split(wgt.format, "{")
	for i, seg := range segments {
		if seg == "" {
			continue
		}
		if i == 0 {
			collection.Add(seg)
		} else {
			p := strings.Index(seg, "}")
			if p > 0 {
				if i, err := strconv.Atoi(seg[:p]); err == nil && i < len(wgt.children) {
					collection.Add(wgt.children[i])
					collection.Add(seg[p+1:])
				} else {
					collection.Add(seg[p+1:])
				}
			} else {
				collection.Add(seg)
			}
		}
	}
	return collection.
		WithID(wgt.ID()).
		HideIf(!wgt.Shown(r)).
		Draw(w, r)
}
