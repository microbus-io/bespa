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

import "math"

/*
In traditional color spaces, a color can be identified solely by the observer's measurement of
the color. Color appearance models such as CAM16 also use information about the environment where
the color was observed, known as the viewing conditions.

For example, white under the traditional assumption of a midday sun white point is accurately
measured as a slightly chromatic blue by CAM16. (roughly, hue 203, chroma 3, lightness 100)

This class caches intermediate values of the CAM16 conversion process that depend only on
viewing conditions, enabling speed ups.
*/
type ViewingConditions struct {
	aw     float64
	nbb    float64
	ncb    float64
	c      float64
	nc     float64
	n      float64
	rgbD   []float64
	fl     float64
	flRoot float64
	z      float64
}

func (vc *ViewingConditions) Aw() float64 {
	return vc.aw
}
func (vc *ViewingConditions) N() float64 {
	return vc.n
}
func (vc *ViewingConditions) Nbb() float64 {
	return vc.nbb
}
func (vc *ViewingConditions) Ncb() float64 {
	return vc.ncb
}
func (vc *ViewingConditions) C() float64 {
	return vc.c
}
func (vc *ViewingConditions) Nc() float64 {
	return vc.nc
}
func (vc *ViewingConditions) RgbD() []float64 {
	return vc.rgbD
}
func (vc *ViewingConditions) Fl() float64 {
	return vc.fl
}
func (vc *ViewingConditions) FlRoot() float64 {
	return vc.flRoot
}
func (vc *ViewingConditions) Z() float64 {
	return vc.z
}

// sRGB-like viewing conditions.
var DEFAULT_VIEWING_CONDITIONS = DefaultWithBackgroundLstar(50.0)

/*
Create ViewingConditions from a simple, physically relevant, set of parameters.

White point, measured in the XYZ color space. default = D65, or sunny day afternoon.

The luminance of the adapting field. Informally, how bright it is in
the room where the color is viewed. Can be calculated from lux by multiplying lux by
0.0586. default = 11.72, or 200 lux.

The lightness of the area surrounding the color. measured by L* in L*a*b*. default = 50.0

A general description of the lighting surrounding the color. 0 is pitch dark,
like watching a movie in a theater. 1.0 is a dimly light room, like watching TV at home at
night. 2.0 means there is no difference between the lighting on the color and around it.
default = 2.0

Whether the eye accounts for the tint of the ambient lighting,
such as knowing an apple is still red in green light. default = false, the eye does not
perform this process on self-luminous objects like displays.
*/
func Make(whitePoint []float64, adaptingLuminance, backgroundLstar, surround float64, discountingIlluminant bool) *ViewingConditions {
	// A background of pure black is non-physical and leads to infinities that represent the idea
	// that any color viewed in pure black can't be seen.
	backgroundLstar = math.Max(0.1, backgroundLstar)

	// Transform white point XYZ to 'cone'/'rgb' responses
	matrix := XYZ_TO_CAM16RGB
	xyz := whitePoint
	rW := (xyz[0] * matrix[0][0]) + (xyz[1] * matrix[0][1]) + (xyz[2] * matrix[0][2])
	gW := (xyz[0] * matrix[1][0]) + (xyz[1] * matrix[1][1]) + (xyz[2] * matrix[1][2])
	bW := (xyz[0] * matrix[2][0]) + (xyz[1] * matrix[2][1]) + (xyz[2] * matrix[2][2])
	f := 0.8 + (surround / 10.0)
	var c float64
	if f >= 0.9 {
		c = Lerp(0.59, 0.69, ((f - 0.9) * 10.0))
	} else {
		c = Lerp(0.525, 0.59, ((f - 0.8) * 10.0))
	}
	var d float64
	if discountingIlluminant {
		d = 1.0
	} else {
		d = f * (1.0 - ((1.0 / 3.6) * math.Exp((-adaptingLuminance-42.0)/92.0)))
	}
	d = ClampDouble(0.0, 1.0, d)
	nc := f
	rgbD := []float64{d*(100.0/rW) + 1.0 - d, d*(100.0/gW) + 1.0 - d, d*(100.0/bW) + 1.0 - d}
	k := 1.0 / (5.0*adaptingLuminance + 1.0)
	k4 := k * k * k * k
	k4F := 1.0 - k4
	fl := (k4 * adaptingLuminance) + (0.1 * k4F * k4F * math.Cbrt(5.0*adaptingLuminance))
	n := YFromLstar(backgroundLstar) / whitePoint[1]
	z := 1.48 + math.Sqrt(n)
	nbb := 0.725 / math.Pow(n, 0.2)
	ncb := nbb
	rgbAFactors := []float64{
		math.Pow(fl*rgbD[0]*rW/100.0, 0.42),
		math.Pow(fl*rgbD[1]*gW/100.0, 0.42),
		math.Pow(fl*rgbD[2]*bW/100.0, 0.42),
	}

	rgbA := []float64{
		(400.0 * rgbAFactors[0]) / (rgbAFactors[0] + 27.13),
		(400.0 * rgbAFactors[1]) / (rgbAFactors[1] + 27.13),
		(400.0 * rgbAFactors[2]) / (rgbAFactors[2] + 27.13),
	}

	aw := ((2.0 * rgbA[0]) + rgbA[1] + (0.05 * rgbA[2])) * nbb
	return &ViewingConditions{n: n, aw: aw, nbb: nbb, ncb: ncb, c: c, nc: nc, rgbD: rgbD, fl: fl, flRoot: math.Pow(fl, 0.25), z: z}
}

// Create sRGB-like viewing conditions with a custom background lstar.
// Default viewing conditions have a lstar of 50, midgray.
func DefaultWithBackgroundLstar(lstar float64) *ViewingConditions {
	return Make(
		WHITE_POINT_D65,
		(200.0 / math.Pi * YFromLstar(50.0) / 100.0),
		lstar,
		2.0,
		false)
}
