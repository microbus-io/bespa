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

import (
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&InputTextWidget{})      // Ensure interface
var _ = InputWidget(&InputTextWidget{}) // Ensure interface

// Compiled once at startup and reused by per-request validators.
var (
	emailRE = regexp.MustCompile(`^.+@.+\..{2,}$`)
	telRE   = regexp.MustCompile(`^[0-9,;\*#\+]*$`)
	colorRE = regexp.MustCompile(`^#[0-9a-f]{6}$`)
)

// InputTextWidget renders a textual input field.
type InputTextWidget struct {
	*widget.InputWidgetBase[*InputTextWidget]
	htmlType     string
	value        string
	placeholder  string
	maxLen       int
	minLen       int
	rows         int
	validators   Predicates
	predicates   Predicates
	pattern      string
	min          string
	max          string
	step         string
	autoComplete string
	width        string
	trimSpaces   bool
	errMsg       string
	autoFocus    bool
	otp          bool
	mode         string
}

// InputText creates a new widget that renders a single-line text input.
// name is the state variable the field reads/writes; value is the initial
// content (used until the user types something). Defaults: leading and
// trailing whitespace are trimmed, browser autocomplete is off, no length
// limits.
func (f FormFactory) InputText(name string, value string) *InputTextWidget {
	x := &InputTextWidget{
		htmlType:     "text",
		value:        value,
		rows:         1,
		minLen:       -1, // Unbound
		maxLen:       -1, // Unbound
		trimSpaces:   true,
		autoComplete: "off",
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// InputPassword creates a new widget for a password field — masked input,
// no whitespace trimming.
func (f FormFactory) InputPassword(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	x.htmlType = "password"
	x.trimSpaces = false
	return x
}

// InputEmail creates a new widget for an email field with built-in
// RFC-5322 + simple TLD validation.
func (f FormFactory) InputEmail(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	x.htmlType = "email"
	x.validators.Add(func(value string) (bool, string) {
		em, err := mail.ParseAddress(value)
		if err != nil {
			return false, "Invalid email address"
		}
		return emailRE.MatchString(em.Address), "Invalid email address"
	})
	return x
}

// InputURL creates a new widget for a URL field. A missing scheme is
// accepted (https:// is assumed for validation); the on-screen keyboard
// shows URL-friendly characters on mobile.
func (f FormFactory) InputURL(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	// HTML type URL requires entering the protocol https:// which most users find confusing
	// x.htmlType = "url"
	x.mode = "url"
	x.validators.Add(func(value string) (bool, string) {
		if !strings.Contains(value, "://") {
			value = "https://" + value
		}
		u, err := url.Parse(value)
		return err == nil && u.Hostname() != "", "Invalid URL"
	})
	return x
}

// InputDate creates a new widget for a date picker. Pass a zero time to
// leave the field empty. Adjust value into the user's time zone before
// passing it in — the field uses the raw Y-M-D from the value.
func (f FormFactory) InputDate(name string, value time.Time) *InputTextWidget {
	v := ""
	if !value.IsZero() {
		v = value.Format("2006-01-02")
	}
	x := f.InputText(name, v)
	x.htmlType = "date"
	x.validators.Add(func(value string) (bool, string) {
		_, err := time.Parse("2006-01-02", value)
		if err != nil {
			return false, "Invalid date"
		}
		return true, ""
	})
	return x
}

// InputMonth creates a new widget for a year-month picker. Pass a zero
// time to leave the field empty. Adjust value into the user's time zone
// before passing it in.
func (f FormFactory) InputMonth(name string, value time.Time) *InputTextWidget {
	v := ""
	if !value.IsZero() {
		v = value.Format("2006-01")
	}
	x := f.InputText(name, v)
	x.htmlType = "month"
	x.validators.Add(func(value string) (bool, string) {
		_, err := time.Parse("2006-01", value)
		if err != nil {
			return false, "Invalid month"
		}
		return true, ""
	})
	return x
}

// InputTime creates a new widget for a time-of-day picker. Pass a zero
// time to leave the field empty. Adjust value into the user's time zone
// before passing it in.
func (f FormFactory) InputTime(name string, value time.Time) *InputTextWidget {
	v := ""
	if !value.IsZero() {
		v = value.Format("15:04:05")
	}
	x := f.InputText(name, v)
	x.htmlType = "time"
	x.validators.Add(func(value string) (bool, string) {
		_, err := time.Parse("15:04:05", value)
		if err != nil {
			return false, "Invalid time"
		}
		return true, ""
	})
	return x
}

// InputInteger creates a new widget for a whole-number field. value may
// be an int or a pre-formatted string; anything else leaves the field
// empty. Defaults to ~16 chars wide.
func (f FormFactory) InputInteger(name string, value any) *InputTextWidget {
	v := ""
	if s, ok := value.(string); ok {
		v = s
	}
	if i, ok := value.(int); ok {
		v = strconv.Itoa(i)
	}
	x := f.InputText(name, v)
	x.htmlType = "number"
	x.WithWidth(16)
	x.validators.Add(func(value string) (bool, string) {
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false, "Invalid number"
		}
		return true, ""
	})
	return x
}

// InputDecimal creates a new widget for a fractional-number field.
// value may be a float64 or a pre-formatted string; anything else leaves
// the field empty. precision is the number of digits shown after the
// decimal point in the initial value.
func (f FormFactory) InputDecimal(name string, value any, precision int) *InputTextWidget {
	v := ""
	if s, ok := value.(string); ok {
		v = s
	}
	if f, ok := value.(float64); ok {
		v = strconv.FormatFloat(f, 'f', precision, 64)
	}
	x := f.InputText(name, v)
	x.htmlType = "number"
	x.WithWidth(16)
	x.validators.Add(func(value string) (bool, string) {
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false, "Invalid number"
		}
		return true, ""
	})
	return x
}

