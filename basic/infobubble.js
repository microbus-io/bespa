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

function infobubble_show(event) {
	const elem=event.currentTarget;
	const box=elem.getBoundingClientRect();
	const pop=elem.nextSibling;

	if (pop.offsetWidth>=window.innerWidth-10) {
		pop.style.width=(window.innerWidth-10)+'px';
	} else {
		pop.style.width="auto";
	}
	
	let top;
	let left;
	if(box.top+pop.offsetHeight>window.innerHeight){
		top=box.top+window.pageYOffset-pop.offsetHeight;
	}else{
		top=box.top+box.height+window.pageYOffset;
	}
	if(box.left+pop.offsetWidth>window.innerWidth){
		left=window.innerWidth-pop.offsetWidth-5+window.pageXOffset;
	}else{
		left=box.left+window.pageXOffset;
	}

	let p=elem;
	while(p){
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
}