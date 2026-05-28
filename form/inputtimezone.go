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
	"io"
	"net/http"
	"strings"

	"github.com/microbus-io/bespa/widget"
)

var _ = Widget(&InputTimeZoneWidget{})      // Ensure interface
var _ = InputWidget(&InputTimeZoneWidget{}) // Ensure interface

// InputTimeZoneWidget renders a time zone selector.
type InputTimeZoneWidget struct {
	*widget.InputWidgetBase[*InputTimeZoneWidget]
	value string
}

// InputTimeZone creates a new widget that renders a paired region / city
// dropdown for picking an IANA time zone name (e.g. "America/Los_Angeles").
// value is the initial selection; pass "" for no preselection. Submitted
// values must match an entry in the bundled IANA name list to pass
// validation.
func (f FormFactory) InputTimeZone(name string, value string) *InputTimeZoneWidget {
	x := &InputTimeZoneWidget{
		value: value,
	}
	x.InputWidgetBase = widget.NewInputWidgetBase(x)
	x.WithName(name)
	return x
}

// Value returns the value of the field.
func (wgt *InputTimeZoneWidget) Value(r *http.Request) string {
	value := wgt.value
	if wgt.Disabled() {
		return value
	}
	state := factory.StateOf(r)
	if state.Has(wgt.Name()) { // || wgt.Submitted(r)
		value = state.Get(wgt.Name())
	}
	return value
}

// Valid validates the field's value against all validators.
func (wgt *InputTimeZoneWidget) Valid(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return true
	}
	value := wgt.Value(r)
	if wgt.Required() && value == "" {
		return false
	}
	if value == "" {
		return true
	}
	for _, v := range timeZoneNames {
		if v == value {
			return true
		}
	}
	return false
}

// Changed indicates if the value of the field changed.
func (wgt *InputTimeZoneWidget) Changed(r *http.Request) bool {
	if wgt.Disabled() || !wgt.Submitted(r) {
		return false
	}
	return wgt.value != wgt.Value(r)
}

// Draw renders the widget's HTML.
func (wgt *InputTimeZoneWidget) Draw(w io.Writer, r *http.Request) (err error) {
	value := wgt.Value(r)
	if value == "" && wgt.Required() {
		value = "US/Pacific"
	}
	invalid := !wgt.Valid(r)
	if invalid {
		if wgt.Required() {
			value = "US/Pacific"
		} else {
			value = ""
		}
	}
	if value == "UTC" {
		value = "UTC/UTC"
	}
	if value == "" {
		value = "/"
	}

	// Region dropdown
	selectIDs := []string{}
	regions := map[string]*widget.TagWidget{}
	randomID := widget.RandomAlphaNumID(8)
	selectIDs = append(selectIDs, randomID)
	tagWorld := Tag("select").
		Attr("id", randomID).
		AttrIf(wgt.Disabled(), "disabled", "1").
		AttrIf(!wgt.Disabled(), "tabindex", "0").
		AttrIf(!wgt.Disabled(), "data-name", wgt.Name()).
		AttrIf(!wgt.Disabled(), "oninput", "inputtimezone_worldInput(event)")
	tzNames := timeZoneNames
	if !wgt.Required() {
		tzNames = append([]string{""}, tzNames...)
	}
	for _, v := range tzNames {
		if v == "UTC" {
			v = "UTC/UTC"
		}
		if v == "" {
			v = "/"
		}
		lastSlash := strings.LastIndex(v, "/")
		if lastSlash < 0 {
			continue
		}
		region := v[:lastSlash]
		tagRegion := regions[region]
		if tagRegion == nil {
			randomID := widget.RandomAlphaNumID(8)
			selectIDs = append(selectIDs, randomID)
			tagRegion = Tag("select").
				Attr("id", randomID)
			if wgt.Disabled() {
				tagRegion.Attr("disabled", "1")
			} else {
				tagRegion.
					Attr("tabindex", "0").
					AttrIf(wgt.AutoSubmit(), "data-autosubmit", "1").
					Attr("oninput", "input_input(event)").
					Attr("oninvalid", "input_invalid(event)").
					Attr("data-region", region)
			}
			regions[region] = tagRegion

			desc := region
			p := strings.LastIndex(desc, "/")
			if p > 0 {
				desc = desc[p+1:]
			}
			desc = strings.ReplaceAll(desc, "_", " ")
			option := Tag("option").
				Attr("value", region).
				Add(desc)
			tagWorld.Add(option)

			if strings.HasPrefix(value, region+"/") {
				option.Attr("selected", "1")
				tagWorld.Attr("value", region)
				tagRegion.Attr("name", wgt.Name())
			} else {
				tagRegion.Class("Off")
			}
			if v == "/" || v == "UTC/UTC" {
				tagRegion.Class("SingleOption")
			}
		}
		desc := v[lastSlash+1:]
		desc = strings.ReplaceAll(desc, "_", " ")
		option := Tag("option").
			Attr("value", strings.TrimSuffix(strings.TrimSuffix(v, "/UTC"), "/")).
			Add(desc)
		if !wgt.Disabled() || v == value {
			tagRegion.Add(option)
		}
		if v == value {
			option.Attr("selected", "1")
			tagRegion.Attr("value", v)
		}
	}
	spanTag := Tag("span").
		Class("InputTimeZone").
		Attr("data-id", wgt.ID()).
		Add(tagWorld)
	for _, region := range regions {
		spanTag.Add(region)
	}
	for _, id := range selectIDs {
		spanTag.Add(
			Tag("script").Add(factory.HTMLUnsafe("input_initBackgroundIconColor('", id, "')")),
		)
	}
	return spanTag.When(wgt.Shown(r)).Draw(w, r)
}

