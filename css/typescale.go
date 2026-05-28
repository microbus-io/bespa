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

package css

import (
	"fmt"
	"io"
	"strings"
)

// https://m3.material.io/styles/typography/type-scale-tokens

// TypeScale is a selection of font styles that can be used across an app,
// ensuring a flexible, yet consistent, style that accommodates a range of purposes.
// The Material Design type scale is a combination of 15 styles, each with an intended
// application and meaning. They’re assigned based on use (such as display or headline),
// and grouped more broadly into categories based on scale (such as large or small).
// Material Design’s default type scale uses Roboto for all titles, labels, and body text,
// creating a cohesive typography experience.
type TypeScale struct {
	DisplayLarge   Font
	DisplayMedium  Font
	DisplaySmall   Font
	HeadlineLarge  Font
	HeadlineMedium Font
	HeadlineSmall  Font
	TitleLarge     Font
	TitleMedium    Font
	TitleSmall     Font
	LabelLarge     Font
	LabelMedium    Font
	LabelSmall     Font
	BodyLarge      Font
	BodyMedium     Font
	BodySmall      Font
}

var DefaultTypeScale = TypeScale{
	DisplayLarge:   Font{DefaultTypeFace.Brand, 48, 42.75, 0, DefaultTypeFace.WeightRegular},
	DisplayMedium:  Font{DefaultTypeFace.Brand, 39, 33.75, 0, DefaultTypeFace.WeightRegular},
	DisplaySmall:   Font{DefaultTypeFace.Brand, 33, 27, 0, DefaultTypeFace.WeightRegular},
	HeadlineLarge:  Font{DefaultTypeFace.Brand, 30, 24, 0, DefaultTypeFace.WeightRegular},
	HeadlineMedium: Font{DefaultTypeFace.Brand, 27, 21, 0, DefaultTypeFace.WeightRegular},
	HeadlineSmall:  Font{DefaultTypeFace.Brand, 24, 18, 0, DefaultTypeFace.WeightRegular},
	TitleLarge:     Font{DefaultTypeFace.Brand, 21, 16.5, 0, DefaultTypeFace.WeightRegular},
	TitleMedium:    Font{DefaultTypeFace.Plain, 18, 12, 0.15, DefaultTypeFace.WeightMedium},
	TitleSmall:     Font{DefaultTypeFace.Plain, 15, 10.5, 0.1, DefaultTypeFace.WeightMedium},
	LabelLarge:     Font{DefaultTypeFace.Plain, 15, 10.5, 0.1, DefaultTypeFace.WeightMedium},
	LabelMedium:    Font{DefaultTypeFace.Plain, 12, 9, 0.5, DefaultTypeFace.WeightMedium},
	LabelSmall:     Font{DefaultTypeFace.Plain, 12, 8.25, 0.5, DefaultTypeFace.WeightMedium},
	BodyLarge:      Font{DefaultTypeFace.Plain, 18, 12, 0.5, DefaultTypeFace.WeightRegular},
	BodyMedium:     Font{DefaultTypeFace.Brand, 15, 10.5, 0.25, DefaultTypeFace.WeightRegular},
	BodySmall:      Font{DefaultTypeFace.Brand, 12, 9, 0.4, DefaultTypeFace.WeightRegular},
}

// WriteCSS writes the 12 fonts of the type scale as CSS variables.
func (ts TypeScale) WriteCSS(w io.Writer) error {
	w.Write([]byte("/* Type scale */\n"))
	w.Write([]byte(":root {\n"))
	ts.writeFont(w, ts.DisplayLarge, "md.sys.typescale.display-large")
	ts.writeFont(w, ts.DisplayMedium, "md.sys.typescale.display-medium")
	ts.writeFont(w, ts.DisplaySmall, "md.sys.typescale.display-small")
	ts.writeFont(w, ts.HeadlineLarge, "md.sys.typescale.headline-large")
	ts.writeFont(w, ts.HeadlineMedium, "md.sys.typescale.headline-medium")
	ts.writeFont(w, ts.HeadlineSmall, "md.sys.typescale.headline-small")
	ts.writeFont(w, ts.TitleLarge, "md.sys.typescale.title-large")
	ts.writeFont(w, ts.TitleMedium, "md.sys.typescale.title-medium")
	ts.writeFont(w, ts.TitleSmall, "md.sys.typescale.title-small")
	ts.writeFont(w, ts.LabelLarge, "md.sys.typescale.label-large")
	ts.writeFont(w, ts.LabelMedium, "md.sys.typescale.label-medium")
	ts.writeFont(w, ts.LabelSmall, "md.sys.typescale.label-small")
	ts.writeFont(w, ts.BodyLarge, "md.sys.typescale.body-large")
	ts.writeFont(w, ts.BodyMedium, "md.sys.typescale.body-medium")
	ts.writeFont(w, ts.BodySmall, "md.sys.typescale.body-small")
	w.Write([]byte("}\n\n"))

	return nil
}

// writeFont writes the properties of a single font as CSS variables.
func (ts TypeScale) writeFont(w io.Writer, f Font, name string) error {
	varName := strings.ReplaceAll(name, ".", "-")
	w.Write([]byte(fmt.Sprintf("\t--%s-font: %s;\n", varName, f.Family)))
	w.Write([]byte(fmt.Sprintf("\t--%s-line-height: %.2fpt;\n", varName, f.LineHeight)))
	w.Write([]byte(fmt.Sprintf("\t--%s-size: %.2fpt;\n", varName, f.Size)))
	w.Write([]byte(fmt.Sprintf("\t--%s-tracking: %.2f;\n", varName, f.Tracking*16.0/f.Size)))
	w.Write([]byte(fmt.Sprintf("\t--%s-weight: %d;\n", varName, f.Weight)))
	return nil
}
