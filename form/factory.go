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

package form

import (
	"embed"
	"encoding/json"
	"strings"

	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/widget"
)

// FormFactory aggregates the widget constructors of this package.
// Use FormFactory{} to construct a new factory
type FormFactory struct{}

// factory is a collection of the dependent factories.
var factory = struct {
	widget.WidgetFactory
	basic.BasicFactory
	FormFactory
}{}

// Function aliases
var Any = factory.Any
var Many = factory.Many
var Text = factory.Text
var HTML = factory.HTML
var HTMLUnsafe = factory.HTMLUnsafe
var Bytes = factory.Bytes
var Tag = factory.Tag

//go:embed *.css *.js
var bundle embed.FS

// init registers the assets of this library with the global registry.
func init() {
	widget.AssetRegistry.RegisterFS(bundle)
}

// Type aliases
type Widget = widget.Widget
type InputWidget = widget.InputWidget
type BytesWidget = widget.BytesWidget

// jsLiteral encodes s as a JSON-quoted JavaScript string literal, then
// unicode-escapes < and > so that an embedded "</script>" in the value
// cannot terminate the surrounding <script> tag.
func jsLiteral(s string) string {
	b, _ := json.Marshal(s)
	out := string(b)
	out = strings.ReplaceAll(out, "<", "\\u003c")
	out = strings.ReplaceAll(out, ">", "\\u003e")
	return out
}

// customValidityScript returns a <script> tag that calls
// input_setCustomValidity(target, message) with both arguments safely
// encoded for interpolation into an inline script.
func customValidityScript(target, msg string) *widget.TagWidget {
	return Tag("script").Add(factory.HTMLUnsafe(
		`input_setCustomValidity(`, jsLiteral(target), `, `, jsLiteral(msg), `)`,
	))
}
