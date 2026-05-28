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

package build

import (
	"net/http"

	"github.com/microbus-io/bespa/website/shared"
)

const tablesBasic = `// store is your data layer; it knows how to query, sort, and paginate.

// 1. Build the table shell first, including its column definitions.
//    Col(visibility, width, alignment) — visibility uses "n"/"w"/"x"
//    letters for narrow/wide/expanded breakpoints.
tbl := wf.Table().
    WithDefaultPageRows(r, 25).
    Add(
        wf.Col("nwx", 12, "left").Add(wf.Sorter("fname", "First name")),
        wf.Col("nwx", 12, "left").Add(wf.Sorter("lname", "Last name")),
        wf.Col("wx", 20, "left").Add("Email"), // not sortable
    )

// 2. Ask the table what its current state is. The table owns the
//    state-variable naming convention; your handler doesn't need to.
rowFrom, rowTo := tbl.DisplayRange(r)
rows, total, _ := store.Query(
    tbl.Query(r),     // current quick-search text
    tbl.SortOrder(r), // current sort column + direction
    rowFrom, rowTo,   // current page window
)

// 3. Tell the table the row count so the paginator can compute the last page.
tbl.WithTotalRows(r, total)

// 4. Add the rows.
for _, p := range rows {
    tbl.Add(wf.Row().Add(
        wf.QuickSearchUnderliner(p.FirstName),
        wf.QuickSearchUnderliner(p.LastName),
        wf.QuickSearchUnderliner(p.Email),
    ).WithHref("/people/" + p.ID))
}
`

const tablesLayout = `wf.Toolbar().
    AddLeft(
        wf.ButtonTonal("").Add(wf.Icon("add"), " New").WithHref("?modal=new"),
        wf.QuickSearch(),
    ).
    AddRight(
        wf.Paginator(),
    ),
tbl,
wf.Toolbar().AddLeft(
    wf.PageSizer(),
),
`

const tablesNamed = `wf.Table().WithName("orders"). // ...
wf.QuickSearch().ForTable("orders")
wf.PageSizer().ForTable("orders")
wf.Paginator().ForTable("orders")
wf.QuickSearchUnderliner(text).ForTable("orders")
`

const tablesSQL = `// store wraps a *sql.DB. It maps the table's accessors onto a single SQL query.
func (s *PeopleStore) Query(q, sort string, from, to int) ([]*Person, int, error) {
    // Whitelist the sort column before interpolating; never trust user input.
    sortSQL := "lname ASC"
    switch sort {
    case "fname":      sortSQL = "fname ASC"
    case "fname desc": sortSQL = "fname DESC"
    case "lname":      sortSQL = "lname ASC"
    case "lname desc": sortSQL = "lname DESC"
    }

    where := ""
    args := []any{}
    if q != "" {
        where = "WHERE fname ILIKE $1 OR lname ILIKE $1 OR email ILIKE $1"
        args = append(args, "%"+q+"%")
    }

    var total int
    if err := s.db.QueryRow("SELECT count(*) FROM people "+where, args...).Scan(&total); err != nil {
        return nil, 0, err
    }

    rows, err := s.db.Query(
        "SELECT id, fname, lname, email FROM people "+where+
            " ORDER BY "+sortSQL+" LIMIT $"+strconv.Itoa(len(args)+1)+" OFFSET $"+strconv.Itoa(len(args)+2),
        append(args, to-from, from)...,
    )
    if err != nil { return nil, 0, err }
    defer rows.Close()

    var out []*Person
    for rows.Next() {
        p := &Person{}
        if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email); err != nil {
            return nil, 0, err
        }
        out = append(out, p)
    }
    return out, total, nil
}
`

const tablesSelection = `// "sel" carries the comma-separated list of selected IDs.
selected := strings.Split(wf.StateOf(r).Get("sel"), ",")
isSelected := func(id string) bool {
    for _, s := range selected {
        if s == id { return true }
    }
    return false
}

// Add a checkbox column as the leftmost column.
tbl.Add(wf.Col("nwx", 1, "center").Add(""))

for _, p := range rows {
    tbl.Add(wf.Row().Add(
        wf.Checkbox("sel_"+p.ID, "").
            WithChecked(isSelected(p.ID)).
            WithAutoSubmit(true),
        wf.QuickSearchUnderliner(p.FirstName),
        // ... other cells ...
    ))
}

// In the handler, on every request: rebuild "sel" from the per-row checkboxes.
sel := []string{}
for _, p := range rows {
    if wf.StateOf(r).Get("sel_"+p.ID) == "true" {
        sel = append(sel, p.ID)
    }
}
wf.StateOf(r).Set("sel", strings.Join(sel, ","))
`