// InputTel creates a new widget for a telephone number field. Validates
// against the IETF tel-dial-string characters (digits and , ; * # +).
// For an unvalidated tel field, use InputPhone.
func (f FormFactory) InputTel(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	x.htmlType = "tel"
	x.validators.Add(func(value string) (bool, string) {
		if !telRE.MatchString(value) {
			return false, "Invalid number"
		}
		return true, ""
	})
	return x
}

// InputRange creates a new widget for a slider that emits an integer.
// Use WithMin/WithMax to set the range and WithStep for the granularity.
// Stretches to fill its container.
func (f FormFactory) InputRange(name string, value any) *InputTextWidget {
	x := f.InputInteger(name, value)
	x.htmlType = "range"
	x.WithWidth(-1)
	return x
}

// InputPhone creates a new widget for a phone-number field — same HTML
// type as InputTel, but with no built-in validator so any characters are
// accepted.
func (f FormFactory) InputPhone(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	x.htmlType = "tel"
	return x
}

// InputColor creates a new widget for a hex color picker. value must be
// in lowercase #RRGGBB form, or empty.
func (f FormFactory) InputColor(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	x.htmlType = "color"
	x.WithWidth(8)
	x.validators.Add(func(value string) (bool, string) {
		if !colorRE.MatchString(value) {
			return false, "Invalid RGB color code"
		}
		return true, ""
	})
	return x
}

// InputOneTimePassword creates a new widget for a one-time password code
// entry, styled with letter spacing for character-by-character readability.
// Combine with WithLength to fix the code length.
func (f FormFactory) InputOneTimePassword(name string, value string) *InputTextWidget {
	x := f.InputText(name, value)
	x.otp = true
	return x
}

// WithPlaceholder sets the placeholder text of the field.
func (wgt *InputTextWidget) WithPlaceholder(placeholder string) *InputTextWidget {
	wgt.placeholder = placeholder
	return wgt
}

