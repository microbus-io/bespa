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

package form

type determination struct {
	ok     bool
	errMsg string
}

// Predicates is a chain of value-validating functions used internally by
// input widgets. App code typically attaches predicates via each input's
// WithPredicate method rather than constructing a Predicates directly.
// Results are memoized by value, so identical values are validated at most
// once per request.
type Predicates struct {
	predicates []func(value string) (bool, string)
	memoized   map[string]determination
}

// Add appends a predicate. nil predicates are skipped during Validate.
func (p *Predicates) Add(predicate func(value string) (ok bool, errMsg string)) {
	p.predicates = append(p.predicates, predicate)
}

// Validate runs each predicate in order against value, short-circuiting
// on the first failure and returning its error message. Returns
// (true, "") if every predicate passes (including the empty case).
func (p *Predicates) Validate(value string) (ok bool, errMsg string) {
	if p.memoized == nil {
		p.memoized = map[string]determination{}
	}
	if cached, ok := p.memoized[value]; ok {
		return cached.ok, cached.errMsg
	}
	for _, predicate := range p.predicates {
		if predicate == nil {
			continue
		}
		var d determination
		d.ok, d.errMsg = predicate(value)
		if !d.ok {
			p.memoized[value] = d
			return false, d.errMsg
		}
	}
	p.memoized[value] = determination{ok: true}
	return true, ""
}
