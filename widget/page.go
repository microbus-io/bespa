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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/microbus-io/bespa/css"
	"github.com/microbus-io/errors"
)

var _ = Widget(&PageWidget{}) // Ensure interface

var cacheBuster = RandomAlphaNumID(8)

// PageWidget renders a top-level HTML page.
type PageWidget struct {
	*WidgetBase[*PageWidget]
	children  []Widget
	nav       Widget
	keyColors css.KeyColors
	theme     string
	title     string
	target    string
	lang      string
}

// Page creates a new widget that renders a top-level HTML page.
func (f WidgetFactory) Page() *PageWidget {
	x := &PageWidget{
		keyColors: css.DefaultKeyColors,
		lang:      "en",
	}
	x.WidgetBase = NewWidgetBase(x)

	pc, file, _, ok := runtime.Caller(1)
	if ok {
		// Ex: /example/package/file.go
		fileDirs := strings.Split(file, string(filepath.Separator))
		if len(fileDirs) >= 2 {
			file = fileDirs[len(fileDirs)-2]
		}

		function := "Unknown"
		runtimeFunc := runtime.FuncForPC(pc)
		if runtimeFunc != nil {
			// Ex: example/package.(*Class).Function
			function = runtimeFunc.Name()
			p := strings.LastIndex(function, "/")
			if p >= 0 {
				function = function[p+1:]
			}
			p = strings.LastIndex(function, ".")
			if p >= 0 {
				function = function[p+1:]
			}
		}
		x.WithID("Page_" + file + "_" + function)
	}

	return x
}

// WithKeyColors sets the color palette of the page.
func (wgt *PageWidget) WithKeyColors(palette css.KeyColors) *PageWidget {
	wgt.keyColors = palette
	return wgt
}

// WithThemeDark forces the page to be rendered in the dark theme.
func (wgt *PageWidget) WithThemeDark() *PageWidget {
	wgt.theme = "DarkTheme"
	return wgt
}

// WithThemeLight forces the page to be rendered in the light theme.
func (wgt *PageWidget) WithThemeLight() *PageWidget {
	wgt.theme = "LightTheme"
	return wgt
}

// WithThemeDefault renders the page using the browser's theme.
func (wgt *PageWidget) WithThemeDefault() *PageWidget {
	wgt.theme = ""
	return wgt
}

// WithTarget sets a default target for links and forms embedded in the page.
func (wgt *PageWidget) WithTarget(target string) *PageWidget {
	wgt.target = target
	return wgt
}

// WithTitle sets the title of the page.
// Titles of top-level pages typically show in the browser.
func (wgt *PageWidget) WithTitle(title string) *PageWidget {
	wgt.title = title
	return wgt
}

// WithLang sets the BCP 47 language tag emitted on the <html> element
// (e.g. "en", "en-US", "fr", "zh-Hans"). Defaults to "en". Screen readers
// use this to pick a pronunciation dictionary; search engines use it for
// language-targeting. Empty string suppresses the attribute.
func (wgt *PageWidget) WithLang(lang string) *PageWidget {
	wgt.lang = lang
	return wgt
}

// Add adds nested widgets.
func (wgt *PageWidget) Add(children ...any) *PageWidget {
	adding := factory.Many(children)
	for _, a := range adding {
		if _, ok := a.(NavAreaMarker); ok {
			wgt.nav = a
		} else if pageTitle, ok := a.(PageTitleMarker); wgt.title == "" && ok {
			wgt.title = pageTitle.PageTitle()
			wgt.children = append(wgt.children, a)
		} else {
			wgt.children = append(wgt.children, a)
		}
	}
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *PageWidget) Children() []Widget {
	return factory.Many(wgt.nav, wgt.children)
}

