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
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/microbus-io/bespa/widget"
	"github.com/microbus-io/errors"
)

// MaxDirectoryRecords is the maximum number of persons allowed in the directory.
const MaxDirectoryRecords = 20

// Per-field maximum lengths enforced at the storage layer as defense-in-depth.
// The form validators already enforce these on input; these hard caps protect
// the server if any future code path bypasses the form (e.g. a direct API).
// Combined with MaxDirectoryRecords (20) and the 100,000 session cap, the
// directory data is hard-bounded at ~400 MB worst-case across all sessions.
const (
	maxFirstNameLen = 20
	maxLastNameLen  = 20
	maxEmailLen     = 40
	maxPhoneLen     = 20
	maxAddressLen   = 40
	maxCityLen      = 20
	maxStateLen     = 3
	maxZipLen       = 5
)

// Person keeps personal information about a person.
// All directories are scoped to a single session (see shared/session.go); persons
// stored in one session's directory are never visible to another session.
type Person struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Birthday  time.Time
	Address1  string
	Address2  string
	City      string
	State     string
	Zip       string
}

// validate rejects a person whose fields exceed the per-field maximum sizes.
// Returns the first violation found, or nil if every field is within bounds.
func (p *Person) validate() error {
	switch {
	case len(p.FirstName) > maxFirstNameLen:
		return errors.New("first name too long")
	case len(p.LastName) > maxLastNameLen:
		return errors.New("last name too long")
	case len(p.Email) > maxEmailLen:
		return errors.New("email too long")
	case len(p.Phone) > maxPhoneLen:
		return errors.New("phone too long")
	case len(p.Address1) > maxAddressLen:
		return errors.New("address line 1 too long")
	case len(p.Address2) > maxAddressLen:
		return errors.New("address line 2 too long")
	case len(p.City) > maxCityLen:
		return errors.New("city too long")
	case len(p.State) > maxStateLen:
		return errors.New("state too long")
	case len(p.Zip) > maxZipLen:
		return errors.New("zip too long")
	}
	return nil
}

