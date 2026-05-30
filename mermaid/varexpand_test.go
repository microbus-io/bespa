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
	"strings"
	"testing"
)

func TestExpandVars_PassThrough(t *testing.T) {
	src := `flowchart LR
    classDef completed fill:#32a7c1,color:#f4f2ef,stroke:#32a7c1
    s1["task"]:::completed
    s1 --> s2`
	got := expandVars(src)
	if got.source != src {
		t.Errorf("source should pass through unchanged:\nwant: %q\n got: %q", src, got.source)
	}
	if len(got.rules) != 0 {
		t.Errorf("expected no rules, got %d: %+v", len(got.rules), got.rules)
	}
}

func TestExpandVars_ClassDefFill(t *testing.T) {
	src := `flowchart LR
    classDef completed fill:var(--md-sys-color-primary)
    s1["task"]:::completed`
	got := expandVars(src)
	if !strings.Contains(got.source, "fill:currentColor") {
		t.Errorf("fill not rewritten: %q", got.source)
	}
	if strings.Contains(got.source, "var(--") {
		t.Errorf("var() should be removed from source: %q", got.source)
	}
	if len(got.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(got.rules))
	}
	r := got.rules[0]
	if r.selector != ".completed" {
		t.Errorf("selector = %q, want .completed", r.selector)
	}
	if r.color != "var(--md-sys-color-primary)" {
		t.Errorf("color = %q", r.color)
	}
	if r.label {
		t.Error("fill rule should not be a label rule")
	}
}

func TestExpandVars_ClassDefAllThree(t *testing.T) {
	src := `classDef completed fill:var(--md-sys-color-primary-container),color:var(--md-sys-color-on-primary-container),stroke:var(--md-sys-color-primary-container)`
	got := expandVars(src)

	// Source should have fill:currentColor,stroke:currentColor; color: dropped.
	if !strings.Contains(got.source, "fill:currentColor") || !strings.Contains(got.source, "stroke:currentColor") {
		t.Errorf("fill/stroke not rewritten: %q", got.source)
	}
	if strings.Contains(got.source, "color:") {
		t.Errorf("color: should be dropped from source: %q", got.source)
	}

	// Rules: fill and stroke dedupe into one (.completed); color is a separate label rule.
	if len(got.rules) != 2 {
		t.Fatalf("expected 2 rules (shape + label), got %d: %+v", len(got.rules), got.rules)
	}
	var shape, label *cssRule
	for i := range got.rules {
		if got.rules[i].label {
			label = &got.rules[i]
		} else {
			shape = &got.rules[i]
		}
	}
	if shape == nil || shape.selector != ".completed" || shape.color != "var(--md-sys-color-primary-container)" {
		t.Errorf("shape rule = %+v", shape)
	}
	if label == nil || label.selector != ".completed .nodeLabel" || label.color != "var(--md-sys-color-on-primary-container)" {
		t.Errorf("label rule = %+v", label)
	}
}

func TestExpandVars_RGBNested(t *testing.T) {
	src := `classDef failed fill:rgb(var(--md-sys-color-error)),color:rgb(var(--md-sys-color-on-error)),stroke:rgb(var(--md-sys-color-error))`
	got := expandVars(src)
	if !strings.Contains(got.source, "fill:currentColor") {
		t.Errorf("nested rgb(var()) not rewritten: %q", got.source)
	}
	if len(got.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(got.rules))
	}
}

func TestExpandVars_RGBAWithComma(t *testing.T) {
	// rgba(var(--x), 0.5) has a comma inside the value; the property splitter
	// must not split on it.
	src := `classDef ghost fill:rgba(var(--md-sys-color-primary), 0.5),stroke:rgba(var(--md-sys-color-primary), 0.5)`
	got := expandVars(src)
	if !strings.Contains(got.source, "fill:currentColor") {
		t.Errorf("rgba(var(), N) not rewritten: %q", got.source)
	}
	// fill and stroke both rewrite; the rules are identical so dedupe to one.
	if len(got.rules) != 1 {
		t.Fatalf("expected 1 rule (deduped), got %d: %+v", len(got.rules), got.rules)
	}
	if got.rules[0].color != "rgba(var(--md-sys-color-primary), 0.5)" {
		t.Errorf("color = %q", got.rules[0].color)
	}
}

