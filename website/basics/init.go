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

package basics

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
	mux.HandleFunc("/basics/overview", HandleOverview)
	mux.HandleFunc("/basics/cheatsheet", HandleCheatsheet)
	mux.HandleFunc("/basics/declarative-views", HandleDeclarativeViews)
	mux.HandleFunc("/basics/incremental", HandleIncremental)
	mux.HandleFunc("/basics/action-url-pattern", HandleActionURLPattern)
	mux.HandleFunc("/basics/embedded-pages", HandleEmbeddedPages)
	mux.HandleFunc("/basics/nesting", HandleNesting)
	mux.HandleFunc("/basics/nested", HandleNested)
	mux.HandleFunc("/basics/frames", HandleFrames)
	mux.HandleFunc("/basics/frame1", HandleFrameOne)
	mux.HandleFunc("/basics/frame2", HandleFrameTwo)
	mux.HandleFunc("/basics/frameempty", HandleFrameEmpty)
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
