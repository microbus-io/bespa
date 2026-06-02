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

package showcase

import (
	"net/http"
	"time"

	"github.com/microbus-io/bespa/website/shared"
)

// handleTextFormatting demonstrates text formatting.
func HandleTextFormatting(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Text formatting"),

		wf.HeadlineLarge("Headline large"),
		wf.HeadlineMedium("Headline medium"),
		wf.MessageBar("All is well"),
		wf.MessageBarError("Oops, it's broken"),
		"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		wf.HeadlineSmall("Headline small"),
		"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		wf.TitleLarge("Title large"),
		"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		wf.TitleMedium("Title medium"),
		"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		wf.TitleSmall("Title small"), " is rendered inline.",
		"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		wf.SpacerParagraph(),
		wf.TextStyle("Default").WithColorDefault(), " ",
		wf.TextStyle("Bold").WithBold(), " ",
		wf.TextStyle("Lightweight").WithLightweight(), " ",
		wf.TextStyle("Italic").WithItalic(), " ",
		wf.TextStyle("Deemphasized").WithColorDeemphasized(), " ",
		wf.TextStyle("Disabled").WithColorDisabled(), " ",
		wf.TextStyle("Monospace").WithMonospace(), " ",
		wf.SpacerNewLine(),
		wf.TextStyle("Primary").WithColorPrimary(), " ",
		wf.TextStyle("Secondary").WithColorSecondary(), " ",
		wf.TextStyle("Tertiary").WithColorTertiary(), " ",
		wf.TextStyle("Error").WithColorError(), " ",
		wf.TextStyle("OK").WithColorOK(), " ",
		wf.SpacerNewLine(),
		wf.TextStyle("Inverse").WithColorInverse(), " ",
		wf.TextStyle("Primary").WithColorOnPrimary(), " ",
		wf.TextStyle("Secondary").WithColorOnSecondary(), " ",
		wf.TextStyle("Tertiary").WithColorOnTertiary(), " ",
		wf.TextStyle("Error").WithColorOnError(), " ",
		wf.TextStyle("OK").WithColorOnOK(), " ",
		wf.SpacerNewLine(),
		wf.TextStyle("1").WithSizeMultiplier(1), " ",
		wf.TextStyle("2").WithSizeMultiplier(2), " ",
		wf.TextStyle("3").WithSizeMultiplier(3), " ",
		wf.TextStyle("4").WithSizeMultiplier(4), " ",
		wf.TextStyle("5").WithSizeMultiplier(5), " ",
		wf.SpacerParagraph(),
		wf.TearOffCalendar(time.Now()),
		wf.SpacerParagraph(),
		"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		wf.SpacerParagraph(),
		"Short value: ", wf.Code("api-key-7f3c2"), " ", wf.CopyToClipboard("api-key-7f3c2"),
		wf.SpacerNewLine(),
		"Long value ", wf.CopyToClipboard(`package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`).WithAltText("Copy code"),
		wf.PlainCodeBlock(`package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`),
		wf.TextAlignCenter(
			wf.HeadlineMedium("Center aligned"),
			"Lorem Ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		),
		wf.SpacerParagraph(),

		wf.Debugger(),
	)

	shared.Render(w, r, page)
}
