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

package mermaid

import (
	"fmt"
	"strings"
)

// expandResult holds the output of expanding CSS var() references in a
// Mermaid source string.
type expandResult struct {
	source string    // rewritten Mermaid with var() refs replaced by currentColor
	rules  []cssRule // bridge rules to be emitted as a <style> block on the page
}

// cssRule is one selector/declaration pair the host page must emit so the
// browser's CSS cascade supplies the color Mermaid's classDef parser cannot.
type cssRule struct {
	selector string // e.g. ".completed", ".completed .nodeLabel", "#fo_s1", ".edgePaths"
	color    string // the original var() / rgb(var()) / rgba(var(), N) value
	label    bool   // true when the rule targets node-label text (needs !important)
}

// expandVars scans the Mermaid source for CSS var() references inside
// classDef, style, and `linkStyle default` directives. Each fill/stroke
// property whose value contains var() is rewritten to currentColor, and a
// bridge rule is recorded that sets `color` on the corresponding class /
// id / edge selector. color: properties (label text) are removed from the
// directive entirely and bridged via a rule targeting Mermaid's documented
// `.nodeLabel` inner hook with !important.
//
// Hex literals, named colors, and bare `currentColor` pass through unchanged.
// Directives the rewriter does not understand (e.g. numeric linkStyle
// indices) pass through unchanged too — they will produce a Mermaid parse
// error if they contain var() references, which is the same behavior the
// author would have hit without this expansion.
func expandVars(source string) expandResult {
	var rewritten strings.Builder
	var rules []cssRule

	lines := strings.Split(source, "\n")
	for i, line := range lines {
		newLine, lineRules := expandLine(line)
		rewritten.WriteString(newLine)
		if i < len(lines)-1 {
			rewritten.WriteString("\n")
		}
		rules = append(rules, lineRules...)
	}

	return expandResult{
		source: rewritten.String(),
		rules:  dedupeRules(rules),
	}
}

// expandLine returns the rewritten line and the bridge rules contributed by
// that line. A line without a recognized directive (or without any var()
// references) is returned unchanged with no rules.
func expandLine(line string) (string, []cssRule) {
	trimmed := strings.TrimSpace(line)
	if !strings.Contains(trimmed, "var(--") {
		return line, nil
	}
	indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]

	var directive, target, body string
	switch {
	case strings.HasPrefix(trimmed, "classDef "):
		directive = "classDef"
		target, body = splitTargetBody(strings.TrimPrefix(trimmed, "classDef "))
	case strings.HasPrefix(trimmed, "style "):
		directive = "style"
		target, body = splitTargetBody(strings.TrimPrefix(trimmed, "style "))
	case strings.HasPrefix(trimmed, "linkStyle "):
		directive = "linkStyle"
		target, body = splitTargetBody(strings.TrimPrefix(trimmed, "linkStyle "))
	default:
		return line, nil
	}
	if target == "" || body == "" {
		return line, nil
	}

	props := splitProperties(body)
	var keptProps []string
	var rules []cssRule
	for _, p := range props {
		key, value := splitKV(p)
		if !strings.Contains(value, "var(--") {
			keptProps = append(keptProps, p)
			continue
		}
		sel, isLabel, ok := buildSelector(directive, target, key)
		if !ok {
			// Unrecognized property/directive shape. Pass the original value
			// through and let Mermaid surface the error.
			keptProps = append(keptProps, p)
			continue
		}
		rules = append(rules, cssRule{selector: sel, color: value, label: isLabel})
		if isLabel {
			// color: directive on classDef/style — drop it; CSS handles text.
			continue
		}
		keptProps = append(keptProps, key+":currentColor")
	}

	if len(keptProps) == 0 {
		// Every property was a bridged color: directive. classDef needs a
		// body to register the class; emit a no-op fill so the :::name
		// reference resolves.
		if directive == "classDef" {
			keptProps = []string{"fill:currentColor"}
		} else {
			// style ID and linkStyle SELECTOR can be omitted entirely.
			return "", rules
		}
	}

	return fmt.Sprintf("%s%s %s %s", indent, directive, target, strings.Join(keptProps, ",")), rules
}

// buildSelector returns the CSS selector for a rewritten property, whether
// it targets node-label text (needs !important), and whether the
// directive/property combination is one the rewriter can bridge. Returns
// ok=false for combinations that should pass through unchanged.
func buildSelector(directive, target, key string) (sel string, label, ok bool) {
	switch directive {
	case "classDef":
		switch key {
		case "fill", "stroke":
			return "." + target, false, true
		case "color":
			return "." + target + " .nodeLabel", true, true
		}
	case "style":
		switch key {
		case "fill", "stroke":
			return "#" + target, false, true
		case "color":
			return "#" + target + " .nodeLabel", true, true
		}
	case "linkStyle":
		// Numeric link indices map awkwardly to CSS (per-edge IDs are not
		// stable across re-renders), so only `linkStyle default` is bridged.
		if target == "default" && key == "stroke" {
			return ".edgePaths", false, true
		}
	}
	return "", false, false
}

// splitTargetBody splits a directive's remainder into the target (class
// name / id / link selector) and the property body.
func splitTargetBody(rest string) (target, body string) {
	sp := strings.IndexAny(rest, " \t")
	if sp < 0 {
		return rest, ""
	}
	return rest[:sp], strings.TrimSpace(rest[sp+1:])
}

// splitProperties splits a classDef/style body into property declarations,
// respecting parenthesized values so commas inside `rgba(var(--x), 0.5)`
// are not treated as separators.
func splitProperties(body string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, c := range body {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(body[start:i]))
				start = i + 1
			}
		}
	}
	parts = append(parts, strings.TrimSpace(body[start:]))
	return parts
}

// splitKV splits "fill:value" into ("fill", "value"). Only the first colon
// is the separator; subsequent colons (e.g. inside a value) are kept verbatim.
func splitKV(prop string) (string, string) {
	idx := strings.Index(prop, ":")
	if idx < 0 {
		return strings.TrimSpace(prop), ""
	}
	return strings.TrimSpace(prop[:idx]), strings.TrimSpace(prop[idx+1:])
}

// dedupeRules collapses rules that target the same (selector, label) slot.
// The first rule wins. This matters when one classDef declares fill and
// stroke as different var() references: the cascade only carries one color,
// so the fill rule is kept and the stroke's currentColor inherits the same
// value. A classDef that needs visually distinct fill and stroke colors
// should hex-code at least one of them rather than var()-reference both.
func dedupeRules(rules []cssRule) []cssRule {
	type key struct {
		selector string
		label    bool
	}
	seen := make(map[key]bool, len(rules))
	out := make([]cssRule, 0, len(rules))
	for _, r := range rules {
		k := key{r.selector, r.label}
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, r)
	}
	return out
}

// formatCSS scopes the rules under a parent selector (typically "#<id>"
// where id is the widget's container) and returns a CSS string suitable for
// embedding in a <style> tag. Returns "" when there are no rules.
func formatCSS(rules []cssRule, scope string) string {
	if len(rules) == 0 {
		return ""
	}
	var b strings.Builder
	for _, r := range rules {
		important := ""
		if r.label {
			important = " !important"
		}
		if scope != "" {
			fmt.Fprintf(&b, "%s %s { color: %s%s; }\n", scope, r.selector, r.color, important)
		} else {
			fmt.Fprintf(&b, "%s { color: %s%s; }\n", r.selector, r.color, important)
		}
	}
	return b.String()
}
