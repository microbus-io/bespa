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

	"github.com/microbus-io/bespa/website/shared"
)

// HandleToolbar demonstrates the toolbar widget.
func HandleToolbar(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Toolbar widget").
			AddLeft(
				wf.ButtonOutlined("").Add("The"),
				wf.ButtonTonal("").Add("App bar"),
			).
			AddRight(
				wf.ButtonOutlined("").Add("Is also a"),
				wf.ButtonTonal("").Add("Toolbar"),
			),

		wf.HeadlineMedium("Horizontal alignment"),
		wf.Printf(
			`A toolbar has two sections that align widgets to the left {0} or to the right {1}.`,
			wf.Icon("align horizontal left"), wf.Icon("align horizontal right"),
		),
		wf.Toolbar().AddLeft(
			wf.Icon("align horizontal left"),
			wf.ButtonTonal("").Add("Left"),
			wf.ButtonOutlined("").Add("Aligned"),
			wf.ButtonOutlined("").Add("Widgets"),
		).AddRight(
			wf.ButtonTonal("").Add("Right"),
			wf.InputText("", "Aligned").WithWidth(8).WithPlaceholder("Aligned"),
			wf.Dropdown("", "").AddOption("", "Widgets"),
			wf.Icon("align horizontal right"),
		),

		wf.HeadlineMedium("Vertical alignment"),
		wf.Printf(
			`Widgets can be aligned to the top {0}, bottom {1} or center {2} when they have varying heights.`,
			wf.Icon("vertical align top"), wf.Icon("vertical align bottom"), wf.Icon("vertical align center"),
		),
		wf.GroupingFrame("").Add(
			wf.Toolbar().WithAlignTop().AddLeft(
				wf.Icon("vertical align top"),
				wf.ButtonFilled("").Add("Top"),
				wf.Avatar("Aligned", "").WithNameLabel(true).WithSize(2.5),
				wf.Checkbox("", true).Add("Toolbar"),
				wf.Icon("vertical align top"),
			),
		),
		wf.GroupingFrame("").Add(
			wf.Toolbar().WithAlignCenter().AddLeft(
				wf.Icon("vertical align center"),
				wf.ButtonFilled("").Add("Center"),
				wf.Avatar("Aligned", "").WithNameLabel(true).WithSize(2.5),
				wf.Checkbox("", true).Add("Toolbar"),
				wf.FilterChip("", "( Default )", true),
				wf.Icon("vertical align center"),
			),
		),
		wf.GroupingFrame("").Add(
			wf.Toolbar().WithAlignBottom().AddLeft(
				wf.Icon("vertical align bottom"),
				wf.ButtonFilled("").Add("Bottom"),
				wf.Avatar("Aligned", "").WithNameLabel(true).WithSize(2.5),
				wf.Checkbox("", true).Add("Toolbar"),
				wf.Icon("vertical align bottom"),
			),
		),

		wf.HeadlineMedium("Wrapping"),
		"If the screen isn't wide enough, the right-aligned section drops below the left-aligned section. ",
		"Resize the window to see how the toolbar adjusts to different widths.",
		wf.Toolbar().AddLeft(
			wf.ButtonOutlined("").Add("The"),
			wf.ButtonOutlined("").Add("Toolbar"),
			wf.FilterChip("", "Wraps", true),
			wf.ButtonOutlined("").Add("If"),
			wf.ButtonOutlined("").Add("There"),
			wf.ButtonOutlined("").Add("Are"),
			wf.ButtonTonal("").Add("Too"),
			wf.ButtonTonal("").Add("Many"),
			wf.ButtonTonal("").Add("Widgets"),
		).AddRight(
			wf.ButtonOutlined("").Add("And"),
			wf.ButtonOutlined("").Add("They"),
			wf.FilterChip("", "Don't", true),
			wf.FilterChip("", "Fit", true),
			wf.ButtonOutlined("").Add("In"),
			wf.ButtonOutlined("").Add("One"),
			wf.ButtonOutlined("").Add("Row"),
		),

		wf.HeadlineMedium("Anchoring to bottom of page"),
		"If the toolbar is the last element on the page or in a form, it is anchored to the bottom of the page.",
		wf.Toolbar().AddLeft(
			wf.ButtonText("").WithHrefBack().Add("Back"),
			wf.ButtonFilled("").WithHrefBack().Add("OK"),
			wf.Collection("This toolbar in anchored"),
		),
	)
	shared.Render(w, r, page)
}
