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

package shared

import (
	"fmt"
	"net/http"

	"github.com/microbus-io/bespa"
	"github.com/microbus-io/bespa/css"
	"github.com/microbus-io/bespa/widget"
)

// wf is the widget factory.
var (
	wf = bespa.DefaultFactory{}
)

// Render sets the user's color preferences before rendering the page.
func Render(w http.ResponseWriter, r *http.Request, page *widget.PageWidget) {
	session := SessionOf(w, r)
	switch session.Theme {
	case "Dark":
		page.WithThemeDark()
	case "Light":
		page.WithThemeLight()
	}
	if session.Palette != "" {
		for _, kc := range css.PresetKeyColors {
			if kc.Name == session.Palette {
				page.WithKeyColors(kc)
				break
			}
		}
	}
	page.Add(navMenu())
	err := page.
		Draw(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%+v", err)))
	}
}
