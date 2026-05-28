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

package hct

// Hct is hue, chroma, and tone. A color system that provides a perceptually accurate color
// measurement system that can also accurately render what colors will appear as in different
// lighting environments.
type Hct struct {
	hue    float64
	chroma float64
	tone   float64
	argb   int
}

func From(hue, chroma, tone float64) *Hct {
	argb := SolveToInt(hue, chroma, tone)
	hct := &Hct{}
	hct.setInternalState(argb)
	return hct
}

// FromInteger creates an HCT color from a color.
func FromInteger(argb int) *Hct {
	hct := &Hct{}
	hct.setInternalState(argb)
	return hct
}

func (hct *Hct) setInternalState(argb int) {
	hct.argb = argb
	cam := FromInt(argb)
	hct.hue = cam.Hue()
	hct.chroma = cam.Chroma()
	hct.tone = LstarFromArgb(argb)
}

// Hue returns the hue component of the color. 0 <= hue < 360
func (hct *Hct) Hue() float64 {
	return hct.hue
}

// Chroma returns the chroma component of the color. 0 <= chroma < ?
func (hct *Hct) Chroma() float64 {
	return hct.chroma
}

// Tone returns the tone component of the color. 0 <= tone <= 100
func (hct *Hct) Tone() float64 {
	return hct.tone
}

// ToInt returns the ARGB representation of the color.
func (hct *Hct) ToInt() int {
	return hct.argb
}

// SetHue sets the hue of this color. Chroma may decrease because chroma has a different maximum for any
// given hue and tone. 0 <= newHue < 360
func (hct *Hct) SetHue(newHue float64) {
	hct.setInternalState(SolveToInt(newHue, hct.chroma, hct.tone))
}

// SetChroma sets the chroma of this color. Chroma may decrease because chroma has a different maximum for
// any given hue and tone. 0 <= newChroma < ?
func (hct *Hct) SetChroma(newChroma float64) {
	hct.setInternalState(SolveToInt(hct.hue, newChroma, hct.tone))
}

// SetTone sets the tone of this color. Chroma may decrease because chroma has a different maximum for any
// given hue and tone. 0 <= newTone <= 100
func (hct *Hct) SetTone(newTone float64) {
	hct.setInternalState(SolveToInt(hct.hue, hct.chroma, newTone))
}

// Translate a color into different ViewingConditions.
//
// Colors change appearance. They look different with lights on versus off, the same color, as
// in hex code, on white looks different when on black. This is called color relativity, most
// famously explicated by Josef Albers in Interaction of Color.
//
// In color science, color appearance models can account for this and calculate the appearance
// of a color in different settings. HCT is based on CAM16, a color appearance model, and uses it
// to make these calculations.
//
// See ViewingConditions.make for parameters affecting color appearance.
func (hct *Hct) InViewingConditions(vc *ViewingConditions) *Hct {
	// 1. Use CAM16 to find XYZ coordinates of color in specified VC.
	cam16 := FromInt(hct.ToInt())
	viewedInVc := cam16.XyzInViewingConditions(vc)

	// 2. Create CAM16 of those XYZ coordinates in default VC.
	recastInVc := FromXyzInViewingConditions(viewedInVc[0], viewedInVc[1], viewedInVc[2], DEFAULT_VIEWING_CONDITIONS)

	// 3. Create HCT from:
	// - CAM16 using default VC with XYZ coordinates in specified VC.
	// - L* converted from Y in XYZ coordinates in specified VC.
	return From(recastInVc.Hue(), recastInVc.Chroma(), LstarFromY(viewedInVc[1]))
}
