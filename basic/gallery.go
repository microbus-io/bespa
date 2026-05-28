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

var _ = Widget(&GalleryWidget{}) // Ensure interface

// GalleryWidget renders a gallery.
type GalleryWidget struct {
	*widget.WidgetBase[*GalleryWidget]
	children []Widget
}

// Gallery creates a new widget that lays out its children as wrap-around
// tiles — typically images or cards. It renders nothing when empty.
func (f BasicFactory) Gallery(children ...any) *GalleryWidget {
	x := &GalleryWidget{
		children: Many(children),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *GalleryWidget) Add(children ...any) *GalleryWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *GalleryWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *GalleryWidget) Draw(w io.Writer, r *http.Request) (err error) {
	divTag := Tag("div").
		Attr("data-id", wgt.ID()).
		Class("Gallery", "Block")
	for _, c := range wgt.children {
		divTag.Add(Tag("div").Add(c))
	}
	return divTag.
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
