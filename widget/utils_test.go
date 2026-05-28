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
	"testing"
)

func TestWidget_SafeHTML(t *testing.T) {
	cases := []string{
		"", "",
		"x", "x",
		"<div>1<p>2<p>3</div></div>", "<div>1<p>2</p><p>3</p></div>",
		"<div>111<script>alert(1);</script>222</div>", "<div>111222</div>",
		"<div>111<SCRIPT>alert(1);</script>222</div>", "<div>111222</div>",
		"<div>111<img onClick='alert(1)'>222</div>", `<div>111<img/>222</div>`,
	}
	for i := 0; i < len(cases); i += 2 {
		safe, err := SafeHTML(cases[i])
		if err != nil {
			t.Error(err)
		}
		if safe != cases[i+1] {
			t.Error(cases[i], safe)
		}
	}
}
