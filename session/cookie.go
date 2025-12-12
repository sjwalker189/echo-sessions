package session

import (
	"net/http"
	"time"
)

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
