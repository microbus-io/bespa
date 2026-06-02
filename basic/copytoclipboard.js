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

function copytoclipboard_click(event) {
	const elem = event.currentTarget;
	const text = elem.getAttribute("data-copy-text") || "";
	const showCopied = function() {
		const icon = elem.querySelector("i.material-symbols-outlined");
		if (!icon) return;
		const pending = elem.getAttribute("data-restore");
		if (pending) window.clearTimeout(Number(pending));
		icon.textContent = "check";
		elem.classList.add("Copied");
		const timeoutID = window.setTimeout(function() {
			icon.textContent = "content_copy";
			elem.classList.remove("Copied");
			elem.removeAttribute("data-restore");
		}, 2000);
		elem.setAttribute("data-restore", timeoutID);
	};
	if (navigator.clipboard && navigator.clipboard.writeText) {
		navigator.clipboard.writeText(text).then(showCopied, function() {});
	}
}
function copytoclipboard_keydown(event) {
	if (event.key === "Enter" || event.key === " ") {
		event.preventDefault();
		copytoclipboard_click(event);
	}
}
