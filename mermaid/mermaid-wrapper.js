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

// The official mermaid bundle exposes itself on the IIFE we wrapped it in.
// jsdelivr's "production" build is namespaced as __esbuild_esm_mermaid_nm.mermaid;
// pluck a stable `mermaid` global out so the rest of this file can ignore that.
(function() {
	if (typeof mermaid !== "undefined") {
		return;
	}
	if (typeof __esbuild_esm_mermaid_nm !== "undefined" && __esbuild_esm_mermaid_nm.mermaid) {
		window.mermaid = __esbuild_esm_mermaid_nm.mermaid;
	}
})();

// mermaid_instances tracks every diagram on the page so all can be re-themed when
// the user switches between light and dark mode without reloading.
const mermaid_instances = [];

let mermaid_idCounter = 0;
function mermaid_nextID() {
	mermaid_idCounter++;
	return "mermaid-svg-" + mermaid_idCounter;
}

// mermaid_render is the entry point invoked from the inline script emitted by
// MermaidWidget.Draw. `id` is the container element id, `zoomPan` is 1 to wire
// up zoom/pan interactions or 0 to leave the diagram static.
function mermaid_render(id, zoomPan) {
	window.setTimeout(function() {
		const elem = document.getElementById(id);
		if (!elem) {
			return;
		}
		const sourceEl = elem.querySelector(".MermaidSource");
		const source = sourceEl ? sourceEl.textContent : elem.textContent;
		const inst = {
			elem: elem,
			source: source,
			zoomPan: !!zoomPan,
			scale: 1, tx: 0, ty: 0,
		};
		mermaid_instances.push(inst);
		mermaid_renderInst(inst);
	}, 1);
}

// mermaid_renderInst runs mermaid against `inst.source` and swaps the resulting
// SVG into `inst.elem`. Called once on creation and again on every theme switch.
function mermaid_renderInst(inst) {
	if (!inst.elem.isConnected) {
		return;
	}
	if (typeof mermaid === "undefined") {
		// Library hasn't loaded yet — try again in a moment. The bundle is async-fetched
		// on the first page that imports the mermaid package.
		window.setTimeout(function() { mermaid_renderInst(inst); }, 50);
		return;
	}
	mermaid.initialize(mermaid_buildConfig());
	const svgID = mermaid_nextID();
	mermaid.render(svgID, inst.source).then(function(result) {
		if (!inst.elem.isConnected) {
			return;
		}
		inst.elem.innerHTML = result.svg;
		// Re-add the (hidden) source so subsequent re-renders find the diagram text.
		const pre = document.createElement("pre");
		pre.className = "MermaidSource";
		pre.textContent = inst.source;
		inst.elem.appendChild(pre);
		const svg = inst.elem.querySelector("svg");
		if (svg) {
			// Mermaid sets a fixed pixel width on the SVG that overflows narrow containers.
			// Strip it so CSS sizing wins; the viewBox keeps the aspect ratio.
			svg.removeAttribute("width");
			svg.removeAttribute("height");
			svg.style.maxWidth = "100%";
			svg.style.height = "100%";
		}
		if (result.bindFunctions) {
			result.bindFunctions(inst.elem);
		}
		if (inst.zoomPan) {
			mermaid_wireZoomPan(inst, svg);
		}
	}).catch(function(err) {
		inst.elem.innerHTML = '<div class="MermaidError">Diagram error: ' + (err && err.message ? err.message : String(err)) + '</div>';
	});
}

// mermaid_cssVar reads a raw CSS custom property value off the document root.
function mermaid_cssVar(name, fallback) {
	const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
	return v || (fallback ?? "");
}

// mermaid_color resolves a Material design CSS custom property (e.g.
// "--md-sys-color-primary") to a literal "rgb(r, g, b)" or "rgba(...)" color
// string. Returns the fallback if the variable is not defined. Mermaid's SVG
// renderer accepts the resulting strings directly.
function mermaid_color(varName, alpha, fallback) {
	const v = mermaid_cssVar(varName, "");
	if (!v) {
		return fallback || "#000";
	}
	const parts = v.split(/[\s,]+/);
	if (parts.length < 3) {
		return fallback || v;
	}
	if (typeof alpha === "number" && alpha < 1) {
		return "rgba(" + parts[0] + "," + parts[1] + "," + parts[2] + "," + alpha + ")";
	}
	return "rgb(" + parts[0] + "," + parts[1] + "," + parts[2] + ")";
}

