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

package basic

import (
	"html"
	"io"
	"net/http"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&ScriptInjectorWidget{}) // Ensure interface

// ScriptInjectorWidget injects a script to the head of the page.
type ScriptInjectorWidget struct {
	*widget.WidgetBase[*ScriptInjectorWidget]
	src string
}

// ScriptInjector creates a new widget that appends a <script src=…> to
// document.head on first render. It deduplicates by src so reusing it across
// partial redraws is safe — the browser only loads the script once. No SRI
// or async/defer attributes are emitted; for tighter control inject your
// own <script> tag via Tag("script") or use AssetRegistry.RegisterIsolatedScript.
func (f BasicFactory) ScriptInjector(src string) *ScriptInjectorWidget {
	x := &ScriptInjectorWidget{
		src: src,
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// Draw renders the widget's HTML.
func (wgt *ScriptInjectorWidget) Draw(w io.Writer, r *http.Request) (err error) {
	return Tag("script").
		Attr("data-id", wgt.ID()).
		Add(HTMLUnsafe(`
if (!document.head.querySelector('script[src="`, html.EscapeString(wgt.src), `"]')) {
	const s = document.createElement("script");
	s.type = "text/javascript";
	s.src = '`, html.EscapeString(wgt.src), `';
	document.head.append(s);
}
`)).
		When(wgt.Shown(r)).
		Draw(w, r)
}
