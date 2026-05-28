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

// Package code provides syntax-highlighted source code widgets.
//
// Highlighting is performed server-side via Chroma (Apache 2.0). The generated
// markup uses CSS classes that map to Material design color tokens, so themes
// switch with the rest of the page without re-rendering.
//
// This package is not part of the default widget library because Chroma's
// lexer set adds several MB to the binary. Import it explicitly:
//
//	import "github.com/microbus-io/bespa/code"
//
// See ATTRIBUTIONS.md for upstream licensing.
package code
