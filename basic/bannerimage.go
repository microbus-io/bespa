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

var _ = Widget(&BannerImageWidget{}) // Ensure interface

// BannerImageWidget renders a banner image.
type BannerImageWidget struct {
	*widget.WidgetBase[*BannerImageWidget]
	src     string
	altText string
	height  string
	align   string
}

// BannerImage creates a new widget that renders an edge-to-edge image as a
// background within its container. Defaults: height 100% of the container,
// anchored center, no alt text.
func (f BasicFactory) BannerImage(src string) *BannerImageWidget {
	x := &BannerImageWidget{
		src:    src,
		height: "100%",
		align:  "center",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithAltText sets the alternate text of the image.
// This text is used by accessibility tools for the blind.
func (wgt *BannerImageWidget) WithAltText(altText string) *BannerImageWidget {
	wgt.altText = altText
	return wgt
}

// WithAnchorTop anchors the image to the top when the container is resized.
func (wgt *BannerImageWidget) WithAnchorTop() *BannerImageWidget {
	wgt.align = "top"
	return wgt
}

// WithAnchorBottom anchors the image to the bottom when the container is resized.
func (wgt *BannerImageWidget) WithAnchorBottom() *BannerImageWidget {
	wgt.align = "bottom"
	return wgt
}

// WithAnchorCenter anchors the image to the center when the container is resized.
// This is the default.
func (wgt *BannerImageWidget) WithAnchorCenter() *BannerImageWidget {
	wgt.align = "center"
	return wgt
}

// WithHeight scales the image to the given height. The width is always
// determined by the container; only the visible vertical slice changes.
// Pass any CSS length, e.g. "100px", "50%" or "calc(100vh - 50px)".
// Empty resets to the 100% default.
func (wgt *BannerImageWidget) WithHeight(css string) *BannerImageWidget {
	if css != "" {
		wgt.height = css
	} else {
		wgt.height = "100%"
	}
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *BannerImageWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("div").
		Class("BannerImage").
		Attr("data-id", wgt.ID()).
		Style(
			"background-image: url('"+wgt.src+"')",
			"height: "+wgt.height,
			"background-position: "+wgt.align,
		).
		Attr("alt", wgt.altText).
		When(wgt.Shown(r)).
		Draw(w, r)
}
