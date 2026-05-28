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
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&EmbedWidget{}) // Ensure interface

// EmbedWidget renders an embedded page.
type EmbedWidget struct {
	*widget.WidgetBase[*EmbedWidget]
	webHandler func() (res *http.Response, err error)
	name       string
}

// Embed creates a new widget that splices the body of an HTTP response into
// the page. Only the content between <body> and </body> is inserted (or the
// whole response if those tags are absent). The fetcher is invoked every
// time the embed needs to render, including during partial redraws — keep
// it cheap or cache its result. For embedding an in-process handler, prefer
// EmbedHandler.
func (f BasicFactory) Embed(fetcher func() (res *http.Response, err error)) *EmbedWidget {
	x := &EmbedWidget{
		webHandler: fetcher,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// EmbedHandler creates a new widget that embeds an in-process HTTP handler
// as a sub-page. The handler is called with a synthetic request: method,
// URL (resolved relative to original.URL), and body are supplied here;
// headers, context, and RemoteAddr are inherited from original. Pass your
// app's mux for handler so the embedded URL routes through the same
// middleware stack. Compressed responses (gzip / deflate / brotli) are
// transparently decoded before the <body> is extracted.
func (f BasicFactory) EmbedHandler(handler http.HandlerFunc, original *http.Request, method string, url string, body io.ReadCloser) *EmbedWidget {
	x := &EmbedWidget{
		webHandler: func() (res *http.Response, err error) {
			rec := httptest.NewRecorder()
			u, err := original.URL.Parse(url)
			if err != nil {
				return nil, err
			}
			r, err := http.NewRequest(method, u.String(), body)
			if err != nil {
				return nil, err
			}
			r = r.WithContext(original.Context())
			for k, v := range original.Header {
				r.Header[k] = v
			}
			r.RemoteAddr = original.RemoteAddr
			handler(rec, r)
			return rec.Result(), nil
		},
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithName names this embed as a frame. Links and forms elsewhere on the
// page can target this frame via the `target` attribute or the `_target`
// state variable, and the framework will route their responses here.
// A named empty embed renders an empty page placeholder so the frame
// remains a valid swap target.
func (wgt *EmbedWidget) WithName(name string) *EmbedWidget {
	wgt.name = name
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *EmbedWidget) Draw(w io.Writer, r *http.Request) (err error) {
	var body []byte
	if wgt.Shown(r) && wgt.webHandler != nil {
		res, err := wgt.webHandler()
		if err != nil {
			return err
		}
		if res != nil && res.Body != nil {
			body, err = readDecoded(res)
			if err != nil {
				return err
			}
			q := 0
			p := bytes.Index(body, []byte("<body"))
			if p >= 0 {
				q = bytes.Index(body[p:], []byte(">"))
				if q >= 0 {
					q += p + 1
				} else {
					q = 0
				}
			}
			p = bytes.LastIndex(body, []byte("</body"))
			if p < 0 {
				p = len(body)
			}
			body = body[q:p]
		}
	}

	tag := Tag("div").
		Class("Embed").
		Attr("data-id", wgt.ID()).
		Attr("data-name", wgt.name)
	hasContent := false
	if len(body) > 0 {
		tag.Add(Bytes(body))
		hasContent = true
	} else if wgt.name != "" {
		tag.Add(factory.Page())
		hasContent = true
	}
	return tag.
		When(wgt.Shown(r) && hasContent).
		Draw(w, r)
}

// readDecoded reads the response body, transparently decompressing it if the
// upstream handler emitted a Content-Encoding header. This lets EmbedHandler
// safely call a mux that sits behind a brotli/gzip/deflate compression layer
// — the embedded fragment is decoded before being parsed for its <body>.
func readDecoded(res *http.Response) ([]byte, error) {
	var src io.Reader = res.Body
	switch strings.ToLower(res.Header.Get("Content-Encoding")) {
	case "gzip", "x-gzip":
		gz, err := gzip.NewReader(src)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		src = gz
	case "deflate":
		zr := flate.NewReader(src)
		defer zr.Close()
		src = zr
	case "br":
		src = brotli.NewReader(src)
	}
	return io.ReadAll(src)
}
