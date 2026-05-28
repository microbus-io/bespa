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

const deploymentCSPCurrent = `Content-Security-Policy:
  default-src 'self';
  script-src  'self' 'unsafe-inline';
  style-src   'self' 'unsafe-inline';
  font-src    'self' data:;
  img-src     'self' data:;
`

const deploymentEmbed = `// Your app's main package can embed everything it needs to serve.
//go:embed assets/*.css assets/*.js assets/fonts/*
var appAssets embed.FS

func main() {
    widget.AssetRegistry.RegisterFS(appAssets)

    mux := http.NewServeMux()
    mux.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP)
    mux.HandleFunc("/", home)
    http.ListenAndServe(":8080", nil)
}
`

const deploymentProxyHeaders = `# nginx
location / {
    proxy_pass         http://app:8080;
    proxy_set_header   Host              $host;
    proxy_set_header   X-Forwarded-Host  $host;
    proxy_set_header   X-Forwarded-Proto $scheme;
    # Optional — set if BESPA is mounted under a path prefix:
    # proxy_set_header X-Forwarded-Prefix /myapp;
}
`

// HandleDeployment covers what changes when a BESPA app moves out of local
// dev: single-binary packaging, CSP posture, automatic cache busting,
// reverse-proxy headers, and compression.
func HandleDeployment(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Deployment"),

		wf.Markdown(
			"A BESPA app deploys as a single Go binary. No companion ",
			"frontend service, no asset pipeline that has to run before ",
			"shipping, no build artifacts that live outside the binary. ",
			"`go build`, copy the output to a server, run it. The five things ",
			"on this page are what you should still think about once you're ",
			"past that.",
		),

		wf.HeadlineMedium("The single-binary shape"),
		wf.Markdown(
			"The framework's own CSS and JS are `//go:embed`-ed inside the ",
			"`widget` package and registered with `AssetRegistry` on import. ",
			"You don't need to ship them separately. The same mechanism is ",
			"how you ship your own assets:",
		),
		wf.CodeBlock(deploymentEmbed).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The resulting binary is fully self-contained. The only files it ",
			"reads from disk at runtime are whatever your handlers ask for — ",
			"a database file, a config file, uploaded user content. The HTML ",
			"and asset surface served to the browser all comes from inside ",
			"the binary.",
		),

		wf.HeadlineMedium("CSP"),
		wf.Markdown(
			"The honest story: **BESPA currently requires `'unsafe-inline'` ",
			"for both `script-src` and `style-src`.** A stricter policy will ",
			"silently break widget initialization and theming.",
		),
		wf.CodeBlock(deploymentCSPCurrent).WithLanguage("plaintext"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Why: widgets bootstrap themselves with small per-instance inline ",
			"scripts like `<script>snackbar_init('abc123')</script>` that ",
			"bind an externally-loaded function to a specific DOM node by ID. ",
			"Theme tokens (the resolved Material color palette for the ",
			"current scheme) ship as an inline `<style>` block in the page ",
			"head. Both patterns are blocked by `script-src 'self'` / ",
			"`style-src 'self'`.",
		),
		wf.HeadlineMedium("CSP — workarounds, in order"),
		wf.Markdown(
			"Pick based on your threat model:",
			"\n\n",
			"1. **`'unsafe-inline'`** (above). What BESPA needs today. The ",
			"cost: a stored-XSS bug becomes executable JavaScript again — ",
			"you fall back to depending on output escaping being perfect.\n",
			"2. **Per-response nonces.** Generate a random nonce in ",
			"middleware, put it in the CSP header (`script-src 'self' ",
			"'nonce-xyz'`), and add `nonce=\"xyz\"` to every inline ",
			"`<script>` and `<style>`. The right answer in principle — but ",
			"BESPA does not currently plumb nonces through `Tag(\"script\")`, ",
			"so you can't enable this without framework changes. This is a ",
			"known gap.\n",
			"3. **`'sha256-...'` hashes.** Not practical here. The per-",
			"instance scripts contain a random ID per render, so the hash ",
			"changes every response. You'd have to recompute and re-list ",
			"every hash per response — at which point nonces are cheaper.\n",
			"4. **Eliminate inline scripts.** A future framework refactor ",
			"could move per-instance init from inline `<script>` to ",
			"`data-init=\"snackbar\"` attributes discovered by ",
			"`MutationObserver`. That's a structural change, not a deploy-",
			"time knob.",
			"\n\n",
			"In the meantime, treat the framework as requiring ",
			"`'unsafe-inline'` for scripts and styles. Keep the rest of your ",
			"CSP strict (`default-src 'self'`, no `'unsafe-eval'`, no ",
			"wildcard hosts) so the loosened script-src is the only ",
			"weakening.",
		),

		wf.HeadlineMedium("Cache busting"),
		wf.Markdown(
			"Asset URLs include an automatic cache-buster: BESPA picks a ",
			"random 8-character ID when the process starts, and appends it ",
			"to every framework asset URL as `?id=XYZ`. The asset handler ",
			"sets `Cache-Control: max-age=2592000` (30 days), so browsers ",
			"hold onto framework CSS and JS aggressively.",
			"\n\n",
			"When you deploy a new binary, the new process generates a new ",
			"ID, the page's `<link>` and `<script>` tags reference the new ",
			"URL, and browsers fetch fresh. No manual hash-naming, no ",
			"asset-manifest file, no \"remember to bump the version.\"",
			"\n\n",
			"There's no story for assets you serve directly (your own ",
			"images, fonts you registered, dynamic-asset handlers). ",
			"`AssetRegistry` sets the 30-day cache header on those too, so ",
			"if their content can change without the binary changing, you ",
			"need your own cache-busting — append a hash to the URL, or ",
			"shorten the max-age via your own handler.",
		),

		wf.HeadlineMedium("Behind a reverse proxy"),
		wf.Markdown(
			"BESPA reads three `X-Forwarded-*` headers when building URLs ",
			"that point back at itself (sourcemap comments in served CSS / ",
			"JS, anything that needs a fully-qualified self-reference):",
			"\n\n",
			"- `X-Forwarded-Host` — overrides `r.Host`.\n",
			"- `X-Forwarded-Proto` — `http` or `https`. Falls back to ",
			"`r.TLS != nil` if absent.\n",
			"- `X-Forwarded-Prefix` — path prefix the app is mounted under. ",
			"Set this if your proxy strips a prefix (e.g. `/myapp/`) before ",
			"passing the request through.",
			"\n\n",
			"Set them at the proxy. Example for nginx:",
		),
		wf.CodeBlock(deploymentProxyHeaders).WithLanguage("plaintext"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Caveat: BESPA doesn't have a notion of \"trusted proxy ",
			"addresses\" — if these headers can reach your app from an ",
			"untrusted source, the values are taken at face value. Either ",
			"terminate at a proxy you control, or strip the headers at the ",
			"edge before they reach the app.",
		),

		wf.HeadlineMedium("Compression"),
		wf.Markdown(
			"BESPA does not compress its own responses — wrap the mux with ",
			"the gzip/brotli middleware of your choice, or terminate ",
			"compression at the proxy. Both work.",
			"\n\n",
			"One subtlety for `EmbedHandler`: if the inner handler's ",
			"response is compressed (because compression middleware wrapped ",
			"the mux you're embedding from), `EmbedHandler` notices the ",
			"`Content-Encoding` header and decompresses transparently. You ",
			"don't have to pass it an inner unwrapped mux to avoid ",
			"double-encoding — see ",
			"[Handlers & routing → Middleware](/build/handlers-and-routing).",
		),

		wf.HeadlineMedium("What this page does not cover"),
		wf.Markdown(
			"Things that are not BESPA-specific and that BESPA doesn't ",
			"opinionate: TLS termination, HSTS, cookie `Secure` / ",
			"`HttpOnly` / `SameSite` flags, load balancing, blue/green ",
			"rollouts, log shipping. Pick the patterns your platform uses; ",
			"the framework is unaware.",
		),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Handlers & routing](/build/handlers-and-routing) — ",
			"the mux wiring this page builds on.\n",
			"\n",
			"[Sessions & auth](/build/sessions-and-auth) — cookies, ",
			"CSRF, the auth posture you're deploying.\n",
			"\n",
			"[Extend → Assets & CSS](/extend/assets) — how widgets ",
			"register the CSS/JS that ends up in the single binary.",
		),
	)
	shared.Render(w, r, page)
}
