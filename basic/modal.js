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

function modal_open(id) {
	const modal = document.getElementById(id);
	const page = page_parentPage(modal);
	if (!page.getAttribute("data-frozen")) {
		page.setAttribute("data-frozen", "1");
		if (document.body.style.position!="fixed") {
			const scrollTop = window.scrollY;
			document.body.style.position = "fixed";
			document.body.style.top = `-${scrollTop}px`;
		}
		let tabStops = page.querySelectorAll('[tabindex="0"]');
		for (const ts of tabStops) {
			ts.setAttribute("tabindex", "-1");
			if (ts===document.activeElement) {
				ts.setAttribute("data-active", "1");
				ts.blur();
			}
		}
		tabStops = modal.querySelectorAll('[tabindex="-1"]');
		for (const ts of tabStops) {
			ts.setAttribute("tabindex", "0");
		}
	}
	// Can't nest blur effect
	let m = modal;
	while (m) {
		m = m.parentElement?.closest(".Modal");
		if (m) {
			m.style.backdropFilter = "initial"
		}
	}
}

function modal_close(id) {
	const modal = document.getElementById(id);
	const page = page_parentPage(modal);
	if (page.getAttribute("data-frozen")) {
		page.removeAttribute("data-frozen");
		if (document.querySelectorAll('[data-frozen="1"]').length==0) {
			const scrollY = document.body.style.top;
			document.body.style.position = "";
			document.body.style.top = "";
			window.scrollTo(0, parseInt(scrollY || '0') * -1);
		}
		const tabStops = page.querySelectorAll('[tabindex="-1"]');
		for (const ts of tabStops) {
			ts.setAttribute("tabindex", "0");
			if (ts.getAttribute("data-active")) {
				page.removeAttribute("data-active");
				ts.focus();
			}
		}
	}
	// Restore blur effect to top modal
	let m = modal;
	while (m) {
		m = m.parentElement?.closest(".Modal");
		if (m) {
			m.style.backdropFilter = "";
			break;
		}
	}
}
