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

package build

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

const formsBasic = `form := wf.Form().Add(
    wf.Field().AddLeft("Name").AddRight(
        wf.InputText("name", "").
            WithRequired(true).
            WithLength(2, 40),
    ),
    wf.Field().AddLeft("Email").AddRight(
        wf.InputEmail("email", "").WithRequired(true),
    ),
    wf.ButtonFilled("save").Add("Save"),
)
`

const formsAutoSubmit = `// Posts the form on every change — useful for live search/filter inputs.
wf.InputText("query", "").
    WithAutoSubmit(true).
    RedrawIfChanged(r, "query"),
`

const formsPredicate = `wf.InputText("password", "").
    WithRequired(true).
    WithLength(8, 64).
    WithPredicate(func(v string) (bool, string) {
        if !strings.ContainsAny(v, "0123456789") {
            return false, "Must contain at least one digit"
        }
        return true, ""
    }),
`

const formsCommit = `// directory and storage.Person stand in for your app's data layer.
if form.ReadyToCommit(r) {
    values := form.Values(r)
    person := &storage.Person{
        FirstName: values.Get("fname"),
        Email:     values.Get("email"),
        // ...
    }
    if err := directory.Insert(person); err != nil {
        // Show error in a snackbar / redraw the relevant field — see below.
    }
}
`

const formsRedirect = `if form.ReadyToCommit(r) {
    if err := directory.Insert(person); err != nil {
        // Fall through to re-render with an inline error.
    } else {
        // Success — return the user to where they came from.
        wf.RedirectBack(w, r)
        return
    }
}
`

const formsModalClose = `// orders is your data layer.
func handleNewOrder(w http.ResponseWriter, r *http.Request) {
    form := wf.Form().Add(/* ... fields ... */)
    if form.ReadyToCommit(r) {
        // Persist, then close the modal by clearing the parent's "modal" key.
        if err := orders.Create(form.Values(r)); err == nil {
            wf.Redirect(w, r, "^?modal=")
            return
        }
    }
    wf.Page().Add(wf.AppBar("New order"), form).Draw(w, r)
}
`

const formsButtons = `wf.Form().Add(
    /* ... fields ... */,
    wf.Toolbar().AddLeft(
        wf.ButtonText("cancel").Add("Cancel"),
        wf.ButtonFilled("save").Add("Save"),
        wf.ButtonFilled("save_and_new").Add("Save and add another"),
    ),
)

// In the handler:
switch wf.StateOf(r).Get("_submit") {
case "save":
    persist(); wf.RedirectBack(w, r); return
case "save_and_new":
    persist(); wf.Redirect(w, r, "?"); return  // reload current page, cleared
case "cancel":
    wf.RedirectBack(w, r); return
}
`

const formsPersistError = `if form.ReadyToCommit(r) {
    if err := directory.Insert(person); err != nil {
        // Re-render the form. A Snackbar widget in the page tree, opted in to
        // a "flash" state variable, surfaces the message without losing input.
        return // falls through to Page().Draw below
    }
    wf.RedirectBack(w, r)
    return
}
page := wf.Page().Add(
    wf.AppBar("New person"),
    wf.Snackbar().
        WithMessage("Couldn't save: email already exists").
        ShowIf(persistErr != nil),
    form,
)
`