// Draw renders the widget's HTML.
func (wgt *PageWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if rr, ok := w.(http.ResponseWriter); ok {
		rr.Header().Set("Content-Type", "text/html; charset=utf8")
	}
	wrt := NewWriterAssistant(w)

	// Read state
	state := factory.StateOf(r)
	redraw := state.HasChanges()

	// Traverse the widget tree structure
	var stats traversalStats
	err = wgt.traverse(wgt, r, !redraw, true, &stats)
	if err != nil {
		return errors.Trace(err)
	}
	// Set the initial values of input widgets into the state
	for n, v := range stats.Inputs {
		if state.Get(n) != v {
			state.Set(n, v)
		}
	}

	// Serialize state
	stateJson, err := json.Marshal(state)
	if err != nil {
		return errors.Trace(err)
	}

	if redraw {
		// Partial page redraw (may redraw nothing). The minimal <html><body> wrapper
		// is here only so the client can parse the response with new HTMLElement().
		wrt.WriteString("<html><body>")
		for _, widget := range stats.Draw {
			wrt.WriteString("\n<!-- ")
			wrt.WriteString(html.EscapeString(widget.ID()))
			wrt.WriteString(" -->")
			err = widget.Draw(w, r)
			if err != nil {
				return errors.Trace(err)
			}
		}
		wrt.WriteString("\n<!-- State -->")
		err = factory.Tag("div").Class("State").Add(stateJson).Draw(wrt, r)
		if err != nil {
			return errors.Trace(err)
		}
		wrt.WriteString("</body></html>\n")
		return wrt.Err()
	}

	// Full page initial drawing
	wrt.WriteString(`<html`)
	if wgt.lang != "" {
		wrt.WriteString(` lang="`, html.EscapeString(wgt.lang), `"`)
	}
	wrt.WriteString(` class="`)
	wrt.WriteString(wgt.userAgent(r))
	if wgt.theme != "" {
		wrt.WriteString(" ", wgt.theme)
	}
	wrt.WriteString(`">`)
	wrt.WriteString("<head>\n")
	wrt.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=yes" />`, "\n")
	wrt.WriteString(`<meta name="color-scheme" content="light dark">`)
	wrt.WriteString(`<link rel="stylesheet" href="/bespa/tones.css?keycolors=`, url.QueryEscape(wgt.keyColors.String()), `">`, "\n")
	wrt.WriteString(`<link rel="stylesheet" href="/bespa/style.css?id=`, cacheBuster, `">`, "\n")
	wrt.WriteString(`<script src="/bespa/script.js?id=`, cacheBuster, `"></script>`, "\n")
	for _, k := range AssetRegistry.IsolatedScriptsOrder() {
		wrt.WriteString(`<script src="/bespa/`+html.EscapeString(k)+`.js?id=`, cacheBuster, `"></script>`, "\n")
	}
	wrt.WriteString("<title>", html.EscapeString(wgt.title), "</title>\n")
	wrt.WriteString(`</head><body class="Top">`, "\n")
	wrt.WriteString(`<div class="FetchError" role="alert"><div class="ErrMsg"></div><i tabindex="0" role="button" aria-label="Dismiss error" class="Icon material-symbols-outlined" onclick="this.parentElement.style.display='none'">close</i></div>`, "\n")

	if wgt.nav != nil {
		err = factory.Tag("nav").Class("TopNav").Add(wgt.nav).Draw(wrt, r)
		if err != nil {
			return errors.Trace(err)
		}
	}

	u := *r.URL
	u.Host = strings.TrimSuffix(u.Host, ":443")
	u.RawQuery = ""
	location := u.String()
	location = strings.TrimPrefix(location, "https:/")
	// Strip the trailing slash so relative URLs resolve against the parent
	// path — but only if there is one. For a bare "/" the trim would yield
	// an empty string, which gets dropped from the HTML attributes; the JS
	// would then read getAttribute("data-location") as null and try to
	// resolve subsequent fetches against "/null".
	if location != "/" {
		location = strings.TrimSuffix(location, "/")
	}
	err = factory.Tag("div").
		Class("Page").
		Attr("data-id", wgt.ID()).
		Attr("data-location", location).
		Attr("data-target", wgt.target).
		Attr("onclick", "page_click(event)").
		Attr("onsubmit", "page_submit(event)").
		ClassIf(wgt.nav != nil, "HasNav").
		Add(
			factory.Tag("div").Class("State").Add(stateJson),
			wgt.children,
		).
		Draw(wrt, r)
	if err != nil {
		return errors.Trace(err)
	}

	wrt.WriteString("</body></html>\n")
	return wrt.Err()
}

// userAgent generates additional CSS classes for the HTML tag based on the user-agent header.
func (wgt *PageWidget) userAgent(r *http.Request) string {
	ua := strings.ToLower(r.Header.Get("User-Agent"))
	if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") {
		return " Safari"
	}
	return ""
}

type traversalStats struct {
	Lookup map[string]Widget
	Draw   []Widget
	Inputs map[string]string
}

// traverse traverses the widget tree of the page to gather information necessary for rendering.
func (wgt *PageWidget) traverse(w Widget, r *http.Request, drawing bool, shown bool, stats *traversalStats) error {
	// Index all widgets
	if w.ID() != "" {
		if stats.Lookup == nil {
			stats.Lookup = map[string]Widget{}
		}
		if stats.Lookup[w.ID()] != nil {
			return fmt.Errorf("duplicate ID '%s' for '%s'", w.ID(), typeOfWidget(w))
		}
		stats.Lookup[w.ID()] = w
	} else if _, ok := w.(*BytesWidget); !ok {
		return fmt.Errorf("missing ID for '%s'", typeOfWidget(w))
	}

	// Collect widgets that need to be drawn
	if !drawing && w.Drawn(r) {
		stats.Draw = append(stats.Draw, w)
		drawing = true // Applies to all descendants
	}
	if !w.Shown(r) {
		shown = false // Applies to all descendants
	}

	// Collect values from input widgets
	if input, ok := w.(InputWidget); ok && shown {
		name := input.Name()
		value := input.Value(r)
		if stats.Inputs == nil {
			stats.Inputs = map[string]string{}
		}
		stats.Inputs[name] = value
	}

	typeCounts := map[string]int{}
	for _, child := range w.Children() {
		// Create a stable ID for children
		typeName := typeOfWidget(child)
		typeCounts[typeName]++
		if child.ID() == "" {
			hash := md5.New()
			fmt.Fprintf(hash, "%s|%s|%d", w.ID(), typeName, typeCounts[typeName])
			stableID := hex.EncodeToString(hash.Sum(nil))[:12]
			if stats.Lookup[stableID] == nil {
				child.SetID(stableID)
			}
		}

		// Recurse
		err := wgt.traverse(child, r, drawing, shown, stats)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func typeOfWidget(widget Widget) string {
	typeOfWidget := reflect.TypeOf(widget)
	typeName := "?"
	if typeOfWidget.Kind() == reflect.Pointer {
		typeName = typeOfWidget.Elem().Name()
	}
	return typeName
}
