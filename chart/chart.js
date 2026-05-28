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

// chart_instances tracks every chart created on the page so they can be re-themed when
// the user switches between light and dark mode without reloading.
const chart_instances = [];

function chart_chart(id, config, renderer) {
	window.setTimeout(function() {
		const elem = document.getElementById(id);
		if (!elem) {
			return;
		}
		const inst = {elem: elem, original: config, renderer: renderer || "canvas"};
		chart_instances.push(inst);
		chart_render(inst);
		if (typeof ResizeObserver !== "undefined") {
			const ch = echarts.getInstanceByDom(elem);
			if (ch) {
				new ResizeObserver(function() { ch.resize(); }).observe(elem);
			}
		}
	}, 1);
}

// chart_render builds a themed copy of the original config and pushes it to the chart instance.
// Called once on creation and again whenever the page theme changes. The original config is
// kept pristine so successive themings start from a clean slate. If the config references a
// named map (e.g. series.map = "USA"), the map's GeoJSON is fetched and registered before
// setOption — ECharts requires the map to be registered up front.
function chart_render(inst) {
	if (!inst.elem.isConnected) {
		return;
	}
	const cfg = chart_deepClone(inst.original);
	chart_applyStyle(cfg);
	chart_resolveConfig(cfg);
	const maps = new Set();
	chart_collectMaps(cfg, maps);
	Promise.all([...maps].map(chart_ensureMap)).then(function() {
		if (!inst.elem.isConnected) {
			return;
		}
		let ch = echarts.getInstanceByDom(inst.elem);
		if (!ch) {
			ch = echarts.init(inst.elem, null, {renderer: inst.renderer});
		}
		ch.setOption(cfg, true);
	});
}

// chart_mapPromises memoizes in-flight or completed map fetches so each map's GeoJSON is
// only downloaded once per page session.
const chart_mapPromises = {};

// chart_ensureMap fetches the GeoJSON for a named ECharts map and registers it. Returns
// a Promise that resolves once the map is available. Already-registered maps short-circuit.
function chart_ensureMap(name) {
	if (echarts.getMap && echarts.getMap(name)) {
		return Promise.resolve();
	}
	if (chart_mapPromises[name]) {
		return chart_mapPromises[name];
	}
	const url = "/bespa/maps/" + name.toLowerCase() + ".json";
	chart_mapPromises[name] = fetch(url).then(function(r) {
		if (!r.ok) {
			throw new Error("failed to load map " + name + ": " + r.status);
		}
		return r.json();
	}).then(function(geo) {
		echarts.registerMap(name, geo);
	});
	return chart_mapPromises[name];
}

// chart_collectMaps walks a resolved config and accumulates every named map referenced.
// Any object with a string "map" property is treated as a map reference — covers series of
// type "map" and top-level geo blocks. False positives are harmless since chart_ensureMap
// is idempotent and short-circuits on already-registered maps.
function chart_collectMaps(obj, out) {
	if (Array.isArray(obj)) {
		for (const x of obj) {
			chart_collectMaps(x, out);
		}
		return;
	}
	if (obj && typeof obj === "object") {
		if (typeof obj.map === "string") {
			out.add(obj.map);
		}
		for (const k in obj) {
			chart_collectMaps(obj[k], out);
		}
	}
}

function chart_deepClone(obj) {
	if (typeof structuredClone === "function") {
		try { return structuredClone(obj); } catch (e) { /* fall through for functions/etc. */ }
	}
	return JSON.parse(JSON.stringify(obj));
}

// chart_onThemeChange re-themes every known chart. Exposed so the framework or app code
// can also trigger re-rendering programmatically if needed.
function chart_onThemeChange() {
	chart_resolveVarCache.clear();
	// Garbage-collect instances whose DOM element was removed (e.g. by a partial redraw).
	for (let i = chart_instances.length - 1; i >= 0; i--) {
		if (!chart_instances[i].elem.isConnected) {
			const stale = echarts.getInstanceByDom(chart_instances[i].elem);
			if (stale) { stale.dispose(); }
			chart_instances.splice(i, 1);
		}
	}
	for (const inst of chart_instances) {
		chart_render(inst);
	}
}

// Wire up theme-change triggers. Both fire chart_onThemeChange:
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
			chart_onThemeChange();
		}).observe(html, {attributes: true, attributeFilter: ["class"]});
	}
	if (typeof window.matchMedia === "function") {
		const mq = window.matchMedia("(prefers-color-scheme: dark)");
		if (mq.addEventListener) {
			mq.addEventListener("change", chart_onThemeChange);
		}
	}
})();

// chart_resolveVar resolves any CSS var() expression in a color string to its computed value.
// ECharts' color parser (zrender) does not evaluate CSS variables — they only resolve when
// the browser applies them to an element. We compute the value against a hidden probe element
// and cache the result. Non-string inputs and strings without var() are passed through unchanged.
const chart_resolveVarCache = new Map();
function chart_resolveVar(value) {
	if (typeof value !== "string" || value.indexOf("var(") < 0) {
		return value;
	}
	const cached = chart_resolveVarCache.get(value);
	if (cached !== undefined) {
		return cached;
	}
	const probe = document.createElement("span");
	probe.style.cssText = "display:none;color:" + value;
	document.body.appendChild(probe);
	const resolved = getComputedStyle(probe).color;
	document.body.removeChild(probe);
	chart_resolveVarCache.set(value, resolved);
	return resolved;
}

