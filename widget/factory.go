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

import "embed"

// WidgetFactory aggregates the widget constructors of this package.
// Use WidgetFactory{} to construct a new factory
type WidgetFactory struct{}

// factory is a collection of the dependent factories.
var factory WidgetFactory

//go:embed *.css *.js
var bundle embed.FS

// init registers the assets of this library with the global registry.
func init() {
	AssetRegistry.RegisterFS(bundle)
}
