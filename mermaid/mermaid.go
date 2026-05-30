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

package mermaid

import (
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

// MermaidWidget is a diagram rendered by the Mermaid JS library.
// See https://mermaid.js.org/intro/syntax-reference.html for the diagram
// grammar. Mermaid is distributed under the MIT License.
type MermaidWidget struct {
	*widget.WidgetBase[*MermaidWidget]
	source  string
	width   string
	height  string
	align   string
	zoomPan bool
}

/*
Mermaid creates a new Mermaid diagram widget. The source argument is the
mermaid diagram text — exactly what you would put inside a
```mermaid fenced code block in markdown. For example:

	graph TD
	    A[Start] --> B{Decision}
	    B -->|Yes| C[Done]
	    B -->|No|  D[Retry]

Rendering happens client-side on every page load. The diagram is themed
against the current Material design tokens and re-renders automatically
on theme switches.
*/
func (f MermaidFactory) Mermaid(source string) *MermaidWidget {
	x := &MermaidWidget{
		source: source,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithWidth sets the diagram container's width.
// The default is 100% of the parent. Pass any CSS length, e.g. "800px",
// "100%" or "calc(100vw - 2em)". Empty falls back to the CSS default.
func (wgt *MermaidWidget) WithWidth(css string) *MermaidWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithHeight sets the diagram container's height.
// Without zoom/pan the diagram renders at the natural pixel dimensions
// of its viewBox, with max-width:100% scaling it down (preserving
// aspect ratio) when the container is narrower than the content; a height
// is generally not needed there.
// When WithZoomPan is on, the container needs a bounded canvas to pan
// within; it defaults to 400px, and WithHeight overrides that. Pass any
// CSS length, e.g. "400px", "100%" or "calc(100vh - 50px)". Empty falls
// back to the CSS default.
func (wgt *MermaidWidget) WithHeight(css string) *MermaidWidget {
	if css != "" {
		wgt.height = "height:" + css
	} else {
		wgt.height = ""
	}
	return wgt
}

// WithAlign positions the SVG horizontally within the canvas when the
// diagram is narrower than the available width. Pass "left", "center"
// (default), or "right". Any other value is treated as "center". Has no
// effect when WithZoomPan is on, since the SVG then fills the canvas.
func (wgt *MermaidWidget) WithAlign(direction string) *MermaidWidget {
	switch direction {
	case "left", "right", "center":
		wgt.align = direction
	default:
		wgt.align = "center"
	}
	return wgt
}

// WithZoomPan enables interactive zoom (mouse wheel / pinch) and pan
// (click-and-drag) on the rendered diagram. Off by default. Double-click
// the diagram to recenter and reset the zoom.
func (wgt *MermaidWidget) WithZoomPan(on bool) *MermaidWidget {
	wgt.zoomPan = on
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *MermaidWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)
	styleParts := []string{}
	if wgt.width != "" {
		styleParts = append(styleParts, wgt.width)
	}
	if wgt.height != "" {
		styleParts = append(styleParts, wgt.height)
	}
	style := strings.Join(styleParts, ";")
	zoomPan := "0"
	if wgt.zoomPan {
		zoomPan = "1"
	}
	canvasClass := "MermaidCanvas"
	switch wgt.align {
	case "left":
		canvasClass += " MermaidAlignLeft"
	case "right":
		canvasClass += " MermaidAlignRight"
	}
	if wgt.zoomPan {
		canvasClass += " MermaidCanvasZoomPan"
	}

	// Mermaid's flowchart parser rejects CSS function syntax in classDef and
	// style values, so var(--x) references the author wrote naturally cannot
	// be passed through. expandVars rewrites those references to currentColor
	// in place and produces bridge CSS rules that resolve the original values
	// from the host page's cascade, scoped to this widget's container id.
	expanded := expandVars(wgt.source)
	bridgeCSS := formatCSS(expanded.rules, "#"+randomID)

	children := []any{}
	if bridgeCSS != "" {
		children = append(children, Tag("style").Add(HTMLUnsafe(bridgeCSS)))
	}
	children = append(children,
		Tag("div").
			Class(canvasClass).
			Attr("id", randomID).
			Style(style).
			Add(
				Tag("pre").
					Class("MermaidSource").
					Add(Text(expanded.source)),
			),
		Tag("script").Add(HTMLUnsafe(
			"\n(() => {\n",
			`mermaid_render('`, randomID, `',`, zoomPan, `);`, "\n",
			"})()\n",
		)),
	)
	return Tag("div").
		Class("Mermaid", "Block").
		Attr("data-id", wgt.ID()).
		Add(children...).
		When(wgt.Shown(r)).
		Draw(w, r)
}
