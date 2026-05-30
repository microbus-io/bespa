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

async function page_click(event) {
	let link = event.target;
	while (link.tagName?.toUpperCase()!="A" && link!=event.currentTarget) {
		link = link.parentElement;
	}
	if (link.tagName?.toUpperCase()=="A") {
		event.stopPropagation()
		event.preventDefault();
		const page = event.currentTarget;
		const action = link.getAttribute("href") ?? link.getAttribute("xlink:href");
		if (action==null) {
			return;
		}
		const actionURL = new URL(action, page_location(page));
		const payload = new URLSearchParams(actionURL.searchParams);
		const target = link.getAttribute("target");
		await page_nav(page, "GET", action, payload, target);
	}
}

async function page_submit(event) {
	if (event.target.tagName?.toUpperCase()=="FORM") {
		event.stopPropagation()
		event.preventDefault();
		const page = event.currentTarget;
		const form = event.target;
		const action = form.getAttribute("action");
		if (action==null) {
			return;
		}
		const actionURL = new URL(action, page_location(page));
		const payload = new URLSearchParams(actionURL.searchParams);
		const target = form.getAttribute("target");
		const inputs = form.querySelectorAll(`INPUT, TEXTAREA, SELECT`)
		for (const i of inputs) {
			if (i.name && page_parentPage(i)===page) {
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
			}
		}
		if (event.submitter.tagName?.toUpperCase()=="BUTTON" && event.submitter.name) {
			payload.set(event.submitter.name, event.submitter.value);
		}
		await page_nav(page, form.method.toUpperCase(), action, payload, target);
	}
}

function page_relativePath(base, path) {
	let pBase = base.split("/");
	pBase = pBase.slice(0, pBase.length-1);
	let pPath = path.split("/");
	let shared = 0;
	while (pBase.length>0 && pPath.length>0 && pBase[0]==pPath[0]) {
		pBase = pBase.slice(1);
		pPath = pPath.slice(1);
		shared++;
	}
	if (shared==1) {
		return path;
	}
	let ellipsis = "../".repeat(pBase.length);
	const remainder = pPath.join("/");
	let result = ellipsis + remainder;
	if (result=="") {
		const p = path.lastIndexOf("/");
		result = ".." + path.substring(p);
	}
	if (result.length>=path.length) {
		return path;
	}
	return result;
}

