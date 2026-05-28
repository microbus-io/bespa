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
	"math"
	"strconv"
	"strings"
)

// Color is an immutable structure that holds an ARGB color value.
type Color int

// NewColor creates a color from its ARGB components.
func NewColor(a, r, g, b int) Color {
	a = a & 0xff
	r = r & 0xff
	g = g & 0xff
	b = b & 0xff
	return Color((a << 24) + (r << 16) + (g << 8) + b)
}

// ParseColor parses a color string in the format "#FF0000", "FF0000" or "rgba(255,0,0,1.0)".
func ParseColor(spec string) Color {
	var r, g, b int
	var a float64
	if strings.HasPrefix(spec, "rgb(") {
		values := strings.Split(spec[4:len(spec)-1], ",")
		if len(values) == 4 {
			r, _ = strconv.Atoi(values[0])
			g, _ = strconv.Atoi(values[1])
			b, _ = strconv.Atoi(values[2])
			a, _ = strconv.ParseFloat(values[3], 64)
			a *= 255.0
			if a > 255.0 {
				a = 255.0
			}
			if a < 0 {
				a = 0
			}
		}
	} else {
		spec = strings.TrimPrefix(spec, "#")
		rgb, _ := strconv.ParseInt(spec, 16, 0)
		b = int(rgb) & 0xFF
		g = (int(rgb) >> 8) & 0xFF
		r = (int(rgb) >> 16) & 0xFF
		a = 255.0
	}
	return NewColor(int(a), r, g, b)
}

// String prints the color as #FF0000, or rgba(255,0,0,1.0) if it has an alpha channel.
func (color Color) String() string {
	c := int(color)
	if c&0xff000000 == 0xff000000 { // 100% opaque
		return fmt.Sprintf("#%06X", c&0xffffff)
	} else {
		a, r, g, b := color.Channels()
		return fmt.Sprintf("rgba(%d,%d,%d,%.2f)", r, g, b, float64(a)/255.0)
	}
}

// Channels returns the ARGB values of the color.
func (color Color) Channels() (a, r, g, b int) {
	c := int(color)
	return c >> 24 & 0xff,
		(c >> 16) & 0xff,
		(c >> 8) & 0xff,
		c & 0xff
}

// Darken darkens the color. Each unit is approx 5% decrease in brightness.
func (color Color) Darken(units int) Color {
	amount := units * 12
	if amount == 0 {
		return color
	}
	a, r, g, b := color.Channels()
	r -= amount
	g -= amount
	b -= amount
	return NewColor(a, r, g, b)
}

// Lighten lightens the color. Each unit is approx 5% increase in brightness.
func (color Color) Lighten(units int) Color {
	amount := units * 12
	if amount == 0 {
		return color
	}
	a, r, g, b := color.Channels()
	r += amount
	g += amount
	b += amount
	return NewColor(a, r, g, b)
}

// IsBright indicates if the color is bright.
func (color Color) IsBright() bool {
	_, r, g, b := color.Channels()
	return (r+g+b >= 96*3)
}

// IsDark indicates if the color is dark.
func (color Color) IsDark() bool {
	return !color.IsBright()
}

// Hover adjusts the color to highlight an element that is hovered.
func (color Color) Hover() Color {
	if color.IsBright() {
		return color.Darken(3)
	}
	return color.Lighten(2)
}

// Hover adjusts the color to highlight an element that is hovered.
func (color Color) HoverActive() Color {
	if color.IsBright() {
		return color.Darken(4)
	}
	return color.Lighten(3)
}

// Rotate rotates the hue of the color by a number of degrees (of 360).
func (color Color) Rotate(degrees int) Color {
	if degrees%360 == 0 {
		return color
	}
	a, r, g, b := color.Channels()
	h, s, l := rgb2hsl(r, g, b)
	h += (256 * degrees / 360)
	h = h % 360
	rr, gg, bb := hsl2rgb(h, s, l)
	return NewColor(a, rr, gg, bb)
}

// Mix creates a new color that is a weighted average of two colors.
// The percentage gives weight to the mixin color.
// A percentage greater than 50 gives more weight to the mixin color.
// A percentage smaller than 50 gives more weight to the original color.
// A percentage of 50 averages the two colors equally.
func (color Color) Mix(mixin Color, mixinWeightPct int) Color {
	pct := mixinWeightPct
	if pct <= 0 {
		return color
	}
	if pct >= 100 {
		return mixin
	}
	a1, r1, g1, b1 := color.Channels()
	a2, r2, g2, b2 := mixin.Channels()
	return NewColor(
		(a2*pct+a1*(100-pct))/100,
		(r2*pct+r1*(100-pct))/100,
		(g2*pct+g1*(100-pct))/100,
		(b2*pct+b1*(100-pct))/100,
	)
}

// Tone returns an approximation of the color at a specified HCT tone.
// Tone 0 is black and tone 100 is white. 40 is considered the base tone.
func (color Color) Tone(tone int) Color {
	return Color(0xff000000 + ApproxTone(int(color), tone))
	// h := hct.FromInteger(int(color) | 0xff000000) // Remove transparency
	// h.SetTone(float64(tone))
	// return Color(h.ToInt())
}

// rgb2hsl returns the equivalent HSL values coresponding to the color.
func rgb2hsl(ri, gi, bi int) (int, int, int) {
	var h, s, l float64

	r := float64(ri) / 255.0
	g := float64(gi) / 255.0
	b := float64(bi) / 255.0

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)

	// Luminosity is the average of the max and min rgb color intensities.
	l = (max + min) / 2

	// saturation
	delta := max - min
	if delta == 0 {
		// it's gray
		return 0, 0, int(l * 255.0)
	}

	// it's not gray
	if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}

	// hue
	r2 := (((max - r) / 6) + (delta / 2)) / delta
	g2 := (((max - g) / 6) + (delta / 2)) / delta
	b2 := (((max - b) / 6) + (delta / 2)) / delta
	switch {
	case r == max:
		h = b2 - g2
	case g == max:
		h = (1.0 / 3.0) + r2 - b2
	case b == max:
		h = (2.0 / 3.0) + g2 - r2
	}

	// fix wraparounds
	switch {
	case h < 0:
		h += 1
	case h > 1:
		h -= 1
	}

	return int(h * 255.0), int(s * 255.0), int(l * 255.0)
}

func hsl2rgb(hi, si, li int) (int, int, int) {
	h := float64(hi) / 255.0
	s := float64(si) / 255.0
	l := float64(li) / 255.0

	if s == 0 {
		// it's gray
		return int(l * 255.0), int(l * 255.0), int(l * 255.0)
	}

	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = (l + s) - (s * l)
	}

	v1 = 2*l - v2

	hueToRGB := func(v1, v2, h float64) float64 {
		if h < 0 {
			h += 1
		}
		if h > 1 {
			h -= 1
		}
		switch {
		case 6*h < 1:
			return (v1 + (v2-v1)*6*h)
		case 2*h < 1:
			return v2
		case 3*h < 2:
			return v1 + (v2-v1)*((2.0/3.0)-h)*6
		}
		return v1
	}
	r := hueToRGB(v1, v2, h+(1.0/3.0))
	g := hueToRGB(v1, v2, h)
	b := hueToRGB(v1, v2, h-(1.0/3.0))

	return int(r * 255.0), int(g * 255.0), int(b * 255.0)
}
