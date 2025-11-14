package session

import (
	"github.com/labstack/echo/v4"
)

type AuthUser struct {
	ID    string
	Name  string
	Email string
}

type FlashMessage struct {
	Type    string
	Message string
}

type FormErrors map[string]string

type SessionData struct {
	Value   int
	User    *AuthUser
	Flashes []FlashMessage
	Errors  FormErrors
}

func NewSessionData() SessionData {
	return SessionData{
		Flashes: make([]FlashMessage, 0),
		Errors:  make(FormErrors),
	}
}

func (s *SessionData) Flash(msg FlashMessage) {
	s.Flashes = append(s.Flashes, msg)
}

func (s *SessionData) Authenticated() bool {
	return s.User != nil
}

func DefaultSession(c echo.Context) *Session[SessionData] {
	return UseSessionByName[SessionData](c, DefaultCookieName)
}
