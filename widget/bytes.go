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

package widget

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var _ = Widget(&BytesWidget{}) // Ensure interface

// BytesWidget renders bytes and strings.
type BytesWidget struct {
	b []byte
}

// Bytes creates a new widget that renders bytes to the page without escaping.
//
// Warning! Use of this type presents a security risk: the content should come from a trusted source,
// as it will be rendered verbatim.
func (f WidgetFactory) Bytes(b []byte) *BytesWidget {
	x := &BytesWidget{b: b}
	return x
}

// HTMLUnsafe creates a new widget that renders HTML to the page without escaping it.
//
// Warning! Use of this type presents a security risk: the content should come from a trusted source,
// as it will be rendered verbatim.
func (f WidgetFactory) HTMLUnsafe(html ...string) *BytesWidget {
	combined := strings.Join(html, "")
	if combined == "" {
		return f.Void()
	}
	return f.Bytes([]byte(combined))
}

// HTML creates a new widget that renders HTML to the page without escaping it.
// The HTML is parsed and rerendered to make sure that it is well-formed.
// Scripts and event handlers are eliminated from the HTML.
func (f WidgetFactory) HTML(html ...string) *BytesWidget {
	combined := strings.Join(html, "")
	if combined == "" {
		return f.Void()
	}
	safe, _ := SafeHTML(combined)
	return f.Bytes([]byte(safe))
}

// Markdown creates a new widget that converts markdown to HTML and
// renders it to the page without escaping it.
// The HTML is parsed and rerendered to make sure that it is well-formed.
// Scripts and event handlers are eliminated from the HTML.
func (f WidgetFactory) Markdown(md ...string) *BytesWidget {
	combined := strings.Join(md, "")
	if combined == "" {
		return f.Void()
	}
	combined = strings.ReplaceAll(combined, "\r\n", "\n")
	b := markdown.ToHTML([]byte(combined), nil, nil)
	return f.HTML(string(b))
}

// Text creates a new widget that renders a string to the page after escaping it.
func (f WidgetFactory) Text(str ...string) *BytesWidget {
	combined := strings.Join(str, "")
	if combined == "" {
		return f.Void()
	}
	return f.HTMLUnsafe(html.EscapeString(combined))
}

// Date creates a new widget that renders a date in the US format "1/2/2006".
func (f WidgetFactory) Date(dateTime time.Time) *BytesWidget {
	if dateTime.IsZero() {
		return f.Text("")
	}
	return f.Text(dateTime.Format("1/2/06"))
}

// DateTime creates a new widget that renders a date and time in the US format "1/2/2006 3:04:05 PM".
func (f WidgetFactory) DateTime(dateTime time.Time) *BytesWidget {
	if dateTime.IsZero() {
		return f.Text("")
	}
	return f.Text(dateTime.Format("1/2/06 3:04 PM"))
}

// Time creates a new widget that renders a time in the US format "3:04:05 PM".
func (f WidgetFactory) Time(dateTime time.Time) *BytesWidget {
	if dateTime.IsZero() {
		return f.Text("")
	}
	return f.Text(dateTime.Format("3:04 PM"))
}

// Time creates a new widget that renders a date in the US format "1/2/2006" or
// time in the US format "3:04:05 PM" if the given time is today.
func (f WidgetFactory) DateOrTime(dateTime time.Time) *BytesWidget {
	if dateTime.IsZero() {
		return f.Text("")
	}
	now := time.Now().In(dateTime.Location())
	if dateTime.Year() == now.Year() && dateTime.Month() == now.Month() && dateTime.Day() == now.Day() {
		return f.Time(dateTime)
	}
	if dateTime.Year() == now.Year() {
		return f.Text(dateTime.Format("Jan 2"))
	}
	return f.Date(dateTime)
}

// TimeZone creates a new widget that renders a time zone in the format "MST (GMT-0700)".
func (f WidgetFactory) TimeZone(loc *time.Location) *BytesWidget {
	if loc == nil {
		return f.Text("")
	}
	return f.Text(time.Now().In(loc).Format("MST"))
}

