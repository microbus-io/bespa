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

import (
	"io"
	"net/http"
)

// WidgetBase facilitates the implementation of the Widget interface.
type WidgetBase[T Widget] struct {
	id      string
	redrawn bool
	shown   bool
	owner   T
}

// NewWidgetBase creates a new widget base holding a pointer to its owner subclass.
func NewWidgetBase[T Widget](owner T) *WidgetBase[T] {
	return &WidgetBase[T]{
		owner: owner,
		shown: true,
	}
}

// ID is an identifier of this widget that is unique in the scope of its page.
func (wb *WidgetBase[T]) ID() string {
	return wb.id
}

// SetID sets an identifier of this widget that is unique in the scope of its page.
func (wb *WidgetBase[T]) SetID(id string) {
	wb.id = id
}

// WithID sets an identifier of this widget that is unique in the scope of its page.
func (wb *WidgetBase[T]) WithID(id string) T {
	wb.id = id
	return wb.owner
}

// Children are the widgets nested under this widget.
// The default implementation returns no children.
func (wb *WidgetBase[T]) Children() []Widget {
	return nil
}

// Draw renders the widget's HTML.
// The default implementation renders an empty placeholder span.
// Widgets should override this method to render their HTML.
func (wb *WidgetBase[T]) Draw(w io.Writer, r *http.Request) (err error) {
	return factory.Tag("span").Attr("data-id", wb.ID()).When(false).Draw(w, r)
}

/*
RedrawIf forces the redrawing of a widget if the condition is satisfied.
Consecutive calls are OR-ed together. It is enough that any condition is true for the widget to be redrawn.
By default the widget is drawn only on the initial page drawing.
Redrawing a widget may either show or hide it, depending on the status of its shown flag.
*/
func (wb *WidgetBase[T]) RedrawIf(condition bool) T {
	wb.redrawn = wb.redrawn || condition
	return wb.owner
}

// RedrawIfChanged is equivalent to multiple RedrawIf(StateOf(r).Changed(stateVarName)).
func (wb *WidgetBase[T]) RedrawIfChanged(r *http.Request, stateVarName ...string) T {
	state := factory.StateOf(r)
	changed := false
	for _, v := range stateVarName {
		changed = changed || state.Changed(v)
	}
	return wb.RedrawIf(changed)
}

// RedrawIfChangedTo is equivalent to RedrawIf(StateOf(r).Changed(stateVarName) && StateOf(r).Get(stateVarName)==value).
func (wb *WidgetBase[T]) RedrawIfChangedTo(r *http.Request, stateVarName string, value string) T {
	state := factory.StateOf(r)
	changedTo := state.Changed(stateVarName) && state.Get(stateVarName) == value
	return wb.RedrawIf(changedTo)
}

// RedrawIfEq is equivalent to RedrawIf(StateOf(r).Get(stateVarName)==value).
func (wb *WidgetBase[T]) RedrawIfEq(r *http.Request, stateVarName string, value string) T {
	state := factory.StateOf(r)
	return wb.RedrawIf(state.Get(stateVarName) == value)
}

// RedrawIfNotEq is equivalent to RedrawIf(StateOf(r).Get(stateVarName)!=value).
func (wb *WidgetBase[T]) RedrawIfNotEq(r *http.Request, stateVarName string, value string) T {
	state := factory.StateOf(r)
	return wb.RedrawIf(state.Get(stateVarName) != value)
}

// RedrawIfEmpty is equivalent to RedrawIf(StateOf(r).Get(stateVarName)=="").
func (wb *WidgetBase[T]) RedrawIfEmpty(r *http.Request, stateVarName string) T {
	state := factory.StateOf(r)
	return wb.RedrawIf(state.Get(stateVarName) == "")
}

// Drawn indicates whether this widget needs to be drawn.
func (wb *WidgetBase[T]) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wb.redrawn || !state.HasChanges()
}

// HideIf is used to indicate if the widget needs to be shown or not.
// Consecutive calls are OR-ed together. It is enough that any condition is true for the widget to be hidden.
// By default the widget is shown.
func (wb *WidgetBase[T]) HideIf(condition bool) T {
	wb.shown = wb.shown && !condition
	return wb.owner
}

// HideIfEmpty is equivalent to HideIf(StateOf(r).Get(stateVarName)=="").
func (wb *WidgetBase[T]) HideIfEmpty(r *http.Request, stateVarName string) T {
	state := factory.StateOf(r)
	return wb.HideIf(state.Get(stateVarName) == "")
}

// HideIfNotEq is equivalent to Hide(StateOf(r).Get(stateVarName)!=value).
func (wb *WidgetBase[T]) HideIfNotEq(r *http.Request, stateVarName string, value string) T {
	state := factory.StateOf(r)
	return wb.HideIf(state.Get(stateVarName) != value)
}

// HideIfEq is equivalent to Hide(StateOf(r).Get(stateVarName)==value).
func (wb *WidgetBase[T]) HideIfEq(r *http.Request, stateVarName string, value string) T {
	state := factory.StateOf(r)
	return wb.HideIf(state.Get(stateVarName) == value)
}

// Shown returns whether the widget is shown or hidden.
// A widget that is hidden is not rendered.
func (wb *WidgetBase[T]) Shown(r *http.Request) bool {
	return wb.shown
}
