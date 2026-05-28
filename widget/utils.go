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
	"math/rand/v2"
	"reflect"
	"strings"

	"github.com/microbus-io/errors"
	"golang.org/x/net/html"
)

// IsNil checks if the interface is nil or points to nil.
func IsNil(object any) bool {
	if object == nil {
		return true
	}
	switch reflect.TypeOf(object).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(object).IsNil()
	}
	return false
}

// RandomAlphaNumID generates a random alphanumeric string of an indicated length.
// The ID is guaranteed to start with a letter so it is safe to use as an identifier of an HTML element
func RandomAlphaNumID(length int) string {
	if length <= 0 {
		return ""
	}
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = letters[rand.IntN(len(letters))]
	}
	if bytes[0] >= '0' && bytes[0] <= '9' {
		bytes[0] = 'Q'
	}
	return string(bytes)
}

// SafeHTML parses and rerenders the unsafe HTML to make sure that it is well-formed.
// It also eliminates scripts and event handlers.
func SafeHTML(unsafeHTML string) (safeHTML string, err error) {
	// Parse
	root, err := html.Parse(strings.NewReader(unsafeHTML))
	if err != nil {
		return "", errors.Trace(err)
	}
	// Eliminate scripts
	var removeScripts func(x *html.Node)
	removeScripts = func(x *html.Node) {
		child := x.FirstChild
		for child != nil {
			next := child.NextSibling
			if child.Type == html.ElementNode && child.Data == "script" {
				x.RemoveChild(child)
			} else {
				kept := child.Attr[:0]
				for _, a := range child.Attr {
					if !strings.HasPrefix(a.Key, "on") {
						kept = append(kept, a)
					}
				}
				child.Attr = kept
				removeScripts(child)
			}
			child = next
		}
	}
	removeScripts(root)
	// Rerender to be sure all tags are closed
	var sb strings.Builder
	err = html.Render(&sb, root)
	if err != nil {
		return "", errors.Trace(err)
	}
	rendered := sb.String()
	rendered = strings.TrimPrefix(rendered, "<html><head></head><body>")
	rendered = strings.TrimSuffix(rendered, "</body></html>")
	return rendered, nil
}
