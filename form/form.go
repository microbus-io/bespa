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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/microbus-io/bespa/basic"
	"github.com/microbus-io/bespa/widget"
)

// FormWidget renders an input form.
type FormWidget struct {
	*widget.WidgetBase[*FormWidget]
	method       string
	action       string
	target       string
	children     []Widget
	autoComplete string
	width        string
	reset        bool
	name         string
}

// Form creates a new widget that renders an input form. Defaults: POST to
// the current URL, autocomplete off, max-width 800px. Submission writes
// `_submit=<form name>` to state so the handler can detect it via
// Submitted or ReadyToCommit; nested input widgets attach to this form
// automatically when added via Form.Add.
func (f FormFactory) Form() *FormWidget {
	x := &FormWidget{
		name:         "form",
		method:       "POST",
		autoComplete: "off",
		action:       "?",
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithName sets the form's name. The name is written into `_submit` on
// submission, which is how Submitted/ReadyToCommit discriminate between
// multiple forms on the same page. Default is "form".
func (wgt *FormWidget) WithName(name string) *FormWidget {
	wgt.name = name
	return wgt
}

// Add appends children to the form. Buttons added at any position are
// automatically collected into a trailing toolbar (right-aligned), so you
// don't have to wrap them yourself. Input widgets in the subtree are
// linked to this form by name so Values/Fields/Submitted can find them.
func (wgt *FormWidget) Add(children ...any) *FormWidget {
	var buttonStrip *basic.ToolbarWidget
	adding := Many(children)
	for _, a := range adding {
		if _, ok := a.(*ButtonWidget); ok {
			if buttonStrip == nil {
				buttonStrip = factory.Toolbar()
				wgt.children = append(wgt.children, buttonStrip)
			}
			buttonStrip.AddRight(a)
		} else {
			buttonStrip = nil
			wgt.children = append(wgt.children, a)
		}
	}

	// Set the parent form name of any nested input widgets
	var f func(w Widget)
	f = func(w Widget) {
		if iw, ok := w.(InputWidget); ok {
			iw.SetFormName(wgt.name)
		}
		for _, c := range w.Children() {
			f(c)
		}
	}
	for _, w := range adding {
		f(w)
	}

	return wgt
}

// WithAutoComplete sets the form's HTML autocomplete attribute. Default
// is "off"; common alternative is "on". Use "" to let the browser inherit
// from the page. Individual text inputs can override via their own
// WithAutoComplete. See the HTML autocomplete spec for the full list of
// tokens.
func (wgt *FormWidget) WithAutoComplete(autoComplete string) *FormWidget {
	wgt.autoComplete = autoComplete
	return wgt
}

// WithWidth caps the form's max-width. Default is 800px; the form still
// compresses below 600px regardless. Pass 0 to clear the cap.
// Allowed CSS units are "px", "%", "ch", "em", "vw", "vh", etc.
func (wgt *FormWidget) WithWidth(width float32, unit string) *FormWidget {
	if width > 0 {
		wgt.width = fmt.Sprintf("max-width:%f%s", width, unit)
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithAction sets the submit method and URL. Default is POST to "?",
// which writes the form's fields into the current page's state — the
// idiomatic choice. Use GET for query/search forms; POST for forms that
// persist data. Only "GET" and "POST" are accepted; other methods are
// silently ignored.
func (wgt *FormWidget) WithAction(method string, action string) *FormWidget {
	method = strings.ToUpper(method)
	if method == "GET" || method == "POST" {
		wgt.method = method
	}
	wgt.action = action
	return wgt
}

// WithTarget sets the HTML target for the form's submission. Defaults to
// the page's `_target` state variable when unset, so submissions route
// back to the active frame.
func (wgt *FormWidget) WithTarget(target string) *FormWidget {
	wgt.target = target
	return wgt
}

// ReadyToCommit returns true when the user just submitted this form and
// every field passes validation — the typical guard around persisting
// values:
//
//	if form.ReadyToCommit(r) { … persist … }
func (wgt *FormWidget) ReadyToCommit(r *http.Request) bool {
	return wgt.Submitted(r) && wgt.Valid(r)
}

// Valid runs every visible field's validator and returns true if all
// pass. Returns true for unsubmitted forms (nothing to validate yet).
// Validators run concurrently — safe for expensive custom predicates.
func (wgt *FormWidget) Valid(r *http.Request) bool {
	if wgt.reset || !wgt.Submitted(r) {
		return true
	}
	// Perform validations in parallel in case there are custom validator that take time
	fields := wgt.Fields(r)
	ch := make(chan bool, len(fields))
	for i := range fields {
		field := fields[i]
		go func() {
			ch <- field.Valid(r)
		}()
	}
	result := true
	for range fields {
		if !<-ch {
			result = false
		}
	}
	return result
}

// Changed reports whether any field's posted value differs from its
// initial value. Returns false for unsubmitted or reset forms.
func (wgt *FormWidget) Changed(r *http.Request) bool {
	if wgt.reset || !wgt.Submitted(r) {
		return false
	}
	for _, field := range wgt.Fields(r) {
		if field.Changed(r) {
			return true
		}
	}
	return false
}

// Reset clears every field's posted value so the form re-renders with its
// initial values. Call this after a successful save when you want to show
// a fresh empty form rather than the user's just-committed input.
func (wgt *FormWidget) Reset(r *http.Request) {
	state := factory.StateOf(r)
	if state.Get("_submit") == wgt.name {
		state.Del("_submit")
	}
	wgt.reset = true
	for _, field := range wgt.Fields(r) {
		state.Del(field.Name())
	}
}

// Submitted returns true when the current request is a submission of
// this specific form (matched by HTTP method and the `_submit` state
// variable). Prefer ReadyToCommit, which also checks validation.
func (wgt *FormWidget) Submitted(r *http.Request) bool {
	return r.Method == wgt.method && factory.StateOf(r).Get("_submit") == wgt.name
}

// Values collects every nested input's current value, keyed by field
// name. Hidden inputs and inputs in hidden subtrees are included; use
// Fields if you need to inspect the InputWidgets themselves.
func (wgt *FormWidget) Values(r *http.Request) Values {
	values := Values(map[string]string{})
	return wgt.traverseValues(wgt, r, values)
}

func (wgt *FormWidget) traverseValues(w Widget, r *http.Request, collector Values) Values {
	if field, ok := w.(InputWidget); ok && field.Name() != "" {
		collector[field.Name()] = field.Value(r)
	}
	for _, child := range w.Children() {
		wgt.traverseValues(child, r, collector)
	}
	return collector
}

// Fields returns the visible InputWidgets currently nested in the form,
// keyed by field name. Hidden subtrees (HideIf*) are excluded — useful
// when validation should ignore fields the user can't see.
func (wgt *FormWidget) Fields(r *http.Request) map[string]InputWidget {
	fields := Fields(map[string]InputWidget{})
	return wgt.traverseFields(r, wgt, fields)
}

func (wgt *FormWidget) traverseFields(r *http.Request, w Widget, collector Fields) Fields {
	if !w.Shown(r) {
		return collector
	}
	if field, ok := w.(InputWidget); ok && field.Name() != "" {
		collector[field.Name()] = field
	}
	for _, child := range w.Children() {
		wgt.traverseFields(r, child, collector)
	}
	return collector
}

// Children are the widgets nested under this widget.
func (wgt *FormWidget) Children() []Widget {
	return wgt.children
}

// Draw renders the widget's HTML.
func (wgt *FormWidget) Draw(w io.Writer, r *http.Request) (err error) {
	// Set the default action to post to itself. Use the relative path.
	action := wgt.action
	if action == "" {
		action = r.URL.Path
		p := strings.LastIndex(action, "/")
		if p >= 0 {
			action = action[p+1:]
		}
	}
	var queryArgs url.Values
	if wgt.method == "GET" {
		// Move the query arguments to be hidden inputs because they cannot be read
		// from the action URL when the form is submitted with the GET method.
		// But do not include query arguments that are overwritten by other fields.
		u, err := url.Parse(action)
		if err != nil {
			return err
		}
		queryArgs = u.Query()
		if len(queryArgs) > 0 {
			// Remove the query arguments from the original action URL
			u.RawQuery = ""
			action = u.String()
		}
		for fieldName := range wgt.Fields(r) {
			queryArgs.Del(fieldName)
		}
	}
	state := factory.StateOf(r)
	target := wgt.target
	if target == "" {
		target = state.Get("_target")
	}
	formTag := Tag("form").
		Attr("method", wgt.method).
		Attr("enctype", "application/x-www-form-urlencoded").
		Attr("action", action).
		Attr("autocomplete", wgt.autoComplete).
		Attr("target", target).
		Add(factory.InputHidden("_submit", wgt.name)).
		Add(Tag("div").Add(wgt.children))
	if wgt.method == "GET" {
		for k, v := range queryArgs {
			formTag.Add(factory.InputHidden(k, v[0]))
		}
	}
	return Tag("div").
		Class("InputForm", "Block").
		Attr("data-id", wgt.ID()).
		Style(wgt.width).
		Add(formTag).
		When(wgt.Shown(r) && len(wgt.children) > 0).
		Draw(w, r)
}

// Drawn indicates whether this widget needs to be drawn.
func (wgt *FormWidget) Drawn(r *http.Request) bool {
	state := factory.StateOf(r)
	return wgt.WidgetBase.Drawn(r) ||
		wgt.reset ||
		state.Get("_submit") == wgt.name
}
