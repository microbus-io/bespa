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

// InputWidgetBase facilitates the implementation of the InputWidget interface.
type InputWidgetBase[T Widget] struct {
	*WidgetBase[T]
	disabled   bool
	required   bool
	autoSubmit bool
	name       string
	owner      T
	formName   string
}

// NewInputWidgetBase creates a new input widget base holding a pointer to its owner subclass.
func NewInputWidgetBase[T Widget](owner T) *InputWidgetBase[T] {
	x := &InputWidgetBase[T]{
		owner: owner,
	}
	x.WidgetBase = NewWidgetBase(owner)
	return x
}

// SetFormName sets the name of the parent form that contains the input widget.
// This name is used to detect if the form was submitted.
func (wb *InputWidgetBase[T]) SetFormName(formName string) {
	wb.formName = formName
}

// Name is the name of the state variable affected by this input widget.
func (wb *InputWidgetBase[T]) Name() string {
	return wb.name
}

// WithName sets the name of the state variable affected by this input widget.
func (wb *InputWidgetBase[T]) WithName(name string) T {
	wb.name = name
	return wb.owner
}

// Disabled indicates if the input widget is disabled.
func (wb *InputWidgetBase[T]) Disabled() bool {
	return wb.disabled
}

// WithDisabled sets whether this input widget is enabled or disabled.
// A disabled widget renders greyed out and its value is not posted back.
// A widget is enabled by default.
func (wb *InputWidgetBase[T]) WithDisabled(disabled bool) T {
	wb.disabled = disabled
	return wb.owner
}

// Required indicates whether the input widget requires a value or not.
func (wb *InputWidgetBase[T]) Required() bool {
	return wb.required
}

// WithRequired sets whether the input widget requires a value or not.
// Values are not required by default.
func (wb *InputWidgetBase[T]) WithRequired(required bool) T {
	wb.required = required
	return wb.owner
}

// AutoSubmit returns a URL that is auto-submitted along with the input's value.
func (wb *InputWidgetBase[T]) AutoSubmit() bool {
	return wb.autoSubmit
}

// WithAutoSubmit sets a URL that is auto-submitted along with the input's value.
func (wb *InputWidgetBase[T]) WithAutoSubmit(autoSubmit bool) T {
	wb.autoSubmit = autoSubmit
	return wb.owner
}

// Submitted indicates if the parent form of the input widget was submitted.
func (wb *InputWidgetBase[T]) Submitted(r *http.Request) bool {
	if wb.formName == "" {
		return false
	}
	return factory.StateOf(r).Get("_submit") == wb.formName
}
