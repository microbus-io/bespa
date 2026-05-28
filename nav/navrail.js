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

function navrail_mouseover(event) {
	const target = event.target.closest(".NavTarget");
	if (!target) {
		return;
	}
	const rail = target.parentElement.parentElement;
	if (rail!=event.currentTarget) {
		return;
	}
	const slider = rail.firstChild;
	const next = target.getAttribute("data-next");
	if (next) {
		event.preventDefault();
		event.stopPropagation();
		slider.style.width = "280px";
		slider.style.left = "100%";
		for (const panel of slider.children) {
			if (panel.id==next) {
				for (const ch of panel.children) {
					if (ch.className=="NavDrawer") {
						ch.firstChild.style.marginLeft = "0";
					}
				}
				panel.style.left = "12px";
				if (event.type=="keydown") {
					panel.querySelector('[tabindex="0"]')?.focus();
				}
			} else {
				panel.style.left = "-280px";
			}
		}
	}
	else {
		slider.style.width = "100%";
		slider.style.left = "0";
	}
}

function navrail_mouseleave(event) {
	const rail = event.currentTarget;
	const slider = rail.firstChild;
	slider.style.width = "100%";
	slider.style.left = "0";
}

function navrail_click(event) {
	navrail_mouseover(event);
}

function navrail_keydown(event) {
	if (event.keyCode==13) {
		navrail_mouseover(event);
	}
}
