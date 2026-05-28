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

// HandleCode demonstrates the code-block widget across a few popular languages.
func HandleCode(w http.ResponseWriter, r *http.Request) {
	goSource := `package main

import (
	"fmt"
	"net/http"
)

// HandleHello greets the requester.
func HandleHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	http.HandleFunc("/hello", HandleHello)
	http.ListenAndServe(":8080", nil)
}
`

	pythonSource := `from dataclasses import dataclass
from typing import Iterable


@dataclass
class State:
    abbrev: str
    name: str
    population: int

    def density(self, land_sq_mi: float) -> float:
        """People per square mile."""
        return self.population / land_sq_mi


def by_population(states: Iterable[State]) -> list[State]:
    return sorted(states, key=lambda s: s.population, reverse=True)


if __name__ == "__main__":
    sample = [
        State("CA", "California", 39_538_223),
        State("TX", "Texas",      29_145_505),
        State("FL", "Florida",    21_538_187),
    ]
    for s in by_population(sample):
        print(f"{s.abbrev}: {s.population:,}")
`

	htmlSource := `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Hello</title>
    <style>
        body { font-family: sans-serif; color: #222; }
        .greeting { font-size: 2rem; font-weight: 600; }
    </style>
</head>
<body>
    <h1 class="greeting">Hello, world!</h1>
    <p>This is a <a href="/about">simple page</a>.</p>
    <script>
        document.querySelector(".greeting").addEventListener("click", () => {
            alert("clicked");
        });
    </script>
</body>
</html>
`

	page := wf.Page().Add(
		wf.AppBar("Code block widget"),
		`The code block widget renders source code with server-side syntax highlighting via Chroma.
		Token classes map to Material design color tokens so the same markup recolors on theme change
		without re-running the highlighter.`,

		wf.HeadlineMedium("Go"),
		`Explicit language via `, wf.Code("WithLanguage(\"go\")"), `. Comments are dimmed-italic, keywords
		use the primary color, types use the secondary, and string and numeric literals use the tertiary.`,
		wf.SpacerBreak(),
		wf.CodeBlock(goSource).WithLanguage("go"),

		wf.HeadlineMedium("Python"),
		`Same widget, different lexer. Decorators (`, wf.Code("@dataclass"), `), built-ins (`, wf.Code("print"), `,
		`, wf.Code("sorted"), `, `, wf.Code("len"), `) and class names each get their own token class. Line
		numbers are enabled here via `, wf.Code("WithLineNumbers(true)"), `.`,
		wf.SpacerBreak(),
		wf.CodeBlock(pythonSource).WithLanguage("python").WithLineNumbers(true),

		wf.HeadlineMedium("HTML"),
		`Mixed-language highlighting: tags, attributes, embedded CSS inside `, wf.Code("<style>"), `, and
		embedded JavaScript inside `, wf.Code("<script>"), ` all get their respective lexers chained
		together by Chroma's HTML lexer. This block is capped at 10 visible rows via `, wf.Code("WithMaxRows(10)"),
		`; scroll the inner pane to see the rest.`,
		wf.SpacerBreak(),
		wf.CodeBlock(htmlSource).WithLanguage("html").WithMaxRows(10),

		wf.SpacerParagraph(),
	)
	shared.Render(w, r, page)
}
