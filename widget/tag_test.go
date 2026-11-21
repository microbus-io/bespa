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
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestWidget_TagAttributeOrder(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	for range 20 {
		tag := WidgetFactory{}.Tag("button").
			Attr("type", "submit").
			Attr("name", "continue").
			Attr("disabled", "1").
			Attr("class", "Text").
			Attr("data-id", "x")
		var buf bytes.Buffer
		err := tag.Draw(&buf, r)
		if err != nil {
			t.Fatal(err)
		}
		expected := `<button class="Text" data-id="x" disabled="1" name="continue" type="submit"></button>`
		if buf.String() != expected {
			t.Fatalf("expected %s, got %s", expected, buf.String())
		}
	}
}