// timeZoneNames is an alphabetically sorted list of most time zone names.
var timeZoneNames = []string{
	"Africa/Abidjan",
	"Africa/Accra",
	"Africa/Addis_Ababa",
	"Africa/Algiers",
	"Africa/Asmara",
	"Africa/Asmera",
	"Africa/Bamako",
	"Africa/Bangui",
	"Africa/Banjul",
	"Africa/Bissau",
	"Africa/Blantyre",
	"Africa/Brazzaville",
	"Africa/Bujumbura",
	"Africa/Cairo",
	"Africa/Casablanca",
	"Africa/Ceuta",
	"Africa/Conakry",
	"Africa/Dakar",
	"Africa/Dar_es_Salaam",
	"Africa/Djibouti",
	"Africa/Douala",
	"Africa/El_Aaiun",
	"Africa/Freetown",
	"Africa/Gaborone",
	"Africa/Harare",
	"Africa/Johannesburg",
	"Africa/Juba",
	"Africa/Kampala",
	"Africa/Khartoum",
	"Africa/Kigali",
	"Africa/Kinshasa",
	"Africa/Lagos",
	"Africa/Libreville",
	"Africa/Lome",
	"Africa/Luanda",
	"Africa/Lubumbashi",
	"Africa/Lusaka",
	"Africa/Malabo",
	"Africa/Maputo",
	"Africa/Maseru",
	"Africa/Mbabane",
	"Africa/Mogadishu",
	"Africa/Monrovia",
	"Africa/Nairobi",
	"Africa/Ndjamena",
	"Africa/Niamey",
	"Africa/Nouakchott",
	"Africa/Ouagadougou",
	"Africa/Porto-Novo",
	"Africa/Sao_Tome",
	"Africa/Timbuktu",
	"Africa/Tripoli",
	"Africa/Tunis",
	"Africa/Windhoek",
	"America/Adak",
	"America/Anchorage",
	"America/Anguilla",
	"America/Antigua",
	"America/Araguaina",
	// "America/Argentina/Buenos_Aires",
	// "America/Argentina/Catamarca",
	// "America/Argentina/ComodRivadavia",
	// "America/Argentina/Cordoba",
	// "America/Argentina/Jujuy",
	// "America/Argentina/La_Rioja",
	// "America/Argentina/Mendoza",
	// "America/Argentina/Rio_Gallegos",
	// "America/Argentina/Salta",
	// "America/Argentina/San_Juan",
	// "America/Argentina/San_Luis",
	// "America/Argentina/Tucuman",
	// "America/Argentina/Ushuaia",
	"America/Aruba",
	"America/Asuncion",
	"America/Atikokan",
	"America/Atka",
	"America/Bahia",
	"America/Bahia_Banderas",
	"America/Barbados",
	"America/Belem",
	"America/Belize",
	"America/Blanc-Sablon",
	"America/Boa_Vista",
	"America/Bogota",
	"America/Boise",
	"America/Buenos_Aires",
	"America/Cambridge_Bay",
	"America/Campo_Grande",
	"America/Cancun",
	"America/Caracas",
	"America/Catamarca",
	"America/Cayenne",
	"America/Cayman",
	"America/Chicago",
	"America/Chihuahua",
	"America/Coral_Harbour",
	"America/Cordoba",
	"America/Costa_Rica",
	"America/Creston",
	"America/Cuiaba",
	"America/Curacao",
	"America/Danmarkshavn",
	"America/Dawson",
	"America/Dawson_Creek",
	"America/Denver",
	"America/Detroit",
	"America/Dominica",
	"America/Edmonton",
	"America/Eirunepe",
	"America/El_Salvador",
	"America/Ensenada",
	"America/Fort_Nelson",
	"America/Fort_Wayne",
	"America/Fortaleza",
	"America/Glace_Bay",
	"America/Godthab",
	"America/Goose_Bay",
	"America/Grand_Turk",
	"America/Grenada",
	"America/Guadeloupe",
	"America/Guatemala",
	"America/Guayaquil",
	"America/Guyana",
	"America/Halifax",
	"America/Havana",
	"America/Hermosillo",
	// "America/Indiana/Indianapolis",
	// "America/Indiana/Knox",
	// "America/Indiana/Marengo",
	// "America/Indiana/Petersburg",
	// "America/Indiana/Tell_City",
	// "America/Indiana/Vevay",
	// "America/Indiana/Vincennes",
	// "America/Indiana/Winamac",
	"America/Indianapolis",
	"America/Inuvik",
	"America/Iqaluit",
	"America/Jamaica",
	"America/Jujuy",
	"America/Juneau",
	// "America/Kentucky/Louisville",
	// "America/Kentucky/Monticello",
	"America/Knox_IN",
	"America/Kralendijk",
	"America/La_Paz",
	"America/Lima",
	"America/Los_Angeles",
	"America/Louisville",
	"America/Lower_Princes",
	"America/Maceio",
	"America/Managua",
	"America/Manaus",
	"America/Marigot",
	"America/Martinique",
	"America/Matamoros",
	"America/Mazatlan",
	"America/Mendoza",
	"America/Menominee",
	"America/Merida",
	"America/Metlakatla",
	"America/Mexico_City",
	"America/Miquelon",
	"America/Moncton",
	"America/Monterrey",
	"America/Montevideo",
	"America/Montreal",
	"America/Montserrat",
	"America/Nassau",
	"America/New_York",
	"America/Nipigon",
	"America/Nome",
	"America/Noronha",
	// "America/North_Dakota/Beulah",
	// "America/North_Dakota/Center",
	// "America/North_Dakota/New_Salem",
	"America/Nuuk",
	"America/Ojinaga",
	"America/Panama",
	"America/Pangnirtung",
	"America/Paramaribo",
	"America/Phoenix",
	"America/Port-au-Prince",
	"America/Port_of_Spain",
	"America/Porto_Acre",
	"America/Porto_Velho",
	"America/Puerto_Rico",
	"America/Punta_Arenas",
	"America/Rainy_River",
	"America/Rankin_Inlet",
	"America/Recife",
	"America/Regina",
	"America/Resolute",
	"America/Rio_Branco",
	"America/Rosario",
	"America/Santa_Isabel",
	"America/Santarem",
	"America/Santiago",
	"America/Santo_Domingo",
	"America/Sao_Paulo",
	"America/Scoresbysund",
	"America/Shiprock",
	"America/Sitka",
	"America/St_Barthelemy",
	"America/St_Johns",
	"America/St_Kitts",
	"America/St_Lucia",
	"America/St_Thomas",
	"America/St_Vincent",
	"America/Swift_Current",
	"America/Tegucigalpa",
	"America/Thule",
	"America/Thunder_Bay",
	"America/Tijuana",
	"America/Toronto",
	"America/Tortola",
	"America/Vancouver",
	"America/Virgin",
	"America/Whitehorse",
	"America/Winnipeg",
	"America/Yakutat",
	"America/Yellowknife",
	"Antarctica/Casey",
	"Antarctica/Davis",
	"Antarctica/DumontDUrville",
	"Antarctica/Macquarie",
	"Antarctica/Mawson",
	"Antarctica/McMurdo",
	"Antarctica/Palmer",
	"Antarctica/Rothera",
	"Antarctica/South_Pole",
	"Antarctica/Syowa",
	"Antarctica/Troll",
	"Antarctica/Vostok",
	"Arctic/Longyearbyen",
	"Asia/Aden",
	"Asia/Almaty",
	"Asia/Amman",
	"Asia/Anadyr",
	"Asia/Aqtau",
	"Asia/Aqtobe",
	"Asia/Ashgabat",
	"Asia/Ashkhabad",
	"Asia/Atyrau",
	"Asia/Baghdad",
	"Asia/Bahrain",
	"Asia/Baku",
	"Asia/Bangkok",
	"Asia/Barnaul",
	"Asia/Beirut",
	"Asia/Bishkek",
	"Asia/Brunei",
	"Asia/Calcutta",
	"Asia/Chita",
	"Asia/Choibalsan",
	"Asia/Chongqing",
	"Asia/Chungking",
	"Asia/Colombo",
	"Asia/Dacca",
	"Asia/Damascus",
	"Asia/Dhaka",
	"Asia/Dili",
	"Asia/Dubai",
	"Asia/Dushanbe",
	"Asia/Famagusta",
	"Asia/Gaza",
	"Asia/Harbin",
	"Asia/Hebron",
	"Asia/Ho_Chi_Minh",
	"Asia/Hong_Kong",
	"Asia/Hovd",
	"Asia/Irkutsk",
	"Asia/Istanbul",
	"Asia/Jakarta",
	"Asia/Jayapura",
	"Asia/Jerusalem",
	"Asia/Kabul",
	"Asia/Kamchatka",
	"Asia/Karachi",
	"Asia/Kashgar",
	"Asia/Kathmandu",
	"Asia/Katmandu",
	"Asia/Khandyga",
	"Asia/Kolkata",
	"Asia/Krasnoyarsk",
	"Asia/Kuala_Lumpur",
	"Asia/Kuching",
	"Asia/Kuwait",
	"Asia/Macao",
	"Asia/Macau",
	"Asia/Magadan",
	"Asia/Makassar",
	"Asia/Manila",
	"Asia/Muscat",
	"Asia/Nicosia",
	"Asia/Novokuznetsk",
	"Asia/Novosibirsk",
	"Asia/Omsk",
	"Asia/Oral",
	"Asia/Phnom_Penh",
	"Asia/Pontianak",
	"Asia/Pyongyang",
	"Asia/Qatar",
	"Asia/Qostanay",
	"Asia/Qyzylorda",
	"Asia/Rangoon",
	"Asia/Riyadh",
	"Asia/Saigon",
	"Asia/Sakhalin",
	"Asia/Samarkand",
	"Asia/Seoul",
	"Asia/Shanghai",
	"Asia/Singapore",
	"Asia/Srednekolymsk",
	"Asia/Taipei",
	"Asia/Tashkent",
	"Asia/Tbilisi",
	"Asia/Tehran",
	"Asia/Tel_Aviv",
	"Asia/Thimbu",
	"Asia/Thimphu",
	"Asia/Tokyo",
	"Asia/Tomsk",
	"Asia/Ujung_Pandang",
	"Asia/Ulaanbaatar",
	"Asia/Ulan_Bator",
	"Asia/Urumqi",
	"Asia/Ust-Nera",
	"Asia/Vientiane",
	"Asia/Vladivostok",
	"Asia/Yakutsk",
	"Asia/Yangon",
	"Asia/Yekaterinburg",
	"Asia/Yerevan",
	"Atlantic/Azores",
	"Atlantic/Bermuda",
	"Atlantic/Canary",
	"Atlantic/Cape_Verde",
	"Atlantic/Faeroe",
	"Atlantic/Faroe",
	"Atlantic/Jan_Mayen",
	"Atlantic/Madeira",
	"Atlantic/Reykjavik",
	"Atlantic/South_Georgia",
	"Atlantic/St_Helena",
	"Atlantic/Stanley",
	"Australia/ACT",
	"Australia/Adelaide",
	"Australia/Brisbane",
	"Australia/Broken_Hill",
	"Australia/Canberra",
	"Australia/Currie",
	"Australia/Darwin",
	"Australia/Eucla",
	"Australia/Hobart",
	"Australia/LHI",
	"Australia/Lindeman",
	"Australia/Lord_Howe",
	"Australia/Melbourne",
	"Australia/NSW",
	"Australia/North",
	"Australia/Perth",
	"Australia/Queensland",
	"Australia/South",
	"Australia/Sydney",
	"Australia/Tasmania",
	"Australia/Victoria",
	"Australia/West",
	"Australia/Yancowinna",
	// "Brazil/Acre",
	// "Brazil/DeNoronha",
	// "Brazil/East",
	// "Brazil/West",
	// "Canada/Atlantic",
	// "Canada/Central",
	// "Canada/Eastern",
	// "Canada/Mountain",
	// "Canada/Newfoundland",
	// "Canada/Pacific",
	// "Canada/Saskatchewan",
	// "Canada/Yukon",
	// "Chile/Continental",
	// "Chile/EasterIsland",
	"Europe/Amsterdam",
	"Europe/Andorra",
	"Europe/Astrakhan",
	"Europe/Athens",
	"Europe/Belfast",
	"Europe/Belgrade",
	"Europe/Berlin",
	"Europe/Bratislava",
	"Europe/Brussels",
	"Europe/Bucharest",
	"Europe/Budapest",
	"Europe/Busingen",
	"Europe/Chisinau",
	"Europe/Copenhagen",
	"Europe/Dublin",
	"Europe/Gibraltar",
	"Europe/Guernsey",
	"Europe/Helsinki",
	"Europe/Isle_of_Man",
	"Europe/Istanbul",
	"Europe/Jersey",
	"Europe/Kaliningrad",
	"Europe/Kiev",
	"Europe/Kirov",
	"Europe/Lisbon",
	"Europe/Ljubljana",
	"Europe/London",
	"Europe/Luxembourg",
	"Europe/Madrid",
	"Europe/Malta",
	"Europe/Mariehamn",
	"Europe/Minsk",
	"Europe/Monaco",
	"Europe/Moscow",
	"Europe/Nicosia",
	"Europe/Oslo",
	"Europe/Paris",
	"Europe/Podgorica",
	"Europe/Prague",
	"Europe/Riga",
	"Europe/Rome",
	"Europe/Samara",
	"Europe/San_Marino",
	"Europe/Sarajevo",
	"Europe/Saratov",
	"Europe/Simferopol",
	"Europe/Skopje",
	"Europe/Sofia",
	"Europe/Stockholm",
	"Europe/Tallinn",
	"Europe/Tirane",
	"Europe/Tiraspol",
	"Europe/Ulyanovsk",
	"Europe/Uzhgorod",
	"Europe/Vaduz",
	"Europe/Vatican",
	"Europe/Vienna",
	"Europe/Vilnius",
	"Europe/Volgograd",
	"Europe/Warsaw",
	"Europe/Zagreb",
	"Europe/Zaporozhye",
	"Europe/Zurich",
	"Indian/Antananarivo",
	"Indian/Chagos",
	"Indian/Christmas",
	"Indian/Cocos",
	"Indian/Comoro",
	"Indian/Kerguelen",
	"Indian/Mahe",
	"Indian/Maldives",
	"Indian/Mauritius",
	"Indian/Mayotte",
	"Indian/Reunion",
	// "Mexico/BajaNorte",
	// "Mexico/BajaSur",
	// "Mexico/General",
	"Pacific/Apia",
	"Pacific/Auckland",
	"Pacific/Bougainville",
	"Pacific/Chatham",
	"Pacific/Chuuk",
	"Pacific/Easter",
	"Pacific/Efate",
	"Pacific/Enderbury",
	"Pacific/Fakaofo",
	"Pacific/Fiji",
	"Pacific/Funafuti",
	"Pacific/Galapagos",
	"Pacific/Gambier",
	"Pacific/Guadalcanal",
	"Pacific/Guam",
	"Pacific/Honolulu",
	"Pacific/Johnston",
	"Pacific/Kiritimati",
	"Pacific/Kosrae",
	"Pacific/Kwajalein",
	"Pacific/Majuro",
	"Pacific/Marquesas",
	"Pacific/Midway",
	"Pacific/Nauru",
	"Pacific/Niue",
	"Pacific/Norfolk",
	"Pacific/Noumea",
	"Pacific/Pago_Pago",
	"Pacific/Palau",
	"Pacific/Pitcairn",
	"Pacific/Pohnpei",
	"Pacific/Ponape",
	"Pacific/Port_Moresby",
	"Pacific/Rarotonga",
	"Pacific/Saipan",
	"Pacific/Samoa",
	"Pacific/Tahiti",
	"Pacific/Tarawa",
	"Pacific/Tongatapu",
	"Pacific/Truk",
	"Pacific/Wake",
	"Pacific/Wallis",
	"Pacific/Yap",
	// "US/Alaska",
	// "US/Aleutian",
	// "US/Arizona",
	// "US/Central",
	// "US/East-Indiana",
	// "US/Eastern",
	// "US/Hawaii",
	// "US/Indiana-Starke",
	// "US/Michigan",
	// "US/Mountain",
	// "US/Pacific",
	// "US/Samoa",
	"UTC",
}
