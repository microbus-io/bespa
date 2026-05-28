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

package nav

var _ = Widget(&NavStripWidget{}) // Ensure interface

// NavStripWidget renders a compact horizontal strip of nav targets,
// embedding NavDrawerWidget so the same Add* methods apply.
type NavStripWidget struct {
	*NavDrawerWidget
}

// NavStrip creates a new widget that renders a horizontal strip of nav
// targets — the always-visible bar of a MainMenu on narrow screens. It's
// a NavDrawer styled inline, so the same AddTop/AddBottom usage applies.
func (f NavFactory) NavStrip() *NavStripWidget {
	drawer := f.NavDrawer()
	drawer.style = "Strip"
	return &NavStripWidget{drawer}
}

// AddTop appends entries to the leading side of the strip.
func (wgt *NavStripWidget) AddTop(topChildren ...any) *NavStripWidget {
	wgt.NavDrawerWidget.AddTop(topChildren...)
	return wgt
}

// AddBottom appends entries to the trailing side of the strip.
func (wgt *NavStripWidget) AddBottom(bottomChildren ...any) *NavStripWidget {
	wgt.NavDrawerWidget.AddBottom(bottomChildren...)
	return wgt
}
