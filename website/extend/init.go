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

package extend

import (
	"net/http"

	"github.com/microbus-io/bespa"
	"github.com/microbus-io/bespa/code"
)

var (
	wf = struct {
		bespa.DefaultFactory
		code.CodeFactory
	}{}
	mux *http.ServeMux
)

// Init registers the web handlers with the global mux.
func Init(m *http.ServeMux) error {
	mux = m
	mux.HandleFunc("/extend/overview", HandleOverview)
	mux.HandleFunc("/extend/anatomy", HandleAnatomy)
	mux.HandleFunc("/extend/composing", HandleComposing)
	mux.HandleFunc("/extend/assets", HandleAssets)
	mux.HandleFunc("/extend/state-aware", HandleStateAware)
	mux.HandleFunc("/extend/form-input-widgets", HandleFormInputs)
	mux.HandleFunc("/extend/packaging", HandlePackaging)
	mux.HandleFunc("/extend/theming", HandleTheming)
	return nil
}

// topicCard builds an outlined card whose entire surface links to a sub-page.
func topicCard(icon, heading, path, desc string) any {
	return wf.CardOutlined().WithHref(path).Add(
		wf.TitleLarge(wf.Icon(icon), " ", heading),
		wf.Spacer(0.25),
		desc,
	)
}
