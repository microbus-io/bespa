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
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// Resettable allows resetting the reader of the request body stream.
type Resettable interface {
	Reset()
}

// State is the collection of the key-value pairs transferred between the frontend and the backend.
// It is a thin wrapper over the request's form.
type State url.Values

// StateOf returns the state variables of the HTTP request. The returned
// State is never nil — if the request body could not be parsed as an HTML
// form, an empty State is returned so callers can safely Get/Has without
// a nil check.
func (f WidgetFactory) StateOf(r *http.Request) State {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			return State{}
		}
		if resettable, ok := r.Body.(Resettable); ok {
			resettable.Reset()
		}
	}
	return State(r.Form)
}

// Get returns the value of a state variable or the empty string if not found.
func (s State) Get(key string) string {
	return url.Values(s).Get(key)
}

// Set sets the value of a state variable, replacing any existing values.
func (s State) Set(key string, value string) {
	url.Values(s).Set(key, value)
}

// Del deletes the value of a state variable.
func (s State) Del(key string) {
	url.Values(s).Del(key)
	ch := s.Get("_changed")
	if ch != "" {
		chNew := ""
		for _, k := range strings.Split(ch, ",") {
			if key != k {
				if chNew != "" {
					chNew += ","
				}
				chNew += k
			}
		}
		s.Set("_changed", chNew)
		// _changed may be "" at this point, in particular when resetting a form.
		// It's important not to remove _changed completely because it serves as
		// indication of a full drawing vs partial page redrawing.
	}
}

// Has indicates if the state variable is found.
func (s State) Has(key string) bool {
	return url.Values(s).Has(key) || s.Changed(key)
}

// Changed returns the list of state variables that were reported changed by the frontend.
func (s State) Changed(key string) bool {
	ch := s.Get("_changed")
	if ch == "" {
		return false
	}
	for _, k := range strings.Split(ch, ",") {
		if key == k {
			return true
		}
	}
	return false
}

// HasChanges indicates if the request includes changes to the state.
// This indicates a partial redrawing request.
func (s State) HasChanges() bool {
	return s.Has("_changed")
}

// MarshalJSON marshals the state as a JSON map, excluding any variables that begin
// with an underscore, with the exception of the _back variable that is marshaled.
func (s State) MarshalJSON() ([]byte, error) {
	m := map[string]string{}
	for k, v := range url.Values(s) {
		if k == "_back" || !strings.HasPrefix(k, "_") {
			m[k] = v[0]
		}
	}
	b, err := json.Marshal(m)
	return b, err
}
