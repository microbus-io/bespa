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

package storage

import (
	"bytes"
	"encoding/csv"
	"sort"
	"strconv"
	"strings"

	"github.com/microbus-io/bespa/website/resources"
)

// UnitedState is one of the 50 states of the USA.
type UnitedState struct {
	Abbrev            string
	Name              string
	Population        int
	Land              float64
	GDP               int
	PopulationDensity float64
	GDPPerCapita      float64
}

// usRecToStruct converts the CSV format to a struct.
func usRecToStruct(rec []string) *UnitedState {
	us := &UnitedState{
		Name:   rec[0],
		Abbrev: rec[1],
	}
	us.Population, _ = strconv.Atoi(rec[2])
	us.Land, _ = strconv.ParseFloat(rec[3], 64)
	us.GDP, _ = strconv.Atoi(rec[4])
	us.PopulationDensity = float64(us.Population) / us.Land
	us.GDPPerCapita = float64(us.GDP) * 1e6 / float64(us.Population)
	return us
}

/*
USLookup performs a query to the states data source to fetch one record.
This is analogous to a query to a database:

	SELECT * FROM States WHERE Abbrev='abbrev'
*/
func USLookup(abbrev string) (*UnitedState, bool, error) {
	dataFile, _ := resources.Bundle.ReadFile("states.csv")
	records, _ := csv.NewReader(bytes.NewReader(dataFile)).ReadAll()
	for _, rec := range records {
		if abbrev == rec[1] {
			return usRecToStruct(rec), true, nil
		}
	}
	return nil, false, nil // Not found
}

/*
USQuery performs a query to the states data source to fetch, filter, sort and slice
the records that match the user's request. This is analogous to a query to a database:

	SELECT * FROM States
	WHERE Name LIKE '%query%' OR Abbrev LIKE '%query%'
	ORDER BY sortOrder
	LIMIT (rowTo-rowFrom) OFFSET rowFrom

If the total number of records matching the query can be discerned (as is in this case),
the table paginator is able to display the total number of pages.
When a database is involved, this often necessitates another query.

	SELECT COUNT(*) FROM States
	WHERE Name LIKE '%query%' OR Abbrev LIKE '%query%'
*/
func USQuery(query string, sortOrder string, rowFrom int, rowTo int) (states []*UnitedState, totalRecords int, err error) {
	dataFile, _ := resources.Bundle.ReadFile("states.csv")
	records, _ := csv.NewReader(bytes.NewReader(dataFile)).ReadAll()
	records = records[1:] // Discard header line
	states = make([]*UnitedState, len(records))
	for i := range records {
		states[i] = usRecToStruct(records[i])
	}

	// Filter
	if len(query) > 0 {
		query = strings.ToLower(query)
		exact := []*UnitedState{}
		hasPrefix := []*UnitedState{}
		contains := []*UnitedState{}
		for _, state := range states {
			name := strings.ToLower(state.Name)
			abbrev := strings.ToLower(state.Abbrev)
			if name == query || abbrev == query {
				exact = append(exact, state)
			} else if strings.HasPrefix(name, query) || strings.HasPrefix(abbrev, query) {
				hasPrefix = append(hasPrefix, state)
			} else if strings.Contains(name, query) || strings.Contains(abbrev, query) {
				contains = append(contains, state)
			}
		}
		states = []*UnitedState{}
		states = append(states, exact...)
		states = append(states, hasPrefix...)
		states = append(states, contains...)
	}

	// Sort
	switch sortOrder {
	case "name":
		sort.Slice(states, func(i, j int) bool {
			return strings.ToUpper(states[i].Name) < strings.ToUpper(states[j].Name)
		})
	case "-name":
		sort.Slice(states, func(j, i int) bool {
			return strings.ToUpper(states[i].Name) < strings.ToUpper(states[j].Name)
		})
	case "abbrev":
		sort.Slice(states, func(i, j int) bool {
			return strings.ToUpper(states[i].Abbrev) < strings.ToUpper(states[j].Abbrev)
		})
	case "-abbrev":
		sort.Slice(states, func(j, i int) bool {
			return strings.ToUpper(states[i].Abbrev) < strings.ToUpper(states[j].Abbrev)
		})
	case "pop":
		sort.Slice(states, func(i, j int) bool {
			return states[i].Population < states[j].Population
		})
	case "-pop":
		sort.Slice(states, func(j, i int) bool {
			return states[i].Population < states[j].Population
		})
	case "land":
		sort.Slice(states, func(i, j int) bool {
			return states[i].Land < states[j].Land
		})
	case "-land":
		sort.Slice(states, func(j, i int) bool {
			return states[i].Land < states[j].Land
		})
	case "density":
		sort.Slice(states, func(i, j int) bool {
			return states[i].PopulationDensity < states[j].PopulationDensity
		})
	case "-density":
		sort.Slice(states, func(j, i int) bool {
			return states[i].PopulationDensity < states[j].PopulationDensity
		})
	case "gdp":
		sort.Slice(states, func(i, j int) bool {
			return states[i].GDP < states[j].GDP
		})
	case "-gdp":
		sort.Slice(states, func(j, i int) bool {
			return states[i].GDP < states[j].GDP
		})
	case "gdpcapita":
		sort.Slice(states, func(i, j int) bool {
			return states[i].GDPPerCapita < states[j].GDPPerCapita
		})
	case "-gdpcapita":
		sort.Slice(states, func(j, i int) bool {
			return states[i].GDPPerCapita < states[j].GDPPerCapita
		})
	}

	// Range
	n := len(states)
	if rowTo > n {
		rowTo = n
	}
	if rowTo < 0 {
		rowTo = n
	}
	if rowFrom > n {
		rowFrom = n
	}
	return states[rowFrom:rowTo], len(states), nil
}
