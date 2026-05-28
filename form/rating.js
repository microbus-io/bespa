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

function starrating_keydown(event) {
	if (event.keyCode==32) {
		event.preventDefault();
		const clickEvent = new Event('click', {
			bubbles: true,
			cancelable: true,
		});
		event.currentTarget.dispatchEvent(clickEvent);
	} else if (event.keyCode==37) {
		event.currentTarget.previousElementSibling?.focus();
	} else if (event.keyCode==39) {
		event.currentTarget.nextElementSibling?.focus();
	}
}

function starrating_click(event) {
	event.stopPropagation();
	const star = event.currentTarget;
	const value = star.getAttribute("value");
	let rater = star.closest(".Rating");
	const stars = rater.querySelectorAll(".Star")
	for (const s of stars) {
		let full = Number(s.getAttribute("value"))<=Number(value);
		if (rater.getAttribute("data-style")=="sentiment") {
			full = Number(s.getAttribute("value"))==Number(value);
		}
		if (full) {
			s.classList.add("Full");
		} else {
			s.classList.remove("Full");
		}
	}
	const input = rater.lastElementChild;
	input.value = value
	const inputEvent = new Event('input', {
		bubbles: true,
		cancelable: true,
	});
	input.dispatchEvent(inputEvent);
}
