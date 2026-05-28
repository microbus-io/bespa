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

const themingMode = `// Always render dark.
page := wf.Page().WithThemeDark()

// Always render light.
page := wf.Page().WithThemeLight()

// Follow the OS / browser preference (the default).
page := wf.Page().WithThemeDefault()
`

const themingPalette = `// Pick from the framework's preset palettes.
page := wf.Page().WithKeyColors(css.PresetKeyColors[0])

// Or build your own from a single source color.
custom := css.KeyColorsFromString("violet")
page := wf.Page().WithKeyColors(custom)
`

const themingPersist = `func Render(w http.ResponseWriter, r *http.Request, page *widget.PageWidget) {
    session := SessionOf(w, r)
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

// HandleTheming covers theming a BESPA page.
func HandleTheming(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Theming"),

		wf.Markdown(
			"BESPA bakes Material Design 3 into the framework: every built-in ",
			"widget references the design tokens (`--md-sys-color-*`, ",
			"`--md-sys-typescale-*`) so a page automatically picks up your theme. ",
			"This page is about the app-level controls — picking dark vs. light, ",
			"choosing a key color palette, and persisting the user's preference.",
		),
		wf.HeadlineMedium("Light, dark, or system"),
		wf.Markdown("Three methods on `Page` set the appearance mode:"),
		wf.CodeBlock(themingMode).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The default follows the browser's `prefers-color-scheme` media ",
			"query, so the page tracks the OS setting and even switches live if ",
			"the user toggles the system theme. `WithThemeDark` / `WithThemeLight` ",
			"force one or the other regardless of the OS.",
		),
		wf.HeadlineMedium("Key color palettes"),
		wf.Markdown(
			"Material 3 generates an entire token palette — primary / secondary ",
			"/ tertiary / error / surface, with their on-color counterparts — ",
			"from a single \"source\" color. The `css` package exposes presets ",
			"plus a constructor that derives a palette from any color string:",
		),
		wf.CodeBlock(themingPalette).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"All built-in widgets and any widget that references ",
			"`--md-sys-color-*` recolors automatically; no widget needs to know ",
			"which palette is in use.",
		),
		wf.HeadlineMedium("Persisting the preference"),
		wf.Markdown(
			"User-specific preferences belong in the session, not the URL. The ",
			"website's `shared.Render` wrapper reads the active session's saved ",
			"theme/palette and applies them to every page before drawing:",
		),
		wf.CodeBlock(themingPersist).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"For a real app, swap the in-memory session for whatever persistence ",
			"you use; the page-side API stays the same.",
		),
		wf.HeadlineMedium("How the tokens get to the browser"),
		wf.Markdown(
			"On the first request the framework writes the resolved palette and ",
			"typography scale to `/bespa/tones.css` and `/bespa/style.css`. ",
			"Browsers cache those CSS files; a 24-hour cache header keeps revisits ",
			"cheap. When the user switches themes, the new key colors are picked ",
			"up on the next page navigation — there's no client-side palette ",
			"computation.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Extend → Material theming](/extend/theming) ",
			"— how custom widgets should reference tokens so they recolor for free.",
			"\n\n",
			"[Profile page](/profile) — a working theme/palette switcher backed ",
			"by session storage.",
		),
	)
	shared.Render(w, r, page)
}
