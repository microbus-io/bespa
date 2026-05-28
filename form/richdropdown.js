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

function richdropdown_click(event) {
	const rich = event.currentTarget;
	const box=rich.getBoundingClientRect();
	const pop=rich.lastChild;

	let top;
	let left=box.left+window.pageXOffset;
	if(box.top+pop.offsetHeight>window.innerHeight){
		top=box.top-pop.offsetHeight+window.pageYOffset;
	} else {
		top=box.top+box.height+window.pageYOffset;
	}

	let p=rich;
	while (p){
		const pos=window.getComputedStyle(p,null).getPropertyValue('position');
		if(pos!=='static'){
			const refBox=p.getBoundingClientRect();
			top-=refBox.top-p.scrollTop+window.pageYOffset;
			left-=refBox.left-p.scrollLeft+window.pageXOffset;
			break;
		}
		p=p.parentElement;
	}

	pop.scrollTop=0;
	pop.style.top=(top)+'px';
	pop.style.left=(left)+'px';
	pop.style.width=(box.width)+'px';
	pop.style.maxHeight=(box.height*3.5)+'px';
	if (pop.style.display=='block'){
		pop.style.display='none';
	} else {
		pop.style.display='block';
	}
}

function richdropdown_mouseleave(event) {
	const rich = event.currentTarget;
	rich.lastChild.style.display='none'
}

function richdropdown_optionClick(event) {
	const option = event.currentTarget;
	option.parentNode.parentNode.firstChild.innerHTML=option.innerHTML;
	option.parentNode.parentNode.firstChild.style.visibility = "visible";
	option.parentNode.parentNode.previousSibling.value=option.getAttribute('value');
}
