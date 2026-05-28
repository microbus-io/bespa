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
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/errors"
)

var _ = Widget(&TagWidget{}) // Ensure interface

// TagWidget renders an HTML tag to the page.
type TagWidget struct {
	name       string
	endTag     bool
	attributes map[string]string
	children   []Widget
	err        error
	when       bool
}

// Tag creates a new tag writer with the specified <name>.
// If the tag name is empty, the tag will not render.
func (f WidgetFactory) Tag(name string) *TagWidget {
	return &TagWidget{
		name:       name,
		endTag:     true,
		when:       true,
		attributes: map[string]string{},
	}
}

// NoEnd indicates not to end the tag with </name>.
// By default the tag is ended.
func (tag *TagWidget) NoEnd() *TagWidget {
	tag.endTag = false
	return tag
}

// AttrIf sets the attribute only if the condition is met.
func (tag *TagWidget) AttrIf(condition bool, name string, value string) *TagWidget {
	if condition {
		tag.Attr(name, value)
	}
	return tag
}

// Attr adds an attribute to the tag, i.e. <tag name="value">.
// The attribute is not rendered if the value is empty,
// except in the case of the "value" attribute.
func (tag *TagWidget) Attr(name string, value string) *TagWidget {
	switch {
	case name == "style" || strings.HasPrefix(name, "on"):
		// Concat using semicolons
		if value != "" {
			if tag.attributes[name] != "" {
				tag.attributes[name] += "; "
			}
			tag.attributes[name] += value
		}
	case name == "class":
		// Concat using spaces
		if value != "" {
			if tag.attributes[name] != "" {
				tag.attributes[name] += " "
			}
			tag.attributes[name] += value
		}
	case name == "value" || name == "href":
		// Accept empty value
		tag.attributes[name] = value
	default:
		// Replace
		if value != "" {
			tag.attributes[name] = value
		} else {
			delete(tag.attributes, name)
		}
	}
	return tag
}

// Hide causes the tag not to be displayed, but still rendered.
func (tag *TagWidget) Hide(hide bool) *TagWidget {
	if hide {
		tag.Style("display:none")
	}
	return tag
}

// ClassIf adds to the class attribute only if the condition is met.
func (tag *TagWidget) ClassIf(condition bool, classes ...string) *TagWidget {
	if condition {
		tag.Class(classes...)
	}
	return tag
}

// Class adds to the class attribute of the tag.
func (tag *TagWidget) Class(classes ...string) *TagWidget {
	for _, i := range classes {
		tag.Attr("class", i)
	}
	return tag
}

// Style adds to the style attribute of the tag.
func (tag *TagWidget) Style(styles ...string) *TagWidget {
	for _, i := range styles {
		tag.Attr("style", i)
	}
	return tag
}

// StyleIf adds to the style attribute of the tag only if the condition is met.
func (tag *TagWidget) StyleIf(condition bool, styles ...string) *TagWidget {
	if condition {
		tag.Style(styles...)
	}
	return tag
}

// Add nests other widgets inside this tag.
func (tag *TagWidget) Add(widgets ...any) *TagWidget {
	tag.children = factory.Many(tag.children, widgets)
	return tag
}

// When sets a flag that must be satisfied for the tag to be rendered.
// If the flag is false, an empty placeholder span is rendered instead.
func (tag *TagWidget) When(flag bool) *TagWidget {
	tag.when = flag
	return tag
}

// ID is a unique identifier within the scope of the page.
// Tag widgets disregard their ID.
func (tag *TagWidget) ID() string {
	return ""
}

// SetID sets a unique identifier within the scope of the page.
// Tag widgets disregard their ID.
func (tag *TagWidget) SetID(id string) {
	// Noop
}

// Children are the widgets nested under this widget, or nil if none.
// Bytes widgets cannot have widgets nested under them.
func (tag *TagWidget) Children() []Widget {
	return tag.children
}

// Draw renders the tag, its attributes, and children to the writer.
func (tag *TagWidget) Draw(w io.Writer, r *http.Request) (err error) {
	if tag.name == "" {
		return nil
	}
	if !tag.when {
		id := tag.attributes["data-id"]
		if id != "" {
			tag.write(w, `<span class="Empty`)
			if tag.attributes["class"] != "" {
				tag.write(w, " ")
				tag.write(w, html.EscapeString(tag.attributes["class"]))
			}
			tag.write(w, `" data-id="`)
			tag.write(w, html.EscapeString(id))
			tag.write(w, `"></span>`)
		}
		return errors.Trace(tag.err)
	}
	tag.write(w, "<")
	tag.write(w, html.EscapeString(tag.name))
	for k, v := range tag.attributes {
		tag.write(w, " ")
		tag.write(w, html.EscapeString(k))
		tag.write(w, `="`)
		tag.write(w, html.EscapeString(v))
		tag.write(w, `"`)
	}
	tag.write(w, ">")
	if tag.err != nil {
		return errors.Trace(tag.err)
	}
	if tag.endTag {
		for _, c := range tag.children {
			err = c.Draw(w, r)
			if err != nil {
				return errors.Trace(err)
			}
		}
		tag.write(w, "</")
		tag.write(w, html.EscapeString(tag.name))
		tag.write(w, ">")
	}
	return errors.Trace(tag.err)
}

// Drawn indicates whether this widget needs to be drawn in either a full or partial page rendering.
// Tag widgets are always drawn.
func (tag *TagWidget) Drawn(r *http.Request) bool {
	return true
}

// Shown indicates whether this widget is shown or hidden.
// Tag widgets are always shown.
func (tag *TagWidget) Shown(r *http.Request) bool {
	return true
}

// Render renders the tag, its attributes, and children to a string.
func (tag *TagWidget) String(r *http.Request) string {
	var buf bytes.Buffer
	tag.Draw(&buf, r)
	return buf.String()
}

// write a string to the writer and keep track of the first error.
func (tag *TagWidget) write(w io.Writer, str string) error {
	if tag.err == nil {
		_, err := w.Write([]byte(str))
		if err != nil {
			tag.err = err
		}
	}
	return errors.Trace(tag.err)
}
