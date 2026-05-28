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

const composingSearchBox = `// factory is the package-local composite factory — see the next section.
type SearchBoxWidget struct {
    *widget.WidgetBase[*SearchBoxWidget]
    name        string
    placeholder string
}

func (f MyFactory) SearchBox(name string) *SearchBoxWidget {
    s := &SearchBoxWidget{name: name, placeholder: "Search…"}
    s.WidgetBase = widget.NewWidgetBase(s)
    return s
}

func (s *SearchBoxWidget) WithPlaceholder(p string) *SearchBoxWidget {
    s.placeholder = p
    return s
}

func (s *SearchBoxWidget) Draw(w io.Writer, r *http.Request) error {
    return factory.Form().Add(
        factory.InputText(s.name, "").
            WithPlaceholder(s.placeholder).
            WithAutoSubmit(true),
    ).Draw(w, r)
}
`

const composingChildren = `func (s *SearchBoxWidget) Children() []widget.Widget {
    return []widget.Widget{ s.composed() }
}

func (s *SearchBoxWidget) Draw(w io.Writer, r *http.Request) error {
    return s.composed().Draw(w, r)
}

// composed builds the sub-tree once so Children and Draw agree.
func (s *SearchBoxWidget) composed() widget.Widget {
    return factory.Form().Add(
        factory.InputText(s.name, "").
            WithPlaceholder(s.placeholder).
            WithAutoSubmit(true),
    )
}
`

const composingFactoryAccess = `// At the top of your package:
import "github.com/microbus-io/bespa/basic"
import "github.com/microbus-io/bespa/form"

var factory = struct {
    widget.WidgetFactory
    basic.BasicFactory
    form.FormFactory
}{}
`

// HandleComposing covers composing existing widgets.
func HandleComposing(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Composing existing widgets"),

		wf.Markdown(
			"Most useful widgets aren't new HTML — they're a recognizable ",
			"combination of existing widgets you want to drop in by name. A ",
			"search box is a form with one input. A toast is a snackbar with ",
			"an icon. A status pill is a chip with a colored dot. These don't ",
			"need `Tag` or hand-written HTML; they just delegate to a sub-tree ",
			"built from the framework's existing widgets.",
		),
		wf.HeadlineMedium("Delegating to a sub-tree"),
		wf.Markdown(
			"A `SearchBox` that wraps a `Form` + `InputText` + `WithAutoSubmit` ",
			"is twelve lines:",
		),
		wf.CodeBlock(composingSearchBox).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The `Draw` method builds the sub-tree on the fly and delegates to ",
			"its `Draw`. From the consumer's side it looks like any other widget: ",
			"`wf.SearchBox(\"q\").WithPlaceholder(\"Find people\")`.",
		),
		wf.HeadlineMedium("Implementing Children too"),
		wf.Markdown(
			"The version above works for static rendering, but the framework's ",
			"traversal pass also walks `Children()` to collect input values, ",
			"assign stable IDs, and propagate redraw. For a widget that exposes ",
			"form inputs to the outside world, return the same sub-tree from ",
			"`Children()`:",
		),
		wf.CodeBlock(composingChildren).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Factor the sub-tree out into a helper so `Draw` and `Children` ",
			"return the same shape — otherwise a state walk and a render walk ",
			"may disagree on which widgets exist, and stable ID generation ",
			"falls apart.",
		),
		wf.HeadlineMedium("Getting at the framework's widgets"),
		wf.Markdown(
			"Inside your widget you'll want to construct `Form`, `InputText`, ",
			"etc. Build a package-local factory struct that composes whichever ",
			"families you need:",
		),
		wf.CodeBlock(composingFactoryAccess).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This is the same pattern every widget package in the framework ",
			"uses — `basic/factory.go` composes `widget.WidgetFactory` + itself, ",
			"`form/factory.go` adds `basic` on top, and so on.",
		),
		wf.HeadlineMedium("When to compose vs. write Draw"),
		wf.Markdown(
			"If everything you need is in the framework already, compose. If you ",
			"need a specific HTML element (a `<canvas>`, a SVG, a third-party JS ",
			"init script), drop down to `Tag` — see Extend → Widget ",
			"anatomy. The two approaches mix freely; a complex widget often uses ",
			"`Tag` for one outer element and composed widgets for everything ",
			"inside.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Widget anatomy](/extend/anatomy) — for the lower-level case ",
			"where you do write Draw.",
			"\n\n",
			"Read `table/quicksearch.go` in the framework source — it's a real ",
			"composed widget that wraps `InputText` with table-aware defaults.",
		),
	)
	shared.Render(w, r, page)
}
