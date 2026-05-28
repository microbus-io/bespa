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

const themingCSSGood = `/* mywidget.css */
BODY.Top .MyWidget {
    background: rgb(var(--md-sys-color-surface-container));
    color: rgb(var(--md-sys-color-on-surface));
    border: solid 1px rgb(var(--md-sys-color-outline-variant));
    border-radius: 4px;
    padding: 0.5em 1ch;
    font-family: var(--md-sys-typescale-body-medium-font);
    font-size: var(--md-sys-typescale-body-medium-size);
    line-height: var(--md-sys-typescale-body-medium-line-height);
}
BODY.Top .MyWidget.Selected {
    background: rgb(var(--md-sys-color-primary-container));
    color: rgb(var(--md-sys-color-on-primary-container));
}
`

const themingCSSBad = `/* DON'T — hard-coded colors don't recolor on theme switch. */
.MyWidget {
    background: #f5f5f5;
    color: #111;
    border: solid 1px #ddd;
}
`

const themingResolve = `// chart/chart.js — resolves CSS var() to RGB at draw time for canvas use.
function chart_color(varName, alpha) {
    const v = getComputedStyle(document.documentElement)
        .getPropertyValue(varName).trim();
    if (!v) return null;
    const parts = v.split(/[\s,]+/);
    if (typeof alpha === "number" && alpha < 1) {
        return "rgba(" + parts[0] + "," + parts[1] + "," + parts[2] + "," + alpha + ")";
    }
    return "rgb(" + parts[0] + "," + parts[1] + "," + parts[2] + ")";
}
`

const themingThemeChange = `// Re-resolve and re-render when the page theme changes.
function onThemeChange() {
    // Invalidate any cached resolved colors and re-render every live instance.
    // ...
}

// Class changes on <html> — flipping LightTheme / DarkTheme.
new MutationObserver(onThemeChange).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
});

// OS-level color scheme changes when the page is in default (system-follows) mode.
window.matchMedia("(prefers-color-scheme: dark)").addEventListener(
    "change", onThemeChange,
);
`

// HandleTheming covers Material theming for widget CSS and runtime use.
func HandleTheming(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Material theming"),

		wf.Markdown(
			"A widget that references the Material design tokens themes itself ",
			"for free: light vs. dark, palette switches, custom user themes all ",
			"flow through. A widget that hardcodes colors becomes a maintenance ",
			"hole. This page is the convention for getting the first kind.",
		),
		wf.HeadlineMedium("The token names"),
		wf.Markdown(
			"BESPA exposes the Material 3 token set as CSS custom properties ",
			"on the `:root` element. The families you'll reach for most often:",
			"\n\n",
			"- Colors: `--md-sys-color-{primary,secondary,tertiary,error}` plus ",
			"the matching `on-` pairs; ",
			"`--md-sys-color-{surface,surface-container,surface-tint}`, ",
			"`--md-sys-color-{background,on-background}`, ",
			"`--md-sys-color-{outline,outline-variant}`.\n",
			"- Typography: ",
			"`--md-sys-typescale-{display,headline,title,body,label}-{large,medium,small}-{font,size,line-height,weight,tracking}`. ",
			"Mix and match — e.g. `--md-sys-typescale-body-medium-size`.\n",
			"- State opacities: ",
			"`--md-sys-state-{disabled,hover,pressed,focus}-state-layer-opacity` ",
			"— multiply against a color when you need a 30%-opacity hover overlay.",
		),
		wf.HeadlineMedium("CSS form"),
		wf.Markdown(
			"Color tokens are stored as `R G B` triplets (no parentheses, no ",
			"rgb()) so you can use them in `rgb()` or `rgba()` with any alpha ",
			"you want:",
		),
		wf.CodeBlock(themingCSSGood).WithLanguage("css"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Compare to the antipattern — hardcoded colors that break the moment ",
			"the user picks dark mode:",
		),
		wf.CodeBlock(themingCSSBad).WithLanguage("css"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Rule of thumb: if your CSS contains a hex color or a literal ",
			"`rgb(123, …)`, you're hard-coding the theme. Replace it with the ",
			"matching token.",
		),
		wf.HeadlineMedium("The BODY.Top prefix"),
		wf.Markdown(
			"You'll notice every framework stylesheet prefixes its selectors ",
			"with `BODY.Top `. This isn't decoration — it's a specificity trick ",
			"so the framework styles outrank generic CSS from third-party ",
			"libraries that share class names. Use the same prefix in your own ",
			"widget CSS for the same reason.",
		),
		wf.HeadlineMedium("Canvas-rendered widgets"),
		wf.Markdown(
			"For widgets that render via JavaScript canvas (charts, custom ",
			"drawing), CSS `var()` inside Canvas colors doesn't work — the ",
			"canvas API doesn't resolve CSS variables. Resolve them in JS first:",
		),
		wf.CodeBlock(themingResolve).WithLanguage("javascript"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The chart widget does exactly this — see `chart/chart.js`. The ",
			"pattern: read the computed CSS value, parse, hand the literal RGB ",
			"to the renderer. Cache the result keyed by token name so you only ",
			"pay the lookup once per token per page.",
		),
		wf.HeadlineMedium("Reacting to theme changes"),
		wf.Markdown(
			"CSS-themed widgets recolor automatically because CSS variables ",
			"flow down on every recompute. Canvas widgets need a hook — listen ",
			"for class changes on the `<html>` element and for ",
			"`prefers-color-scheme` media changes, clear your resolved-color ",
			"cache, and re-render:",
		),
		wf.CodeBlock(themingThemeChange).WithLanguage("javascript"),
		wf.SpacerBreak(),
		wf.Markdown(
			"This is the pattern the chart widget uses in production — see ",
			"`chart/chart.js` for the full implementation including the ",
			"instance registry.",
		),
		wf.HeadlineMedium("Tonal palettes"),
		wf.Markdown(
			"For widgets that need many distinct colors (charts, color-coded ",
			"badges, multi-series visuals), the framework also exposes ",
			"`--md-sys-color-primary-{0..330}deg` at 30° hue rotations of the ",
			"primary tone. Picking your series colors from these gives you a ",
			"harmonized rainbow that retracks any palette switch.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Build → Theming](/build/theming) — the consumer-side ",
			"controls (light/dark, palette).",
			"\n\n",
			"[Assets & CSS](/extend/assets) — where the widget's CSS ",
			"lives and how it gets to the browser.",
			"\n\n",
			"Read `css/keycolors.go` for the palette generation and ",
			"`chart/chart.js` for the canvas-side theme handling.",
		),
	)
	shared.Render(w, r, page)
}
