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

package showcase

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/microbus-io/bespa/website/shared"
	"github.com/microbus-io/bespa/website/storage"
)

// handleGallery demonstrates the gallery widget.
func HandleGallery(w http.ResponseWriter, r *http.Request) {
	galFixed := wf.Gallery()
	galUneven := wf.Gallery()
	sz, _ := strconv.Atoi(wf.StateOf(r).Get("sizer"))
	if sz < 2 || sz > 8 {
		sz = 4
	}
	consistentRnd := rand.New(rand.NewSource(11))
	states, _, _ := storage.USQuery("", "abbrev", 0, -1)
	for i := 0; i < len(states); i++ {
		r := float32(consistentRnd.Intn(4) + 1) // [1,4]
		abbrev := strings.ToLower(states[i].Abbrev)
		galFixed.Add(wf.Avatar(states[i].Name, "/images/stateflags/"+abbrev+".webp").WithSize(float32(sz) / 2))
		galUneven.Add(wf.Avatar(states[i].Name, "/images/stateflags/"+abbrev+".webp").WithSize(float32(sz) * r / 4).WithNameLabel(true))
	}

	page := wf.Page().Add(
		wf.AppBar("Gallery widget"),
		`The gallery is used to display a collection of nested widgets in a flexible layout that adjusts
		to the width of the screen and the nested widgets. Use the slider to change the size of the state avatars.`,

		wf.SpacerParagraph(),
		wf.InputRange("sizer", 4).WithMin(2).WithMax(8).WithAutoSubmit(true),

		wf.HeadlineMedium("Evenly sized items"),
		galFixed.RedrawIfChanged(r, "sizer"),
		wf.HeadlineMedium("Unevenly sized items"),
		galUneven.RedrawIfChanged(r, "sizer"),
	)
	shared.Render(w, r, page)
}

// HandleDeck demonstrates the deck of cards widget.
func HandleDeck(w http.ResponseWriter, r *http.Request) {
	deck := wf.Deck(2, 4, 5)
	states, _, _ := storage.USQuery("", "abbrev", 0, -1)
	for i := 0; i < len(states); i++ {
		abbrev := strings.ToLower(states[i].Abbrev)
		card := wf.CardFilled().Add(
			wf.BannerImage("/images/stateflags/"+abbrev+".webp").WithHeight(100, "px"),
			wf.TitleMedium(states[i].Name),
			"Population: ", states[i].Population, wf.SpacerNewLine(),
			"Land: ", states[i].Land,
		)
		deck.Add(card)
	}
	page := wf.Page().Add(
		wf.AppBar("Deck of cards widget"),
		`The deck is used to display a collection of nested card widgets in a flexible layout that adjusts
		to the width of the screen and the nested widgets.`,
		deck,
	)
	shared.Render(w, r, page)
}
