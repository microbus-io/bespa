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

function input_input(event) {
	const input = event.currentTarget;
	input.setCustomValidity("");
	if (input.validity.valid) {
		let x = input;
		while (!x.getAttribute("data-id")) {
			x.classList.remove('Invalid');
			x = x.parentElement;
		}
		x.classList.remove('Invalid');
	}
	input_autoSubmit(input,
		input.type=="text" ||
		input.type=="password" ||
		input.type=="email" ||
		input.type=="url" ||
		input.type=="number" ||
		input.type=="tel" ||
		input.tagName=="TEXTAREA");
}

function input_change(event) {
	const input = event.currentTarget;
	const prevTimeoutID = input.getAttribute("data-debouncer");
	if (prevTimeoutID) {
		input_autoSubmit(input, false);
	}
}

function input_autoSubmit(input, debounce) {
	const autoSubmit = input.getAttribute("data-autosubmit");
	if (!autoSubmit) {
		return;
	}
	const f = function () {
		if (input.value && !input.validity.valid) {
			return;
		}
		const page = page_parentPage(input);
		const actionURL = page_location(page);
		const payload = new URLSearchParams(actionURL.searchParams);
		const i = input;
		if (i.type=="checkbox") {
			if (i.checked) {
				payload.set(i.name, "1");
			} else if (!payload.get(i.name)) {
				payload.set(i.name, "0");
			}
		} else if (i.type=="radio") {
			if (i.checked) {
				payload.set(i.name, i.value);
			}
		} else {
			payload.set(i.name, i.value);
		}
		page_nav(page, "PATCH", "?", payload, null);
	};
	const prevTimeoutID = input.getAttribute("data-debouncer");
	if (prevTimeoutID) {
		window.clearTimeout(prevTimeoutID);
		input.removeAttribute("data-debouncer");
	}
	if (debounce) {
		const timeoutID = window.setTimeout(f, 250);
		input.setAttribute("data-debouncer", timeoutID);
	} else {
		f();
	}
}

function input_invalid(event) {
	const input = event.currentTarget;
	let x = input;
	while (!x.getAttribute("data-id")) {
		x.classList.add('Invalid');
		x = x.parentElement;
	}
	x.classList.add('Invalid');
}

function input_initBackgroundIconColor(id) {
	const elem = document.getElementById(id);
	const compStyle = getComputedStyle(elem);
	if (compStyle.backgroundImage?.indexOf("currentColor")>=0) {
		elem.style.backgroundImage = compStyle.backgroundImage.replace("currentColor", compStyle.color);
	}
}

function input_setCustomValidity(ids, msg) {
	for (let id of ids.split(' ')) {
		const elem = document.getElementById(id);
		elem.setCustomValidity(msg);
		const form = elem.closest("FORM");
		if (form) {
			try {
				form.reportValidity();
			} catch(e) {
				// reportValidity throws on detached form nodes; ignored.
			}
		}
	}
}