// chart_resolveConfig walks an ECharts config in place and resolves CSS var() in any string
// value. Run before setOption so zrender's color parser sees literal color values.
function chart_resolveConfig(obj) {
	if (Array.isArray(obj)) {
		for (let i = 0; i < obj.length; i++) {
			if (typeof obj[i] === "string") {
				obj[i] = chart_resolveVar(obj[i]);
			} else if (obj[i] && typeof obj[i] === "object") {
				chart_resolveConfig(obj[i]);
			}
		}
		return;
	}
	if (obj && typeof obj === "object") {
		for (const k in obj) {
			const v = obj[k];
			if (typeof v === "string") {
				obj[k] = chart_resolveVar(v);
			} else if (v && typeof v === "object") {
				chart_resolveConfig(v);
			}
		}
	}
}

// Patch CanvasRenderingContext2D so direct canvas draws also resolve CSS var() in colors.
// This catches any rendering path that bypasses the config-walker above (e.g. animations,
// gradients computed at draw time). Idempotent — guards against double-patching if chart.js
// is somehow loaded twice.
(function() {
	if (typeof CanvasRenderingContext2D === "undefined") {
		return;
	}
	const proto = CanvasRenderingContext2D.prototype;
	if (proto.__chart_varPatched) {
		return;
	}
	for (const prop of ["fillStyle", "strokeStyle", "shadowColor"]) {
		const desc = Object.getOwnPropertyDescriptor(proto, prop);
		if (!desc || !desc.set) {
			continue;
		}
		Object.defineProperty(proto, prop, {
			configurable: true,
			get: desc.get,
			set: function(v) { desc.set.call(this, chart_resolveVar(v)); },
		});
	}
	proto.__chart_varPatched = true;
})();

// chart_color resolves a Material design CSS custom property (e.g. "--md-sys-color-primary")
// to a literal "rgb(r, g, b)" color string. Returns null if the variable is not defined.
// ECharts' canvas renderer cannot parse CSS var() expressions itself, so values must be
// pre-resolved to literal colors here.
function chart_color(varName, alpha) {
	const v = getComputedStyle(document.documentElement).getPropertyValue(varName).trim();
	if (!v) {
		return null;
	}
	const parts = v.split(/[\s,]+/);
	if (parts.length < 3) {
		return null;
	}
	if (typeof alpha === "number" && alpha < 1) {
		return "rgba(" + parts[0] + "," + parts[1] + "," + parts[2] + "," + alpha + ")";
	}
	return "rgb(" + parts[0] + "," + parts[1] + "," + parts[2] + ")";
}

// chart_cssVar reads a raw CSS custom property value off the document root.
function chart_cssVar(name, fallback) {
	const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
	return v || (fallback ?? "");
}

// chart_cssSize parses a CSS length token like "22px" to a number ECharts can consume.
function chart_cssSize(name, fallback) {
	const v = chart_cssVar(name, "");
	const parsed = parseFloat(v);
	return isNaN(parsed) ? fallback : parsed;
}

