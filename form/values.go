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

// Values are the values posted for all the widgets of the form.
// The underlying structure is a map[string]string that can be iterated.
type Values map[string]string

// Get returns the named value.
func (v Values) Get(name string) (value string) {
	return map[string]string(v)[name]
}

// Has indicates if the named value was posted with the form.
func (v Values) Has(name string) bool {
	_, ok := map[string]string(v)[name]
	return ok
}

// Fields are the input widgets posted with the form.
// The underlying structure is a map[string]InputWidget that can be iterated.
type Fields map[string]InputWidget

// Get returns the named field.
func (f Fields) Get(name string) (field InputWidget) {
	return map[string]InputWidget(f)[name]
}

// Has indicates if the named field was posted with the form.
func (f Fields) Has(name string) bool {
	_, ok := map[string]InputWidget(f)[name]
	return ok
}
