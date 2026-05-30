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

package showcase

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

func HandleMermaid(w http.ResponseWriter, r *http.Request) {
	flowchart := wf.Mermaid(`flowchart LR
    Start([Request arrives]) --> Auth{Authenticated?}
    Auth -- No --> Login[Redirect to login]
    Auth -- Yes --> Handler[Run handler]
    Handler --> Render[Render widgets]
    Render --> Diff[Diff fragments]
    Diff --> Reply([Send response])`)

	sequence := wf.Mermaid(`sequenceDiagram
    participant Browser
    participant Server
    participant DB
    Browser->>Server: GET /page?x=1
    Server->>DB: Query
    DB-->>Server: Rows
    Server-->>Browser: HTML
    Browser->>Server: POST ?x=2 (Bespa-Fetch: 1)
    Server-->>Browser: Partial fragments
    Note over Browser: Swap by data-id`)

	state := wf.Mermaid(`stateDiagram-v2
    [*] --> Idle
    Idle --> Loading: submit
    Loading --> Success: 2xx
    Loading --> Error: 4xx/5xx
    Success --> Idle: reset
    Error --> Idle: retry`)

	gantt := wf.Mermaid(`gantt
    title Release plan
    dateFormat YYYY-MM-DD
    section Design
    Mockups          :done,  des1, 2026-05-01, 7d
    Review           :active, des2, 2026-05-08, 3d
    section Build
    Widgets          :        b1, after des2, 10d
    Docs             :        b2, after b1, 5d`).WithHeight("300px")

	classDiagram := wf.Mermaid(`classDiagram
    class Widget {
        +ID() string
        +Draw(w, r) error
    }
    class WidgetBase~T~ {
        -shown bool
        +RedrawIfChanged(r, keys) T
        +HideIf(pred) T
    }
    class ChartWidget
    class MermaidWidget
    Widget <|.. WidgetBase
    WidgetBase <|-- ChartWidget
    WidgetBase <|-- MermaidWidget`)

	zoomable := wf.Mermaid(`flowchart TB
    A[Request] --> B[Router]
    B --> C{Method}
    C -->|GET| D[Handler GET]
    C -->|POST| E[Handler POST]
    D --> F[Read state]
    E --> F
    F --> G[Build tree]
    G --> H[Draw]
    H --> I[Compress]
    I --> J[Response]
    J --> K[Client]
    K --> L{Bespa-Fetch?}
    L -->|Yes| M[Swap fragments]
    L -->|No| N[Replace page]`).WithHeight("500px").WithZoomPan(true)

	page := wf.Page().Add(
		wf.AppBar("Mermaid widgets"),
		wf.HeadlineMedium("Flowchart"),
		flowchart,
		wf.HeadlineMedium("Sequence"),
		sequence,
		wf.HeadlineMedium("State"),
		state,
		wf.HeadlineMedium("Gantt"),
		gantt,
		wf.HeadlineMedium("Class"),
		classDiagram,
		wf.HeadlineMedium("Zoom & pan"),
		"Wheel or pinch to zoom, drag to pan, double-click to recenter.",
		zoomable,
	)
	shared.Render(w, r, page)
}
