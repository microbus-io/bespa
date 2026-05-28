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

var WHITE_POINT_D65 = []float64{95.047, 100.0, 108.883}

var SRGB_TO_XYZ = [][]float64{
	{0.41233895, 0.35762064, 0.18051042},
	{0.2126, 0.7152, 0.0722},
	{0.01932141, 0.11916382, 0.95034478},
}

var XYZ_TO_SRGB = [][]float64{
	{3.2413774792388685, -1.5376652402851851, -0.49885366846268053},
	{-0.9691452513005321, 1.8758853451067872, 0.04156585616912061},
	{0.05562093689691305, -0.20395524564742123, 1.0571799111220335},
}

// Converts an L* value to a Y value.
// L* in L*a*b* and Y in XYZ measure the same quantity, luminance.
// L* measures perceptual luminance, a linear scale. Y in XYZ measures relative luminance, a
// logarithmic scale.
func YFromLstar(lstar float64) float64 {
	return 100.0 * LabInvf((lstar+16.0)/116.0)
}

func LabInvf(ft float64) float64 {
	e := 216.0 / 24389.0
	kappa := 24389.0 / 27.0
	ft3 := ft * ft * ft
	if ft3 > e {
		return ft3
	} else {
		return (116*ft - 16) / kappa
	}
}

// Linearized linearizes an RGB component.
func Linearized(rgbComponent int) float64 {
	normalized := float64(rgbComponent) / 255.0
	if normalized <= float64(0.040449936) {
		return normalized / 12.92 * 100.0
	} else {
		return math.Pow((normalized+0.055)/1.055, 2.4) * 100.0
	}
}

// Delinearized delinearizes an RGB component.
func Delinearized(rgbComponent float64) int {
	normalized := rgbComponent / 100.0
	delinearized := 0.0
	if normalized <= 0.0031308 {
		delinearized = normalized * 12.92
	} else {
		delinearized = 1.055*math.Pow(normalized, 1.0/2.4) - 0.055
	}
	return ClampInt(0, 255, int(math.Round(delinearized*255.0)))
}

// ArgbFromXyz converts a color from ARGB to XYZ.
func ArgbFromXyz(x, y, z float64) int {
	matrix := XYZ_TO_SRGB
	linearR := matrix[0][0]*x + matrix[0][1]*y + matrix[0][2]*z
	linearG := matrix[1][0]*x + matrix[1][1]*y + matrix[1][2]*z
	linearB := matrix[2][0]*x + matrix[2][1]*y + matrix[2][2]*z
	r := Delinearized(linearR)
	g := Delinearized(linearG)
	b := Delinearized(linearB)
	return ArgbFromRgb(r, g, b)
}

// ArgbFromRgb Converts a color from RGB components to ARGB format.
func ArgbFromRgb(red, green, blue int) int {
	return (255 << 24) | ((red & 255) << 16) | ((green & 255) << 8) | (blue & 255)
}

// ArgbFromLinrgb converts a color from linear RGB components to ARGB format.
func ArgbFromLinrgb(linrgb []float64) int {
	r := Delinearized(linrgb[0])
	g := Delinearized(linrgb[1])
	b := Delinearized(linrgb[2])
	return ArgbFromRgb(r, g, b)
}

// ArgbFromLstar converts an L* value to an ARGB representation.
func ArgbFromLstar(lstar float64) int {
	y := YFromLstar(lstar)
	component := Delinearized(y)
	return ArgbFromRgb(component, component, component)
}

// LstarFromArgb computes the L* value of a color in ARGB representation.
func LstarFromArgb(argb int) float64 {
	y := XyzFromArgb(argb)[1]
	return 116.0*LabF(y/100.0) - 16.0
}

// XyzFromArgb converts a color from XYZ to ARGB
func XyzFromArgb(argb int) []float64 {
	r := Linearized(RedFromArgb(argb))
	g := Linearized(GreenFromArgb(argb))
	b := Linearized(BlueFromArgb(argb))
	return MatrixMultiply([]float64{r, g, b}, SRGB_TO_XYZ)
}

// RedFromArgb returns the red component of a color in ARGB format.
func RedFromArgb(argb int) int {
	return (argb >> 16) & 255
}

// GreenFromArgb returns the green component of a color in ARGB format. */
func GreenFromArgb(argb int) int {
	return (argb >> 8) & 255
}

// BlueFromArgb returns the blue component of a color in ARGB format. */
func BlueFromArgb(argb int) int {
	return argb & 255
}

func LabF(t float64) float64 {
	e := 216.0 / 24389.0
	kappa := 24389.0 / 27.0
	if t > e {
		return math.Pow(t, 1.0/3.0)
	} else {
		return (kappa*t + 16) / 116
	}
}

// LstarFromY converts a Y value to an L* value.
// L* in L*a*b* and Y in XYZ measure the same quantity, luminance.
// L* measures perceptual luminance, a linear scale. Y in XYZ measures relative luminance, a
// logarithmic scale.
func LstarFromY(y float64) float64 {
	return LabF(y/100.0)*116.0 - 16.0
}
