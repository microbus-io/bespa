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

// HandleTabSwitcher demonstrates the tab widget.
func HandleTabSwitcher(w http.ResponseWriter, r *http.Request) {
	state := wf.StateOf(r)
	top := state.Get("tab")
	if top == "" {
		top = "top1"
	}

	page := wf.Page().Add(
		wf.AppBar("Tab switcher widget").
			AddBottom(
				// Add a tab switcher inside the app bar
				wf.TabSwitcher().
					WithName("tab"). // If not specified, "tab" is the default name
					AddLeft(
						// The labels of the app bar tab switcher
						wf.TabLabel("top1").Add("Intro"),
						wf.TabLabel("top2").Add("State"),
						wf.TabLabel("top3").Add("In page"),
					).
					AddRight(
						// The labels of the app bar tab switcher
						wf.TabLabel("top4").Add("Right aligned"),
						wf.TabLabel("top5").Add(wf.Icon("lightbulb")),
					),
			),

		// The bodies of the app bar tab switcher
		wf.TabSwitcher().
			WithName("tab").
			AddLeft(
				wf.TabBody("top1").Add(
					`Tabs added to the app bar are best suited for grouping together pages that are
					closely related where the content of the page changes in its entirety.`,
					wf.SpacerParagraph(),
					`Switch to the next tab to continue...`,
				),
				wf.TabBody("top2").Add(
					wf.Markdown(`Tab widgets affect a state variable that is then used to show or hide content.
						The tab switcher widget in the app bar sets the value of the state variable `, "`tab`",
						`to the key of the selected tab: `, "`top1`, `top2` or `top3`.",
						"\n\n",
						`Click on the next tab label ("In page") to switch tabs.`,
						"\n\n",
						"Slide open the debugger to see the impact on the state of the page.",
					),
				),
				wf.TabBody("top3").Add(
					wf.Markdown(`Tab switcher widgets can be placed in the middle of the page as well.
						Tabs are shown or hidden according to the user selection purely on the frontend.
						Nevertheless, a request is still sent to the backend when the value of
						the tab widget's state variable changes.
						This tab switcher is tied to the state variable `, "`inpage`",
						`to avoid conflict with the state variable `, "`tab`", ` that's used for the top tab switcher.`,
					),

					// An in-page tab switcher
					wf.TabSwitcher().
						WithName("inpage").
						AddLeft(
							wf.TabLabel("a").Add("Nested page"),
							wf.TabLabel("b").Add("Lorem ipsum"),
							wf.TabLabel("c").Add("Proin lacinia"),
							wf.TabBody("a").Add(
								`Tabs can contain anything. In this examples the tab body embeds another page inside one of its tabs.`,
								wf.SpacerParagraph(),
								wf.GroupingFrame("").Add(
									wf.EmbedHandler(mux.ServeHTTP, r, "GET", "/basics/nested", nil),
								),
							),
							wf.TabBody("b").Add(
								`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ut purus sed purus vestibulum lacinia ut et purus.
								Praesent porttitor fringilla nunc, ut varius tortor lacinia maximus. Nam luctus sollicitudin justo, id aliquet
								velit condimentum vitae. Vivamus scelerisque varius commodo. Vestibulum interdum leo leo, id finibus odio
								pulvinar sit amet. Morbi quis neque velit. Integer euismod elementum nunc, id egestas nisl posuere id.
								Sed egestas felis eget vulputate lobortis.`,
							),
							wf.TabBody("c").Add(
								`Proin lacinia, diam ut condimentum fringilla, turpis nisl egestas mi, a pulvinar magna velit ut quam.
								Aenean nec nisl scelerisque, congue justo ut, consequat felis. Suspendisse consequat maximus ex
								sollicitudin elementum. Ut sed ipsum blandit, vulputate nulla at, faucibus felis. Donec ac tristique sapien.`,
							),
						),
				),
				wf.TabBody("top4").Add(
					`Tabs labels can be aligned to the right.`,
				),
				wf.TabBody("top5").Add(
					`Tab labels do not have to be textual. In this case an icon is used.`,
				),
			),
		wf.Debugger(),
	)

	// Prevent redrawing the nested page when switching tabs after it was initially fetched
	if state.Get("nested") == "page" {
		state.Set("fetched", "1")
	}

	shared.Render(w, r, page)
}
