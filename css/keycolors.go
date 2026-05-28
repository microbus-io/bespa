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

package css

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// https://m3.material.io/theme-builder#/dynamic

var PresetKeyColors = []KeyColors{
	{
		Name:           "Green",
		Primary:        ParseColor("#426915"),
		Secondary:      ParseColor("#57624a"),
		Tertiary:       ParseColor("#386663"),
		Neutral:        ParseColor("#5e5f59"),
		NeutralVariant: ParseColor("#5c6054"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Blue",
		Primary:        ParseColor("#215fa6"),
		Secondary:      ParseColor("#555f71"),
		Tertiary:       ParseColor("#6e5676"),
		Neutral:        ParseColor("#5e5e62"),
		NeutralVariant: ParseColor("#5b5e66"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Yellow",
		Primary:        ParseColor("#636100"),
		Secondary:      ParseColor("#616042"),
		Tertiary:       ParseColor("#3e6656"),
		Neutral:        ParseColor("#605e58"),
		NeutralVariant: ParseColor("#605f51"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Red",
		Primary:        ParseColor("#a03f28"),
		Secondary:      ParseColor("#77574f"),
		Tertiary:       ParseColor("#6d5d2e"),
		Neutral:        ParseColor("#655c5a"),
		NeutralVariant: ParseColor("#6c5b57"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Violet",
		Primary:        ParseColor("#715573"),
		Secondary:      ParseColor("#005ac1"),
		Tertiary:       ParseColor("#575e71"),
		Neutral:        ParseColor("#5e5e62"),
		NeutralVariant: ParseColor("#5c5e66"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Olive",
		Primary:        ParseColor("#705d00"),
		Secondary:      ParseColor("#675e40"),
		Tertiary:       ParseColor("#44664e"),
		Neutral:        ParseColor("#615e57"),
		NeutralVariant: ParseColor("#635e50"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Meadow",
		Primary:        ParseColor("#006496"),
		Secondary:      ParseColor("#50606f"),
		Tertiary:       ParseColor("#66587b"),
		Neutral:        ParseColor("#5d5e61"),
		NeutralVariant: ParseColor("#595f65"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Roses",
		Primary:        ParseColor("#9b4429"),
		Secondary:      ParseColor("#77574d"),
		Tertiary:       ParseColor("#6b5e2f"),
		Neutral:        ParseColor("#655c5a"),
		NeutralVariant: ParseColor("#6c5b56"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Byzantine",
		Primary:        ParseColor("#8e437e"),
		Secondary:      ParseColor("#6f5767"),
		Tertiary:       ParseColor("#815341"),
		Neutral:        ParseColor("#635d60"),
		NeutralVariant: ParseColor("#675b62"),
		Error:          ParseColor("#ba1a1a"),
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "New Mexico",
		Primary:        ParseColor("#705d00"), // Gold
		Secondary:      ParseColor("#675e40"),
		Tertiary:       ParseColor("#44664e"),
		Neutral:        ParseColor("#615e57"),
		NeutralVariant: ParseColor("#635e50"),
		Error:          ParseColor("#bd082f"), // Red
		OK:             ParseColor("#1aba1a"),
	}, {
		Name:           "Miami Dolphins",
		Primary:        ParseColor("#006970"), // Aqua
		Secondary:      ParseColor("#00658b"), // Blue
		Tertiary:       ParseColor("#4f5e7d"),
		Neutral:        ParseColor("#5d5f5a"),
		NeutralVariant: ParseColor("#5b6056"),
		Error:          ParseColor("#ae3100"), // Orange
		OK:             ParseColor("#1aba1a"),
	},
}

var DefaultKeyColors = PresetKeyColors[0]

// KeyColors are the Material key colors of the color scheme from which
// all the color tones are derived.
// https://m3.material.io/styles/color/the-color-system/key-colors-tones .
type KeyColors struct {
	Name           string
	Primary        Color
	Secondary      Color
	Tertiary       Color
	Neutral        Color
	NeutralVariant Color
	Error          Color
	OK             Color
}

// String serializes the key colors to a string.
func (kc KeyColors) String() string {
	var s strings.Builder
	s.WriteString(kc.Name)
	s.WriteString(";")
	s.WriteString(kc.Primary.String())
	s.WriteString(";")
	s.WriteString(kc.Secondary.String())
	s.WriteString(";")
	s.WriteString(kc.Tertiary.String())
	s.WriteString(";")
	s.WriteString(kc.Neutral.String())
	s.WriteString(";")
	s.WriteString(kc.NeutralVariant.String())
	s.WriteString(";")
	s.WriteString(kc.Error.String())
	s.WriteString(";")
	s.WriteString(kc.OK.String())
	return s.String()
}

// KeyColorsFromString deserializes the key colors from a string.
func KeyColorsFromString(s string) KeyColors {
	var kc KeyColors
	parts := strings.Split(s, ";")
	if len(parts) == 8 {
		kc.Name = parts[0]
		kc.Primary = ParseColor(parts[1])
		kc.Secondary = ParseColor(parts[2])
		kc.Tertiary = ParseColor(parts[3])
		kc.Neutral = ParseColor(parts[4])
		kc.NeutralVariant = ParseColor(parts[5])
		kc.Error = ParseColor(parts[6])
		kc.OK = ParseColor(parts[7])
	}
	return kc
}

// WriteCSSTones writes the 13 Material tones of the key color as CSS variables.
func (kc KeyColors) WriteCSSTones(w io.Writer) error {
	w.Write([]byte("/* Key color tones */\n"))
	w.Write([]byte(":root {\n"))
	kc.writeTones(w, kc.Primary, "md-ref-palette-primary")
	kc.writeTones(w, kc.Secondary, "md-ref-palette-secondary")
	kc.writeTones(w, kc.Tertiary, "md-ref-palette-tertiary")
	kc.writeTones(w, kc.Neutral, "md-ref-palette-neutral")
	kc.writeTones(w, kc.NeutralVariant, "md-ref-palette-neutral-variant")
	kc.writeTones(w, kc.Error, "md-ref-palette-error")
	kc.writeTones(w, kc.OK, "md-ref-palette-ok")

	// Rotated colors
	for deg := 0; deg < 360; deg += 30 {
		d := strconv.Itoa(deg) + "deg"
		kc.writeTones(w, kc.Primary.Rotate(deg), "md-ref-palette-primary-"+d)
		kc.writeTones(w, kc.Secondary.Rotate(deg), "md-ref-palette-secondary-"+d)
		kc.writeTones(w, kc.Tertiary.Rotate(deg), "md-ref-palette-tertiary-"+d)
	}

	w.Write([]byte("}\n\n"))
	return nil
}

// WriteCSSThemes writes the Material variables of the light and dark themes.
func (kc KeyColors) WriteCSSThemes(w io.Writer) error {
	w.Write([]byte("/* Elevation and opacities */\n"))
	w.Write([]byte(":root {\n"))
	// https://m3.material.io/styles/color/the-color-system/color-roles
	// https://m3.material.io/styles/elevation/applying-elevation
	w.Write([]byte("\t--md-sys-elevation-level1-tint-layer-opacity: .05;\n"))
	w.Write([]byte("\t--md-sys-elevation-level2-tint-layer-opacity: .08;\n"))
	w.Write([]byte("\t--md-sys-elevation-level3-tint-layer-opacity: .11;\n"))
	w.Write([]byte("\t--md-sys-elevation-level4-tint-layer-opacity: .12;\n"))
	w.Write([]byte("\t--md-sys-elevation-level5-tint-layer-opacity: .14;\n"))

	// https://m3.material.io/styles/elevation/tokens
	w.Write([]byte("\t--md-sys-elevation-level0: 0;\n"))
	w.Write([]byte("\t--md-sys-elevation-level1: 1px;\n"))
	w.Write([]byte("\t--md-sys-elevation-level2: 3px;\n"))
	w.Write([]byte("\t--md-sys-elevation-level3: 6px;\n"))
	w.Write([]byte("\t--md-sys-elevation-level4: 8px;\n"))
	w.Write([]byte("\t--md-sys-elevation-level5: 12px;\n"))

	w.Write([]byte("\t--md-sys-elevation-level0-shadow: 0;\n"))
	w.Write([]byte("\t--md-sys-elevation-level1-shadow: 0 2px calc(var(--md-sys-elevation-level1) * 2 + 4px) calc(var(--md-sys-elevation-level1)) rgba(var(--md-sys-color-shadow), .35);\n"))
	w.Write([]byte("\t--md-sys-elevation-level2-shadow: 0 2px calc(var(--md-sys-elevation-level2) * 2 + 4px) calc(var(--md-sys-elevation-level2)) rgba(var(--md-sys-color-shadow), .35);\n"))
	w.Write([]byte("\t--md-sys-elevation-level3-shadow: 0 2px calc(var(--md-sys-elevation-level3) * 2 + 4px) calc(var(--md-sys-elevation-level3)) rgba(var(--md-sys-color-shadow), .35);\n"))
	w.Write([]byte("\t--md-sys-elevation-level4-shadow: 0 2px calc(var(--md-sys-elevation-level4) * 2 + 4px) calc(var(--md-sys-elevation-level4)) rgba(var(--md-sys-color-shadow), .35);\n"))
	w.Write([]byte("\t--md-sys-elevation-level5-shadow: 0 2px calc(var(--md-sys-elevation-level5) * 2 + 4px) calc(var(--md-sys-elevation-level5)) rgba(var(--md-sys-color-shadow), .35);\n"))

	// https://m3.material.io/foundations/interaction-states
	w.Write([]byte("\t--md-sys-state-hover-state-layer-opacity: .08;\n"))
	w.Write([]byte("\t--md-sys-state-focus-state-layer-opacity: .12;\n"))
	w.Write([]byte("\t--md-sys-state-pressed-state-layer-opacity: .12;\n"))
	w.Write([]byte("\t--md-sys-state-dragged-state-layer-opacity: .16;\n"))
	w.Write([]byte("\t--md-sys-state-disabled-state-layer-opacity: .38;\n"))
	w.Write([]byte("\t--md-sys-state-disabled-container-state-layer-opacity: .12;\n"))

	w.Write([]byte("}\n\n"))

	// https://m3.material.io/styles/color/the-color-system/tokens#7fd4440e-986d-443f-8b3a-4933bff16646
	tokens := []string{
		"Primary", "md-sys-color-primary", "md-ref-palette-primary40", "md-ref-palette-primary80",
		"Primary container", "md-sys-color-primary-container", "md-ref-palette-primary90", "md-ref-palette-primary30",
		"Secondary", "md-sys-color-secondary", "md-ref-palette-secondary40", "md-ref-palette-secondary80",
		"Secondary container", "md-sys-color-secondary-container", "md-ref-palette-secondary90", "md-ref-palette-secondary30",
		"Tertiary", "md-sys-color-tertiary", "md-ref-palette-tertiary40", "md-ref-palette-tertiary80",
		"Tertiary container", "md-sys-color-tertiary-container", "md-ref-palette-tertiary90", "md-ref-palette-tertiary30",
		"Surface", "md-sys-color-surface", "md-ref-palette-neutral99", "md-ref-palette-neutral10",
		"Surface dim", "md-sys-color-surface-dim", "md-ref-palette-neutral87", "md-ref-palette-neutral6",
		"Surface bright", "md-sys-color-surface-bright", "md-ref-palette-neutral98", "md-ref-palette-neutral24",
		"Surface container lowest", "md-sys-color-surface-container-lowest", "md-ref-palette-neutral100", "md-ref-palette-neutral4",
		"Surface container low", "md-sys-color-surface-container-low", "md-ref-palette-neutral96", "md-ref-palette-neutral10",
		"Surface container", "md-sys-color-surface-container", "md-ref-palette-neutral94", "md-ref-palette-neutral12",
		"Surface container high", "md-sys-color-surface-container-high", "md-ref-palette-neutral92", "md-ref-palette-neutral17",
		"Surface container highest", "md-sys-color-surface-container-highest", "md-ref-palette-neutral90", "md-ref-palette-neutral22",
		"Surface variant", "md-sys-color-surface-variant", "md-ref-palette-neutral-variant90", "md-ref-palette-neutral-variant30",
		"Background", "md-sys-color-background", "md-ref-palette-neutral98", "md-ref-palette-neutral6",
		"Error", "md-sys-color-error", "md-ref-palette-error40", "md-ref-palette-error80",
		"Error container", "md-sys-color-error-container", "md-ref-palette-error90", "md-ref-palette-error30",
		"On primary", "md-sys-color-on-primary", "md-ref-palette-primary100", "md-ref-palette-primary20",
		"On primary container", "md-sys-color-on-primary-container", "md-ref-palette-primary10", "md-ref-palette-primary90",
		"On secondary", "md-sys-color-on-secondary", "md-ref-palette-secondary100", "md-ref-palette-secondary20",
		"On secondary container", "md-sys-color-on-secondary-container", "md-ref-palette-secondary10", "md-ref-palette-secondary90",
		"On tertiary", "md-sys-color-on-tertiary", "md-ref-palette-tertiary100", "md-ref-palette-tertiary20",
		"On tertiary container", "md-sys-color-on-tertiary-container", "md-ref-palette-tertiary10", "md-ref-palette-tertiary90",
		"On surface", "md-sys-color-on-surface", "md-ref-palette-neutral10", "md-ref-palette-neutral90",
		"On surface variant", "md-sys-color-on-surface-variant", "md-ref-palette-neutral-variant30", "md-ref-palette-neutral-variant80",
		"On error", "md-sys-color-on-error", "md-ref-palette-error100", "md-ref-palette-error20",
		"On error container", "md-sys-color-on-error-container", "md-ref-palette-error10", "md-ref-palette-error90",
		"On background", "md-sys-color-on-background", "md-ref-palette-neutral10", "md-ref-palette-neutral90",
		"Outline", "md-sys-color-outline", "md-ref-palette-neutral-variant50", "md-ref-palette-neutral-variant60",
		"Outline variant", "md-sys-color-outline-variant", "md-ref-palette-neutral-variant80", "md-ref-palette-neutral-variant30",
		"Shadow", "md-sys-color-shadow", "md-ref-palette-neutral0", "md-ref-palette-neutral0",
		"Surface tint", "md-sys-color-surface-tint", "md-sys-color-primary", "md-sys-color-primary",
		"Inverse surface", "md-sys-color-inverse-surface", "md-ref-palette-neutral20", "md-ref-palette-neutral90",
		"Inverse on surface", "md-sys-color-inverse-on-surface", "md-ref-palette-neutral95", "md-ref-palette-neutral20",
		"Inverse primary", "md-sys-color-inverse-primary", "md-ref-palette-primary80", "md-ref-palette-primary40",
		"Scrim", "md-sys-color-scrim", "md-ref-palette-neutral0", "md-ref-palette-neutral0",

		// OK color is green
		"OK", "md-sys-color-ok", "md-ref-palette-ok40", "md-ref-palette-ok80",
		"OK container", "md-sys-color-ok-container", "md-ref-palette-ok90", "md-ref-palette-ok30",
		"On OK", "md-sys-color-on-ok", "md-ref-palette-ok100", "md-ref-palette-ok20",
		"On OK container", "md-sys-color-on-ok-container", "md-ref-palette-ok10", "md-ref-palette-ok90",
	}

	// Rotated colors
	for deg := 0; deg < 360; deg += 30 {
		d := strconv.Itoa(deg) + "deg"
		t := []string{
			"Primary", "md-sys-color-primary", "md-ref-palette-primary40", "md-ref-palette-primary80",
			"Primary container", "md-sys-color-primary-container", "md-ref-palette-primary90", "md-ref-palette-primary30",
			"On primary", "md-sys-color-on-primary", "md-ref-palette-primary100", "md-ref-palette-primary20",
			"On primary container", "md-sys-color-on-primary-container", "md-ref-palette-primary10", "md-ref-palette-primary90",
		}
		for i := range t {
			t[i] = strings.ReplaceAll(t[i], "Primary", "Primary "+d)
			t[i] = strings.ReplaceAll(t[i], "On primary", "On primary "+d)
			t[i] = strings.ReplaceAll(t[i], "-primary", "-primary-"+d)
		}
		tokens = append(tokens, t...)
		for i := range t {
			t[i] = strings.ReplaceAll(t[i], "Primary", "Secondary")
			t[i] = strings.ReplaceAll(t[i], "primary", "secondary")
		}
		tokens = append(tokens, t...)
		for i := range t {
			t[i] = strings.ReplaceAll(t[i], "Secondary", "Tertiary")
			t[i] = strings.ReplaceAll(t[i], "secondary", "tertiary")
		}
		tokens = append(tokens, t...)
	}

	w.Write([]byte("/* Light theme tokens */\n"))
	w.Write([]byte("@media (prefers-color-scheme: light) {\n"))
	w.Write([]byte(":root {\n"))
	for i := 0; i < len(tokens); i += 4 {
		rule := fmt.Sprintf("\t--%s: var(--%s);\n", tokens[i+1], tokens[i+2])
		_, err := w.Write([]byte(rule))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("}\n"))
	w.Write([]byte("}\n"))
	w.Write([]byte("HTML.LightTheme {\n"))
	for i := 0; i < len(tokens); i += 4 {
		rule := fmt.Sprintf("\t--%s: var(--%s);\n", tokens[i+1], tokens[i+2])
		_, err := w.Write([]byte(rule))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("}\n\n"))

	w.Write([]byte("/* Dark theme tokens */\n"))
	w.Write([]byte("@media (prefers-color-scheme: dark) {\n"))
	w.Write([]byte(":root {\n"))
	for i := 0; i < len(tokens); i += 4 {
		rule := fmt.Sprintf("\t--%s: var(--%s);\n", tokens[i+1], tokens[i+3])
		_, err := w.Write([]byte(rule))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("}\n"))
	w.Write([]byte("}\n"))
	w.Write([]byte("HTML.DarkTheme {\n"))
	for i := 0; i < len(tokens); i += 4 {
		rule := fmt.Sprintf("\t--%s: var(--%s);\n", tokens[i+1], tokens[i+3])
		_, err := w.Write([]byte(rule))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("}\n\n"))

	return nil
}

// writeTones writes the 13 Material tones of the color as CSS variables.
func (kc KeyColors) writeTones(w io.Writer, color Color, name string) error {
	tones := []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 99, 100}
	if name == "md-ref-palette-neutral" {
		tones = []int{0, 4, 6, 10, 12, 17, 20, 22, 24, 30, 40, 50, 60, 70, 80, 87, 90, 92, 94, 95, 96, 98, 99, 100}
	}
	for _, tone := range tones {
		toned := color.Tone(tone)
		_, r, g, b := toned.Channels()
		rule := fmt.Sprintf("\t--%s%d: %d,%d,%d;\n", name, tone, r, g, b)
		_, err := w.Write([]byte(rule))
		if err != nil {
			return err
		}
	}
	return nil
}
