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

package chart

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/bespa/widget"
	"github.com/microbus-io/errors"
)

// ChartWidget is a chart from the Apache ECharts library.
// See https://echarts.apache.org/en/option.html for the configuration reference.
// ECharts is distributed under the Apache License 2.0.
type ChartWidget struct {
	*widget.WidgetBase[*ChartWidget]
	configTemplate string
	data           any
	svg            bool
	height         string
}

/*
Chart creates a new ECharts widget.

configTemplate is a Go html/template that, once executed against data, must
produce a JavaScript object literal matching the ECharts option spec at
https://echarts.apache.org/en/option.html. The template is parsed and
executed on every Draw, so keep it small or precompute heavy strings.

Any user-controlled string interpolated into the JS object literal MUST be
passed through the `escape` template function, which escapes single and
double quotes:

	title: {
	    text: '{{ .Title | escape }}',
	}

Note: `escape` only protects against breaking out of a JS string literal —
it is not sufficient if you concatenate untrusted text into JS structure or
HTML. Prefer numbers and pre-validated values for non-string fields.

Use AsSVG(true) for print-friendly output. To use named maps (e.g. world or
US states) blank-import "github.com/microbus-io/bespa/chart/maps".
*/
func (f ChartFactory) Chart(configTemplate string, data any) *ChartWidget {
	x := &ChartWidget{
		configTemplate: configTemplate,
		data:           data,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// AsSVG switches the chart to ECharts' SVG renderer instead of the default
// canvas renderer. SVG output scales cleanly for print and high-DPI
// displays at some cost in interactive performance for very large datasets.
func (wgt *ChartWidget) AsSVG(svg bool) *ChartWidget {
	wgt.svg = svg
	return wgt
}

// WithHeight sets the chart container's height.
// Pass any CSS length, e.g. "400px", "100%" or "calc(100vh - 50px)".
// The default (set via CSS) is 400px. Empty falls back to the CSS default.
func (wgt *ChartWidget) WithHeight(css string) *ChartWidget {
	wgt.height = css
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *ChartWidget) Draw(w io.Writer, r *http.Request) (err error) {
	funcMap := template.FuncMap{
		"escape": func(s string) string {
			s = strings.ReplaceAll(s, `'`, `\'`)
			s = strings.ReplaceAll(s, `"`, `\"`)
			return s
		},
	}
	htmlTmpl, err := template.New("chart").Funcs(funcMap).Parse(wgt.configTemplate)
	if err != nil {
		return errors.Trace(err)
	}
	var buf bytes.Buffer
	err = htmlTmpl.ExecuteTemplate(&buf, "chart", wgt.data)
	if err != nil {
		return errors.Trace(err)
	}
	config := buf.String()
	randomID := widget.RandomAlphaNumID(8)
	renderer := "canvas"
	if wgt.svg {
		renderer = "svg"
	}
	heightStyle := ""
	if wgt.height != "" {
		heightStyle = "height:" + wgt.height
	}
	return Tag("div").
		Class("Chart", "Block").
		Attr("data-id", wgt.ID()).
		Add(
			Tag("div").
				Class("ChartCanvas").
				Attr("id", randomID).
				Style(heightStyle),
			Tag("script").Add(HTMLUnsafe(
				"\n(() => {\n",
				`const config = `, config, `;`, "\n",
				`chart_chart('`, randomID, `',config,'`, renderer, `');`, "\n",
				"})()\n",
			)),
		).
		When(wgt.Shown(r)).
		Draw(w, r)
}
