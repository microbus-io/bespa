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

// XYZ_TO_CAM16RGB transforms XYZ color space coordinates to 'cone'/'RGB' responses in CAM16.
var XYZ_TO_CAM16RGB = [][]float64{
	{0.401288, 0.650173, -0.051461},
	{-0.250268, 1.204414, 0.045854},
	{-0.002079, 0.048952, 0.953127},
}

// CAM16RGB_TO_XYZ transforms 'cone'/'RGB' responses in CAM16 to XYZ color space coordinates.
var CAM16RGB_TO_XYZ = [][]float64{
	{1.8620678, -1.0112547, 0.14918678},
	{0.38752654, 0.62144744, -0.00897398},
	{-0.01584150, -0.03412294, 1.0499644},
}

/**
 * Cam16, a color appearance model. Colors are not just defined by their hex code, but rather, a hex
 * code and viewing conditions.
 *
 * <p>CAM16 instances also have coordinates in the CAM16-UCS space, called J*, a*, b*, or jstar,
 * astar, bstar in code. CAM16-UCS is included in the CAM16 specification, and should be used when
 * measuring distances between colors.
 *
 * <p>In traditional color spaces, a color can be identified solely by the observer's measurement of
 * the color. Color appearance models such as CAM16 also use information about the environment where
 * the color was observed, known as the viewing conditions.
 *
 * <p>For example, white under the traditional assumption of a midday sun white point is accurately
 * measured as a slightly chromatic blue by CAM16. (roughly, hue 203, chroma 3, lightness 100)
 */
type Cam16 struct {
	// CAM16 color dimensions, see getters for documentation.
	hue    float64
	chroma float64
	j      float64
	q      float64
	m      float64
	s      float64

	// Coordinates in UCS space. Used to determine color distance, like delta E equations in L*a*b*.
	jstar float64
	astar float64
	bstar float64
}

// Hue in CAM16.
func (c *Cam16) Hue() float64 {
	return c.hue
}

// Chroma in CAM16.
func (c *Cam16) Chroma() float64 {
	return c.chroma
}

// J is lightness in CAM16.
func (c *Cam16) J() float64 {
	return c.j
}

// Q is brightness in CAM16.
// Prefer lightness, brightness is an absolute quantity. For example, a sheet of white paper is
// much brighter viewed in sunlight than in indoor light, but it is the lightest object under any
// lighting.
func (c *Cam16) Q() float64 {
	return c.q
}

// M is colorfulness in CAM16.
// Prefer chroma, colorfulness is an absolute quantity. For example, a yellow toy car is much
// more colorful outside than inside, but it has the same chroma in both environments.
func (c *Cam16) M() float64 {
	return c.m
}

// S is saturation in CAM16.
// Colorfulness in proportion to brightness. Prefer chroma, saturation measures colorfulness
// relative to the color's own brightness, where chroma is colorfulness relative to white.
func (c *Cam16) S() float64 {
	return c.s
}

// Jstar is lightness coordinate in CAM16-UCS.
func (c *Cam16) Jstar() float64 {
	return c.jstar
}

// Astar is a* coordinate in CAM16-UCS.
func (c *Cam16) Astar() float64 {
	return c.astar
}

// Bstar is b* coordinate in CAM16-UCS.
func (c *Cam16) Bstar() float64 {
	return c.bstar
}

// Distance measures the distance between colors.
// CAM16 instances also have coordinates in the CAM16-UCS space, called J*, a*, b*, or jstar,
// astar, bstar in code. CAM16-UCS is included in the CAM16 specification, and is used to measure
// distances between colors.
func (c *Cam16) Distance(other *Cam16) float64 {
	dJ := c.Jstar() - other.Jstar()
	dA := c.Astar() - other.Astar()
	dB := c.Bstar() - other.Bstar()
	dEPrime := math.Sqrt(dJ*dJ + dA*dA + dB*dB)
	dE := 1.41 * math.Pow(dEPrime, 0.63)
	return dE
}

// FromInt creates a CAM16 color from an ARGB representation of a color,
// assuming the color was viewed in default viewing conditions.
func FromInt(argb int) *Cam16 {
	return FromIntInViewingConditions(argb, DEFAULT_VIEWING_CONDITIONS)
}

// FromIntInViewingConditions creates a CAM16 color from an ARGB representation of a color,
// defined viewing conditions.
func FromIntInViewingConditions(argb int, viewingConditions *ViewingConditions) *Cam16 {
	// Transform ARGB int to XYZ
	red := (argb & 0x00ff0000) >> 16
	green := (argb & 0x0000ff00) >> 8
	blue := (argb & 0x000000ff)
	redL := Linearized(red)
	greenL := Linearized(green)
	blueL := Linearized(blue)
	x := 0.41233895*redL + 0.35762064*greenL + 0.18051042*blueL
	y := 0.2126*redL + 0.7152*greenL + 0.0722*blueL
	z := 0.01932141*redL + 0.11916382*greenL + 0.95034478*blueL

	return FromXyzInViewingConditions(x, y, z, viewingConditions)
}

