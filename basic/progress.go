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
	"time"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&ProgressWidget{}) // Ensure interface

// ProgressWidget renders a progress bar.
type ProgressWidget struct {
	*widget.WidgetBase[*ProgressWidget]
	refreshURL string
	refresh    time.Duration
	max        int
	value      int
	width      string
	height     string
}

// Progress creates a new widget that renders a progress bar.
// Set the maximum (WithMax) to make it visible. Combine WithValue for
// static progress, or WithRefreshURL for a live-updating bar that polls
// a JSON endpoint.
func (f BasicFactory) Progress() *ProgressWidget {
	x := &ProgressWidget{
		max:     0,
		value:   0,
		refresh: time.Millisecond * 250,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithValue sets the bar's current value (clamped to 0..max by the
// browser). Pass a negative value to render an indeterminate / "infinite"
// progress animation instead.
func (wgt *ProgressWidget) WithValue(value int) *ProgressWidget {
	wgt.value = value
	return wgt
}

// WithMax sets the bar's maximum. The bar is hidden until this is set to
// a positive value (the default is 0). Negative values are ignored.
func (wgt *ProgressWidget) WithMax(max int) *ProgressWidget {
	if max >= 0 {
		wgt.max = max
	}
	return wgt
}

/*
WithRefreshURL sets a remote resource that updates a dynamic progress bar.
The response is expected in the form of a JSON object in the form:

	{
		"value": 123,			// A negative value results in an infinite progress bar
		"stop": true,			// Once stopped, a dynamic progress bar will no longer refresh
		"action": "?done=1",	// An action URL to call
	}
*/
func (wgt *ProgressWidget) WithRefreshURL(refreshURL string) *ProgressWidget {
	wgt.refreshURL = refreshURL
	return wgt
}

// WithRefreshInterval sets the polling cadence for the refresh URL.
// Default is 250ms. Non-positive values are ignored.
func (wgt *ProgressWidget) WithRefreshInterval(interval time.Duration) *ProgressWidget {
	if interval > 0 {
		wgt.refresh = interval
	}
	return wgt
}

// WithWidth sets the width of the progress bar.
// Pass any CSS length, e.g. "200px", "100%" or "calc(100% - 1em)". Empty clears it.
func (wgt *ProgressWidget) WithWidth(css string) *ProgressWidget {
	if css != "" {
		wgt.width = "width:" + css
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithHeight sets the height of the progress bar.
// Pass any CSS length, e.g. "8px", "1em". Empty clears it.
func (wgt *ProgressWidget) WithHeight(css string) *ProgressWidget {
	if css != "" {
		wgt.height = "height:" + css
	} else {
		wgt.height = ""
	}
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *ProgressWidget) Draw(w io.Writer, r *http.Request) (err error) {
	randomID := widget.RandomAlphaNumID(8)
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Class("LiveProgress").
		ClassIf(wgt.value < 0, "Infinite").
		Attr("id", randomID).
		Attr("data-interval", strconv.Itoa(int(wgt.refresh.Milliseconds()))).
		Attr("data-refresh", wgt.refreshURL).
		Add(
			Tag("progress").
				Attr("max", strconv.Itoa(wgt.max)).
				AttrIf(wgt.value >= 0, "value", strconv.Itoa(wgt.value)).
				Style(wgt.width, wgt.height),
			Tag("script").
				Add(factory.HTMLUnsafe("progress_start('", randomID, "')")),
			Tag("a").
				Attr("href", ""),
		).
		When(wgt.Shown(r) && wgt.max > 0).
		Draw(w, r)
}