// mermaid_buildConfig produces a mermaid initialize() config wired to the
// current Material design tokens. Picks the `base` theme so themeVariables
// fully control colors.
function mermaid_buildConfig() {
	const primary = mermaid_color("--md-sys-color-primary-container", 1, "#ddd");
	const onPrimary = mermaid_color("--md-sys-color-on-primary-container", 1, "#000");
	const primaryBorder = mermaid_color("--md-sys-color-primary", 1, "#888");
	const secondary = mermaid_color("--md-sys-color-secondary-container", 1, "#eee");
	const onSecondary = mermaid_color("--md-sys-color-on-secondary-container", 1, "#000");
	const tertiary = mermaid_color("--md-sys-color-tertiary-container", 1, "#eee");
	const onTertiary = mermaid_color("--md-sys-color-on-tertiary-container", 1, "#000");
	const background = mermaid_color("--md-sys-color-background", 1, "#fff");
	const surface = mermaid_color("--md-sys-color-surface", 1, "#fff");
	const text = mermaid_color("--md-sys-color-on-background", 1, "#000");
	const line = mermaid_color("--md-sys-color-outline", 1, "#888");
	const note = mermaid_color("--md-sys-color-surface-variant", 1, "#f7f7f7");
	const onNote = mermaid_color("--md-sys-color-on-surface-variant", 1, "#333");
	const fontFamily = mermaid_cssVar("--md-sys-typescale-body-medium-font", "Roboto, sans-serif");

	return {
		startOnLoad: false,
		securityLevel: "strict",
		theme: "base",
		fontFamily: fontFamily,
		themeVariables: {
			background: background,
			primaryColor: primary,
			primaryTextColor: onPrimary,
			primaryBorderColor: primaryBorder,
			secondaryColor: secondary,
			secondaryTextColor: onSecondary,
			secondaryBorderColor: primaryBorder,
			tertiaryColor: tertiary,
			tertiaryTextColor: onTertiary,
			tertiaryBorderColor: primaryBorder,
			lineColor: line,
			textColor: text,
			mainBkg: primary,
			secondBkg: secondary,
			tertiaryBkg: tertiary,
			nodeBorder: primaryBorder,
			clusterBkg: surface,
			clusterBorder: line,
			defaultLinkColor: line,
			edgeLabelBackground: surface,
			noteBkgColor: note,
			noteTextColor: onNote,
			noteBorderColor: line,
			titleColor: text,
			fontFamily: fontFamily,
		},
	};
}

// mermaid_onThemeChange re-renders every known diagram with the new theme tokens.
function mermaid_onThemeChange() {
	for (let i = mermaid_instances.length - 1; i >= 0; i--) {
		if (!mermaid_instances[i].elem.isConnected) {
			mermaid_instances.splice(i, 1);
		}
	}
	for (const inst of mermaid_instances) {
		inst.scale = 1;
		inst.tx = 0;
		inst.ty = 0;
		mermaid_renderInst(inst);
	}
}

// Wire up theme-change triggers. Both fire mermaid_onThemeChange:
//   1. MutationObserver on <html>'s class attribute — for explicit DarkTheme/LightTheme toggling.
//   2. matchMedia "prefers-color-scheme" — for OS-level theme changes when the page is in
//      default (system-follows) mode.
(function() {
	if (typeof MutationObserver !== "undefined") {
		const html = document.documentElement;
		let lastClass = html.className;
		new MutationObserver(function() {
			if (html.className === lastClass) {
				return;
			}
			lastClass = html.className;
			mermaid_onThemeChange();
		}).observe(html, {attributes: true, attributeFilter: ["class"]});
	}
	if (typeof window.matchMedia === "function") {
		const mq = window.matchMedia("(prefers-color-scheme: dark)");
		if (mq.addEventListener) {
			mq.addEventListener("change", mermaid_onThemeChange);
		}
	}
})();

