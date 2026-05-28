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
	"time"

	"github.com/microbus-io/bespa/widget"
)

// TearOffCalendarWidget renders a tear-off calendar with month, day and weekday.
type TearOffCalendarWidget struct {
	*widget.WidgetBase[*TearOffCalendarWidget]
	date time.Time
}

// TearOffCalendar creates a new widget that renders a stylized date block
// showing month, day, and weekday — the look of a desk calendar page.
// A zero-value date renders nothing.
func (f BasicFactory) TearOffCalendar(date time.Time) *TearOffCalendarWidget {
	x := &TearOffCalendarWidget{
		date: date,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Children are the widgets nested under this widget.
func (wgt *TearOffCalendarWidget) Children() []Widget {
	return nil
}

// Draw renders the widget's HTML.
func (wgt *TearOffCalendarWidget) Draw(w io.Writer, r *http.Request) (err error) {
	month := wgt.date.Format("Jan")
	day := wgt.date.Day()
	weekday := wgt.date.Format("Mon")

	return Tag("div").
		Class("TearOffCalendar").
		Attr("data-id", wgt.ID()).
		Add(
			Tag("div").Add((factory.TitleSmall(month))),
			Tag("div").Add(factory.HeadlineLarge(day)),
			Tag("div").Add(factory.TextLightweight(weekday)),
		).
		When(wgt.Shown(r) && !wgt.date.IsZero()).
		Draw(w, r)
}