// FullName combines the person's first and last names.
func (p *Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// PersonDirectory keeps a list of persons.
type PersonDirectory struct {
	data []*Person
	sync.Mutex
}

// List returns all persons in the directory.
// The data returned is a fresh copy so it is thread-safe.
func (ds *PersonDirectory) List() []*Person {
	ds.Lock()
	defer ds.Unlock()
	// Add 3 records to get started
	if ds.data == nil {
		uniqueEmails := map[string]bool{}
		for len(uniqueEmails) < 3 {
			p := ds.RandomPerson()
			if !uniqueEmails[p.Email] {
				p.ID = widget.RandomAlphaNumID(8)
				ds.data = append(ds.data, p)
				uniqueEmails[p.Email] = true
			}
		}
	}
	result := []*Person{}
	for _, i := range ds.data {
		clone := *i
		result = append(result, &clone)
	}
	return result
}

// Update updates the data of a person in the directory.
func (ds *PersonDirectory) Update(person *Person) (err error) {
	if err := person.validate(); err != nil {
		return errors.Trace(err)
	}
	ds.Lock()
	defer ds.Unlock()
	for _, i := range ds.data {
		if i.Email == person.Email && i.ID != person.ID {
			return errors.New("email already registered")
		}
	}
	for k, i := range ds.data {
		if i.ID == person.ID {
			clone := *person
			ds.data[k] = &clone
			return nil
		}
	}
	return errors.New("record not found")
}

// Insert add a person to the directory.
// If a person with the same email address already exists, a new record is not created.
func (ds *PersonDirectory) Insert(person *Person) (id string, err error) {
	if err := person.validate(); err != nil {
		return "", errors.Trace(err)
	}
	ds.Lock()
	defer ds.Unlock()
	for _, i := range ds.data {
		if i.Email == person.Email {
			return i.ID, errors.New("email already registered")
		}
	}
	if len(ds.data) < MaxDirectoryRecords {
		clone := *person
		clone.ID = widget.RandomAlphaNumID(8)
		ds.data = append(ds.data, &clone)
		return clone.ID, nil
	}
	return "", errors.New("directory size limit reached")
}

// CanInsert indicates if there is room in the directory to add more persons.
func (ds *PersonDirectory) CanInsert(r *http.Request) bool {
	ds.Lock()
	defer ds.Unlock()
	return len(ds.data) < MaxDirectoryRecords
}

// Lookup returns a person by ID.
// The data returned is a fresh copy so it is thread-safe.
func (ds *PersonDirectory) Lookup(personID string) (person *Person, ok bool) {
	ds.Lock()
	defer ds.Unlock()
	for _, i := range ds.data {
		if i.ID == personID {
			clone := *i
			return &clone, true
		}
	}
	return nil, false
}

// Lookup returns a person by email.
// The data returned is a fresh copy so it is thread-safe.
func (ds *PersonDirectory) LookupByEmail(email string) (person *Person, ok bool) {
	ds.Lock()
	defer ds.Unlock()
	for _, i := range ds.data {
		if i.Email == email {
			clone := *i
			return &clone, true
		}
	}
	return nil, false
}

// Delete removes the data of a person from the directory.
func (ds *PersonDirectory) Delete(personID string) (ok bool) {
	ds.Lock()
	defer ds.Unlock()
	for k, i := range ds.data {
		if i.ID == personID {
			ds.data = append(ds.data[:k], ds.data[k+1:]...)
			return true
		}
	}
	return false
}

func (ds *PersonDirectory) RandomPerson() *Person {
	firstNames := []string{
		"James", "John", "Robert", "Michael", "William", "David",
		"Richard", "Charles", "Joseph", "Thomas", "Christopher",
		"Daniel", "Paul", "Mark", "Donald", "George", "Kenneth",
		"Steven", "Edward", "Brian",
		"Mary", "Patricia", "Linda", "Barbara", "Elizabeth", "Jennifer",
		"Maria", "Susan", "Margaret", "Dorothy", "Lisa", "Nancy",
		"Karen", "Betty", "Helen", "Sandra", "Donna", "Carol", "Ruth",
		"Sharon",
	}
	lastNames := []string{
		"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis",
		"Garcia", "Rodriguez", "Wilson", "Martinez", "Anderson", "Taylor",
		"Thomas", "Hernandez", "Moore", "Martin", "Jackson", "Thompson",
		"White", "Lopez", "Lee", "Gonzalez", "Harris", "Clark", "Lewis",
		"Robinson", "Walker", "Perez", "Hall",
	}
	cities := []string{
		"New York", "NY",
		"Los Angeles", "CA",
		"Chicago", "IL",
		"Houston", "TX",
		"Phoenix", "AZ",
		"Philadelphia", "PA",
		"San Antonio", "TX",
		"San Diego", "CA",
		"Dallas", "TX",
		"San Jose", "CA",
		"Austin", "TX",
		"Jacksonville", "FL",
		"Forth Worth", "TX",
		"Columbus", "OH",
		"Charlotte", "NC",
		"Indianapolis", "IN",
		"San Francisco", "CA",
		"Seattle", "WA",
		"Denver", "CO",
		"Nashville", "TN",
	}

	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	email := strings.ToLower(firstName + "." + lastName + "@example.com")
	phone := ""
	for i := 0; i < 12; i++ {
		if i == 3 || i == 7 {
			phone += " "
		} else {
			phone += strconv.Itoa(rand.Intn(10))
		}
	}

	age := rand.Intn(60) + 19
	yyyy := time.Now().Year() - age
	mm := rand.Intn(12) + 1
	dd := rand.Intn(28) + 1
	birthday := time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)

	x := 2 * rand.Intn(len(cities)/2)
	city := cities[x]
	state := cities[x+1]
	zip := ""
	for i := 0; i < 5; i++ {
		zip += strconv.Itoa(rand.Intn(10))
	}
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Birthday:  birthday,
		City:      city,
		State:     state,
		Zip:       zip,
	}
}
