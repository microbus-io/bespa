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

package main

import (
	"net/http"
	"sort"

	"github.com/microbus-io/bespa/css"
	"github.com/microbus-io/bespa/website/shared"
)

// HandleProfile enables the user to manage their preferences.
func HandleProfile(w http.ResponseWriter, r *http.Request) {
	session := shared.SessionOf(w, r)

	// Palette dropdown
	keyColorsDropdown := wf.RichDropdown("palette", session.Palette)
	palettes := []css.KeyColors{}
	palettes = append(palettes, css.PresetKeyColors...)
	sort.Slice(palettes, func(i, j int) bool {
		return palettes[i].Name < palettes[j].Name
	})
	for _, kc := range palettes {
		tb := wf.Toolbar().
			WithAlignCenter().
			WithWrap(false).
			AddLeft(kc.Name)
		colors := []css.Color{kc.Primary, kc.Secondary, kc.Tertiary, kc.Neutral, kc.NeutralVariant, kc.Error}
		for _, c := range colors {
			tb.AddRight(wf.Swatch(c))
		}
		keyColorsDropdown.AddOption(kc.Name, tb)
	}

	form := wf.Form()
	page := wf.Page().Add(
		wf.AppBar("Profile settings"),
		form.
			WithAction("POST", "").
			Add(
				wf.Field().
					AddLeft("Palette").
					AddRight(keyColorsDropdown),
				wf.Field().
					AddLeft("Theme").
					AddRight(
						wf.Dropdown("theme", session.Theme).
							AddOption("", "Browser's default").
							AddOption("Light", "Light").
							AddOption("Dark", "Dark"),
					),
				wf.ButtonText("cancel").Add("Cancel").WithHrefBack(),
				wf.ButtonFilled("save").Add("Save"),
			),
	)
	if form.ReadyToCommit(r) {
		values := form.Values(r)
		session.Palette = values.Get("palette")
		session.Theme = values.Get("theme")
	}
	shared.Render(w, r, page)
}
