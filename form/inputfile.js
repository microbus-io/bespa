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

function inputfile_keydown(event) {
	const elem = event.currentTarget;
	if (event.keyCode==32 && elem.classList.contains("Upload")) {
		elem.querySelector('INPUT[type="file"]').click();
		event.preventDefault();
		return;
	}
	if ((event.keyCode==8 || event.keyCode==46) && elem.classList.contains("Uploaded")) {
		elem.classList.remove("Uploaded");
		elem.classList.add("Upload");
		event.preventDefault();
		return;
	}
}

function inputfile_click(event) {
	const elem = event.target;
	// Browse for file
	if (elem.classList.contains("DropZone")) {
		elem.parentElement.querySelector('INPUT[type="file"]').click();
		event.stopPropagation();
	}
	// Cancel button
	if (elem.classList.contains("Icon") && elem.parentElement.className=="FileName") {
		elem.previousSibling.innerText = "";
		const e = elem.parentElement.parentElement;
		e.classList.remove("Invalid");
		e.classList.remove("Uploaded");
		e.classList.add("Upload");
		const input = e.querySelector('INPUT[type="hidden"]');
		input.value = "";
		input_autoSubmit(input, false);
	}
}

function inputfile_allowed(e, f) {
	const accept = e.querySelector('INPUT[type="file"]').getAttribute("accept");
	if (!accept) {
		return true;
	}
	let result = false;
	accept.split(",").forEach(function (type) {
		type = type.trim();
		if (!type) {
			return;
		}
		let typePrefix = type;
		if (typePrefix.endsWith("*")) {
			typePrefix = typePrefix.slice(0, -1);
		}
		if (f.name.endsWith(type) || f.type==type || f.type.startsWith(typePrefix)) {
			result = true;
		}
	});
	return result;	
}

function inputfile_draggedFile(event) {
	let f;
	if (event.dataTransfer.items) {
		if (event.dataTransfer.items.length==1 && event.dataTransfer.items[0].kind==='file') {
			f=event.dataTransfer.items[0].getAsFile();
		}
	} else {
		if (event.dataTransfer.files.length==1) {
			f=event.dataTransfer.files[0];
		}
	}
	return f;
}

function inputfile_drop(event) {
	event.preventDefault();
	const dropZone = event.target;
	const f = inputfile_draggedFile(event);
	dropZone.parentElement.classList.remove("Dragging")
	if (!inputfile_allowed(dropZone.parentElement, f)) {
		const e = dropZone.parentElement;
		e.querySelector(".FileName").firstChild.innerText = "Unacceptable file type";
		e.classList.remove("Upload");
		e.classList.add("Uploaded");
		e.classList.add("Invalid");
		return;
	}
	dropZone.parentElement.classList.remove("Invalid")
	var fr = new FileReader();
	fr.onloadend=function(){
		inputfile_upload(dropZone.parentElement, f.name, f.type, fr.result);
	};
	if(f){fr.readAsArrayBuffer(f);}
}

function inputfile_dragover(event) {
	event.preventDefault();
	event.dataTransfer.dropEffect='copy';
}

function inputfile_dragenter(event) {
	event.target.parentElement.classList.add("Dragging")
}
function inputfile_dragleave(event) {
	event.target.parentElement.classList.remove("Dragging")
}

function inputfile_fileSelected(event) {
	event.preventDefault();
	const btn = event.target;
	var f = btn.files[0];
	var fr = new FileReader();
	fr.onloadend=function(){
		inputfile_upload(btn.parentElement, f.name, f.type, fr.result);
		btn.value='';
	};
	if(f){fr.readAsArrayBuffer(f);}
}

async function inputfile_upload(e, fileName, fileType, content) {
	const maxSize = e.getAttribute("data-max-size");
	if (maxSize && content.byteLength>maxSize) {
		e.querySelector(".FileName").firstChild.innerText = "File too large";
		e.classList.remove("Upload");
		e.classList.add("Uploaded");
		e.classList.add("Invalid");
		return;
	}

	const chunkSize = 512*1024; // 512KB
	if (fileType=="") {
		fileType = "application/octet-stream";
	}	
	const url = e.getAttribute("data-receiver");
	let upload;
	try {
		const chunk = content.slice(0, Math.min(chunkSize, content.byteLength));
		const response = await fetch(url + "?name="+fileName, {
			method: "POST",
			body: chunk,
			headers: {"Content-Type": fileType}
		})
		if (response.status!=200) {
			throw new Error(response.statusText);
		}
		upload = await response.json();
	} catch (exc) {
		e.querySelector(".FileName").firstChild.innerText = exc.message;
		e.classList.remove("Upload");
		e.classList.add("Uploaded");
		e.classList.add("Invalid");
		return;
	};

	if (content.byteLength > chunkSize) {
		const progress = e.querySelector("PROGRESS");
		progress.setAttribute("max", content.byteLength)
		progress.setAttribute("value", chunkSize)
		e.classList.remove("Upload");
		e.classList.add("Uploading");
	
		for (let i=chunkSize; i<content.byteLength; i+=chunkSize) {
			progress.setAttribute("value", i);
			try {
				const chunk = content.slice(i, Math.min(i+chunkSize, content.byteLength));
				const response = await fetch(url + "?key="+upload.key, {
					method: "POST",
					body: chunk,
					headers: {"Content-Type": fileType}
				});
				if (response.status!=200) {
					throw new Error(response.statusText);
				}
			} catch (exc) {
				e.querySelector(".FileName").firstChild.innerText = exc.message;
				e.classList.remove("Uploading");
				e.classList.add("Uploaded");
				e.classList.add("Invalid");
				return;
			};
		}
	}

	e.querySelector(".FileName").firstChild.innerText = fileName;
	e.classList.remove("Upload");
	e.classList.remove("Uploading");
	e.classList.add("Uploaded");

	const input = e.querySelector('INPUT[type="hidden"]');
	input.value = upload.key;
	input_autoSubmit(input, false);
}
