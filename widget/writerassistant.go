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

package widget

import (
	"io"
)

// WriterAssistant wraps an underlying writer,
// adding helper functions for writing strings and keeping track of the first error.
type WriterAssistant struct {
	w   io.Writer
	err error
}

// NewWriterAssistant creates a new assistant for writing to an underlying writer.
func NewWriterAssistant(w io.Writer) *WriterAssistant {
	return &WriterAssistant{w: w}
}

// Write writes bytes to the underlying writer.
func (hw *WriterAssistant) Write(b []byte) (int, error) {
	if hw.err != nil {
		return 0, hw.err
	}
	n, err := hw.w.Write(b)
	if err != nil {
		hw.err = err
	}
	return n, err
}

// WriteString writes one or more strings to the underlying writer.
func (hw *WriterAssistant) WriteString(items ...string) (int, error) {
	if hw.err != nil {
		return 0, hw.err
	}
	nn := 0
	for _, i := range items {
		n, err := hw.w.Write([]byte(i))
		if err != nil {
			hw.err = err
			return nn, hw.err
		}
		nn += n
	}
	return nn, hw.err
}

// Err is the first error encountered during a write operation.
func (hw *WriterAssistant) Err() error {
	return hw.err
}
