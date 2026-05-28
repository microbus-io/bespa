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
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/microbus-io/bespa/css"
)

// AssetRegistry is a singleton that keeps all the assets necessary to render widgets.
// These include CSS (templatized), JavaScript and HTML assets.
// All widgets must register their assets with the asset registry.
var AssetRegistry = &assets{}

type assets struct {
	mu                   sync.RWMutex
	styles               map[string]string
	scripts              map[string]string
	files                map[string][]byte
	isolatedScriptsOrder []string
	isolatedScripts      map[string]bool
	handlers             map[string]http.Handler
}

// FileSystem supports reading files and directories.
type FileSystem interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

// RegisterFS locates *.css, *.js, *.js.map and *.woff2 in the file system and
// adds them to the assets used to render the page.
func (a *assets) RegisterFS(fileSystem FileSystem) error {
	dir, err := fileSystem.ReadDir(".")
	if err != nil {
		return err
	}
	for _, file := range dir {
		fileName := file.Name()
		b, err := fileSystem.ReadFile(fileName)
		if err != nil {
			return err
		}
		err = a.Register(fileName, b)
		if err != nil {
			return err
		}
	}
	return nil
}

// Register adds *.css, *.js, *.js.map, *.woff2, *.png, *.jpg to the assets used to render the page.
func (a *assets) Register(fileName string, b []byte) error {
	switch {
	case strings.HasSuffix(fileName, ".css"):
		AssetRegistry.RegisterStyle(strings.TrimSuffix(fileName, ".css"), string(b))
	case strings.HasSuffix(fileName, ".js") && !strings.HasSuffix(fileName, ".src.js"):
		if bytes.Contains(b, []byte("//# source")) {
			AssetRegistry.RegisterIsolatedScript(strings.TrimSuffix(fileName, ".js"), string(b))
		} else {
			AssetRegistry.RegisterScript(strings.TrimSuffix(fileName, ".js"), string(b))
		}
	case strings.HasSuffix(fileName, ".woff2"):
		AssetRegistry.RegisterFile(fileName, b)
	case strings.HasSuffix(fileName, ".js.map"):
		AssetRegistry.RegisterFile(fileName, b)
	case strings.HasSuffix(fileName, ".src.js"):
		AssetRegistry.RegisterFile(fileName, b)
	case strings.HasSuffix(fileName, ".png"):
		AssetRegistry.RegisterFile(fileName, b)
	case strings.HasSuffix(fileName, ".jpg"):
		AssetRegistry.RegisterFile(fileName, b)
	case strings.HasSuffix(fileName, ".svg"):
		AssetRegistry.RegisterFile(fileName, b)
	}
	return nil
}

// RegisterStyle adds or replaces a CSS asset used to render the page.
func (a *assets) RegisterStyle(key string, css string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.styles == nil {
		a.styles = map[string]string{}
	}
	a.styles[key] = a.stripCopyrightComment(css) + "\n"
}

// RegisterScript adds or replaces a script asset used to render the page.
func (a *assets) RegisterScript(key string, js string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.scripts == nil {
		a.scripts = map[string]string{}
	}
	a.scripts[key] = a.stripCopyrightComment(js) + "\n"
}

// RegisterIsolatedScript adds or replaces a script asset used to render the page.
// Isolated scripts are served separately from the rest of the scripts.
func (a *assets) RegisterIsolatedScript(key string, js string) {
	a.mu.Lock()
	if a.isolatedScripts == nil {
		a.isolatedScripts = map[string]bool{}
		a.isolatedScriptsOrder = []string{}
	}
	if a.isolatedScripts[key] || js == "" {
		a.mu.Unlock()
		return
	}
	a.isolatedScripts[key] = true
	a.isolatedScriptsOrder = append(a.isolatedScriptsOrder, key)
	a.mu.Unlock()
	a.RegisterFile(key, []byte(js))
}

// RegisterHandler registers an http.Handler that owns a sub-path of the /bespa/ namespace.
// A request for /bespa/<prefix>/... is delegated to the handler if no other entry matches.
// Among overlapping registrations, the longest prefix wins. The handler sees the request
// unchanged — it can inspect r.URL.Path to dispatch within its own namespace.
func (a *assets) RegisterHandler(prefix string, handler http.Handler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.handlers == nil {
		a.handlers = map[string]http.Handler{}
	}
	a.handlers[prefix] = handler
}

// IsolatedScriptsOrder returns a snapshot of the registered isolated script keys in registration order.
func (a *assets) IsolatedScriptsOrder() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]string, len(a.isolatedScriptsOrder))
	copy(out, a.isolatedScriptsOrder)
	return out
}

