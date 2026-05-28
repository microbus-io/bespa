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

// richedit_init is called inline from each rendered widget. It hides the form's
// <textarea>, mounts a Quill editor next to it, and mirrors content changes back
// into the textarea so the standard bespa auto-submit / form-post path picks them up.
function richedit_init(opts) {
	const textarea = document.getElementById(opts.textareaID);
	if (!textarea) {
		return;
	}
	if (textarea.dataset.richeditInitialized === "1") {
		return; // Guard against double-init when widgets re-render.
	}
	textarea.dataset.richeditInitialized = "1";

	// Build container: editor mounts above the now-hidden textarea, inside the
	// wrapper <span> the Go side rendered.
	const wrapper = document.getElementById("richedit" + opts.textareaID);
	const editorDiv = document.createElement("div");
	editorDiv.className = "RichEditEditor";
	wrapper.insertBefore(editorDiv, textarea);
	textarea.style.display = "none";

	// Pre-fill content from the textarea before Quill mounts so the editor inherits it.
	const initialHTML = textarea.value || textarea.textContent || "";
	editorDiv.innerHTML = initialHTML;

	const modules = {
		toolbar: richedit_toolbar(opts.toolbar || []),
	};
	if (opts.mentions && opts.mentions.length > 0) {
		modules.mention = richedit_mentionConfig(opts.mentions);
	}

	const quill = new Quill(editorDiv, {
		theme: "snow",
		placeholder: opts.placeholder || "",
		readOnly: !!opts.disabled,
		modules: modules,
	});

	if (opts.autoFocus) {
		quill.focus();
	}

	// Mirror Quill content changes into the hidden textarea. The editor's HTML is
	// available at quill.root.innerHTML. Quill emits a "text-change" event on every
	// edit, including paste, formatting, and undo/redo.
	quill.on("text-change", function() {
		const html = quill.root.innerHTML;
		// Quill leaves an "<p><br></p>" placeholder for an empty editor; normalize.
		textarea.value = (html === "<p><br></p>") ? "" : html;
		// Reuse the standard input auto-submit pipeline.
		if (typeof input_autoSubmit === "function") {
			input_autoSubmit(textarea, true);
		}
	});
}

// richedit_toolbar translates the legacy ckeditor-style button name list into the
// nested array Quill expects. A "|" token starts a new toolbar group (row).
// Unknown tokens are skipped silently.
function richedit_toolbar(tokens) {
	const groups = [[]];
	for (const t of tokens) {
		if (t === "|") {
			groups.push([]);
			continue;
		}
		const mapped = richedit_mapToken(t);
		if (mapped == null) {
			continue;
		}
		// Some tokens (e.g. "alignment") expand to multiple controls.
		if (Array.isArray(mapped)) {
			groups[groups.length - 1].push(...mapped);
		} else {
			groups[groups.length - 1].push(mapped);
		}
	}
	// Drop empty trailing groups.
	return groups.filter(function(g) { return g.length > 0; });
}

function richedit_mapToken(t) {
	switch (t) {
		case "bold":                return "bold";
		case "italic":              return "italic";
		case "underline":           return "underline";
		case "strikethrough":       return "strike";
		case "code":                return "code";
		case "codeBlock":           return "code-block";
		case "subscript":           return {script: "sub"};
		case "superscript":         return {script: "super"};
		case "bulletedList":        return {list: "bullet"};
		case "numberedList":        return {list: "ordered"};
		case "todoList":            return {list: "check"};
		case "outdent":             return {indent: "-1"};
		case "indent":              return {indent: "+1"};
		case "blockQuote":          return "blockquote";
		case "removeFormat":        return "clean";
		case "link":                return "link";
		case "image":               return "image";
		case "mediaEmbed":          return "video";
		case "heading":             return [{header: [1, 2, 3, 4, 5, 6, false]}];
		case "fontColor":           return {color: []};
		case "fontBackgroundColor": return {background: []};
		case "fontFamily":          return {font: []};
		case "fontSize":            return {size: ["small", false, "large", "huge"]};
		case "alignment":           return {align: []};
		case "alignment:left":      return {align: ""};
		case "alignment:center":    return {align: "center"};
		case "alignment:right":     return {align: "right"};
		case "alignment:justify":   return {align: "justify"};
		default:
			// Unsupported tokens (undo/redo/specialCharacters/horizontalLine/findAndReplace/
			// sourceEditing/pageBreak/insertTable/tableColumn/tableRow/mergeTableCells/
			// tableCellProperties/tableProperties/selectAll) — silently skip.
			return null;
	}
}

// richedit_mentionConfig assembles the quill-mention configuration from the
// list of feeds provided by the server. Each feed has its own marker and
// item list; quill-mention dispatches by the typed character.
function richedit_mentionConfig(feeds) {
	const byMarker = {};
	const markers = [];
	for (const f of feeds) {
		byMarker[f.marker] = {minChars: f.minimumCharacters || 0, items: f.feed || []};
		markers.push(f.marker);
	}
	return {
		mentionDenotationChars: markers,
		allowedChars: /^[\p{L}\p{N}\s'-]*$/u,
		minChars: 0,
		source: function(searchTerm, renderList, mentionChar) {
			const cfg = byMarker[mentionChar];
			if (!cfg) {
				renderList([], searchTerm);
				return;
			}
			if (searchTerm.length < cfg.minChars) {
				renderList([], searchTerm);
				return;
			}
			const needle = searchTerm.toLowerCase();
			const matches = [];
			for (let i = 0; i < cfg.items.length; i++) {
				const v = cfg.items[i];
				if (!needle || v.toLowerCase().includes(needle)) {
					matches.push({id: i, value: v});
				}
			}
			renderList(matches, searchTerm);
		},
	};
}
