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
	"net/http"

	"github.com/microbus-io/bespa"
	"github.com/microbus-io/bespa/chart"
	_ "github.com/microbus-io/bespa/chart/maps"
	"github.com/microbus-io/bespa/code"
	"github.com/microbus-io/bespa/mermaid"
	"github.com/microbus-io/bespa/richedit"
)

var (
	wf = struct {
		bespa.DefaultFactory
		richedit.RichEditFactory
		chart.ChartFactory
		code.CodeFactory
		mermaid.MermaidFactory
	}{}
	mux *http.ServeMux
)

// Init registers the web handlers with the global mux.
func Init(m *http.ServeMux) error {
	mux = m
	mux.HandleFunc("/showcase/overview", HandleHome)
	mux.HandleFunc("/showcase/tab-switcher", HandleTabSwitcher)
	mux.HandleFunc("/showcase/gallery", HandleGallery)
	mux.HandleFunc("/showcase/deck", HandleDeck)
	mux.HandleFunc("/showcase/text-formatting", HandleTextFormatting)
	mux.HandleFunc("/showcase/states", HandleStates)
	mux.HandleFunc("/showcase/state", HandleState)
	mux.HandleFunc("/showcase/query-states", HandleQueryStates)
	mux.HandleFunc("/showcase/form-input", HandleFormInput)
	mux.HandleFunc("/showcase/form-validation", HandleFormValidation)
	mux.HandleFunc("/showcase/dir-list", HandleDirList)
	mux.HandleFunc("/showcase/dir-edit", HandleDirEdit)
	mux.HandleFunc("/showcase/toolbar", HandleToolbar)
	mux.HandleFunc("/showcase/navigation", HandleNavigation)
	mux.HandleFunc("/showcase/receiver", HandleReceiver)
	mux.HandleFunc("/showcase/progress", HandleProgress)
	mux.HandleFunc("/showcase/finite-progress-status", HandleFiniteProgressStatus)
	mux.HandleFunc("/showcase/infinite-progress-status", HandleInfiniteProgressStatus)
	mux.HandleFunc("/showcase/charts", HandleCharts)
	mux.HandleFunc("/showcase/code", HandleCode)
	mux.HandleFunc("/showcase/mermaid", HandleMermaid)
	return nil
}