// RegisterFile adds or replaces an arbitrary asset used to render the page.
func (a *assets) RegisterFile(key string, content []byte) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.files == nil {
		a.files = map[string][]byte{}
	}
	a.files[key] = content
}

// stripCopyrightComment removes the copyright notice at the top of JS and CSS files.
func (a *assets) stripCopyrightComment(code string) string {
	code = strings.TrimSpace(code)
	p := strings.Index(code, "/*")
	if p < 0 {
		return code
	}
	q := strings.Index(code[p+2:], "*/")
	if q < 0 {
		return code
	}
	if !strings.Contains(code[p+2:p+2+q], "Copyright") {
		return code
	}
	if p == 0 {
		return strings.TrimSpace(code[q+2+2:])
	}
	return code
}

// writeStyle writes the aggregated stylesheet assets to the writer.
func (a *assets) writeStyle(w io.Writer) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	keys := make([]string, 0, len(a.styles)+4)
	for k := range a.styles {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		_, err := w.Write([]byte(a.styles[key]))
		if err != nil {
			return err
		}
	}
	return nil
}

// writeScript writes the aggregated script assets to the writer.
func (a *assets) writeScript(w io.Writer) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	keys := make([]string, 0, len(a.scripts)+4)
	for k := range a.scripts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		_, err := w.Write([]byte(a.scripts[key]))
		if err != nil {
			return err
		}
	}
	return nil
}

// ServeHTTP serves the CSS and JavaScript assets.
// The handler should be added to the HTTP mux at /bespa/ .
func (a *assets) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}
	prefix := r.Header.Get("X-Forwarded-Prefix")
	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		if r.TLS != nil {
			proto = "https"
		} else {
			proto = "http"
		}
	}
	uri := r.URL.Path
	qm := strings.Index(uri, "?")
	if qm >= 0 {
		uri = uri[:qm]
	}

	// Cache resources for 30 days
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 30*24*60*60))

	if strings.HasSuffix(uri, "/style.css") {
		// Remove cache buster
		u := proto + "://" + host + prefix + "/bespa" + uri

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "/*# sourceURL=%s */\n\n", u)
		css.DefaultTypeScale.WriteCSS(&buf)
		css.DefaultKeyColors.WriteCSSThemes(&buf)
		a.writeStyle(&buf)
		w.Header().Set("Content-Type", "text/css")
		w.Write(buf.Bytes())
		return
	}

	if strings.HasSuffix(uri, "/tones.css") {
		arg := r.URL.Query().Get("keycolors")
		kc := css.KeyColorsFromString(arg)
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		kc.WriteCSSTones(w)
		return
	}

	if strings.HasSuffix(uri, "/script.js") {
		// Remove cache buster
		u := proto + "://" + host + prefix + "/bespa" + uri

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "//# sourceURL=%s\n\n", u)
		a.writeScript(&buf)
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(buf.Bytes())
		return
	}

	if strings.HasSuffix(uri, ".js") {
		// Remove cache buster
		u := proto + "://" + host + prefix + "/bespa" + uri

		p := strings.LastIndex(uri, "/")
		fileName := uri[p+1:]
		fileName = strings.TrimSuffix(fileName, ".js")
		a.mu.RLock()
		data, ok := a.files[fileName]
		a.mu.RUnlock()
		if ok {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "//# sourceURL=%s\n\n", u)
			buf.Write(data)
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Write(buf.Bytes())
			return
		}
	}

	// Files may be registered with sub-path keys (e.g. "maps/usa.json") so callers can
	// group assets under a logical namespace served at /bespa/<group>/<name>. The lookup
	// uses everything after /bespa/ as the key; legacy flat keys without slashes still
	// match because they were always served at /bespa/<name> in the first place.
	if key := strings.TrimPrefix(uri, "/bespa/"); key != "" && key != uri {
		a.mu.RLock()
		data, ok := a.files[key]
		a.mu.RUnlock()
		if ok {
			w.Header().Set("Content-Type", http.DetectContentType(data))
			w.Write(data)
			return
		}
		// Fall back to a registered handler whose prefix matches the longest portion of
		// the path. Handlers are how packages serve dynamically-generated assets.
		a.mu.RLock()
		var bestPrefix string
		var bestHandler http.Handler
		for prefix, h := range a.handlers {
			if strings.HasPrefix(key, prefix) && len(prefix) > len(bestPrefix) {
				bestPrefix = prefix
				bestHandler = h
			}
		}
		a.mu.RUnlock()
		if bestHandler != nil {
			bestHandler.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
