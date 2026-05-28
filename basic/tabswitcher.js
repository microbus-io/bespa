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

function tabswitcher_keydown(event) {
	if (event.keyCode!=32 && event.keyCode!=13) {
		return;
	}
	const clickEvent = new Event('click', {
		bubbles: true,
		cancelable: true,
	});
	event.currentTarget.dispatchEvent(clickEvent);
}

function tabswitcher_click(event) {
	const label = event.currentTarget;
	if (label.classList.contains("selected") || event.target==label.lastElementChild) {
		return;
	}
	const clickEvent = new Event('click', {
		bubbles: true,
		cancelable: true,
	});
	const link = label.lastElementChild;
	link.dispatchEvent(clickEvent);

	const page = page_parentPage(label);
	const switcherName = label.getAttribute("data-name");
	const selectedTab = label.getAttribute("data-tab");
	const dynamic = label.getAttribute("data-dynamic");
	if (!dynamic) {
		const siblings = label.closest(".TabLabels").querySelectorAll(".TabLabel");
		for (const s of siblings) {
			if (s!==label) {
				s.classList.remove("selected");
				const body = page.querySelector('.TabBody[data-name="' + switcherName + '"][data-tab="' + s.getAttribute("data-tab") + '"]');
				if (body) {
					body.classList.remove("selected");
					body.setAttribute("hidden", "true");
				}
			}
		}
		label.classList.add("selected");
		const body = page.querySelector('.TabBody[data-name="' + switcherName + '"][data-tab="' + selectedTab + '"]');
		if (body) {
			body.classList.add("selected");
			body.removeAttribute("hidden");
		}
	}
}