async function page_nav(page, method, action, searchParams, target) {
	// Delegate to parent page
	if (action.startsWith("^")) {
		await page_nav(page_parentPage(page), method, action.substring(1), searchParams, null);
		return;
	}
	// Delegate to top page
	if (action.startsWith("~")) {
		await page_nav(page_topPage(page), method, action.substring(1), searchParams, null);
		return;
	}

	// Apply to state and redraw
	if (action.startsWith("?")) {
		const changed = [];
		const stateData = page_stateData(page);
		for(const [key, value] of searchParams) {
			if (!value) {
				if (!stateData.has(key) || stateData.get(key)!="") {
					changed.push(key);
				}
				stateData.delete(key);
			} else {
				if (!stateData.has(key) || stateData.get(key)!=value) {
					changed.push(key);
				}
				stateData.set(key, value);
			}
		}
		if (changed.length==0) {
			return;
		}
		stateData.append("_changed", changed);
	
		const fetchURL = new URL(action, page_location(page));
		const response = await page_fetch(method, fetchURL, stateData);
		if (Number(response.location)<0) {
			history.go(Number(response.location));
			return;
		}
		if (response.location) {
			const action = response.location;
			const actionURL = new URL(action, page_location(page));
			const payload = new URLSearchParams(actionURL.searchParams);
			await page_nav(page, "GET", action, payload, null);
			return;
		}
		const errElem = document.body.querySelector(".FetchError");
		if (response.ok) {
			const htmlElem = document.createElement("HTML");
			htmlElem.innerHTML = response.text;
			const bodyElem = htmlElem.querySelector("BODY");
			const redrawn = [...bodyElem.children];
			page_applyRedrawnElements(page, redrawn);
			errElem.style.display = "none";
			if (debugger_refresh) {
				debugger_refresh(response.ok, page, redrawn, method, fetchURL, stateData);
			}
		} else {
			errElem.querySelector(".ErrMsg").innerText = page_errorMessage(response);
			errElem.style.display = "flex";
			if (debugger_refresh) {
				debugger_refresh(response.ok, page, response.text, method, fetchURL, stateData);
			}
		}
		return;
	}

	let targetPage;
	let parentRefPage = page;
	for (let p = page; p && !target; p = page_parentPage(p)) {
		target = target || p?.getAttribute("data-target");
		if (target=="_parent") {
			parentRefPage = p;
		}
	}
	if (!target) {
		targetPage = page;
	} else if (target=="_top") {
		targetPage = page_topPage();
	} else if (target=="_parent") {
		targetPage = page_parentPage(parentRefPage);
	} else {
		// Locate named page
		targetPage = document.body.querySelector('.Embed[data-name="'+target+'"] .Page')

		// Open in named window
		if (!targetPage) {
			if (method=="GET") {
				const fetchURL = new URL(action, page_location(page));
				fetchURL.search = searchParams.toString();
				window.open(fetchURL.href, target);
			} else {
				const fetchURL = new URL(action, page_location(page));
				fetchURL.search = "";
				const formElem = document.createElement("FORM");
				formElem.method = method;
				formElem.action = fetchURL.href;
				formElem.target = target;
				for(const [key, value] of searchParams) {
					const inputElem = document.createElement("INPUT");
					inputElem.type = "hidden";
					inputElem.name = key;
					inputElem.value = value;
					formElem.appendChild(inputElem);
				}
				document.body.appendChild(formElem);
				formElem.submit();
				window.setTimeout(function() {
					document.body.removeChild(formElem);
				}, 5000);
			}
			return;
		}
	}
	if (!targetPage) {
		targetPage = page;
	}

	// Replace entire nested page
	const parentPage = page_parentPage(targetPage);
	if (parentPage) { // Not the top page
		// Auto-calculate the back link, when navigating inside the same page
		if (!searchParams.has("_back")) {
			const actionPath = new URL(action, page_location(page)).pathname;
			let backURL = new URL(page_location(page)).pathname;
			backURL = page_relativePath(actionPath, backURL);
			const state = page_stateData(page).toString();
			if (state!="") {
				backURL = backURL + "?" + state;
			}
			searchParams.set("_back", backURL);
			if (targetPage!=page) {
				searchParams.set("_back", "0");
			}
		}
		const fetchURL = new URL(action, page_location(page));
		const response = await page_fetch(method, fetchURL, searchParams);
		if (Number(response.location)<0) {
			history.go(Number(response.location));
			return;
		}
		if (response.location) {
			const action = response.location;
			const actionURL = new URL(action, page_location(page));
			const payload = new URLSearchParams(actionURL.searchParams);
			await page_nav(page, "GET", action, payload, target);
			return;
		}
		const errElem = document.body.querySelector(".FetchError");
		if (response.ok) {
			const htmlElem = document.createElement("HTML");
			htmlElem.innerHTML = response.text;
			const pageElem = htmlElem.querySelector(".Page");
			targetPage.parentElement.replaceChild(pageElem, targetPage);
			page_observeSwap(targetPage, pageElem);
			const scripts = pageElem.querySelectorAll("SCRIPT");
			for (const s of scripts) {
				eval(s.innerText);
			}
			const autoFocus = pageElem.querySelectorAll('[autofocus="1"');
			if (autoFocus?.length) {
				autoFocus[0].focus();
			}
			errElem.style.display = "none";
			if (debugger_refresh) {
				debugger_refresh(response.ok, pageElem, [pageElem], method, fetchURL, searchParams);
			}
		} else {
			errElem.querySelector(".ErrMsg").innerText = page_errorMessage(response);
			errElem.style.display = "flex";
			if (debugger_refresh) {
				debugger_refresh(response.ok, page, response.text, method, fetchURL, searchParams);
			}
		}
		return;
	}

	// Replace the top page
	const docURL = new URL(document.location);
	docURL.search = page_stateData(targetPage).toString();
	history.replaceState(null, null, docURL.toString());
	if (method=="GET") {
		const fetchURL = new URL(action, page_location(page));
		fetchURL.search = searchParams.toString();
		document.location.assign(fetchURL.href);
	} else {
		const fetchURL = new URL(action, page_location(page));
		fetchURL.search = "";
		const formElem = document.createElement("FORM");
		formElem.method = method;
		formElem.action = fetchURL.href;
		for(const [key, value] of searchParams) {
			const inputElem = document.createElement("INPUT");
			inputElem.type = "hidden";
			inputElem.name = key;
			inputElem.value = value;
			formElem.appendChild(inputElem);
		}
		document.body.appendChild(formElem);
		formElem.submit();
	}
}