// WithLength bounds the field length to [minChars, maxChars] in characters.
// Pass a negative value for either bound to leave it unbounded. Note: a
// non-zero minimum does NOT make the field required on its own — pair
// with WithRequired if the field must be non-empty.
func (wgt *InputTextWidget) WithLength(minChars int, maxChars int) *InputTextWidget {
	wgt.minLen = minChars
	wgt.maxLen = maxChars
	return wgt
}

// WithWidth sets the visible width of the field in characters. Pass 0
// (or any non-positive value) to let it stretch to fill the container —
// the default.
func (wgt *InputTextWidget) WithWidth(chars int) *InputTextWidget {
	if chars > 0 {
		// Extra width to accommodate wider letters and the padding
		wgt.width = fmt.Sprintf("width: calc(%fch + 2px)", float32(chars)+1)
		if wgt.otp {
			// Accommodate the letter spacing
			wgt.width = fmt.Sprintf("width: calc(%fch + 2px + %fch)", float32(chars)+1, float32(chars)*0.5)
		}
	} else {
		wgt.width = ""
	}
	return wgt
}

// WithRows turns the field into a multi-line <textarea> of the given
// height. Default is 1 (single-line <input>). Multi-row fields skip
// client-side HTML validation — server-side validators still run.
func (wgt *InputTextWidget) WithRows(rows int) *InputTextWidget {
	wgt.rows = rows
	return wgt
}

// WithAutoComplete determines the browser's auto complete behavior of the field.
// See https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/autocomplete for valid values.
// The common options are "on" and "off", with "off" being the default.
// Use the empty string "" to inherit the behavior of the parent form.
func (wgt *InputTextWidget) WithAutoComplete(autoComplete string) *InputTextWidget {
	wgt.autoComplete = autoComplete
	return wgt
}

// WithTrimSpaces toggles leading/trailing whitespace trimming on the
// posted value. Default is true except for password fields, which are
// always preserved verbatim.
func (wgt *InputTextWidget) WithTrimSpaces(trim bool) *InputTextWidget {
	wgt.trimSpaces = trim
	return wgt
}

// WithAutoFocus automatically focuses the cursor on the text input field.
// Auto-focus is off by default.
func (wgt *InputTextWidget) WithAutoFocus(autoFocus bool) *InputTextWidget {
	wgt.autoFocus = autoFocus
	return wgt
}

// WithMode sets the HTML `inputmode` attribute, hinting which on-screen
// keyboard mobile browsers should show (e.g. "decimal", "tel", "url",
// "email", "search", "numeric").
func (wgt *InputTextWidget) WithMode(inputMode string) *InputTextWidget {
	wgt.mode = inputMode
	return wgt
}

// WithMin sets an inclusive lower bound. Pass an int for integer/range
// fields or a time.Time for date/month/time fields; types that don't
// match the field's HTML type are silently ignored. Adjust times into
// the user's time zone first. The current value is clamped to the new
// bound.
func (wgt *InputTextWidget) WithMin(lowerBound any) *InputTextWidget {
	if min, ok := lowerBound.(int); ok && (wgt.htmlType == "number" || wgt.htmlType == "range") {
		wgt.min = strconv.Itoa(min)
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return false, "Not a valid number"
			}
			if int(v) < min {
				return false, fmt.Sprintf("Must be %d or greater", min)
			}
			return true, ""
		})
	}
	if min, ok := lowerBound.(time.Time); ok && wgt.htmlType == "date" && !min.IsZero() {
		wgt.min = min.Format("2006-01-02")
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := time.ParseInLocation("2006-01-02", value, min.Location())
			if err != nil {
				return false, "Not a valid date"
			}
			if v.Before(min) {
				return false, fmt.Sprintf("Must be %s or later", min.Format("1/2/06"))
			}
			return true, ""
		})
	}
	if min, ok := lowerBound.(time.Time); ok && wgt.htmlType == "month" && !min.IsZero() {
		wgt.min = min.Format("2006-01")
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := time.ParseInLocation("2006-01", value, min.Location())
			if err != nil {
				return false, "Not a valid month"
			}
			if v.Before(min) {
				return false, fmt.Sprintf("Must be %s or later", min.Format("January 2006"))
			}
			return true, ""
		})
	}
	if min, ok := lowerBound.(time.Time); ok && wgt.htmlType == "time" && !min.IsZero() {
		wgt.min = min.Format("15:04:05")
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := time.ParseInLocation("15:04:05", value, min.Location())
			if err != nil {
				return false, "Not a valid time"
			}
			if v.Before(min) {
				return false, fmt.Sprintf("Must be %s or later", min.Format("15:04:05"))
			}
			return true, ""
		})
	}
	wgt.value = wgt.clamp(wgt.value)
	return wgt
}