// Duration creates a new widget that renders a duration in the format "2d 1h 0m 1s"
func (f WidgetFactory) Duration(dur time.Duration) *BytesWidget {
	seconds := int(dur.Seconds())
	days := seconds / (24 * 3600)
	seconds = seconds % (24 * 3600)
	hours := seconds / 3600
	seconds = seconds % 3600
	minutes := seconds / 60
	seconds = seconds % 60
	switch {
	case days == 0 && hours == 0 && minutes == 0:
		return f.Text(fmt.Sprintf("%ds", seconds))
	case days == 0 && hours == 0:
		return f.Text(fmt.Sprintf("%dm %ds", minutes, seconds))
	case days == 0:
		return f.Text(fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds))
	default:
		return f.Text(fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds))
	}
}

// Integer creates a new widget that renders an integer in US format "1,234".
func (f WidgetFactory) Integer(i int) *BytesWidget {
	return f.Text(message.NewPrinter(language.English).Sprintf("%d", i))
}

// Float creates a new widget that renders a float in US format "1,234.5678".
func (f WidgetFactory) Float(flt float64, precision int) *BytesWidget {
	return f.Text(message.NewPrinter(language.English).Sprintf("%.0"+strconv.Itoa(precision)+"f", flt))
}

// Void creates a new widget that does nothing.
func (f WidgetFactory) Void() *BytesWidget {
	return f.Bytes([]byte{})
}

// Any creates a new widget that best fit the type of the argument.
// Nils are tolerated. Unrecognized types are formatted into a string by default.
func (f WidgetFactory) Any(object any) Widget {
	if IsNil(object) {
		return f.Void()
	}
	if widget, ok := object.(Widget); ok {
		return widget
	}

	switch v := object.(type) {
	case string:
		return f.Text(v)
	case []byte:
		return f.Bytes(v)
	case int:
		return f.Integer(v)
	case int64:
		return f.Integer(int(v))
	case int32:
		return f.Integer(int(v))
	case float64:
		return f.Float(float64(v), 2)
	case float32:
		return f.Float(float64(float32(v)), 2)
	case time.Time:
		return f.DateTime(v)
	case *time.Location:
		return f.TimeZone(v)
	default:
		s := fmt.Sprintf("%v", v)
		return f.Text(s)
	}
}

// Many creates a collection of widgets that best fit the type of the arguments.
// Nils are tolerated. Unrecognized types are formatted into a string by default.
func (f WidgetFactory) Many(objects ...any) []Widget {
	widgets := make([]Widget, 0, len(objects))
	for _, obj := range objects {
		if IsNil(obj) {
			continue
		}
		switch v := obj.(type) {
		case []Widget:
			widgets = append(widgets, v...)
		case []any:
			// Call recursively
			widgets = append(widgets, f.Many(v...)...)
		default:
			widgets = append(widgets, f.Any(obj))
		}
	}
	return widgets
}

// ID is a unique identifier within the scope of the page.
// Bytes widgets disregard their ID.
func (wgt *BytesWidget) ID() string {
	return ""
}

// SetID sets a unique identifier within the scope of the page.
// Bytes widgets disregard their ID.
func (wgt *BytesWidget) SetID(id string) {
	// Noop
}

// Children are the widgets nested under this widget, or nil if none.
// Bytes widgets cannot have widgets nested under them.
func (wgt *BytesWidget) Children() []Widget {
	return nil
}

// Draw renders the widget's HTML.
func (wgt *BytesWidget) Draw(w io.Writer, r *http.Request) (err error) {
	// if wgt.ID() != "" {
	// 	_, err = w.Write([]byte(`<span data-id="` + html.EscapeString(wgt.ID()) + `">`))
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	_, err = w.Write(wgt.b)
	if err != nil {
		return err
	}
	// if wgt.ID() != "" {
	// 	_, err = w.Write([]byte("</span>"))
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// Drawn indicates whether this widget needs to be drawn in either a full or partial page rendering.
func (wgt *BytesWidget) Drawn(r *http.Request) bool {
	return false
}

// Shown indicates whether this widget is shown or hidden.
func (wgt *BytesWidget) Shown(r *http.Request) bool {
	return true
}

// String returns the text contained in the widget.
func (wgt *BytesWidget) String() string {
	return html.UnescapeString(string(wgt.b))
}
