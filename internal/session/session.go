package session

import (
	"app/internal/types"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

const sessionCookieName = "session"

type Session struct {
	ID         string                         `json:"id"`
	Cookie     *http.Cookie                   `json:"cookie"`
	User       *AuthUser                      `json:"user"`
	Flashes    []FlashMessage                 `json:"flashes"`
	FormErrors *types.HashMap[string, string] `json:"formErrors"`
	FormValues *url.Values                    `json:"formValues"`
}

func New() Session {
	cookie := NewSessionCookie(sessionCookieName)
	errors := make(types.HashMap[string, string])
	return Session{
		ID:         cookie.Value,
		Cookie:     cookie,
		FormErrors: &errors,
	}
}

func (s *Session) Authenticated() bool {
	return s.User != nil
}

func (s *Session) Logout() {
	s.User = nil
}

func (s *Session) Flash(msg FlashMessage) {
	s.Flashes = append(s.Flashes, msg)
}

func (s *Session) RegenerateID() {
	id := NewSessionId()
	s.ID = id
	if s.Cookie != nil {
		s.Cookie.Value = id
	}
}

func (s *Session) Flush() {
	errors := make(types.HashMap[string, string])
	s.FormErrors = &errors
	s.FormValues = nil
	s.Flashes = make([]FlashMessage, 0)
}

func NewSessionId() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func Default(c echo.Context) *Session {
	store, ok := c.Get(sessionContextKey(sessionCookieName)).(*Session)
	if !ok {
		panic("Session is not present on the request. Did you forget to attach the session middleware?")
	}
	return store
}