const tablesBulk = `// A bulk-action button reads the selected IDs out of state.
wf.Toolbar().AddLeft(
    wf.QuickSearch(),
    wf.ButtonOutlined("delete").
        Add(wf.Icon("delete"), " Delete selected").
        WithDisabled(wf.StateOf(r).Get("sel") == "").
        RedrawIfChanged(r, "sel"),
)

// In the handler:
if wf.StateOf(r).Get("_submit") == "delete" {
    for _, id := range strings.Split(wf.StateOf(r).Get("sel"), ",") {
        store.Delete(id)
    }
    wf.StateOf(r).Set("sel", "") // clear the selection
}
`

const tablesPerColumnFilter = `// Each filterable column gets its own state variable.
filterFname := wf.StateOf(r).Get("filter_fname")
filterEmail := wf.StateOf(r).Get("filter_email")

tbl.Add(
    wf.Col("nwx", 12, "left").Add(
        wf.Sorter("fname", "First name"),
        wf.InputText("filter_fname", filterFname).
            WithPlaceholder("filter").
            WithAutoSubmit(true),
    ),
    wf.Col("nwx", 20, "left").Add(
        "Email",
        wf.InputText("filter_email", filterEmail).
            WithPlaceholder("filter").
            WithAutoSubmit(true),
    ),
)

// Pass the per-column filters down to the store.
rows, total, _ := store.QueryFiltered(
    map[string]string{"fname": filterFname, "email": filterEmail},
    tbl.SortOrder(r), rowFrom, rowTo,
)
`

