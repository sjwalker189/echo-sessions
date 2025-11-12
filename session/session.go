package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const DefaultCookieName = "session"

type Session[T any] struct {
	ID     string       `json:"id"`
	Cookie *http.Cookie `json:"cookie"`
	Data   *T           `json:"data"`
}

func (s Session[T]) Empty() bool {
	return s.Data == nil
}

func (s *Session[T]) RegenerateID() {
	id := NewSessionId()
	s.ID = id
	if s.Cookie != nil {
		s.Cookie.Value = id
	}
}

func NewSession[T any]() Session[T] {
	cookie := NewSessionCookie(DefaultCookieName)
	return Session[T]{
		ID:     cookie.Value,
		Cookie: cookie,
		Data:   new(T),
	}
}

func NewSessionWithData[T any](data T) Session[T] {
	cookie := NewSessionCookie(DefaultCookieName)
	return Session[T]{
		ID:     cookie.Value,
		Cookie: cookie,
		Data:   &data,
	}
}

func NewSessionId() string {
	// Generate 32 random bytes
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := base64.URLEncoding.EncodeToString(b)
	return sessionID
}

func NewSessionCookie(name string) *http.Cookie {
	id := NewSessionId()

	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Name = name
	cookie.Value = id
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true

	return cookie
}