// page_location returns the fully qualified location of the page element, as a string.
function page_location(page) {
	return new URL(page.getAttribute("data-location"), document.location).href;
}

async function page_fetch(method, url, searchParams) {
	let response;
	if (method=="GET" && searchParams.toString().length>2000) {
		method = "PATCH";
	}
	try {
		if (method=="GET") {
			url.search = searchParams.toString();
			response = await fetch(url.toString(), {
				method: "GET",
				headers: {
					"Bespa-Fetch": "1"
				}
			});
		} else {
			url.search = "";
			response = await fetch(url.toString(), {
				method: method,
				body: searchParams,
				headers: {
					"Content-Type": "application/x-www-form-urlencoded",
					"Bespa-Fetch": "1"
				}
			});
		}
	} catch (e) {
		let msg = e.message;
		if (msg=="Failed to fetch") {
			msg = "Failed to connect to server";
		}
		return {"ok": false, "text": msg};
	}
	const text = await response.text();
	if (text.startsWith("Location: ")) {
		return {"ok": true, "location": text.substring(10)};
	}
	return {"ok": response.ok, "status": response.status, "statusText": response.statusText, "text": text};
}

function page_applyRedrawnElements(page, elems) {
	for (const elem of elems) {
		if (elem.className=="State") {
			const div = page.querySelector("DIV.State")
			if (div) {
				div.innerText = elem.innerText;
			}
			continue;
		}
		const dataID = elem.getAttribute("data-id")
		if (dataID) {
			const oldElems = page.querySelectorAll('[data-id="' + dataID + '"]');
			let autoFocus;
			for (const oldElem of oldElems) {
				if (page_parentPage(oldElem)===page) {
					const newElem = elem;
					oldElem.parentElement.replaceChild(newElem, oldElem);
					page_observeSwap(oldElem, newElem);
					const scripts = newElem.querySelectorAll("SCRIPT");
					for (const s of scripts) {
						eval(s.innerText);
					}
					if (!autoFocus?.length) {
						autoFocus = newElem.querySelectorAll('[autofocus="1"');
					}
				}
			}
			if (autoFocus?.length) {
				autoFocus[0].focus();
			}
		}
	}
}

function page_stateData(page) {
	const stateData = new URLSearchParams();
	const div = page.querySelector("DIV.State")
	if (div) {
		const state = JSON.parse(div.innerText);
		for (const s in state) {
			stateData.set(s, state[s]);
		}
	}
	return stateData;
}

function page_parentPage(elem) {
	return elem.parentElement.closest(".Page");
}

function page_topPage() {
	return document.body.querySelector(".Page")
}

const page_resizeObserver = new ResizeObserver(page_observe);
window.addEventListener("load", page_initObserver);

function page_initObserver() {
	const observed = document.querySelectorAll(`[data-observe-width]`)
	for (const o of observed) {
		page_resizeObserver.observe(o, {box: "border-box"});
	}
}