function chart_applyStyle(config) {
	const text = chart_color("--md-sys-color-on-background");
	const background = chart_color("--md-sys-color-background");
	const outline = chart_color("--md-sys-color-outline");
	const outlineVariant = chart_color("--md-sys-color-outline-variant");
	const disabledOpacity = parseFloat(chart_cssVar("--md-sys-state-disabled-state-layer-opacity")) || 0.38;
	const disabled = chart_color("--md-sys-color-on-background", disabledOpacity);
	const deemphasized = chart_color("--md-sys-color-on-background", 1 - disabledOpacity);

	// Resolve Material typescale tokens to literal pixel sizes and a literal font-family
	// string, since ECharts' canvas renderer can't evaluate CSS var() at draw time.
	const fontFamily = chart_cssVar("--md-sys-typescale-body-medium-font", "Roboto, sans-serif");
	const titleSize = chart_cssSize("--md-sys-typescale-title-large-size", 22);
	const titleWeight = chart_cssSize("--md-sys-typescale-title-large-weight", 400);
	const subtitleSize = chart_cssSize("--md-sys-typescale-body-medium-size", 14);
	const subtitleWeight = chart_cssSize("--md-sys-typescale-body-medium-weight", 400);
	const bodySize = chart_cssSize("--md-sys-typescale-body-small-size", 12);

	// Global text style
	config.textStyle = config.textStyle ?? {};
	config.textStyle.color = config.textStyle.color ?? text;
	config.textStyle.fontFamily = config.textStyle.fontFamily ?? fontFamily;
	config.textStyle.fontSize = config.textStyle.fontSize ?? bodySize;

	// Background
	config.backgroundColor = config.backgroundColor ?? "transparent";

	// Title
	if (config.title) {
		for (let t of chart_toArray(config.title)) {
			t.textStyle = t.textStyle ?? {};
			t.textStyle.color = t.textStyle.color ?? text;
			t.textStyle.fontSize = t.textStyle.fontSize ?? titleSize;
			t.textStyle.fontWeight = t.textStyle.fontWeight ?? titleWeight;
			t.textStyle.fontFamily = t.textStyle.fontFamily ?? fontFamily;
			t.subtextStyle = t.subtextStyle ?? {};
			t.subtextStyle.color = t.subtextStyle.color ?? deemphasized;
			t.subtextStyle.fontSize = t.subtextStyle.fontSize ?? subtitleSize;
			t.subtextStyle.fontWeight = t.subtextStyle.fontWeight ?? subtitleWeight;
			t.subtextStyle.fontFamily = t.subtextStyle.fontFamily ?? fontFamily;
		}
	}

	// Legend. Default itemGap (10px) is too tight for multi-word labels and the
	// chips end up running together; widen to 24px so neighbouring labels stay
	// visually separated.
	if (config.legend) {
		for (let l of chart_toArray(config.legend)) {
			l.textStyle = l.textStyle ?? {};
			l.textStyle.color = l.textStyle.color ?? text;
			l.inactiveColor = l.inactiveColor ?? disabled;
			l.itemGap = l.itemGap ?? 24;
		}
	}

	// Tooltip
	config.tooltip = config.tooltip ?? {};
	config.tooltip.backgroundColor = config.tooltip.backgroundColor ?? background;
	config.tooltip.borderColor = config.tooltip.borderColor ?? outlineVariant;
	config.tooltip.textStyle = config.tooltip.textStyle ?? {};
	config.tooltip.textStyle.color = config.tooltip.textStyle.color ?? text;

	// Axes (xAxis/yAxis can be a single object or an array)
	const styleAxis = function(axis) {
		axis.axisLine = axis.axisLine ?? {};
		axis.axisLine.lineStyle = axis.axisLine.lineStyle ?? {};
		axis.axisLine.lineStyle.color = axis.axisLine.lineStyle.color ?? text;
		axis.axisLabel = axis.axisLabel ?? {};
		axis.axisLabel.color = axis.axisLabel.color ?? text;
		axis.axisTick = axis.axisTick ?? {};
		axis.axisTick.lineStyle = axis.axisTick.lineStyle ?? {};
		axis.axisTick.lineStyle.color = axis.axisTick.lineStyle.color ?? text;
		axis.splitLine = axis.splitLine ?? {};
		axis.splitLine.lineStyle = axis.splitLine.lineStyle ?? {};
		axis.splitLine.lineStyle.color = axis.splitLine.lineStyle.color ?? outline;
		axis.minorSplitLine = axis.minorSplitLine ?? {};
		axis.minorSplitLine.lineStyle = axis.minorSplitLine.lineStyle ?? {};
		axis.minorSplitLine.lineStyle.color = axis.minorSplitLine.lineStyle.color ?? outlineVariant;
		if (axis.nameTextStyle) {
			axis.nameTextStyle.color = axis.nameTextStyle.color ?? deemphasized;
		} else {
			axis.nameTextStyle = {color: deemphasized};
		}
	};
	for (let a of chart_toArray(config.xAxis)) { styleAxis(a); }
	for (let a of chart_toArray(config.yAxis)) { styleAxis(a); }
	for (let a of chart_toArray(config.radiusAxis)) { styleAxis(a); }
	for (let a of chart_toArray(config.angleAxis)) { styleAxis(a); }

	// Color palette derived from Material tonal primaries at 12 hue rotations (every 30°).
	if (!config.color) {
		config.color = [];
		for (let i = 180; i < 180 + 12*210; i += 210) {
			const deg = i % 360;
			const c = chart_color("--md-sys-color-primary-" + deg + "deg");
			if (c) {
				config.color.push(c);
			}
		}
	}

	// Per-series label defaults. ECharts draws a thick text outline on pie labels by default
	// (textBorderColor "auto" + textBorderWidth 3) which looks chunky against a themed background.
	// Force the label text to the on-background color with no outline; chart authors can override.
	if (config.series) {
		for (let s of chart_toArray(config.series)) {
			s.label = s.label ?? {};
			s.label.color = s.label.color ?? text;
			s.label.textBorderColor = s.label.textBorderColor ?? "transparent";
			s.label.textBorderWidth = s.label.textBorderWidth ?? 0;
			if (s.type === "pie") {
				s.labelLine = s.labelLine ?? {};
				s.labelLine.lineStyle = s.labelLine.lineStyle ?? {};
				s.labelLine.lineStyle.color = s.labelLine.lineStyle.color ?? outline;
			}
		}
	}
}

function chart_toArray(x) {
	if (!x) {
		return [];
	}
	if (Array.isArray(x)) {
		return x;
	}
	return [x];
}
