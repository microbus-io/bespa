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

// Font is used to customize the style of the page.
// LineHeight and Size are in CSS points (pt). The Material design defaults
// use the px-equivalent values scaled by 12/16 (1pt = 1/72 in ≈ 1.333px at
// 96 dpi, so 12pt ≈ 16px).
type Font struct {
	Family     string
	LineHeight float32
	Size       float32
	Tracking   float32
	Weight     int
}
