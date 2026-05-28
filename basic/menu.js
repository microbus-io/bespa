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

function menu_click(event) {
	event.stopPropagation();
}

function menu_mouseenter(event) {
	const menu = event.currentTarget;
	const pop = menu.nextSibling;
	const menuRect = page_getAbsPosRect(menu);
	let top = menuRect.top + menuRect.height;
	let left = menuRect.left;
	const cropElem = page_getCroppingAncestor(menu);
	const cropRect = page_getAbsPosRect(cropElem);
	if (left + pop.offsetWidth + 8 > cropRect.left + cropRect.width) {
		left = cropRect.left + cropRect.width - pop.offsetWidth - 8;
	}
	if (top + pop.offsetHeight + 8 > cropRect.top + cropRect.height) {
		top = cropRect.top + cropRect.height - pop.offsetHeight - 2*menuRect.height + 1;
	}
	pop.style.top = Math.max(top,0) + "px";
	pop.style.left = Math.max(left,0) + "px";
}