func FromXyzInViewingConditions(x, y, z float64, viewingConditions *ViewingConditions) *Cam16 {
	// Transform XYZ to 'cone'/'rgb' responses
	matrix := XYZ_TO_CAM16RGB
	rT := (x * matrix[0][0]) + (y * matrix[0][1]) + (z * matrix[0][2])
	gT := (x * matrix[1][0]) + (y * matrix[1][1]) + (z * matrix[1][2])
	bT := (x * matrix[2][0]) + (y * matrix[2][1]) + (z * matrix[2][2])

	// Discount illuminant
	rD := viewingConditions.RgbD()[0] * rT
	gD := viewingConditions.RgbD()[1] * gT
	bD := viewingConditions.RgbD()[2] * bT

	// Chromatic adaptation
	rAF := math.Pow(viewingConditions.Fl()*math.Abs(rD)/100.0, 0.42)
	gAF := math.Pow(viewingConditions.Fl()*math.Abs(gD)/100.0, 0.42)
	bAF := math.Pow(viewingConditions.Fl()*math.Abs(bD)/100.0, 0.42)
	rA := float64(Signum(rD)) * 400.0 * rAF / (rAF + 27.13)
	gA := float64(Signum(gD)) * 400.0 * gAF / (gAF + 27.13)
	bA := float64(Signum(bD)) * 400.0 * bAF / (bAF + 27.13)

	// redness-greenness
	a := (11.0*rA + -12.0*gA + bA) / 11.0
	// yellowness-blueness
	b := (rA + gA - 2.0*bA) / 9.0

	// auxiliary components
	u := (20.0*rA + 20.0*gA + 21.0*bA) / 20.0
	p2 := (40.0*rA + 20.0*gA + bA) / 20.0

	// hue
	atan2 := math.Atan2(b, a)
	atanDegrees := ToDegrees(atan2)
	var hue float64
	if atanDegrees < 0 {
		hue = atanDegrees + 360.0
	} else if atanDegrees >= 360 {
		hue = atanDegrees - 360.0
	} else {
		hue = atanDegrees
	}
	hueRadians := ToRadians(hue)

	// achromatic response to color
	ac := p2 * viewingConditions.Nbb()

	// CAM16 lightness and brightness
	j := 100.0 * math.Pow(
		ac/viewingConditions.Aw(),
		viewingConditions.C()*viewingConditions.Z())
	q := 4.0 / viewingConditions.C() *
		math.Sqrt(j/100.0) *
		(viewingConditions.Aw() + 4.0) *
		viewingConditions.FlRoot()

	// CAM16 chroma, colorfulness, and saturation.
	var huePrime float64
	if hue < 20.14 {
		huePrime = hue + 360
	} else {
		huePrime = hue
	}
	eHue := 0.25 * (math.Cos(ToRadians(huePrime)+2.0) + 3.8)
	p1 := 50000.0 / 13.0 * eHue * viewingConditions.Nc() * viewingConditions.Ncb()
	t := p1 * math.Hypot(a, b) / (u + 0.305)
	alpha := math.Pow(1.64-math.Pow(0.29, viewingConditions.N()), 0.73) * math.Pow(t, 0.9)
	// CAM16 chroma, colorfulness, saturation
	c := alpha * math.Sqrt(j/100.0)
	m := c * viewingConditions.FlRoot()
	s := 50.0 * math.Sqrt((alpha*viewingConditions.C())/(viewingConditions.Aw()+4.0))

	// CAM16-UCS components
	jstar := (1.0 + 100.0*0.007) * j / (1.0 + 0.007*j)
	mstar := 1.0 / 0.0228 * math.Log1p(0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)

	return &Cam16{hue: hue, chroma: c, j: j, q: q, m: m, s: s, jstar: jstar, astar: astar, bstar: bstar}
}

// FromJch creates a CAM16 color from CAM16 lighness, chroma and hue,
// assuming the color was viewed in default viewing conditions.
func FromJch(j, c, h float64) *Cam16 {
	return FromJchInViewingConditions(j, c, h, DEFAULT_VIEWING_CONDITIONS)
}