// WithMax sets an inclusive upper bound. Same type rules as WithMin.
// The current value is clamped to the new bound.
func (wgt *InputTextWidget) WithMax(upperBound any) *InputTextWidget {
	if max, ok := upperBound.(int); ok && (wgt.htmlType == "number" || wgt.htmlType == "range") {
		wgt.max = strconv.Itoa(max)
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return false, "Not a valid number"
			}
			if int(v) > max {
				return false, fmt.Sprintf("Must be %d or less", max)
			}
			return true, ""
		})
	}
	if max, ok := upperBound.(time.Time); ok && wgt.htmlType == "date" && !max.IsZero() {
		wgt.max = max.Format("2006-01-02")
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := time.ParseInLocation("2006-01-02", value, max.Location())
			if err != nil {
				return false, "Not a valid date"
			}
			if v.After(max) {
				return false, fmt.Sprintf("Must be %s or earlier", max.Format("1/2/06"))
			}
			return true, ""
		})
	}
	if max, ok := upperBound.(time.Time); ok && wgt.htmlType == "month" && !max.IsZero() {
		wgt.max = max.Format("2006-01")
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := time.ParseInLocation("2006-01", value, max.Location())
			if err != nil {
				return false, "Not a valid month"
			}
			if v.After(max) {
				return false, fmt.Sprintf("Must be %s or earlier", max.Format("January 2006"))
			}
			return true, ""
		})
	}
	if max, ok := upperBound.(time.Time); ok && wgt.htmlType == "time" && !max.IsZero() {
		wgt.max = max.Format("15:04:05")
		wgt.validators.Add(func(value string) (bool, string) {
			v, err := time.ParseInLocation("15:04:05", value, max.Location())
			if err != nil {
				return false, "Not a valid time"
			}
			if v.After(max) {
				return false, fmt.Sprintf("Must be %s or earlier", max.Format("15:04:05"))
			}
			return true, ""
		})
	}
	wgt.value = wgt.clamp(wgt.value)
	return wgt
}

