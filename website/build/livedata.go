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

const liveDataWidgetSketch = `// A small custom widget that opens an SSE connection and, on each event,
// programmatically clicks a hidden link inside itself. BESPA's client
// runtime catches the click the same way it catches a user click — same
// action-URL grammar, same state merging, same partial-redraw pass.

type LiveFeed struct {
    *widget.WidgetBase[*LiveFeed]
    feedURL string
}

func (f MyFactory) LiveFeed(feedURL string) *LiveFeed {
    l := &LiveFeed{feedURL: feedURL}
    l.WidgetBase = widget.NewWidgetBase(l)
    return l
}

func (l *LiveFeed) Draw(w io.Writer, r *http.Request) error {
    _, err := fmt.Fprintf(w,
        "<span data-id=%q hidden>"+
        "<a class=\"LiveFeedTrigger\"></a>"+
        "<script>livefeed_init(%q, %q)</script>"+
        "</span>",
        l.ID(), l.ID(), l.feedURL)
    return err
}`

const liveDataClientJS = `// Companion JavaScript, registered via widget.AssetRegistry.RegisterScript.
// On each server event, set the href on the hidden anchor and click it. The
// page-level click handler in page.js picks it up and treats it as ordinary
// navigation — no private runtime APIs needed.

function livefeed_init(id, feedURL) {
    const root = document.getElementById(id);
    const trigger = root.querySelector("a.LiveFeedTrigger");
    const src = new EventSource(feedURL);
    src.onmessage = (e) => {
        trigger.setAttribute("href", e.data);
        trigger.click();
    };
}
`

// HandleLiveData explains BESPA's deliberate non-position on async / live
// data, and the recipe for plugging in SSE or WebSockets as a custom widget.
func HandleLiveData(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Live data"),

		wf.Markdown(
			"BESPA does not ship a built-in mechanism for server-pushed ",
			"updates — no SSE channel, no WebSocket bus, no long-poll helper. ",
			"This is deliberate. The framework is stateless top to bottom: ",
			"every request is independent, no per-connection state lives on ",
			"the server, and you can scale horizontally or restart a node ",
			"without coordination. Persistent connections would break that ",
			"property, so they're left to the app to wire up when needed.",
		),

		wf.HeadlineMedium("The pattern"),
		wf.Markdown(
			"Live updates are a five-minute widget on top of the framework, ",
			"not a missing feature.",
			"\n\n",
			"Write a small custom widget that opens an SSE or WebSocket ",
			"connection from the browser. On each message, the widget hands ",
			"the payload to the client runtime as if the user had clicked a ",
			"link — same action-URL grammar, same state merging, same ",
			"partial-redraw pass. The server pushes a URL fragment; BESPA ",
			"navigates to it.",
			"\n\n",
			"That means the server side of the feature is whatever HTTP ",
			"endpoint produces the events. There's no \"BESPA connection ",
			"layer\" to integrate with — write the SSE handler the way you'd ",
			"write any SSE handler in Go, emit `data:` lines containing URL ",
			"fragments, and you're done.",
		),

		wf.HeadlineMedium("Polling: the built-in option"),
		wf.Markdown(
			"Before reaching for SSE or WebSockets, check whether polling is ",
			"enough. `Progress` ships with a built-in polling loop: give it a ",
			"refresh URL and an interval, and it asks that endpoint for ",
			"`{value, stop, action}` JSON on its own clock. When the job is ",
			"done, the server returns `stop:true` with an `action` like ",
			"`?done=1` — the widget dispatches a synthetic click on a hidden ",
			"anchor with that href, and BESPA navigates as if the user had ",
			"clicked it. (Same trick as the SSE widget below; just packaged.)",
			"\n\n",
			"```go\n",
			"wf.Progress().\n",
			"    WithMax(100).\n",
			"    WithRefreshURL(\"/jobs/\"+jobID+\"/progress\").\n",
			"    WithRefreshInterval(500 * time.Millisecond)\n",
			"```\n",
			"\n",
			"This is the right tool for: long-running jobs with a known end, ",
			"upload/processing progress, anything where \"the answer arrives ",
			"in seconds, not milliseconds\" is fine. It costs you no extra ",
			"infrastructure — it's just HTTP requests on an interval, served ",
			"by your existing handlers, and the interval is bounded by the ",
			"`stop` flag rather than running forever.",
			"\n\n",
			"For everything else — push notifications with no predictable ",
			"end, presence, multi-tab sync — keep reading.",
		),

		wf.HeadlineMedium("Sketch: an SSE-driven widget"),
		wf.CodeBlock(liveDataWidgetSketch).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.CodeBlock(liveDataClientJS).WithLanguage("js"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The widget renders one hidden `<a>` element. On every server ",
			"event, its JS sets the `href` to the pushed payload and calls ",
			"`.click()`. BESPA's page-level click handler catches the bubbled ",
			"event and treats it as a normal link click. If the server pushes ",
			"`?notif=3`, the page's notification counter redraws via its ",
			"existing `RedrawIfChanged(r, \"notif\")` binding. If it pushes ",
			"`^?modal=/alerts/42`, a modal opens. Nothing about the receiving ",
			"page has to know live data is involved.",
		),

		wf.HeadlineMedium("What the server emits"),
		wf.Markdown(
			"Keep the payload to a URL fragment, not HTML. Two reasons:",
			"\n\n",
			"- The receiving page already knows how to render itself from a ",
			"URL. Pushing HTML would require the server to know the page's ",
			"structure — coupling that BESPA otherwise avoids.\n",
			"- A URL fragment is tiny. SSE is cheap when each event is a ",
			"dozen bytes; it stops being cheap when each event is an HTML ",
			"document.",
		),

		wf.HeadlineMedium("When to reach for this — and when not"),
		wf.Markdown(
			"Live data is the right call for: notifications, presence ",
			"indicators, \"someone else just edited this row,\" job-progress ",
			"updates, multi-tab synchronization.",
			"\n\n",
			"It's the wrong call when polling is enough — use `Progress` ",
			"(above) for jobs with a known end, or any state with a ",
			"once-every-few-seconds cadence. No persistent connection, no ",
			"fanout to manage.",
		),

		wf.HeadlineMedium("Operational note"),
		wf.Markdown(
			"Once you open persistent connections, you've taken on what ",
			"BESPA chose to avoid: a connection table, heartbeats, server ",
			"affinity (or a pub/sub layer to escape it), and the failure ",
			"modes that come with long-lived sockets behind load balancers. ",
			"That's fine — those are well-understood problems — but it is a ",
			"real shift in operational shape from a stateless app. Reach for ",
			"it when the UX needs it, not as a default.",
		),

		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Extend → Widget anatomy](/extend/anatomy) — how to ",
			"write the custom widget shown above.",
			"\n\n",
			"[Extend → Assets & CSS](/extend/assets) — how to ",
			"register the companion JavaScript.",
			"\n\n",
			"[Basics → Action-URL pattern](/basics/action-url-pattern) ",
			"— the grammar the server-pushed payload uses.",
		),
	)
	shared.Render(w, r, page)
}
