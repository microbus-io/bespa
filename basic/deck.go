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
	"fmt"
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&DeckWidget{}) // Ensure interface

// DeckWidget renders a deck of cards.
type DeckWidget struct {
	*widget.WidgetBase[*DeckWidget]
	children     []Widget
	colsNarrow   int
	colsWide     int
	colsExpanded int
}

// Deck creates a new widget that arranges its children into equal-width
// columns with an 8px gutter, choosing between three column counts based on
// the container's measured width: colsNarrow (<600px), colsWide (600–1200px),
// and colsExpanded (>1200px). A typical call is Deck(1, 2, 4).
func (f BasicFactory) Deck(colsNarrow, colsWide, colsExpanded int) *DeckWidget {
	x := &DeckWidget{
		colsNarrow:   colsNarrow,
		colsWide:     colsWide,
		colsExpanded: colsExpanded,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Add adds nested widgets.
func (wgt *DeckWidget) Add(children ...any) *DeckWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *DeckWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *DeckWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)
	style := fmt.Sprintf(`
	#%s.Width_600 > DIV {
		width: calc(%.8f%% - %.8fpx);
	}
	#%s.Width600_1200 > DIV {
		width: calc(%.8f%% - %.8fpx);
	}
	#%s.Width1200_ > DIV {
		width: calc(%.8f%% - %.8fpx);
	}`,
		randomID, 100.0/float64(wgt.colsNarrow), 8.0*float64(wgt.colsNarrow-1)/float64(wgt.colsNarrow),
		randomID, 100.0/float64(wgt.colsWide), 8.0*float64(wgt.colsWide-1)/float64(wgt.colsWide),
		randomID, 100.0/float64(wgt.colsExpanded), 8.0*float64(wgt.colsExpanded-1)/float64(wgt.colsExpanded),
	)
	styleTag := Tag("style").Add(factory.HTMLUnsafe(style))
	divTag := Tag("div").
		Attr("id", randomID).
		Attr("data-observe-width", "600,1200")
	for _, c := range wgt.children {
		divTag.Add(Tag("div").Add(c))
	}
	return Tag("div").
		Class("Deck", "Block").
		Attr("data-id", wgt.ID()).
		Add(styleTag, divTag).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}