// WithStep constrains the value to multiples of granularity, measured
// from the field's initial value (or min if no initial value). Units:
// integer/range — units; time — seconds; date — days; month — months.
// Defaults: 1 for numbers/dates/months, 60 for time. Ignored for non-
// numeric/non-temporal fields.
func (wgt *InputTextWidget) WithStep(granularity int) *InputTextWidget {
	supportedTypes := map[string]bool{
		"number": true,
		"range":  true,
		"time":   true,
		"date":   true,
		"month":  true,
	}
	if !supportedTypes[wgt.htmlType] {
		return wgt
	}
	wgt.step = strconv.Itoa(granularity)
	wgt.validators.Add(func(value string) (bool, string) {
		// The step anchor is the initial value if set, otherwise the min.
		if wgt.htmlType == "number" || wgt.htmlType == "range" {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return false, "Not a valid number"
			}
			var base int
			if wgt.value != "" {
				base, _ = strconv.Atoi(wgt.value)
			} else if wgt.min != "" {
				base, _ = strconv.Atoi(wgt.min)
			}
			if (int(v)-base)%granularity != 0 {
				return false, "Not a valid value"
			}
		}
		if wgt.htmlType == "time" {
			v, err := time.Parse("15:04:05", value)
			if err != nil {
				return false, "Not a valid time"
			}
			var base time.Time
			if wgt.value != "" {
				base, _ = time.Parse("15:04:05", wgt.value)
			} else if wgt.min != "" {
				base, _ = time.Parse("15:04:05", wgt.min)
			}
			if int(v.Sub(base).Seconds())%granularity != 0 {
				return false, "Not a valid value"
			}
		}
		if wgt.htmlType == "date" {
			v, err := time.Parse("2006-01-02", value)
			if err != nil {
				return false, "Not a valid date"
			}
			var base time.Time
			if wgt.value != "" {
				base, _ = time.Parse("2006-01-02", wgt.value)
			} else if wgt.min != "" {
				base, _ = time.Parse("2006-01-02", wgt.min)
			}
			if int(v.Sub(base).Hours()/24.0)%granularity != 0 {
				return false, "Not a valid value"
			}
		}
		if wgt.htmlType == "month" {
			v, err := time.Parse("2006-01", value)
			if err != nil {
				return false, "Not a valid month"
			}
			var base time.Time
			if wgt.value != "" {
				base, _ = time.Parse("2006-01", wgt.value)
			} else if wgt.min != "" {
				base, _ = time.Parse("2006-01", wgt.min)
			}
			delta := (v.Year()-base.Year())*12 + (int(v.Month()) - int(base.Month()))
			if delta%granularity != 0 {
				return false, "Not a valid value"
			}
		}
		return true, ""
	})
	return wgt
}

// WithPattern requires the value to match the given regexp. Applies to
// text/email/url/password/tel fields only; silently ignored for others.
// Pattern syntax is RE2 (Go's regexp) on the server; the browser uses its
// own engine for client-side hints.
func (wgt *InputTextWidget) WithPattern(exp string) *InputTextWidget {
	supportedTypes := map[string]bool{
		"text":     true,
		"email":    true,
		"url":      true,
		"password": true,
		"tel":      true,
	}
	if !supportedTypes[wgt.htmlType] {
		return wgt
	}
	wgt.pattern = exp
	wgt.validators.Add(func(value string) (bool, string) {
		match, err := regexp.MatchString(exp, value)
		return match && err == nil, "Must match the requested format"
	})
	return wgt
}

// WithPredicate adds a custom validator. Unlike the built-in validators,
// predicates are also called with an empty value — useful for required-
// like rules with custom messaging. The returned errMsg is shown to the
// user when valid is false.
func (wgt *InputTextWidget) WithPredicate(predicate func(value string) (valid bool, errMsg string)) *InputTextWidget {
	wgt.predicates.Add(predicate)
	return wgt
}

// Value returns the field's effective value: the user's posted input
// when available, or the initial value otherwise. Trimmed (unless
// WithTrimSpaces(false)), clamped to WithMin/WithMax, and normalized for
// times (zero-second suffix). Disabled fields always return the initial
// value.
func (wgt *InputTextWidget) Value(r *http.Request) string {
	value := wgt.value
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) {
		value = state.Get(wgt.Name())
	}
	if wgt.trimSpaces {
		value = strings.TrimSpace(value)
	}
	if wgt.htmlType == "time" && value != "" {
		// The time widget may not post back seconds for round minutes (hh:mm:00)
		parts := strings.Split(value, ":")
		if len(parts) == 2 {
			value += ":00"
		}
	}
	return wgt.clamp(value)
}

// clamp returns value pinned within [wgt.min, wgt.max] according to wgt.htmlType.
// Empty values and unparseable bounds are passed through unchanged.
func (wgt *InputTextWidget) clamp(value string) string {
	if value == "" {
		return value
	}
	switch wgt.htmlType {
	case "number", "range":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return value
		}
		if wgt.min != "" {
			if m, err := strconv.ParseInt(wgt.min, 10, 64); err == nil && v < m {
				return wgt.min
			}
		}
		if wgt.max != "" {
			if m, err := strconv.ParseInt(wgt.max, 10, 64); err == nil && v > m {
				return wgt.max
			}
		}
	case "date":
		return clampTime(value, wgt.min, wgt.max, "2006-01-02")
	case "month":
		return clampTime(value, wgt.min, wgt.max, "2006-01")
	case "time":
		return clampTime(value, wgt.min, wgt.max, "15:04:05")
	}
	return value
}

