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

function inputtimezone_worldInput(event) {
	const world = event.currentTarget;
	let x = world.nextSibling;
	while (x) {
		if (x.tagName != "SELECT") {
			x = x.nextSibling;
			continue;
		}
		if ((x.getAttribute("data-region") ?? "")==world.value) {
			x.name = world.getAttribute("data-name");
			x.classList.remove("Off");
			const inputEvent = new Event('input', {
				bubbles: true,
				cancelable: true,
			});
			x.dispatchEvent(inputEvent);
		} else {
			x.removeAttribute("name");
			x.classList.add("Off");
		}
		x = x.nextSibling;
	}
}
