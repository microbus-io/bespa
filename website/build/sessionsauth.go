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

const saMiddleware = `// session.go in your app — wraps a session library you chose.
var sessions = scs.New() // github.com/alexedwards/scs/v2

// In main():
mux := http.NewServeMux()
mux.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)
mux.HandleFunc("/", handleHome)
http.ListenAndServe(":8080", sessions.LoadAndSave(mux))
`

const saReadInHandler = `// Inside any handler:
func handleProfile(w http.ResponseWriter, r *http.Request) {
    userID := sessions.GetInt(r.Context(), "user_id")
    if userID == 0 {
        // Send the user to login, remember where they were.
        wf.Redirect(w, r, "/login?_back=" + url.QueryEscape(r.URL.String()))
        return
    }
    // ... render the page with userID in scope ...
}
`

const saGate = `// Wrap a single handler:
mux.Handle("/admin/", requireLogin(adminMux))

// Or wrap a whole sub-mux:
adminMux := http.NewServeMux()
adminMux.HandleFunc("/admin/users", handleUsers)
adminMux.HandleFunc("/admin/orders", handleOrders)
mux.Handle("/admin/", http.StripPrefix("/admin", requireLogin(adminMux)))

func requireLogin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if sessions.GetInt(r.Context(), "user_id") == 0 {
            wf.Redirect(w, r, "/login?_back=" + url.QueryEscape(r.URL.String()))
            return
        }
        next.ServeHTTP(w, r)
    })
}
`

const saLogin = `func handleLogin(w http.ResponseWriter, r *http.Request) {
    form := wf.Form().Add(
        wf.Field().AddLeft("Email").AddRight(
            wf.InputEmail("email", "").WithRequired(true),
        ),
        wf.Field().AddLeft("Password").AddRight(
            wf.InputPassword("password", "").WithRequired(true),
        ),
        wf.ButtonFilled("login").Add("Sign in"),
    )
    if form.ReadyToCommit(r) {
        user, err := authenticate(form.Values(r).Get("email"), form.Values(r).Get("password"))
        if err == nil {
            sessions.Put(r.Context(), "user_id", user.ID)
            // RedirectBack returns to _back, or to / if unset.
            if !wf.RedirectBack(w, r) {
                wf.Redirect(w, r, "/")
            }
            return
        }
        // Fall through and re-render with a Snackbar holding the message.
    }
    wf.Page().Add(wf.AppBar("Sign in"), form).Draw(w, r)
}
`

const saCSRF = `wf.Form().Add(
    wf.InputHidden("_csrf", nosurf.Token(r)),
    // ... fields ...
)

// On the receiving side, your CSRF middleware reads the same token and
// rejects mismatches before the handler runs.
`

const saPrefs = `func Render(w http.ResponseWriter, r *http.Request, page *widget.PageWidget) {
    session := SessionOf(w, r)            // your helper
    switch session.Theme {
    case "Dark":
        page.WithThemeDark()
    case "Light":
        page.WithThemeLight()
    }
    if session.Palette != "" {
        for _, kc := range css.PresetKeyColors {
            if kc.Name == session.Palette {
                page.WithKeyColors(kc)
                break
            }
        }
    }
    page.Draw(w, r)
}
`

