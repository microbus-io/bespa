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

package richedit

import (
	"embed"

	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/widget"
)

// RichEditFactory aggregates the widget constructors of this package.
// Use RichEditFactory{} to construct a new factory.
type RichEditFactory struct{}

// factory is a collection of the dependent factories.
var factory = struct {
	widget.WidgetFactory
	basic.BasicFactory
	RichEditFactory
}{}

// Function aliases
var Any = factory.Any
var Many = factory.Many
var Text = factory.Text
var HTML = factory.HTML
var HTMLUnsafe = factory.HTMLUnsafe
var Bytes = factory.Bytes
var Tag = factory.Tag

// Wrapper CSS / JS and the Quill themes auto-join the main bundle.
// Quill core and the quill-mention extension are big enough to warrant
// isolated <script> tags — embed those separately below.
//
//go:embed inputrichtext.css inputrichtext.js quill.snow.css quill.mention.css
var bundle embed.FS

//go:embed quill.js
var quillJS string

//go:embed quill.mention.js
var quillMentionJS string

// init registers the assets of this library with the global registry.
func init() {
	widget.AssetRegistry.RegisterFS(bundle)
	widget.AssetRegistry.RegisterIsolatedScript("quill", quillJS)
	widget.AssetRegistry.RegisterIsolatedScript("quill.mention", quillMentionJS)
}

// Type aliases
type Widget = widget.Widget
type InputWidget = widget.InputWidget
type BytesWidget = widget.BytesWidget
