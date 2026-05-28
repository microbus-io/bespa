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

function mainmenu_reveal(event, drawerID) {
	event.preventDefault();
	const navMenu = event.currentTarget.closest(".MainMenu");
	const vertical = navMenu.querySelector(".VerticalSection");	
	const horizontal = navMenu.querySelector(".HorizontalSection");
	const backdrop = navMenu.querySelector(".Backdrop");
	navMenu.closest("NAV").style.zIndex = 10000;
	vertical.classList.add("Shown");
	vertical.querySelector('[tabindex="0"]')?.focus();
	backdrop.classList.add("On");
}

function mainmenu_conceal(event) {
	event.preventDefault();
	const navMenu = event.currentTarget.closest(".MainMenu");
	const vertical = navMenu.querySelector(".VerticalSection");
	const horizontal = navMenu.querySelector(".HorizontalSection");
	const backdrop = navMenu.querySelector(".Backdrop");
	window.setTimeout(function() {
		navMenu.closest("NAV").style.zIndex = 2;
	}, 150);
	vertical.classList.remove("Shown");
	horizontal.querySelector('[tabindex="0"]')?.focus();
	backdrop.classList.remove("On");
}

function mainmenu_railMousenter(event) {
	event.currentTarget.closest("NAV").style.zIndex = 10000;
}

function mainmenu_mouseleave(event) {
	if (event.currentTarget!=event.target) {
		return;
	}
	const navMenu = event.currentTarget;
	const vertical = navMenu.querySelector(".VerticalSection");
	if (vertical.classList.contains("Shown")) {
		const horizontal = navMenu.querySelector(".HorizontalSection");
		const backdrop = navMenu.querySelector(".Backdrop");
		window.setTimeout(function() {
			navMenu.closest("NAV").style.zIndex = 2;
		}, 150);
		vertical.classList.remove("Shown");
		horizontal.querySelector('[tabindex="0"]')?.focus();
		backdrop.classList.remove("On");
	} else {
		navMenu.closest("NAV").style.zIndex = 2;	
	}
}
