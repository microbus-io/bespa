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

function debugger_click(event) {
	event.currentTarget.classList.toggle("Shown");
}

let debugger_present = false;
function debugger_init(id) {
	debugger_present = true;
	const deb = document.getElementById(id);
	const page = page_parentPage(deb);
	const loc = new URL(document.location);
	debugger_refresh(true, page, [page], "GET", loc.pathname, loc.searchParams);
}

function debugger_refresh(ok, page, content, method, action, params) {
	if (!debugger_present) {
		return;
	}
	// Identify the last active debugger
	const debs = document.querySelectorAll(".Debugger");
	if (debs.length==0) {
		return;
	}
	let deb = debs[0];
	for (const d of debs) {
		if (d.style.display=="inline-block") {
			deb = d;
			break;
		}
	}
	deb.style.display = "inline-block";

	deb.innerHTML = "";
	if (!ok) {
		deb.classList.add("Error");
	} else {
		deb.classList.remove("Error");
	}

	// Request
	if (method) {
		const title = document.createElement("h5");
		title.innerText = "Request";
		const info = document.createElement("div");
		const actionURL = new URL(action, document.location);
		info.innerHTML = debugger_escape(method + " " + actionURL.pathname) + "<br>";
		for(const [key, value] of params) {
			info.innerHTML += "<b>" + debugger_escape(key) + ":</b> " + debugger_escape(value) + "<br>";
		}
		deb.appendChild(title)
		deb.appendChild(info)
	}

	// Page states
	const pages = document.querySelectorAll("DIV.Page");
	for (const page of pages) {
		const randomID = Math.floor(Math.random() * Date.now()).toString(16);
		page.setAttribute("data-debugger-id", randomID);
		const title = document.createElement("h5");
		title.setAttribute("data-ref", randomID);
		title.setAttribute("onmouseenter", "debugger_highlight(this, true)");
		title.setAttribute("onmouseleave", "debugger_highlight(this, false)");
		title.innerText = page.getAttribute("data-location");
		const info = document.createElement("div");
		const state = page_stateData(page);
		let f = false;
		for (let [key,val] of state) {
			info.innerHTML += "<b>" + debugger_escape(key) + ":</b> " + debugger_escape(val) + "<br>";
			f = true;
		}
		if (!f) {
			info.innerHTML += "<i>empty</i>";
		}
		deb.appendChild(title)
		deb.appendChild(info)
	}

	// Redrawn elements
	if (ok && content) {
		const title = document.createElement("h5");
		title.innerText = "Redrawn";
		const info = document.createElement("div");
		for (const i of content) {
			const dataID = i.getAttribute("data-id")
			if (dataID) {
				const randomID = i.getAttribute("data-debugger-id") || Math.floor(Math.random() * Date.now()).toString(16);
				i.setAttribute("data-debugger-id", randomID);
				const span = document.createElement("span");
				span.setAttribute("data-ref", randomID);
				span.setAttribute("onmouseenter", "debugger_highlight(this, true)");
				span.setAttribute("onmouseleave", "debugger_highlight(this, false)");
				span.innerText += i.tagName;
				if (i.tagName=="INPUT" && i.type) {
					span.innerText += '[type="' + i.type + '"]';
				}
				if (i.className) {
					const classNames = i.className.split(" ");
					span.innerText += "." + classNames[0];
				}
				span.innerText += "#" + dataID;
				info.appendChild(span);
				info.appendChild(document.createElement("br"));
			}	
			deb.appendChild(title)
			deb.appendChild(info)
		}
	}

	// Error
	if (!ok && content) {
		const title = document.createElement("h5");
		title.innerText = "Error";
		const info = document.createElement("pre");
		info.innerHTML = content;
		deb.appendChild(title)
		deb.appendChild(info)
	}
}

function debugger_escape(str) {
	return str.replaceAll(">", "&gt").replaceAll("<", "&lt");
}

function debugger_highlight(elem, flag) {
	const debuggerID = elem.getAttribute("data-ref");
	if (debuggerID) {
		const ref = document.querySelector('[data-debugger-id="'+debuggerID+'"]')
		if (ref) {
			const compStyle = getComputedStyle(ref);
			if (compStyle.position=="fixed") {
				ref.classList.toggle("DebuggerHighlightNoAnim", flag);
			} else {
				ref.classList.toggle("DebuggerHighlight", flag);
			}
		}
	}
}
