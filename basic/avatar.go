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

package basic

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&AvatarWidget{}) // Ensure interface

// AvatarWidget renders an avatar image.
type AvatarWidget struct {
	*widget.WidgetBase[*AvatarWidget]
	src      string
	name     string
	showName bool
	size     string
}

// Avatar creates a new widget that renders an avatar.
// When imageSrc is empty, the avatar falls back to the user's initials
// (derived from the first and last name parts), or to a generic person
// icon if name is also empty.
func (f BasicFactory) Avatar(name string, imageSrc string) *AvatarWidget {
	x := &AvatarWidget{
		src:  imageSrc,
		name: strings.TrimSpace(name),
	}
	x.WidgetBase = widget.NewWidgetBase(x)
	return x
}

// WithNameLabel displays the full name to the right of the avatar image.
// By default the name is not shown.
func (wgt *AvatarWidget) WithNameLabel(show bool) *AvatarWidget {
	wgt.showName = show
	return wgt
}

// WithSize scales the avatar relative to its default size (1.0 = default,
// 2.0 = double). Values <= 0 are ignored.
func (wgt *AvatarWidget) WithSize(ratio float32) *AvatarWidget {
	if ratio == 1 {
		wgt.size = ""
	} else if ratio > 0 {
		wgt.size = fmt.Sprintf("font-size:%.2fem", ratio)
	}
	return wgt
}

// Draw renders the widget's HTML.
func (wgt *AvatarWidget) Draw(w io.Writer, r *http.Request) (err error) {
	// Calculate initials
	initials := ""
	if wgt.src == "" {
		nameParts := []string{}
		for _, part := range strings.Split(wgt.name, " ") {
			if strings.TrimSpace(part) != "" {
				nameParts = append(nameParts, strings.TrimSpace(part))
			}
		}
		n := len(nameParts)
		if n == 1 {
			initials = string([]rune(nameParts[0])[0])
		} else if n > 1 {
			initials = string([]rune(nameParts[0])[0]) + string([]rune(nameParts[n-1])[0])
		}
		initials = strings.ToUpper(initials)
	}
	var letters Widget
	letters = factory.Icon("person")
	if initials != "" {
		letters = Text(initials)
	}

	background := ""
	if wgt.src != "" {
		background = `background-image:url('` + strings.ReplaceAll(wgt.src, "'", "%27") + `')`
		letters = nil
	}
	return Tag("div").
		Attr("data-id", wgt.ID()).
		Class("Avatar").
		Style(wgt.size).
		Add(
			Tag("div").
				Style(background).
				Add(
					Tag("div").
						When(letters != nil).
						Add(letters)),
			Tag("span").
				When(wgt.showName && wgt.name != "").
				Add(wgt.name)).
		When(wgt.Shown(r)).
		Draw(w, r)
}