// HandleTables covers building data tables.
func HandleTables(w http.ResponseWriter, r *http.Request) {
	page := wf.Page().Add(
		wf.Modal("modal").Add(
			wf.EmbedHandler(mux.ServeHTTP, r, "GET",
				wf.StateOf(r).Get("modal")+"?_back=^?modal=", nil),
		),

		wf.AppBar("Data tables"),

		wf.Markdown(
			"The `table` package wraps the bespa table primitive with the chrome ",
			"that turns rows into a usable list view. A typical data table is the ",
			"`Table` widget plus three support widgets — quick-search, paginator, ",
			"and page-size selector — laid out in a recognizable toolbar-bracketed ",
			"shape.",
		),
		wf.HeadlineMedium("Build the shell, then ask it for state"),
		wf.Markdown(
			"The idiom is to construct the table widget first — its name, ",
			"columns, defaults — and then *ask* it what the current quick-search ",
			"text, sort order, and page range are. The widget owns its ",
			"state-variable naming convention; your handler stays out of it.",
		),
		wf.CodeBlock(tablesBasic).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"The accessors that matter:",
			"\n\n",
			"- `tbl.Query(r)` — the current quick-search text.\n",
			"- `tbl.SortOrder(r)` — current sort column and direction (e.g. ",
			"`\"Last name\"`, `\"Last name desc\"`).\n",
			"- `tbl.DisplayRange(r)` — the row indexes for the current page, ",
			"as a half-open range you can hand straight to a SQL `LIMIT/OFFSET` ",
			"or a slice expression.\n",
			"- `tbl.WithTotalRows(r, n)` — feed back the total row count so the ",
			"paginator can compute the last page number. This is the one thing ",
			"only your handler knows.",
		),
		wf.HeadlineMedium("Support widgets"),
		wf.Markdown(
			"Four widgets pair with a table. They sit outside it and auto-bind ",
			"to the table by name:",
			"\n\n",
			"- `wf.QuickSearch()` — the live search input. Auto-submits on every ",
			"keystroke. Writes to the table's query state; the table reads it ",
			"via `tbl.Query(r)`.\n",
			"- `wf.Paginator()` — page-number navigation. \"First / prev / 1 2 3 ",
			"… / next / last\". Knows the last page number from `WithTotalRows`.\n",
			"- `wf.PageSizer()` — the rows-per-page selector. Lets the user pick ",
			"from a few common sizes.\n",
			"- `wf.QuickSearchUnderliner(content)` — wrap a cell value with this ",
			"to underline the substring matching the current quick-search term. ",
			"Cell-level helper; not part of the toolbar layout.",
		),
		wf.HeadlineMedium("Idiomatic layout"),
		wf.Markdown(
			"Two toolbars bracket the table. The top toolbar carries action ",
			"buttons (\"+ New\", etc.) and the quick-search on the left, and the ",
			"paginator on the right. The bottom toolbar carries the page-size ",
			"selector. This is the shape used by both the data-table and CRUD ",
			"showcases:",
		),
		wf.CodeBlock(tablesLayout).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"Notice no widget calls `RedrawIfChanged` — the table and its ",
			"support widgets each declare their own redraw conditions internally. ",
			"A change to sort, search, or page redraws just the table and the ",
			"paginator, leaving the rest of the page alone.",
		),
		wf.Markdown(
			"For a working version — search, sort, paginate, page-size, plus ",
			"per-row edit and delete actions — open the CRUD example in a modal:",
		),
		wf.ButtonFilled("").
			Add(wf.Icon("play arrow"), " Try it").
			WithHref("?modal=/showcase/dir-list"),
		wf.HeadlineMedium("Multiple tables on one page"),
		wf.Markdown(
			"By default a table picks up the name `\"table\"` and the support ",
			"widgets pick up the same default. For two tables on the same page, ",
			"name one (or both) and point each companion widget at it with ",
			"`ForTable`:",
		),
		wf.CodeBlock(tablesNamed).WithLanguage("go"),
		wf.SpacerBreak(),
		wf.Markdown(
			"State variables become `orders_q`, `orders_sort`, `orders_page` — ",
			"but your handler still doesn't see those names directly; the ",
			"accessor methods `tbl.Query(r)`, etc., still do the right thing.",
		),
		wf.HeadlineMedium("Sorting"),
		wf.Markdown(
			"Marking a column `Sortable(true)` makes its header a clickable ",
			"control that writes the sort order into the table's state. The ",
			"framework doesn't sort the rows for you — your backing store does, ",
			"given `tbl.SortOrder(r)` as a hint. You're free to ignore ",
			"unsupported sort columns or map the column name to a different ",
			"database field.",
		),
		wf.HeadlineMedium("Paging vs. infinite scroll"),
		wf.Markdown(
			"BESPA's table is page-oriented — the server returns a fixed window ",
			"of rows per request. For very large datasets this is usually what ",
			"you want: bounded memory on the server, bounded payload over the ",
			"wire, and the rendered HTML stays small.",
		),
		wf.HeadlineMedium("Backing with `database/sql`"),
		wf.Markdown(
			"Translate the table's three accessors — `Query`, `SortOrder`, ",
			"`DisplayRange` — into one SQL query that does the filter, sort, ",
			"and limit/offset in the database. The handler hands the values ",
			"straight to the store; the store maps them to columns:",
		),
		wf.CodeBlock(tablesSQL).WithLanguage("go"),
		wf.Markdown(
			"Two points worth highlighting:",
			"\n\n",
			"- **Whitelist the sort column.** The framework gives you whatever ",
			"string the user clicked on as the sort key; don't interpolate it ",
			"into SQL. A small `switch` mapping known keys to known column ",
			"names is the cleanest defense.\n",
			"- **Run `COUNT(*)` in the same statement when you can.** ",
			"`tbl.WithTotalRows(r, n)` wants the unfiltered-by-page count. ",
			"On PostgreSQL, a window function (`COUNT(*) OVER ()`) folds the ",
			"count into the row query and skips the second round-trip.",
		),
		wf.HeadlineMedium("Row selection and bulk actions"),
		wf.Markdown(
			"The table itself doesn't carry a notion of \"selected rows\" — ",
			"that's an app-level pattern using checkboxes and a state variable ",
			"that holds the selected IDs. Add a checkbox column and a state ",
			"variable that aggregates the per-row checkbox values:",
		),
		wf.CodeBlock(tablesSelection).WithLanguage("go"),
		wf.Markdown(
			"With selection captured, a bulk-action button reads the IDs out ",
			"of state and runs the operation. The button discriminates by ",
			"name (the `_submit` state var) and disables itself when nothing ",
			"is selected:",
		),
		wf.CodeBlock(tablesBulk).WithLanguage("go"),
		wf.Markdown(
			"For long-running bulk operations, redirect to a background job ",
			"page rather than blocking the request — and consider a confirm ",
			"modal before destructive actions.",
		),
		wf.HeadlineMedium("Per-column filtering"),
		wf.Markdown(
			"Quick-search is one filter against all columns. For more ",
			"control, add an input *inside* each column header — each writes ",
			"to its own state variable and the handler passes a map to the ",
			"store:",
		),
		wf.CodeBlock(tablesPerColumnFilter).WithLanguage("go"),
		wf.Markdown(
			"Per-column filters compose naturally with quick-search — keep ",
			"both if it makes sense for the dataset. The filter inputs use ",
			"`WithAutoSubmit(true)` so the table re-queries on every ",
			"keystroke, the same way quick-search does.",
		),
		wf.HeadlineMedium("See also"),
		wf.Markdown(
			"[Showcase → Data table](/showcase/states) — the worked example with ",
			"the toolbar-bracketed layout, sortable columns, quick-search, and ",
			"the underliner highlighting matches.",
			"\n\n",
			"[Showcase → CRUD](/showcase/dir-list) — the same layout with a \"+ ",
			"New\" action button in the top toolbar and per-row edit / delete ",
			"actions.",
		),
	)
	shared.Render(w, r, page)
}
