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

const errNotFound = `// store is your data layer; ErrNotFound is its sentinel error.
func handleOrder(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    order, err := store.Find(id)
    if errors.Is(err, ErrNotFound) {
        renderNotFound(w, r, "Order "+id)
        return
    }
    // ... render the order page ...
}

func renderNotFound(w http.ResponseWriter, r *http.Request, what string) {
    w.WriteHeader(http.StatusNotFound)
    wf.Page().Add(
        wf.AppBar("Not found"),
        wf.HeadlineLarge("404 — Not found"),
        wf.Markdown(what, " doesn't exist or has been removed."),
        wf.ButtonText("").Add("Back").WithHrefBack(),
    ).Draw(w, r)
}
`

const errCatchAll = `// In main():
mux.HandleFunc("/", handleRoot)   // matches everything not otherwise routed

func handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        renderNotFound(w, r, "The page")
        return
    }
    // ... real home page ...
}
`

const errRecover = `func recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("panic in %s %s: %v\n%s",
                    r.Method, r.URL.Path, rec, debug.Stack())
                renderServerError(w, r)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func renderServerError(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    wf.Page().Add(
        wf.AppBar("Something went wrong"),
        wf.HeadlineLarge("500 — Internal error"),
        wf.Markdown(
            "An error occurred handling your request. Try again, or ",
            "[contact support](/contact) if it keeps happening.",
        ),
    ).Draw(w, r)
}

// In main():
http.ListenAndServe(":8080", recoverPanic(mux))
`

const errFlash = `// On error inside a form handler, re-render with a Snackbar instead of
// redirecting. The URL still carries the user's input, so the form
// repopulates.

if form.ReadyToCommit(r) {
    if err := persist(values); err != nil {
        page := wf.Page().Add(
            wf.AppBar("New person"),
            wf.Snackbar().
                WithMessage("Couldn't save: " + friendlyMessage(err)).
                ShowIf(true),
            form,
        )
        page.Draw(w, r)
        return
    }
    wf.RedirectBack(w, r)
}
`

const errInsideEmbed = `// mux is your *http.ServeMux (see Handlers & routing).
// In a modal that embeds /orders/{id}/edit:
wf.Modal("modal").Add(
    wf.EmbedHandler(mux.ServeHTTP, r, "GET",
        wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
),

// The embedded handler can render its own error page — the framework just
// reads its <body> and drops it into the modal. A 404 inside the modal
// shows up as a "Not found" page rendered inside the modal frame; the
// parent page stays open underneath, and the user can dismiss the modal.

// If you want the modal to close automatically on error, redirect to the
// parent with the error in a state variable the parent reacts to:
wf.Redirect(w, r, "^?modal=&error=order-not-found")
`

// HandleErrors covers the BESPA-specific aspects of error handling.
func HandleErrors(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Errors"),

		wf.Markdown(
			"Three kinds of errors come up in a real BESPA app: a route that ",
			"matched but the resource doesn't exist (404), a handler that ",
			"crashed or returned an unexpected error (500), and a domain error ",
			"that the user can act on (\"email already exists\"). The first two ",
			"render dedicated error pages. The third is part of the form / ",
			"persist flow — see [Forms & validation](/build/forms) for that ",
			"side of the story.",
		),

		wf.HeadlineMedium("404 — not found"),
		wf.Markdown(
			"A handler that decides the resource doesn't exist writes the ",
			"status and renders a regular BESPA page:",
		),
		wf.CodeBlock(errNotFound).WithLanguage("go"),
		wf.Markdown(
			"The `WriteHeader(404)` must come **before** the page's first byte ",
			"of body is written. Calling `Draw` after `WriteHeader` is the ",
			"correct order — `WriteHeader` writes the status line; `Draw` ",
			"writes the body.",
		),
		wf.Markdown(
			"For URLs that match no registered handler at all, the stdlib mux ",
			"automatically returns a plain 404. To replace it with a friendly ",
			"BESPA page, route `/` to a catch-all handler that checks the path:",
		),
		wf.CodeBlock(errCatchAll).WithLanguage("go"),

		wf.HeadlineMedium("500 — panic recovery"),
		wf.Markdown(
			"A panic inside a handler crashes the goroutine unless something ",
			"recovers. Add a top-level `recover` middleware that logs the ",
			"stack and renders a 500 page:",
		),
		wf.CodeBlock(errRecover).WithLanguage("go"),
		wf.Markdown(
			"Log the full stack to your observability system; render only a ",
			"friendly message to the user. Never echo `recover()`'s value into ",
			"the HTML — it's untrusted, and often discloses internal paths or ",
			"data.",
		),

		wf.HeadlineMedium("Form & persist errors"),
		wf.Markdown(
			"For errors that are part of normal flow — validation, duplicate ",
			"records, transient downstream failures — *don't* redirect to an ",
			"error page. Re-render the form with a `Snackbar` (or attach an ",
			"`InfoBubble` to the offending field):",
		),
		wf.CodeBlock(errFlash).WithLanguage("go"),
		wf.Markdown(
			"Because no redirect happens, the URL still carries the submitted ",
			"values; the inputs come back populated and the user can fix the ",
			"cause and re-submit. Full pattern in ",
			"[Forms & validation → Persist errors](/build/forms).",
		),

		wf.HeadlineMedium("Errors inside an embedded page"),
		wf.Markdown(
			"When a modal or side panel embeds a handler via `EmbedHandler`, ",
			"the embedded handler's response is what gets rendered inside the ",
			"overlay. If that handler returns a 404 or 500, the user sees the ",
			"error page *inside the modal* — which is usually fine, but ",
			"sometimes you'd rather close the modal and surface the message on ",
			"the parent:",
		),
		wf.CodeBlock(errInsideEmbed).WithLanguage("go"),
		wf.Markdown(
			"`EmbedHandler` decompresses `Content-Encoding` automatically and ",
			"passes through the response body — the wire format of the ",
			"embedded error page is the same as any other BESPA response.",
		),

		wf.HeadlineMedium("Logging vs surfacing"),
		wf.Markdown(
			"Two principles worth naming:",
			"\n\n",
			"- **Log technical errors; show generic messages.** ",
			"Database connection refused, decode failure, third-party timeout — ",
			"these go to your logs (with full context). The user sees \"Try ",
			"again in a moment.\"\n",
			"- **Surface user-actionable errors inline.** Bad email format, ",
			"duplicate username, insufficient permissions — render next to the ",
			"input or as a Snackbar, in language the user can act on. Don't ",
			"redirect to an error page.",
		),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Forms & validation → Persist errors](/build/forms) — the ",
			"re-render-with-Snackbar idiom for domain errors.",
			"\n\n",
			"[Handlers & routing](/build/handlers-and-routing) — where the ",
			"recover middleware sits in the middleware stack.",
			"\n\n",
			"[Modals & side panels](/build/modals) — what an embedded handler's ",
			"errors look like to the user.",
		),
	)
	shared.Render(w, r, page)
}
