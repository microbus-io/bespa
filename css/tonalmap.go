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
	"embed"
)

var (
	//go:embed tonalmap.bin
	embeddedFS embed.FS
)

// ApproxTone calculates an approximation of the HCT tonal level for a given color.
// It used a pre-calculated tonal map with a granularity of 0x08.
// Tone must be a number between 0 and 100.
func ApproxTone(rgb int, tone int) int {
	if tone <= 0 || tone%10 == 0 || tone == 95 || tone == 99 || tone >= 100 {
		r, g, b := approxTone(rgb, tone)
		return r<<16 + g<<8 + b
	}

	var lower int
	var upper int
	if tone <= 90 {
		lower = tone / 10 * 10
		upper = lower + 10
	} else if tone <= 95 {
		lower = 90
		upper = 95
	} else if tone <= 99 {
		lower = 95
		upper = 99
	}

	r1, g1, b1 := approxTone(rgb, lower)
	r2, g2, b2 := approxTone(rgb, upper)

	r := (r2*(tone-lower) + r1*(upper-tone)) / (upper - lower)
	g := (g2*(tone-lower) + g1*(upper-tone)) / (upper - lower)
	b := (b2*(tone-lower) + b1*(upper-tone)) / (upper - lower)
	return r<<16 + g<<8 + b
}

// approxTone returns an approximation of the HCT tonal level for a given color.
// It used a pre-calculated tonal map with a granularity of 0x08.
// Tone must be 0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 99 or 100.
func approxTone(rgb int, tone int) (rx int, gx int, bx int) {
	if tone <= 0 {
		return 0, 0, 0
	}
	if tone >= 100 {
		return 0xff, 0xff, 0xff
	}

	r := rgb >> 16 & 0xFF
	ri := (r + 0x4) / 0x8
	g := rgb >> 8 & 0xFF
	gi := (g + 0x4) / 0x8
	b := rgb & 0xFF
	bi := (b + 0x4) / 0x8

	pos := ri*33*33 + gi*33 + bi
	pos *= 3 * 11
	if tone <= 90 && tone%10 == 0 {
		// 10, 20, ... 90
		pos += +3 * (tone/10 - 1)
	} else if tone == 95 {
		// 95
		pos += 3 * 9
	} else if tone == 99 {
		// 99
		pos += 3 * 10
	} else {
		return 0, 0, 0 // Invalid tone
	}
	bin, _ := embeddedFS.ReadFile("tonalmap.bin")
	rgbResult := bin[pos : pos+3]
	return int(rgbResult[0]), int(rgbResult[1]), int(rgbResult[2])
}
