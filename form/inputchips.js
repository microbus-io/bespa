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

function inputchips_focusin(event) {
	const elem = event.currentTarget.closest(".InputChips");
	elem.classList.add("Focus");
}
function inputchips_focusout(event) {
	const elem = event.currentTarget.closest(".InputChips")
	elem.classList.remove("Focus");
	const popup = elem.querySelector("UL");
	if (!popup.getAttribute("data-hover")) {
		popup.innerHTML = "";
		inputchips_resizeObserver.unobserve(elem);
	}
}

function inputchips_remove(event) {
	const elem = event.currentTarget.closest(".InputChips");
	const chip = event.currentTarget.closest(".Chip");
	chip.parentElement.removeChild(chip);
	event.stopPropagation();
	inputchips_recalcValue(elem);
}

function inputchips_keydownRemove(event) {
	if (event.keyCode!=56 && event.keyCode!=8) {
		return;
	}
	inputchips_remove(event);
}

function inputchips_recalcValue(elem) {
	let aggregatedValue = "";
	let aggregatedTitle = "";
	const chips = elem.querySelectorAll(".Chip");
	for (const ch of chips) {
		aggregatedValue += ch.getAttribute("data-value") + "\n";
		aggregatedTitle += ch.firstChild.innerText + "\n";
	}
	const hiddenValue = elem.querySelector('INPUT[type="hidden"]');
	hiddenValue.value = aggregatedValue.trim();
	const hiddenTitle = hiddenValue.nextSibling;
	hiddenTitle.value = aggregatedTitle.trim();
	const maxItems = elem.getAttribute("data-maxitems");
	const input = elem.querySelector('INPUT[type="text"]');
	if (maxItems) {
		if (chips.length>=maxItems) {
			input.classList.add("Saturated");
			chips[chips.length-1].focus();
		} else {
			input.classList.remove("Saturated");
			input.focus();
		}
		if (chips.length<=maxItems) {
			input.setCustomValidity("");
			elem.classList.remove('Invalid');
		}
	} else {
		input.focus();
	}
	input_autoSubmit(hiddenValue, false);
}

async function inputchips_input(event) {
	const input = event.currentTarget;
	input.setCustomValidity("");
	const q = input.value;
	const elem = event.currentTarget.closest(".InputChips");
	elem.classList.remove('Invalid');
	const url = elem.getAttribute("data-url");	
	let queryResults = {};
	if (q) {
		try {
			let u = url;
			if (u.includes("?")) {
				u = u + "&"
			} else {
				u = u + "?"
			}
			const response = await fetch(u + "q=" + encodeURIComponent(q), {
				method: "GET"
			});
			queryResults = await response.json();
		} catch (e) {
		}
	}

	const popup = elem.querySelector("UL");
	popup.innerHTML = "";
	inputchips_resizeObserver.unobserve(elem);
	if (!queryResults || !queryResults.options) {
		popup.style.display = "none";
		return;
	}
	let added = 0;
	for (let i in queryResults.options) {
		if (added>=8) {
			break;
		}
		const option = queryResults.options[i];
		let dup = false;
		const chips = elem.querySelectorAll('.Chip');
		for (let c of chips) {
			if (c.getAttribute("data-value")==option.value) {
				dup = true;
				break;
			}
		}
		if (dup) {
			continue;
		}
		const li = document.createElement("li")
		li.setAttribute("data-value", option.value);
		if (added==0) {
			li.classList.add("Active");
		}
		const titleDiv = document.createElement("div")
		titleDiv.innerText = option.title;
		inputchips_highlight(titleDiv, q);
		li.appendChild(titleDiv);
		if (option.desc) {
			const descDiv = document.createElement("div")
			descDiv.innerText = option.desc;
			inputchips_highlight(descDiv, q);
			li.appendChild(descDiv);
		}
		popup.appendChild(li);
		added++;
	}
	popup.style.display = "block";

	inputchips_positionPopup(input, popup);
	inputchips_resizeObserver.observe(elem);
}