// mermaid_wireZoomPan adds wheel-to-zoom, drag-to-pan, pinch-to-zoom, and a
// toolbar (zoom in / zoom out / reset) to the rendered SVG. State lives on
// `inst` so it survives across re-renders.
//
// Transform is applied via CSS on the <svg> element (not the SVG transform
// attribute on an inner <g>) so the coordinate system is client pixels —
// 1px of mouse movement equals 1px of diagram movement regardless of the
// SVG's internal viewBox scale.
function mermaid_wireZoomPan(inst, svg) {
	if (!svg) {
		return;
	}
	svg.style.transformOrigin = "0 0";

	const apply = function() {
		svg.style.transform = "translate(" + inst.tx + "px," + inst.ty + "px) scale(" + inst.scale + ")";
	};
	apply();

	const minScale = 0.2;
	const maxScale = 8;
	const clamp = function(s) {
		return Math.max(minScale, Math.min(maxScale, s));
	};

	// Wheel zoom anchored at the cursor. Cursor position is measured relative
	// to the (untransformed) container — using svg.getBoundingClientRect would
	// give a moving reference frame because the SVG itself is what we're
	// transforming.
	svg.addEventListener("wheel", function(e) {
		e.preventDefault();
		const rect = inst.elem.getBoundingClientRect();
		const cx = e.clientX - rect.left;
		const cy = e.clientY - rect.top;
		const factor = Math.exp(-e.deltaY * 0.005);
		const newScale = clamp(inst.scale * factor);
		const k = newScale / inst.scale;
		inst.tx = cx - k * (cx - inst.tx);
		inst.ty = cy - k * (cy - inst.ty);
		inst.scale = newScale;
		apply();
	}, {passive: false});

	// Mouse drag to pan.
	let dragging = false;
	let dragX = 0;
	let dragY = 0;
	svg.addEventListener("mousedown", function(e) {
		dragging = true;
		dragX = e.clientX;
		dragY = e.clientY;
		svg.style.cursor = "grabbing";
		e.preventDefault();
	});
	window.addEventListener("mousemove", function(e) {
		if (!dragging) {
			return;
		}
		inst.tx += e.clientX - dragX;
		inst.ty += e.clientY - dragY;
		dragX = e.clientX;
		dragY = e.clientY;
		apply();
	});
	window.addEventListener("mouseup", function() {
		if (dragging) {
			dragging = false;
			svg.style.cursor = "grab";
		}
	});
	svg.style.cursor = "grab";

	// Touch: one finger pans, two fingers pinch-zoom.
	let pinch = null;
	svg.addEventListener("touchstart", function(e) {
		if (e.touches.length === 1) {
			dragging = true;
			dragX = e.touches[0].clientX;
			dragY = e.touches[0].clientY;
		} else if (e.touches.length === 2) {
			dragging = false;
			const dx = e.touches[0].clientX - e.touches[1].clientX;
			const dy = e.touches[0].clientY - e.touches[1].clientY;
			pinch = {dist: Math.hypot(dx, dy), scale: inst.scale};
		}
	}, {passive: true});
	svg.addEventListener("touchmove", function(e) {
		if (e.touches.length === 1 && dragging) {
			inst.tx += e.touches[0].clientX - dragX;
			inst.ty += e.touches[0].clientY - dragY;
			dragX = e.touches[0].clientX;
			dragY = e.touches[0].clientY;
			apply();
			e.preventDefault();
		} else if (e.touches.length === 2 && pinch) {
			const dx = e.touches[0].clientX - e.touches[1].clientX;
			const dy = e.touches[0].clientY - e.touches[1].clientY;
			const dist = Math.hypot(dx, dy);
			inst.scale = clamp(pinch.scale * (dist / pinch.dist));
			apply();
			e.preventDefault();
		}
	}, {passive: false});
	svg.addEventListener("touchend", function(e) {
		if (e.touches.length === 0) {
			dragging = false;
			pinch = null;
		} else if (e.touches.length < 2) {
			pinch = null;
		}
	});

	// Double-click anywhere on the SVG resets to the original framing.
	svg.addEventListener("dblclick", function(e) {
		e.preventDefault();
		inst.scale = 1;
		inst.tx = 0;
		inst.ty = 0;
		apply();
	});
}