// HandleSessionsAuth covers integrating a session library + auth gates into
// BESPA handlers. BESPA itself has no session API — this page is the
// integration pattern.
func HandleSessionsAuth(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Sessions & auth"),

		wf.Markdown(
			"BESPA stays out of sessions and authentication — neither lives in ",
			"the framework. Bring whichever session library you like ",
			"([scs](https://github.com/alexedwards/scs), ",
			"[gorilla/sessions](https://github.com/gorilla/sessions), ",
			"or your own), wrap your mux with its middleware, and read the ",
			"session inside your handlers. This page is the integration pattern.",
		),
		wf.HeadlineMedium("Wire the middleware"),
		wf.Markdown(
			"Standard `net/http` middleware. The session library handles cookies ",
			"and storage; you treat it as an opaque blob on the request context:",
		),
		wf.CodeBlock(saMiddleware).WithLanguage("go"),
		wf.Markdown(
			"Middleware order is up to you. Compression, logging, request ID — ",
			"any of these can wrap or be wrapped by the session middleware; ",
			"BESPA is unaware. See [Handlers & routing](/build/handlers-and-routing) ",
			"for the only middleware-order rule the framework imposes.",
		),

		wf.HeadlineMedium("Read in a handler"),
		wf.Markdown(
			"Anywhere inside a handler, read whatever your library exposes off ",
			"`r.Context()`. If the user isn't logged in, redirect to the login ",
			"page with `_back` so they return to where they tried to go:",
		),
		wf.CodeBlock(saReadInHandler).WithLanguage("go"),

		wf.HeadlineMedium("Gate routes with middleware"),
		wf.Markdown(
			"Standard `net/http` middleware composes. Gate one handler or wrap a ",
			"whole sub-tree:",
		),
		wf.CodeBlock(saGate).WithLanguage("go"),

		wf.HeadlineMedium("Login form"),
		wf.Markdown(
			"A BESPA login flow is just a form whose handler reads credentials ",
			"and stores user state in the session. The `RedirectBack` helper ",
			"returns the user to wherever they came from:",
		),
		wf.CodeBlock(saLogin).WithLanguage("go"),
		wf.Markdown(
			"Open it in a modal (`?modal=/login`) for a soft prompt that ",
			"doesn't lose the user's context, or hit it as a full page for ",
			"first-time sign-in. Both work with the same handler.",
		),

		wf.HeadlineMedium("CSRF"),
		wf.Markdown(
			"Most session libraries include a CSRF helper (`nosurf`, ",
			"`gorilla/csrf`, scs's `Token`). Render the token as a hidden field ",
			"in every state-changing form:",
		),
		wf.CodeBlock(saCSRF).WithLanguage("go"),
		wf.Markdown(
			"BESPA forms post all hidden inputs along with visible ones, so the ",
			"token rides back automatically.",
		),

		wf.HeadlineMedium("OAuth and external redirects"),
		wf.Markdown(
			"The `_back` state variable survives an external OAuth round-trip ",
			"if you embed it in the OAuth `state` parameter and restore it on ",
			"the callback. On success, set the session and `RedirectBack`:",
			"\n\n",
			"- Originating page: store `r.URL` in the OAuth `state` (after ",
			"URL-encoding) along with a random nonce for CSRF.\n",
			"- Callback handler: validate the nonce, exchange the code, set ",
			"the session, parse `_back` out of the state, redirect.",
		),

		wf.HeadlineMedium("Per-user preferences"),
		wf.Markdown(
			"Anything that's *per-user* and should survive navigation — theme, ",
			"locale, layout preferences — belongs in the session, not the URL. ",
			"The website's own theme/palette switcher uses this pattern; its ",
			"`Render` wrapper reads the session and applies the user's choices ",
			"before drawing every page:",
		),
		wf.CodeBlock(saPrefs).WithLanguage("go"),
		wf.Markdown(
			"In contrast, anything that's *per-page* and should be ",
			"bookmarkable — search query, sort order, the currently open ",
			"modal — belongs in the URL. See [State patterns](/build/state) ",
			"for the URL-vs-session rule.",
		),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Handlers & routing](/build/handlers-and-routing) — where the ",
			"session middleware sits relative to the rest.",
			"\n\n",
			"[Forms & validation](/build/forms) — the form / RedirectBack / ",
			"persist-error pattern used in the login example.",
			"\n\n",
			"`website/shared/session.go` in this repo — a small concrete session ",
			"helper used by the website example.",
		),
	)
	shared.Render(w, r, page)
}
