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

// Lerp is the linear interpolation function.
// It returns start if amount = 0 and stop if amount = 1.
func Lerp(start, stop, amount float64) float64 {
	return (1.0-amount)*start + amount*stop
}

// ClampDouble clamps an floating point number between two floating-point numbers.
func ClampDouble(min, max, input float64) float64 {
	if input < min {
		return min
	} else if input > max {
		return max
	}
	return input
}

// ClampDouble clamps an integer between two integers.
func ClampInt(min, max, input int) int {
	if input < min {
		return min
	} else if input > max {
		return max
	}
	return input
}

// Signum return 1 if num > 0, -1 if num < 0, and 0 if num = 0
func Signum(num float64) float64 {
	if num < 0 {
		return -1
	} else if num == 0 {
		return 0
	} else {
		return 1
	}
}

// ToDegrees converts an angle measured in radians to an approximately equivalent angle
// measured in degrees.
func ToDegrees(radians float64) float64 {
	return radians * 180.0 / math.Pi
}

// ToRadians converts an angle measured in degrees to an approximately equivalent angle
// measured in radians.
func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// MatrixMultiply multiplies a 1x3 row vector with a 3x3 matrix.
func MatrixMultiply(row []float64, matrix [][]float64) []float64 {
	a := row[0]*matrix[0][0] + row[1]*matrix[0][1] + row[2]*matrix[0][2]
	b := row[0]*matrix[1][0] + row[1]*matrix[1][1] + row[2]*matrix[1][2]
	c := row[0]*matrix[2][0] + row[1]*matrix[2][1] + row[2]*matrix[2][2]
	return []float64{a, b, c}
}

// SanitizeDegreesDouble sanitizes a degree measure as a floating-point number.
func SanitizeDegreesDouble(degrees float64) float64 {
	degrees = math.Mod(degrees, 360.0)
	if degrees < 0 {
		degrees = degrees + 360.0
	}
	return degrees
}