// FromJchInViewingConditions creates a CAM16 color from CAM16 lighness, chroma and hue,
// defined viewing conditions.
func FromJchInViewingConditions(j, c, h float64, viewingConditions *ViewingConditions) *Cam16 {
	q := 4.0 / viewingConditions.C() *
		math.Sqrt(j/100.0) *
		(viewingConditions.Aw() + 4.0) *
		viewingConditions.FlRoot()
	m := c * viewingConditions.FlRoot()
	alpha := c / math.Sqrt(j/100.0)
	s := 50.0 * math.Sqrt((alpha*viewingConditions.C())/(viewingConditions.Aw()+4.0))

	hueRadians := ToRadians(h)
	jstar := (1.0 + 100.0*0.007) * j / (1.0 + 0.007*j)
	mstar := 1.0 / 0.0228 * math.Log1p(0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)
	return &Cam16{hue: h, chroma: c, j: j, q: q, m: m, s: s, jstar: jstar, astar: astar, bstar: bstar}
}

// FromUcs creates a CAM16 color from CAM16-UCS coordinates
func FromUcs(jstar, astar, bstar float64) *Cam16 {
	return FromUcsInViewingConditions(jstar, astar, bstar, DEFAULT_VIEWING_CONDITIONS)
}

// FromUcsInViewingConditions creates a CAM16 color from CAM16-UCS coordinates in defined viewing conditions.
func FromUcsInViewingConditions(jstar, astar, bstar float64, viewingConditions *ViewingConditions) *Cam16 {
	m := math.Hypot(astar, bstar)
	m2 := math.Expm1(m*0.0228) / 0.0228
	c := m2 / viewingConditions.FlRoot()
	h := math.Atan2(bstar, astar) * (180.0 / math.Pi)
	if h < 0.0 {
		h += 360.0
	}
	j := jstar / (1. - (jstar-100.)*0.007)
	return FromJchInViewingConditions(j, c, h, viewingConditions)
}

// ToInt returns the ARGB representation of the color.
// Assumes the color was viewed in default viewing conditions,
// which are near-identical to the default viewing conditions for sRGB.
func (c *Cam16) ToInt() int {
	return c.Viewed(DEFAULT_VIEWING_CONDITIONS)
}

func (c *Cam16) Viewed(viewingConditions *ViewingConditions) int {
	xyz := c.XyzInViewingConditions(viewingConditions)
	return ArgbFromXyz(xyz[0], xyz[1], xyz[2])
}

func (c *Cam16) XyzInViewingConditions(viewingConditions *ViewingConditions) []float64 {
	var alpha float64
	if c.Chroma() == 0.0 || c.J() == 0.0 {
		alpha = 0.0
	} else {
		alpha = c.Chroma() / math.Sqrt(c.J()/100.0)
	}

	t := math.Pow(alpha/math.Pow(1.64-math.Pow(0.29, viewingConditions.N()), 0.73), 1.0/0.9)
	hRad := ToRadians(c.Hue())

	eHue := 0.25 * (math.Cos(hRad+2.0) + 3.8)
	ac := viewingConditions.Aw() * math.Pow(c.J()/100.0, 1.0/viewingConditions.C()/viewingConditions.Z())
	p1 := eHue * (50000.0 / 13.0) * viewingConditions.Nc() * viewingConditions.Ncb()
	p2 := (ac / viewingConditions.Nbb())

	hSin := math.Sin(hRad)
	hCos := math.Cos(hRad)

	gamma := 23.0 * (p2 + 0.305) * t / (23.0*p1 + 11.0*t*hCos + 108.0*t*hSin)
	a := gamma * hCos
	b := gamma * hSin
	rA := (460.0*p2 + 451.0*a + 288.0*b) / 1403.0
	gA := (460.0*p2 - 891.0*a - 261.0*b) / 1403.0
	bA := (460.0*p2 - 220.0*a - 6300.0*b) / 1403.0

	rCBase := math.Max(0, (27.13*math.Abs(rA))/(400.0-math.Abs(rA)))
	rC := Signum(rA) * (100.0 / viewingConditions.Fl()) * math.Pow(rCBase, 1.0/0.42)
	gCBase := math.Max(0, (27.13*math.Abs(gA))/(400.0-math.Abs(gA)))
	gC := Signum(gA) * (100.0 / viewingConditions.Fl()) * math.Pow(gCBase, 1.0/0.42)
	bCBase := math.Max(0, (27.13*math.Abs(bA))/(400.0-math.Abs(bA)))
	bC := Signum(bA) * (100.0 / viewingConditions.Fl()) * math.Pow(bCBase, 1.0/0.42)
	rF := rC / viewingConditions.RgbD()[0]
	gF := gC / viewingConditions.RgbD()[1]
	bF := bC / viewingConditions.RgbD()[2]

	matrix := CAM16RGB_TO_XYZ
	x := (rF * matrix[0][0]) + (gF * matrix[0][1]) + (bF * matrix[0][2])
	y := (rF * matrix[1][0]) + (gF * matrix[1][1]) + (bF * matrix[1][2])
	z := (rF * matrix[2][0]) + (gF * matrix[2][1]) + (bF * matrix[2][2])

	return []float64{x, y, z}
}