func clampTime(value, min, max, layout string) string {
	v, err := time.Parse(layout, value)
	if err != nil {
		return value
	}
	if min != "" {
		if m, err := time.Parse(layout, min); err == nil && v.Before(m) {
			return min
		}
	}
	if max != "" {
		if m, err := time.Parse(layout, max); err == nil && v.After(m) {
			return max
		}
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *InputTextWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	// Required
	if value == "" && wgt.Required() {
		// wgt.errMsg = "A value is required"
		return false
	}
	// Predicates
	if ok, errMsg := wgt.predicates.Validate(value); !ok {
		wgt.errMsg = errMsg
		return false
	}
	// Accept empty
	if value == "" {
		return true
	}
	// Length
	if wgt.maxLen >= 0 && len([]rune(value)) > wgt.maxLen {
		return false
	}
	if wgt.minLen >= 0 && len([]rune(value)) < wgt.minLen {
		return false
	}
	// Validators
	if ok, errMsg := wgt.validators.Validate(wgt.Value(r)); !ok {
		wgt.errMsg = errMsg
		return false
	}
	// All good
	return true
}

// Changed indicates if the value of the field changed.
func (wgt *InputTextWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *InputTextWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	invalid := !wgt.Valid(r)
	textTag := Tag("input").
		Attr("type", wgt.htmlType).
		NoEnd().
		Attr("value", value)
	if wgt.rows > 1 {
		textTag = Tag("textarea").Attr("rows", strconv.Itoa(wgt.rows))
		if value != "" {
			textTag.Add(value)
		}
	}
	randomID := widget.RandomAlphaNumID(8)
	errID := randomID + "_err"
	textTag.
		Attr("id", randomID).
		Attr("placeholder", wgt.placeholder).
		Attr("min", wgt.min).
		Attr("max", wgt.max).
		Attr("step", wgt.step).
		Style(wgt.width).
		AttrIf(wgt.Disabled(), "disabled", "1").
		AttrIf(invalid, "aria-invalid", "true").
		AttrIf(invalid && wgt.errMsg != "", "aria-describedby", errID).
		ClassIf(invalid, "Invalid")
	if !wgt.Disabled() {
		textTag.
			Attr("name", wgt.Name()).
			Attr("pattern", wgt.pattern).
			Attr("tabindex", "0").
			AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
			Attr("inputmode", wgt.mode).
			Attr("oninput", "input_input(event)").
			Attr("onchange", "input_change(event)").
			Attr("oninvalid", "input_invalid(event)").
			Attr("autoComplete", wgt.autoComplete).
			AttrIf(wgt.minLen >= 0, "minlength", strconv.Itoa(wgt.minLen)).
			AttrIf(wgt.maxLen >= 0, "maxlength", strconv.Itoa(wgt.maxLen)).
			AttrIf(wgt.Required(), "required", "1").
			AttrIf(wgt.autoFocus, "autofocus", "1").
			ClassIf(wgt.otp, "OTP")
	}
	errTag := Tag("")
	errLive := Tag("")
	if invalid && wgt.errMsg != "" {
		errTag = customValidityScript(randomID, wgt.errMsg)
		// Visually-hidden live region so screen readers announce the error
		// when the input becomes invalid, paired with aria-describedby above
		// so they can also read the message on focus.
		errLive = Tag("span").
			Attr("id", errID).
			Attr("role", "alert").
			Class("VisuallyHidden").
			Add(wgt.errMsg)
	}
	return Tag("span").
		Attr("data-id", wgt.ID()).
		Add(textTag, errTag, errLive).
		When(wgt.Shown(r)).
		Draw(w, r)
}
