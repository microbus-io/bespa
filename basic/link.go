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
	"net/url"
	"strconv"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&LinkWidget{}) // Ensure interface

// LinkWidget renders a link.
type LinkWidget struct {
	*widget.WidgetBase[*LinkWidget]
	href     string
	disabled bool
	target   string
	back     bool
	children []Widget
}

// Link creates a new widget that renders an anchor.
// href accepts the full action-URL grammar: `?key=val` (apply state),
// `^?…` (parent page), `~?…` (top page), `/path` (full navigation),
// `path` (relative to the page's data-location). An empty href, or one
// that resolves empty, causes the link to render nothing.
func (f BasicFactory) Link(href string) *LinkWidget {
	x := &LinkWidget{
		href: href,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithHref replaces the link's href. Resets any prior WithHrefBack().
func (wgt *LinkWidget) WithHref(href string) *LinkWidget {
	wgt.back = false
	wgt.href = href
	return wgt
}

// WithHrefBack turns this into a "back" link. It follows the `_back` state
// variable when set; otherwise it falls back to a browser-history back step
// if the referrer is the same host. The link auto-hides when there's no
// history to go back to. Set `_back=0` to force-disable.
func (wgt *LinkWidget) WithHrefBack() *LinkWidget {
	wgt.back = true
	return wgt
}

// WithTarget sets the link's target. When unset, the page's `_target`
// state variable is used so the response routes into the active frame.
func (wgt *LinkWidget) WithTarget(target string) *LinkWidget {
	wgt.target = target
	return wgt
}

// WithDisabled greys out the link and removes its anchor — the children
// still render but no navigation occurs.
func (wgt *LinkWidget) WithDisabled(disabled bool) *LinkWidget {
	wgt.disabled = disabled
	return wgt
}

// Add adds nested widgets.
func (wgt *LinkWidget) Add(children ...any) *LinkWidget {
	wgt.children = Many(wgt.children, children)
	return wgt
}

// Children are the widgets nested under this widget.
func (wgt *LinkWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *LinkWidget) Draw(w io.Writer, r *http.Request) (err error) {
	state := factory.StateOf(r)
	href := wgt.href
	if wgt.back {
		href = state.Get("_back")
		if href == "" { // Default is go back one, if on same domain
			refHost := ""
			if ref, _ := url.Parse(r.Header.Get("Referer")); ref != nil {
				refHost = ref.Host
				if ref.Port() == "" {
					switch ref.Scheme {
					case "https":
						refHost += ":443"
					case "http":
						refHost += ":80"
					}
				}
			}
			host := r.Header.Get("X-Forwarded-Host")
			if host == "" {
				host = r.Host
			}
			if refHost == host {
				href = "-1"
			}
		}
		if href == "0" { // Disable back
			href = ""
		}
	}
	var detectionScript *widget.BytesWidget
	backLinkID := ""
	if backCount, err := strconv.Atoi(href); err == nil && backCount < 0 {
		href = "javascript:history.go(" + href + ")"
		// Detect if there's enough history and hide the link if not
		backLinkID = "back" + widget.RandomAlphaNumID(8)
		detectionScript = factory.HTMLUnsafe("<script>if (history.length<", strconv.Itoa(1-backCount), ") {document.getElementById('"+backLinkID+"').style.display='none';}</script>")
	}
	if wgt.disabled {
		return Tag("span").
			Class("Disabled").
			Attr("data-id", wgt.ID()).
			Add(wgt.children).
			When(wgt.Shown(r) && len(wgt.children) > 0 && href != "").
			Draw(w, r)
	}
	target := wgt.target
	if target == "" {
		target = state.Get("_target")
	}
	return Tag("a").
		Attr("data-id", wgt.ID()).
		Attr("id", backLinkID).
		Class("TextLink").
		Attr("href", href).
		Attr("target", target).
		Attr("tabindex", "0").
		Add(wgt.children, detectionScript).
		When(wgt.Shown(r) && len(wgt.children) > 0 && href != "").
		Draw(w, r)
}