// HandleForms covers forms and validation.
func HandleForms(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Forms & validation"),

		wf.Markdown(
			"BESPA forms are server-rendered HTML forms plus a handful of ",
			"widgets that know how to validate themselves on the client (for ",
			"fast feedback) and the server (for trust). The same widget runs ",
			"both checks.",
		),
		wf.HeadlineMedium("The shape"),
		wf.Markdown(
			"A `Form` widget wraps any number of input widgets. Inputs go ",
			"inside `Field` rows so the label/input pairing is visually ",
			"consistent:",
		),
		wf.CodeBlock(formsBasic).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Submitting this form posts `name` and `email` back into the ",
			"page's state. The `save` button's name becomes `_submit=save` so ",
			"the handler can distinguish multiple buttons.",
		),
		wf.HeadlineMedium("Live updates"),
		wf.Markdown(
			"Auto-submit turns a form into a live search box. Each keystroke ",
			"posts the form (debounced), so the server sees the current value ",
			"and any `RedrawIfChanged` widgets reflect it:",
		),
		wf.CodeBlock(formsAutoSubmit).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Do not put `RedrawIfChanged` on the input itself — see ",
			"[When to redraw](/build/redraw).",
		),
		wf.HeadlineMedium("Built-in validators"),
		wf.Markdown(
			"Every input widget chains a standard set of checks:",
			"\n\n",
			"- `WithRequired(true)` — non-empty.\n",
			"- `WithLength(min, max)` — rune count bounds. `-1` means unbounded.\n",
			"- `WithPattern(regexp)` — regex match on the value.\n",
			"- `WithMin(...)` / `WithMax(...)` — numeric or date bounds (clamps ",
			"the value on read too).\n\n",
			"Format-specific widgets bake their own validator in: `InputEmail` ",
			"checks email shape, `InputURL` checks URL parse, `InputInteger` ",
			"checks numeric parse, etc.",
		),
		wf.HeadlineMedium("Custom predicates"),
		wf.Markdown(
			"For anything the built-ins can't express, attach a predicate. The ",
			"predicate runs on every validation pass; returning `(false, msg)` ",
			"surfaces the message inline next to the input:",
		),
		wf.CodeBlock(formsPredicate).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Predicates are pure Go — use them for cross-field rules, lookups ",
			"against a database, or anything else the standard validators can't do.",
		),
		wf.HeadlineMedium("Committing"),
		wf.Markdown(
			"`form.ReadyToCommit(r)` returns true exactly once per request, when ",
			"the form was submitted AND all validators pass. That's the moment ",
			"to persist:",
		),
		wf.CodeBlock(formsCommit).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Re-run validation on the server, even if the same checks already ",
			"passed on the client — a hostile client can skip the client checks. ",
			"The widget validators are the single source of truth for whether ",
			"the data is acceptable.",
		),
		wf.HeadlineMedium("After a successful commit"),
		wf.Markdown(
			"On success, the handler should *redirect*, not re-render. ",
			"Redirecting clears the submitted state from the URL, prevents a ",
			"browser refresh from re-submitting, and lets the user land on a ",
			"clean page. The framework provides two helpers:",
			"\n\n",
			"- `wf.RedirectBack(w, r)` — sends the user to whatever URL is in ",
			"the `_back` state variable. Returns `false` if there is no ",
			"`_back`, so you can fall through to a default destination.\n",
			"- `wf.Redirect(w, r, location)` — explicit destination. The ",
			"location string respects the [action-URL grammar](/basics/action-url-pattern) ",
			"— `^?modal=` to close a parent modal, `~/somewhere` to navigate ",
			"the top page, or a plain path to navigate the current frame.",
		),
		wf.CodeBlock(formsRedirect).WithLanguage("go"),
		wf.HeadlineMedium("A modal form that posts and closes itself"),
		wf.Markdown(
			"Combine the modal embed pattern with `^?modal=` to get \"open a ",
			"form in a modal, save it, modal disappears, parent page reflects ",
			"the change\":",
		),
		wf.CodeBlock(formsModalClose).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The parent page renders the modal slot with ",
			"`wf.EmbedHandler(mux.ServeHTTP, r, ..., \"?_back=^?modal=\", nil)`. ",
			"When the handler redirects to `^?modal=`, the parent's `modal` ",
			"state variable is cleared, the modal closes, and any widget on ",
			"the parent that opted in to a redraw on relevant state ",
			"(e.g. the orders table) re-renders with the new row.",
		),
		wf.HeadlineMedium("Multiple buttons"),
		wf.Markdown(
			"Every submit button writes its name to the `_submit` state ",
			"variable. Inspect it to branch:",
		),
		wf.CodeBlock(formsButtons).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"`form.ReadyToCommit(r)` is only true when the *named* form's ",
			"submit fired and validators passed — it ignores submits from ",
			"other forms on the same page, so you don't need to guard against ",
			"cross-form events.",
		),
		wf.HeadlineMedium("Persist errors"),
		wf.Markdown(
			"When validation passes but persistence fails — duplicate key, ",
			"external service down, conflict on update — surface the error and ",
			"keep the user's input. The idiom is to *not* redirect on error; ",
			"re-render the form with a `Snackbar` (or an `InfoBubble` attached ",
			"to a specific field):",
		),
		wf.CodeBlock(formsPersistError).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Because no redirect happens, the URL still carries the submitted ",
			"values, so the inputs come back populated. The user fixes the ",
			"cause and resubmits.",
		),
		wf.HeadlineMedium("File uploads"),
		wf.Markdown(
			"`wf.InputFile(name)` posts a multipart form. Read the file in the ",
			"handler with `r.FormFile(name)` — standard `net/http` — after ",
			"`form.ReadyToCommit(r)`. Validate size and content type on the ",
			"server: client-side `accept=\"image/*\"` is a hint, not a ",
			"guarantee. See [Showcase → Form input](/showcase/form-input) for ",
			"a working file picker.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Showcase → Form input widgets](/showcase/form-input) ",
			"— every input type live.",
			"\n\n",
			"[Showcase → Form validation](/showcase/form-validation) ",
			"— validators and predicates in action.",
			"\n\n",
			"[Showcase → CRUD](/showcase/dir-list) ",
			"— a full register/edit/delete flow.",
		),
	)
	shared.Render(w, r, page)
}
