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
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Redirect writes a 307 redirect to the response.
func (f WidgetFactory) Redirect(w http.ResponseWriter, r *http.Request, location string) {
	backCount, _ := strconv.Atoi(location)

	// Add the _back URL of the current page to the redirected location
	if backCount >= 0 && r.URL.Query().Has("_back") {
		u, err := url.Parse(location)
		if err == nil {
			q := u.Query()
			if !q.Has("_back") {
				q.Set("_back", r.URL.Query().Get("_back"))
				u.RawQuery = q.Encode()
				location = u.String()
			}
		}
	}

	if r.Header.Get("Bespa-Fetch") == "1" {
		// Request is coming from JavaScript's fetch and must be handled manually.
		// Strip control characters to prevent the location string from breaking the protocol.
		safeLocation := strings.Map(func(c rune) rune {
			if c == '\r' || c == '\n' || c < 0x20 {
				return -1
			}
			return c
		}, location)
		w.Write([]byte("Location: " + safeLocation))
		return
	}

	if backCount < 0 {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<html><script>history.go(%d)</script></html>`, backCount-1)))
		return
	}

	xForwardedScheme := r.Header.Get("X-Forwarded-Proto")
	xForwardedHost := r.Header.Get("X-Forwarded-Host")
	xForwardedPrefix := r.Header.Get("X-Forwarded-Prefix")
	if xForwardedHost != "" {
		// Use the original request URL as the base URL for the redirection
		rr := *r
		rr.URL.Path = xForwardedPrefix + rr.URL.Path
		rr.URL.Host = xForwardedHost
		if xForwardedScheme != "" {
			rr.URL.Scheme = xForwardedScheme
		}
		http.Redirect(w, &rr, location, http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, location, http.StatusTemporaryRedirect)
	}
}

// RedirectBack redirects to the URL in the _back state variable.
// If a _back variable is not present, it does nothing and returns false.
func (f WidgetFactory) RedirectBack(w http.ResponseWriter, r *http.Request) (ok bool) {
	back := factory.StateOf(r).Get("_back")
	if back == "0" {
		return false
	}
	if back == "" {
		back = "-1"
	}
	f.Redirect(w, r, back)
	return true
}

/*

// RelPathOf returns the relative path of the request.
// The path starts with the sub path and includes the query arguments.
// It excludes both the microservice's and the ingress proxy's host names.
// It does not start with a slash.
// A request to https://localhost:8080/some.example:443/subpath?arg=val yields subpath?arg=val.
func RelPathOf(r *http.Request) string {
	u := r.URL.String()
	p := strings.Index(u, "://")
	if p < 0 {
		return u
	}
	q := strings.Index(u[p+3:], "/")
	if q < 0 {
		return ""
	}
	return u[p+3+q+1:]
}

// AbsPathOf returns the absolute path of the request.
// The path starts with the microservice's host name and then sub path and query.
// It excludes the ingress proxy's host names.
// It does not start with a slash.
// A request to https://localhost:8080/some.example:443/subpath?arg=val yields some.example:443/subpath?arg=val.
func AbsPathOf(r *http.Request) string {
	u := r.URL.String()
	p := strings.Index(u, "://")
	if p < 0 {
		return u
	}
	return u[p+3:]
}

*/
