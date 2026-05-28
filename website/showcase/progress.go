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

package showcase

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

var progress = 0

// HandleProgress demonstrates the progress widget.
func HandleProgress(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.AppBar("Progress widget"),

		wf.HeadlineMedium("Infinite"),
		"This is a static progress bar set to an infinite value.",
		wf.SpacerParagraph(),
		wf.Progress().WithMax(1).WithValue(-1),
		wf.SpacerParagraph(),

		"This dynamic progress bar refreshes its state from the backend. When it completes, it updates its value to complete.",
		wf.SpacerParagraph(),
		wf.Progress().WithMax(100).WithValue(-1).WithRefreshURL("infinite-progress-status"),

		wf.HeadlineMedium("Finite"),
		"This is a static progress bar set to a finite value.",
		wf.SpacerParagraph(),
		wf.Progress().WithMax(100).WithValue(60),
		wf.SpacerParagraph(),

		wf.Printf(
			`This dynamic progress bar refreshes its state from the backend. When it completes, it sets the value of the state variable
			{0} to 1 which results in hiding the progress bar and displaying a completion message instead.`,
			wf.Code("done"),
		),
		wf.SpacerParagraph(),
		wf.Progress().WithMax(100).WithRefreshURL("finite-progress-status").HideIfEq(r, "done", "1").RedrawIfChanged(r, "done"),
		wf.Collection(wf.SpacerBreak(), "Job complete!").HideIfEmpty(r, "done").RedrawIfChanged(r, "done"),
		wf.SpacerParagraph(),

		wf.Debugger(),
	)
	shared.Render(w, r, page)
}

// HandleFiniteProgressStatus provides the status of a progress bar.
func HandleFiniteProgressStatus(w http.ResponseWriter, r *http.Request) {
	if progress == 100 {
		progress = 0
	} else {
		progress += 5
	}
	u := ""
	if progress == 100 {
		u = "?done=1"
	}
	json.NewEncoder(w).Encode(struct {
		Value  int    `json:"value"`
		Stop   bool   `json:"stop"`
		Action string `json:"action"`
	}{
		Value:  progress,
		Stop:   progress == 100,
		Action: u,
	})
}

// HandleInfiniteProgressStatus provides the status of a progress bar.
func HandleInfiniteProgressStatus(w http.ResponseWriter, r *http.Request) {
	value := -1
	if rand.Intn(50) == 0 {
		// Simulate a non-deterministic duration for the job
		value = 100
	}
	// Job complete
	json.NewEncoder(w).Encode(struct {
		Value int `json:"value"`
	}{
		Value: value,
	})
}
