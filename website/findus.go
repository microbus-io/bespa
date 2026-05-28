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

	"github.com/microbus-io/bespa/website/shared"
)

// HandleFindUs echoes the Microbus contact info at
// https://docs.microbus.io/about/contact/.
func HandleFindUs(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Find us"),

		contactRow("language", "Website",
			"https://www.microbus.io", "www.microbus.io"),
		contactRow("mail", "Email",
			"mailto:info@microbus.io", "info@microbus.io"),
		contactRow("code", "GitHub",
			"https://github.com/microbus-io", "github.com/microbus-io"),
		contactRow("business", "LinkedIn",
			"https://linkedin.com/company/microbus-io", "linkedin.com/company/microbus-io"),
	)
	shared.Render(w, r, page)
}

// contactRow renders one labeled contact method as a Field.
func contactRow(icon, label, href, text string) any {
	return wf.Field().
		AddLeft(wf.Icon(icon), " ", label).
		AddRight(wf.Link(href).Add(text))
}
