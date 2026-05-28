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

var _ = Widget(&SplitterWidget{}) // Ensure interface

// SplitterWidget renders content in columns.
type SplitterWidget struct {
	*widget.WidgetBase[*SplitterWidget]
	colChildren [][]Widget
	colWidth    []int
	wrap        string
}

// Splitter creates a new widget that lays out content in columns.
// Each width is a flex weight relative to the others; 0 means auto-size
// to the column's content. Empty columns are omitted entirely.
// Use AddLeft, AddRight, or AddToCol to populate columns.
// Splitter() with no widths is shorthand for a single auto-sized column.
func (f BasicFactory) Splitter(widths ...int) *SplitterWidget {
	if len(widths) == 0 {
		widths = []int{0}
	}
	x := &SplitterWidget{
		colWidth:    widths,
		colChildren: make([][]Widget, len(widths)),
		wrap:        "Wrap",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithWrap controls whether the splitter stacks columns vertically on
// narrow (<600px) viewports. When wrapped, columns stack in declaration
// order: left → middle → … → right. Default is true.
func (wgt *SplitterWidget) WithWrap(wrap bool) *SplitterWidget {
	if !wrap {
		wgt.wrap = ""
	} else {
		wgt.wrap = "Wrap"
	}
	return wgt
}

// AddLeft adds nested widgets to the leftmost column.
func (wgt *SplitterWidget) AddLeft(leftChildren ...any) *SplitterWidget {
	return wgt.AddToCol(0, leftChildren...)
}

// AddRight adds nested widgets to the rightmost column.
func (wgt *SplitterWidget) AddRight(rightChildren ...any) *SplitterWidget {
	return wgt.AddToCol(len(wgt.colChildren)-1, rightChildren...)
}

// AddToCol adds nested widgets to the 0-indexed column.
func (wgt *SplitterWidget) AddToCol(index int, children ...any) *SplitterWidget {
	wgt.colChildren[index] = Many(wgt.colChildren[index], children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *SplitterWidget) Children() []Widget {
	var result []Widget
	for _, col := range wgt.colChildren {
		result = Many(result, col)
	}
	return result
}

// Draw renders the widget's HTML.
func (wgt *SplitterWidget) Draw(w io.Writer, r *http.Request) (err error) {
	splitter := Tag("div").
		Class("Splitter", wgt.wrap).
		Attr("data-observe-width", "600").
		Attr("data-id", wgt.ID())
	hasItems := false
	// totalWidth := 0
	// for i := range wgt.colChildren {
	// 	if len(wgt.colChildren[i]) == 0 {
	// 		continue
	// 	}
	// 	totalWidth += wgt.colWidth[i]
	// }
	for i := range wgt.colChildren {
		if len(wgt.colChildren[i]) == 0 {
			continue
		}
		tag := Tag("div").Add(wgt.colChildren[i])
		if wgt.colWidth[i] > 0 {
			tag.Style(fmt.Sprintf("flex: %d", wgt.colWidth[i]))
			// tag.Style(fmt.Sprintf("flex-shrink: %d", wgt.colWidth[i]))
			// tag.Style(fmt.Sprintf("flex-basis: %.2f%%", float64(wgt.colWidth[i])*100.0/float64(totalWidth)))
		} else {
			tag.Style("flex-grow: 0; flex-shrink: 0")
		}
		splitter.Add(tag)
		hasItems = true
	}
	return splitter.
		When(wgt.Shown(r) && hasItems).
		Draw(w, r)
}
