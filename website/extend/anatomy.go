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

package extend

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

const anatomyInterface = `// widget/widget.go
type Widget interface {
    ID() string
    SetID(id string)
    Children() []Widget
    Draw(w io.Writer, r *http.Request) error
    Drawn(r *http.Request) bool
    Shown(r *http.Request) bool
}
`

const anatomyBase = `type GreetingWidget struct {
    *widget.WidgetBase[*GreetingWidget]
    name string
}

func (f MyFactory) Greeting(name string) *GreetingWidget {
    g := &GreetingWidget{name: name}
    g.WidgetBase = widget.NewWidgetBase(g)
    return g
}

func (g *GreetingWidget) WithName(name string) *GreetingWidget {
    g.name = name
    return g
}

func (g *GreetingWidget) Draw(w io.Writer, r *http.Request) error {
    return widget.NewWriterAssistant(w).
        WriteString("<span data-id=\"", g.ID(), "\">Hello, ", html.EscapeString(g.name), "!</span>").
        Err()
}
`

const anatomyComposed = `func (g *GreetingWidget) Draw(w io.Writer, r *http.Request) error {
    return factory.Tag("span").
        Class("Greeting").
        Attr("data-id", g.ID()).
        Add("Hello, ", g.name, "!").
        When(g.Shown(r)).
        Draw(w, r)
}
`

// HandleAnatomy covers the Widget interface and WidgetBase.
func HandleAnatomy(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Widget anatomy"),

		wf.Markdown(
			"Every widget in BESPA is a Go struct that implements the `Widget` ",
			"interface. In practice you never implement it from scratch — you ",
			"embed `WidgetBase[T]` and write a `Draw` method.",
		),
		wf.HeadlineMedium("The interface"),
		wf.Markdown(
			"The contract is six methods, but only one of them (`Draw`) is ",
			"something you're likely to write:",
		),
		wf.CodeBlock(anatomyInterface).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Everything else has a sensible default in `WidgetBase`. `ID`/`SetID` ",
			"plumb the partial-redraw protocol; `Children` defaults to `nil` (no ",
			"children); `Drawn` and `Shown` are controlled by `RedrawIfChanged` / ",
			"`HideIf` — see Extend → State-aware widgets.",
		),
		wf.HeadlineMedium("The base type"),
		wf.Markdown(
			"The minimal widget — a Go struct, a constructor that wires up ",
			"`WidgetBase`, optional builder methods, and a `Draw`:",
		),
		wf.CodeBlock(anatomyBase).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Three things going on:",
			"\n\n",
			"- The `*widget.WidgetBase[*GreetingWidget]` embed is generic — the ",
			"type parameter tells `WithFoo`-style methods on the base what type ",
			"to return. That's how you can chain `.WithName(...)` on the ",
			"concrete type.\n",
			"- `widget.NewWidgetBase(g)` passes the concrete pointer back to the ",
			"base so its `WithID`, `RedrawIf`, etc. can return `*GreetingWidget`, ",
			"not `*WidgetBase`.\n",
			"- `Draw` writes the widget's HTML. Include `data-id` set to ",
			"`g.ID()` on whatever element should be addressable by the ",
			"partial-redraw swap — otherwise the framework can't replace this ",
			"widget in place.",
		),
		wf.HeadlineMedium("Composing instead of escaping by hand"),
		wf.Markdown(
			"In practice almost no widget writes raw HTML. Use `Tag` — the ",
			"framework's tag builder that handles escaping, attribute merging, ",
			"and conditional rendering:",
		),
		wf.CodeBlock(anatomyComposed).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This is the form you'll see used by every widget under `basic/` ",
			"and `form/`. `Tag(\"span\").Add(g.name, ...)` calls the framework's ",
			"`Any` coercion under the hood, which escapes plain strings; numbers, ",
			"times, and other types convert to widgets too.",
		),
		wf.HeadlineMedium("Where the widget lives"),
		wf.Markdown(
			"By convention each widget gets its own `.go` file plus optional ",
			"`.css` and `.js` siblings — e.g. `basic/heading.go` + ",
			"`basic/heading.css`. The factory file in the package's root registers ",
			"all the assets in one `init()`. See Extend → Assets & CSS ",
			"for that pattern.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Composing existing widgets](/extend/composing) — for most ",
			"cases you don't need to write a Draw at all.",
			"\n\n",
			"[State-aware widgets](/extend/state-aware) — how Drawn / ",
			"Shown / RedrawIfChanged interact with your Draw.",
			"\n\n",
			"Read the source of `basic/heading.go` for a tiny worked example, ",
			"or `basic/menu.go` for one that uses companion CSS and JS.",
		),
	)
	shared.Render(w, r, page)
}
