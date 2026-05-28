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

import "net/http"

// InputWidget is a widget that collects input and updates state variables.
// InputWidgetBase is used to implement this interface.
type InputWidget interface {
	// Name of the input state variable.
	Name() string

	// Value of the input state variable, given an HTTP request.
	Value(r *http.Request) string

	// Valid indicates if the value is valid, given an HTTP request.
	Valid(r *http.Request) bool

	// Changed indicates if the value posted in an HTTP request differs from the initial value.
	Changed(r *http.Request) bool

	// SetFormName sets the name of the parent form that contains the input widget.
	// This name is used to detect if the form was submitted.
	SetFormName(formName string)
}
