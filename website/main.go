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

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/microbus-io/bespa"
	"github.com/microbus-io/bespa/code"
	"github.com/microbus-io/bespa/website/basics"
	"github.com/microbus-io/bespa/website/build"
	"github.com/microbus-io/bespa/website/extend"
	"github.com/microbus-io/bespa/website/showcase"
	"github.com/microbus-io/bespa/widget"
)

var (
	mux = &LoggerMux{http.NewServeMux()}
	wf  = struct {
		bespa.DefaultFactory
		code.CodeFactory
	}{}
)

func main() {
	mux.HandleFunc("/bespa/", widget.AssetRegistry.ServeHTTP) // CSS and JavaScript assets
	mux.HandleFunc("/images/", HandleImages)                  // Trailing slash
	mux.HandleFunc("/", HandleRoot)
	mux.HandleFunc("/llms.txt", HandleLLMs)
	mux.HandleFunc("/start", HandleStart)
	mux.HandleFunc("/profile", HandleProfile)
	mux.HandleFunc("/contact", HandleFindUs)
	mux.HandleFunc("/try/hello", HandleTryHello)
	mux.HandleFunc("/try/counter", HandleTryCounter)
	showcase.Init(mux.ServeMux)
	basics.Init(mux.ServeMux)
	build.Init(mux.ServeMux)
	extend.Init(mux.ServeMux)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

// LoggerMux wraps the standard HTTP mux to log incoming requests.
type LoggerMux struct {
	*http.ServeMux
}

// ServeHTTP delegates the request to the standard HTTP mux after logging it to stdout.
func (lm *LoggerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%-4s %s\n", r.Method, r.URL.String())

	// Brotli compression often yields 75% reduction in page sizes
	var closer io.Closer
	if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
		br := &BrotliResponseWriter{
			ResponseWriter: w,
		}
		w = br
		closer = br
		w.Header().Set("Content-Encoding", "br")
	}
	lm.ServeMux.ServeHTTP(w, r)
	if closer != nil {
		closer.Close()
	}
}

type BrotliResponseWriter struct {
	http.ResponseWriter
	writer *brotli.Writer
}

func (br *BrotliResponseWriter) Write(b []byte) (int, error) {
	if br.writer == nil {
		br.writer = brotli.NewWriter(br.ResponseWriter)
	}
	return br.writer.Write(b)
}

func (br *BrotliResponseWriter) Close() error {
	if br.writer != nil {
		return br.writer.Close()
	}
	return nil
}
