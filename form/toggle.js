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

function toggle_click(event) {
	const toggle = event.currentTarget;
	toggle.classList.toggle("Selected");
	const input = toggle.firstChild;
	input.checked = toggle.classList.contains("Selected");
	const inputEvent = new Event('input', {
		bubbles: false,
		cancelable: true,
	});
	input.dispatchEvent(inputEvent);
}

function toggle_keydown(event) {
	if (event.keyCode!=32) {
		return;
	}
	event.preventDefault();
	toggle_click(event);
}
