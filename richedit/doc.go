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

// Package richedit wraps Quill 2 as a rich-text input widget.
//
// Quill is distributed under the BSD-3-Clause license; the quill-mention add-on used
// to implement @ mentions is MIT. See ATTRIBUTIONS.md.
//
// The widget exposes the same surface as the legacy ckeditor package where the
// underlying capability exists in Quill. Features that have no built-in equivalent
// (table-properties / cell-properties dialogs, special-characters menu, find &
// replace, source view) are intentionally not surfaced.
package richedit
