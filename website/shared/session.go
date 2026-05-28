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

package shared

import (
	"net/http"
	"sync"

	"github.com/microbus-io/bespa/css"
	"github.com/microbus-io/bespa/website/storage"
	"github.com/microbus-io/bespa/widget"
)

var (
	sessionLock sync.Mutex
	sessions    map[string]*Session
	fifo        []string
)

// Session keeps track of user data in-memory.
// It is an in-memory emulation of a database.
type Session struct {
	ID          string
	Theme       string
	Palette     string
	directoryDB *storage.PersonDirectory
	lock        sync.Mutex
}

// DirectoryDB returns the directory database associated with the session.
func (s *Session) DirectoryDB() *storage.PersonDirectory {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.directoryDB == nil {
		s.directoryDB = &storage.PersonDirectory{}
	}
	return s.directoryDB
}

// SessionOf returns or creates a session for the request.
func SessionOf(w http.ResponseWriter, r *http.Request) *Session {
	sessionID := widget.RandomAlphaNumID(32)
	cookie, err := r.Cookie("SessionID")
	if err == nil {
		sessionID = cookie.Value
		if len(sessionID) > 32 {
			sessionID = sessionID[:32]
		}
	} else {
		cookie = &http.Cookie{
			Name:     "SessionID",
			Value:    sessionID,
			MaxAge:   30 * 24 * 60 * 60, // One month
			HttpOnly: true,
			Path:     "/",
		}
		w.Header().Add("Set-Cookie", cookie.String())
	}

	sessionLock.Lock()
	defer sessionLock.Unlock()
	if sessions == nil {
		sessions = map[string]*Session{}
	}
	session, ok := sessions[sessionID]
	if !ok {
		// Create a new session
		session = &Session{
			ID:      sessionID,
			Palette: css.DefaultKeyColors.Name,
		}
		sessions[sessionID] = session
		// Limit to 10000 sessions
		fifo = append(fifo, sessionID)
		if len(fifo) > 10000 {
			delete(sessions, fifo[0])
			fifo = fifo[1:]
		}
	}
	return session
}
