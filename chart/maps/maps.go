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

// Package maps registers a dynamic GeoJSON server for Apache ECharts maps.
// Blank-import this package to make every supported region available to chart widgets:
//
//	import _ "github.com/microbus-io/bespa/chart/maps"
//
// Supported URL patterns under /bespa/maps/:
//
//	world.json           — the full world country dataset
//	us.json              — the US states dataset (50 + DC + PR)
//	us-mainland.json     — the US states minus Alaska, Hawaii, and territories
//	us-<state>.json      — a single US state by two-letter postal abbrev (us-ca, us-tx, …)
//	ca.json              — Canada with provinces and territories
//	au.json              — Australia with states and territories
//	it.json              — Italy with regions
//	de.json              — Germany with states (Länder)
//	jp.json              — Japan with prefectures
//	cn.json              — China with provinces
//	in.json              — India with states and union territories
//	<country>.json       — for any other code, a single-country outline filtered from world
//
// Filtered responses are computed on first request and memoized. See ATTRIBUTIONS.md at
// the repository root for dataset sources and licenses.
package maps

import (
	"bytes"
	"embed"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/microbus-io/bespa/widget"
)

// All bundled GeoJSON datasets. Filtering by extension keeps .go source
// out of the binary.
//
//go:embed *.json
var bundle embed.FS

var (
	world      *indexedFC
	us         *indexedFC
	responseMu sync.RWMutex
	responses  = map[string][]byte{}
)

func init() {
	world = mustIndex(mustRead("world.json"))
	us = mustIndex(mustRead("us.json"))
	// Pre-warm the bundled subdivision datasets so requests for them are served from
	// the response cache without ever entering compute().
	for _, name := range []string{"world", "us", "ca", "au", "it", "de", "jp", "cn", "in"} {
		responses[name] = mustRead(name + ".json")
	}
	widget.AssetRegistry.RegisterHandler("maps/", http.HandlerFunc(serve))
}

func mustRead(name string) []byte {
	b, err := bundle.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return b
}

// indexedFC is a FeatureCollection that has been parsed enough to filter by name without
// re-decoding the whole document on every request. Each feature retains its original raw
// JSON for verbatim re-emission, plus its extracted name for matching.
type indexedFC struct {
	features []indexedFeature
}

type indexedFeature struct {
	raw  json.RawMessage
	name string
}

func mustIndex(data []byte) *indexedFC {
	var top struct {
		Features []json.RawMessage `json:"features"`
	}
	if err := json.Unmarshal(data, &top); err != nil {
		panic(err)
	}
	out := &indexedFC{features: make([]indexedFeature, 0, len(top.Features))}
	for _, raw := range top.Features {
		var probe struct {
			Properties struct {
				Name string `json:"name"`
			} `json:"properties"`
		}
		_ = json.Unmarshal(raw, &probe)
		out.features = append(out.features, indexedFeature{raw: raw, name: probe.Properties.Name})
	}
	return out
}

// filter returns a minimal FeatureCollection containing only the features whose name
// satisfies keep. The result is a fresh byte slice safe to send and cache.
func (fc *indexedFC) filter(keep func(name string) bool) []byte {
	var buf bytes.Buffer
	buf.WriteString(`{"type":"FeatureCollection","features":[`)
	first := true
	for _, f := range fc.features {
		if !keep(f.name) {
			continue
		}
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.Write(f.raw)
	}
	buf.WriteString(`]}`)
	return buf.Bytes()
}

func serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	const prefix = "/bespa/maps/"
	if !strings.HasPrefix(path, prefix) || !strings.HasSuffix(path, ".json") {
		http.NotFound(w, r)
		return
	}
	spec := strings.ToLower(path[len(prefix) : len(path)-len(".json")])
	if spec == "" {
		http.NotFound(w, r)
		return
	}

	responseMu.RLock()
	data, ok := responses[spec]
	responseMu.RUnlock()
	if !ok {
		var err error
		data, err = compute(spec)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		responseMu.Lock()
		responses[spec] = data
		responseMu.Unlock()
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "max-age=86400")
	w.Write(data)
}

func compute(spec string) ([]byte, error) {
	// us-* subdivisions
	if sub, ok := strings.CutPrefix(spec, "us-"); ok {
		if sub == "mainland" {
			return us.filter(func(name string) bool {
				return name != "Alaska" && name != "Hawaii" && name != "Puerto Rico"
			}), nil
		}
		stateName, ok := usStateByCode[sub]
		if !ok {
			return nil, errNotFound
		}
		return us.filter(func(name string) bool { return name == stateName }), nil
	}
	// Country code → country outline from world
	countryName, ok := countryByCode[spec]
	if !ok {
		return nil, errNotFound
	}
	return world.filter(func(name string) bool { return name == countryName }), nil
}

var errNotFound = &notFoundError{}

type notFoundError struct{}

func (*notFoundError) Error() string { return "map not found" }
