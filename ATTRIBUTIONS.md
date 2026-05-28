# Third-Party Attributions

This project bundles third-party software and data. Each component is listed below
with its source, license, and where it lives in the source tree.

---

## JavaScript Libraries

### Apache ECharts

- **Location:** `chart/echarts.js`
- **Version:** 5.5.1
- **Source:** https://echarts.apache.org/
- **License:** Apache License 2.0
- **Copyright:** Copyright 2017 The Apache Software Foundation

Apache ECharts is a charting and visualization library used by the `chart` package
to render interactive charts on the client.

License text: https://www.apache.org/licenses/LICENSE-2.0

### Mermaid

- **Location:** `mermaid/mermaid.js`
- **Version:** 11.15.0
- **Source:** https://mermaid.js.org/ — https://github.com/mermaid-js/mermaid
- **License:** MIT License
- **Copyright:** Copyright (c) 2014 - 2022 Knut Sveidqvist

Mermaid is a JavaScript-based diagramming and charting tool that turns text
definitions into SVG diagrams (flowcharts, sequence, class, state, gantt,
and more). The `mermaid` package wraps it as a BESPA widget and themes the
output against the Material design tokens.

License text:

```
MIT License

Copyright (c) 2014 - 2022 Knut Sveidqvist

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### Quill

- **Location:** `richedit/quill.js`, `richedit/quill.snow.css`
- **Version:** 2.0.3
- **Source:** https://quilljs.com/ — https://github.com/slab/quill
- **License:** BSD 3-Clause
- **Copyright:** Copyright (c) 2017, Slab — Copyright (c) 2014, Jason Chen

Quill is the rich-text editor that powers the `richedit` package. The snow theme
CSS ships alongside; both are styled to fit the Material design system via
`richedit/inputrichtext.css`.

License text:

```
BSD 3-Clause License

Copyright (c) 2017, Slab
Copyright (c) 2014, Jason Chen

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its contributors
  may be used to endorse or promote products derived from this software
  without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```

### quill-mention

- **Location:** `richedit/quill.mention.js`, `richedit/quill.mention.css`
- **Version:** 6.0.2
- **Source:** https://github.com/quill-mention/quill-mention
- **License:** MIT License
- **Copyright:** Copyright (c) 2017 João Pedro Schmitz

quill-mention adds `@` / `#` autocomplete feeds to Quill; consumed by the
`richedit` package's `WithMentionFeed(…)` API.

License text:

```
MIT License

Copyright (c) 2017 João Pedro Schmitz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## Go Libraries (bundled)

### Chroma

- **Imported by:** `code` package
- **Source:** https://github.com/alecthomas/chroma
- **License:** MIT License
- **Copyright:** Copyright (c) 2017 Alec Thomas

Chroma is a server-side syntax highlighter with ~300 embedded lexers. The
`code` package uses it to tokenize source code into class-tagged HTML which
is then styled by `code/codeblock.css` against the Material design tokens.

License text:

```
The MIT License (MIT)

Copyright (c) 2017 Alec Thomas

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## Map Data

### World countries and US states (Apache ECharts examples)

- **Location:** `chart/maps/world.json`, `chart/maps/us.json`
- **Source:** https://echarts.apache.org/examples/data/asset/geo/
- **License:** Apache License 2.0
- **Copyright:** Copyright The Apache Software Foundation

`world.json` is a `FeatureCollection` of country outlines. `us.json` is a
`FeatureCollection` of US state boundaries (50 states + DC + Puerto Rico).
Both are sourced from the Apache ECharts examples directory and inherit the
project's Apache 2.0 license.

### Country subdivision maps (Click That 'Hood)

- **Location:** `chart/maps/ca.json`, `chart/maps/au.json`, `chart/maps/it.json`,
  `chart/maps/de.json`, `chart/maps/jp.json`, `chart/maps/cn.json`,
  `chart/maps/in.json`
- **Source:** https://github.com/codeforgermany/click_that_hood
- **License:** MIT License
- **Copyright:** Copyright (c) 2013-2021 Code for America

These per-country GeoJSON files contain administrative subdivisions:

| File | Country | Features |
|---|---|---|
| `ca.json` | Canada | provinces and territories |
| `au.json` | Australia | states and territories |
| `it.json` | Italy | regions |
| `de.json` | Germany | states (Länder) |
| `jp.json` | Japan | prefectures |
| `cn.json` | China | provinces |
| `in.json` | India | states and union territories |

License text:

```
MIT License

Copyright (c) 2013-2021  Code for America

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## Fonts

### Roboto Flex / Roboto Mono

- **Location:** `basic/roboto-flex-latin.woff2`, `basic/roboto-mono-latin.woff2`
- **License:** SIL Open Font License 1.1
- **See:** `basic/roboto-readme.txt`

### Material Symbols Outlined

- **Location:** `basic/material-symbols-outlined.woff2`
- **License:** Apache License 2.0
- **See:** `basic/material-symbols-readme.txt`

---

## Go Dependencies

Direct module dependencies are listed in `go.mod`. Their licenses can be inspected with:

```
go mod download
go-licenses report ./... > licenses.txt
```

At time of writing, the direct dependencies are:

- `github.com/gomarkdown/markdown` — BSD-2-Clause
- `github.com/microbus-io/copyrighter` — Apache License 2.0
- `golang.org/x/text` — BSD-3-Clause
- `github.com/andybalholm/brotli` — MIT
- `golang.org/x/net` — BSD-3-Clause