function inputchips_positionPopup(input, popup) {
	const inputRect = page_getAbsPosRect(input);
	let top = inputRect.top + inputRect.height/2;
	let left = inputRect.left;
	const cropElem = page_getCroppingAncestor(input);
	const cropRect = page_getAbsPosRect(cropElem);
	if (left + popup.offsetWidth + 8 > cropRect.left + cropRect.width) {
		left = cropRect.left + cropRect.width - popup.offsetWidth - 8;
	}
	// if (top + popup.offsetHeight + 8 > cropRect.top + cropRect.height) {
	// 	top = cropRect.top + cropRect.height - popup.offsetHeight;
	// }
	popup.style.top = Math.max(top,0) + "px";
	popup.style.left = Math.max(left,0) + "px";
}

function inputchips_highlight(elem, query) {
	const esc = query.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&');
	let reg = new RegExp("(\\b"+esc+")", 'i');
	let value = elem.innerText;
	let newValue = value.replace(reg, "<u>$1</u>");
	if (value==newValue) {
		reg = new RegExp("("+esc+")", 'i');
		newValue = value.replace(reg, "<u>$1</u>");	
	}
	elem.innerHTML = newValue;
}

function inputchips_keydown(event) {
	if (event.keyCode!=8 && event.keyCode!=27 && event.keyCode!=38 && event.keyCode!=40 && event.keyCode!=13 && event.keyCode!=9) {
		return;
	}
	const input = event.currentTarget;
	const elem = event.currentTarget.closest(".InputChips");
	const popup = elem.querySelector("UL");
	if (event.keyCode==8) { // Backspace
		if (input.value) {
			return;
		}
		const lastChip = elem.querySelector(".Chip:last-of-type");
		if (lastChip) {
			lastChip.focus();
		}
	}
	if (event.keyCode==27) { // Escape
		popup.innerHTML = "";
		inputchips_resizeObserver.unobserve(elem);
	}
	if (event.keyCode==38) { // Up arrow
		for (let i=popup.children.length-1; i>=1; i--) {
			if (popup.children[i].classList.contains("Active")) {
				popup.children[i].classList.remove("Active");
				popup.children[i-1].classList.add("Active");
				break;
			}
		}
	}
	if (event.keyCode==40) { // Down arrow
		for (let i=0; i<popup.children.length-1; i++) {
			if (popup.children[i].classList.contains("Active")) {
				popup.children[i].classList.remove("Active");
				popup.children[i+1].classList.add("Active");
				break;
			}
		}
	}
	if (event.keyCode==13 || event.keyCode==9) { // Enter or tab
		if (event.shiftKey) {
			return;
		}
		const active = popup.querySelector(".Active");
		if (active) {
			inputchips_newChip(active);
		}
	}
	event.preventDefault();
	event.stopPropagation();
}

function inputchips_newChip(clicked) {
	const elem = clicked.closest(".InputChips");
	const input = elem.querySelector('INPUT[type="text"]');
	const blank = elem.querySelector(".BlankChip");
	const popup = elem.querySelector("UL");
	const newChip = blank.cloneNode(true);
	newChip.classList.remove("BlankChip");
	newChip.classList.add("Chip");
	newChip.firstChild.innerText = clicked.firstChild.innerText;
	newChip.setAttribute("data-value", clicked.getAttribute("data-value"));
	elem.insertBefore(newChip, input)
	input.value = "";
	popup.innerHTML = "";
	inputchips_resizeObserver.unobserve(elem);
	inputchips_recalcValue(elem);
}

function inputchips_mouseenter(event) {
	const popup = event.currentTarget;
	popup.setAttribute("data-hover", "1");
}
function inputchips_mouseleave(event) {
	const popup = event.currentTarget;
	popup.removeAttribute("data-hover");
}

function inputchips_click(event) {
	const clicked = event.target.closest("li");
	inputchips_newChip(clicked);
	event.preventDefault();
	event.stopPropagation();
}

const inputchips_resizeObserver = new ResizeObserver(inputchips_observe);
function inputchips_observe(entries) {
	for (const entry of entries) {
		if (entry.borderBoxSize) {
			const popup = entry.target.querySelector("UL");
			const input = entry.target.querySelector('INPUT[type="text"]');
			inputchips_positionPopup(input, popup);
		}
	}
}
