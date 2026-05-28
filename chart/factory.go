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
	"embed"

	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/widget"
)

// ChartFactory aggregates the widget constructors of this package.
type ChartFactory struct{}

// factory is a collection of the dependent factories.
var factory = struct {
	widget.WidgetFactory
	basic.BasicFactory
	ChartFactory
}{}

// Function aliases
var Any = factory.Any
var Many = factory.Many
var Text = factory.Text
var HTML = factory.HTML
var HTMLUnsafe = factory.HTMLUnsafe
var Bytes = factory.Bytes
var Tag = factory.Tag

// Wrapper CSS and JS auto-join the main bundle via RegisterFS.
//
//go:embed chart.css chart.js
var bundle embed.FS

// ECharts ships without a sourcemap reference, so RegisterFS wouldn't
// auto-isolate it. Embed it separately and register explicitly so it
// becomes its own <script src> tag instead of bloating /bespa/script.js.
//
//go:embed echarts.js
var echartsJS string

// init registers the assets of this library with the global registry.
//
// Map data is not bundled here. To enable named maps (countries, US states,
// US-mainland, the world, etc.) blank-import the maps sub-package:
//
//	import _ "github.com/microbus-io/bespa/chart/maps"
//
// The sub-package registers a single handler that serves any supported region
// from /bespa/maps/<spec>.json — chart.js fetches the URL the first time a
// chart references the map name. See chart/maps for the supported URL grammar.
func init() {
	widget.AssetRegistry.RegisterFS(bundle)
	widget.AssetRegistry.RegisterIsolatedScript("echarts.js", echartsJS)
}

// Type aliases
type Widget = widget.Widget
type InputWidget = widget.InputWidget
type BytesWidget = widget.BytesWidget
