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

function progress_start(id) {
    const elem = document.getElementById(id);
    const progress = elem.firstChild;
    const refresh = elem.getAttribute("data-refresh");
    if (!refresh) {
        return;
    }
    const interval = elem.getAttribute("data-interval");
    const timeoutID = window.setInterval(async function() {
        let status = {
            value: 0,
            stop: false,
            action: ""
        };
        try {
            const response = await fetch(refresh, {
                method: "GET"
            });
            status = await response.json();
        } catch (e) {
        }
        if (status.value>=0) {
            progress.setAttribute("value", status.value);
            elem.classList.remove("Infinite");
            if (progress.value>=progress.max || status.stop) {
                window.clearInterval(timeoutID);
            }
            if (status.action) {
                const clickEvent = new Event('click', {
                    bubbles: true,
                    cancelable: true,
                });
                elem.lastChild.setAttribute("href", status.action);
                elem.lastChild.dispatchEvent(clickEvent);
            }
        } else {
            progress.removeAttribute("value");
            elem.classList.add("Infinite");
        }
    }, interval);
}
