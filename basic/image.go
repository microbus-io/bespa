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

var _ = Widget(&ImageWidget{}) // Ensure interface

// ImageWidget renders an image.
type ImageWidget struct {
	*widget.WidgetBase[*ImageWidget]
	src     string
	altText string
	width   string
	height  string
}

// Image creates a new widget that renders an image.
func (f BasicFactory) Image(src string) *ImageWidget {
	x := &ImageWidget{
		src: src,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithAltText sets the alternate text of the image.
// This text is used by accessibility tools for the blind.
func (wgt *ImageWidget) WithAltText(altText string) *ImageWidget {
	wgt.altText = altText
	return wgt
}

// WithWidth scales the image to the given width.
// Unless explicitly set, the height is adjusted to maintain the aspect ratio.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
func (wgt *ImageWidget) WithWidth(css string) *ImageWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithHeight scales the image to the given height.
// Unless explicitly set, the width is adjusted to maintain the aspect ratio.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
func (wgt *ImageWidget) WithHeight(css string) *ImageWidget {
	if css != "" {
		wgt.height = "height:" + css
	} else {
		wgt.height = ""
	}
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *ImageWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("img").
		Class("PictureImage").
		Attr("data-id", wgt.ID()).
		Attr("src", wgt.src).
		Attr("alt", wgt.altText).
		Style(wgt.width, wgt.height).
		NoEnd().
		When(wgt.Shown(r)).
		Draw(w, r)
}