function page_observeSwap(oldElem, newElem) {
	if (oldElem.getAttribute("data-observe-width")) {
		page_resizeObserver.unobserve(oldElem);
	}
	const unobserved = oldElem.querySelectorAll(`[data-observe-width]`)
	for (const o of unobserved) {
		page_resizeObserver.unobserve(o);
	}
	if (newElem.getAttribute("data-observe-width")) {
		page_resizeObserver.observe(newElem);
	}
	const observed = newElem.querySelectorAll(`[data-observe-width]`)
	for (const o of observed) {
		page_resizeObserver.observe(o, {box: "border-box"});
	}
}

function page_observe(entries) {
	for (const entry of entries) {
		if (entry.borderBoxSize) {
			const observedWidth = entry.borderBoxSize[0].inlineSize;
			page_setObservedWidth(entry.target, observedWidth);
		}
	}
}

function page_setObservedWidth(elem, observedWidth) {
	if (!observedWidth) {
		observedWidth = elem.offsetWidth;
	}
	const steps = elem.getAttribute("data-observe-width").split(",");
	elem.classList.toggle("Width_"+steps[0], observedWidth<Number(steps[0]))
	for (let i=1; i<steps.length; i++) {
		elem.classList.toggle("Width"+steps[i-1]+"_"+steps[i], observedWidth>=Number(steps[i-1]) && observedWidth<Number(steps[i]))
	}
	const lastStep = steps[steps.length-1];
	elem.classList.toggle("Width"+lastStep+"_", observedWidth>=Number(lastStep))
	window.setTimeout(function() {
		elem.classList.add("WidthObserved");
	}, 1);
}

function page_getAbsPosRect(elem) {
	const box = elem.getBoundingClientRect();
	let top = box.top + window.scrollY;
	let left = box.left + window.scrollX
	let p = elem;
	while(p){
		const pos = window.getComputedStyle(p,null).getPropertyValue('position');
		if(pos!=='static'){
			const refBox=p.getBoundingClientRect();
			top -= refBox.top - p.scrollTop + window.scrollY;
			left -= refBox.left - p.scrollLeft + window.scrollX;
			break;
		}
		p = p.parentElement;
	}
	return {top: top, left: left, width: box.width, height: box.height};
}

function page_getCroppingAncestor(elem) {
	let p = elem;
	while (p) {
		const computedStyle = window.getComputedStyle(p);
		if (computedStyle.overflowX=="hidden" || computedStyle.overflowX=="scroll" ||
			computedStyle.overflowY=="hidden" || computedStyle.overflowY=="scroll") {
			return p;
		}
		p = p.parentElement;
	}
	return document.querySelector("BODY");
}

// page_errorMessage extracts a human-readable message from a non-OK page_fetch
// response. The argument is the response object returned by page_fetch and
// includes ok, status, statusText, and text. The body is preferred — JSON
// {err:{...}}/{error:{...}}/{error:...}/{message:...} shapes are unwrapped — and
// the HTTP status line is used as a final fallback so the user always sees
// something more useful than an empty banner.
function page_errorMessage(response) {
	const status = response && response.status ? response.status : 0;
	const statusText = response && response.statusText ? response.statusText : "";
	const body = response && response.text ? response.text : "";
	const fallback = status > 0
		? ("Server error " + status + (statusText ? " " + statusText : ""))
		: "Server error";
	if (body) {
		try {
			const j = JSON.parse(body);
			if (j.err && (j.err.error || j.err.message)) {
				return j.err.error || j.err.message;
			}
			if (j.error && (j.error.error || j.error.message)) {
				return j.error.error || j.error.message;
			}
			const msg = j.error || j.message;
			if (msg) {
				return msg;
			}
		} catch (e) {
			const firstLine = body.split("\n")[0].trim();
			if (firstLine) {
				return firstLine;
			}
		}
	}
	return fallback;
}