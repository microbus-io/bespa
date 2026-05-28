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

function button_click(event) {
	const btn = event.currentTarget;
	if (btn.previousSibling) {
		btn.previousSibling.setCustomValidity("");
	}
	if (btn.type!="submit" && btn.lastChild.tagName=="A" && btn.lastChild.getAttribute("href") && event.target.tagName!="A") {
		event.preventDefault();
		event.stopPropagation();
		btn.lastChild.click();
	}
}