func TestExpandVars_MixedVarAndHex(t *testing.T) {
	// Properties without var() must pass through unchanged.
	src := `classDef alert fill:var(--md-sys-color-error),color:#000,stroke-dasharray:4 2`
	got := expandVars(src)
	if !strings.Contains(got.source, "color:#000") {
		t.Errorf("hex color: should be preserved: %q", got.source)
	}
	if !strings.Contains(got.source, "stroke-dasharray:4 2") {
		t.Errorf("stroke-dasharray should be preserved: %q", got.source)
	}
	if !strings.Contains(got.source, "fill:currentColor") {
		t.Errorf("var() fill should be rewritten: %q", got.source)
	}
	if len(got.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(got.rules))
	}
}

func TestExpandVars_StyleDirective(t *testing.T) {
	// `style ID ...` is per-element and bridges via an id selector.
	src := `style fo_s1 fill:var(--md-sys-color-primary),fill-opacity:0.05,stroke:var(--md-sys-color-primary),stroke-dasharray:4 2`
	got := expandVars(src)
	if !strings.Contains(got.source, "fill:currentColor") {
		t.Errorf("style fill not rewritten: %q", got.source)
	}
	if !strings.Contains(got.source, "fill-opacity:0.05") {
		t.Errorf("fill-opacity should be preserved: %q", got.source)
	}
	if len(got.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(got.rules))
	}
	if got.rules[0].selector != "#fo_s1" {
		t.Errorf("style selector should use id form: %q", got.rules[0].selector)
	}
}

func TestExpandVars_LinkStyleDefault(t *testing.T) {
	src := `linkStyle default stroke:var(--md-sys-color-primary)`
	got := expandVars(src)
	if !strings.Contains(got.source, "stroke:currentColor") {
		t.Errorf("linkStyle stroke not rewritten: %q", got.source)
	}
	if len(got.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(got.rules))
	}
	if got.rules[0].selector != ".edgePaths" {
		t.Errorf("linkStyle selector should target .edgePaths, got %q", got.rules[0].selector)
	}
}

func TestExpandVars_LinkStyleNumericPassesThrough(t *testing.T) {
	// Numeric link indices are not bridged (per-edge mapping is brittle).
	src := `linkStyle 0,1 stroke:var(--md-sys-color-error)`
	got := expandVars(src)
	if !strings.Contains(got.source, "var(--md-sys-color-error)") {
		t.Errorf("numeric linkStyle should pass through unchanged: %q", got.source)
	}
	if len(got.rules) != 0 {
		t.Errorf("expected no rules for numeric linkStyle, got %d", len(got.rules))
	}
}

func TestExpandVars_MultipleClassDefsDedupe(t *testing.T) {
	src := `classDef completed fill:var(--md-sys-color-primary),stroke:var(--md-sys-color-primary)
classDef running fill:var(--md-sys-color-primary),stroke:var(--md-sys-color-primary),stroke-dasharray:4 2`
	got := expandVars(src)
	// .completed has fill+stroke (dedupes to one rule), .running same — but
	// .completed and .running are different selectors, so still two rules total.
	if len(got.rules) != 2 {
		t.Fatalf("expected 2 rules across 2 classDefs, got %d: %+v", len(got.rules), got.rules)
	}
	// Verify each selector appears exactly once.
	selectors := map[string]int{}
	for _, r := range got.rules {
		selectors[r.selector]++
	}
	if selectors[".completed"] != 1 || selectors[".running"] != 1 {
		t.Errorf("expected one rule per classDef: %+v", selectors)
	}
}

func TestExpandVars_FormatCSSScoping(t *testing.T) {
	rules := []cssRule{
		{selector: ".completed", color: "var(--p)", label: false},
		{selector: ".completed .nodeLabel", color: "var(--op)", label: true},
	}
	got := formatCSS(rules, "#abc123")
	wants := []string{
		"#abc123 .completed { color: var(--p); }",
		"#abc123 .completed .nodeLabel { color: var(--op) !important; }",
	}
	for _, w := range wants {
		if !strings.Contains(got, w) {
			t.Errorf("missing %q in:\n%s", w, got)
		}
	}
}

func TestExpandVars_NoVarRefsNoSideEffects(t *testing.T) {
	// A line with var(--x) in a node label (not a directive) should pass through.
	src := `s1["set --foo=var(--x)"]`
	got := expandVars(src)
	if got.source != src {
		t.Errorf("non-directive line should pass through unchanged:\nwant %q\n got %q", src, got.source)
	}
	if len(got.rules) != 0 {
		t.Errorf("non-directive var() should not produce rules: %+v", got.rules)
	}
}
