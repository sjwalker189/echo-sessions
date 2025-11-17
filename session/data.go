package session

import (
	"app/types"

	"github.com/labstack/echo/v4"
)

type AuthUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FlashMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type FormErrors map[string]string

type Data struct {
	User       *AuthUser                     `json:"user"`
	Flashes    []FlashMessage                `json:"flashes"`
	FormErrors types.HashMap[string, string] `json:"formErrors"`
}

func NewData() Data {
	return Data{
		Flashes:    make([]FlashMessage, 0),
		FormErrors: make(types.HashMap[string, string]),
	}
}

func (s *Data) Flash(msg FlashMessage) {
	s.Flashes = append(s.Flashes, msg)
}

func Default(c echo.Context) *Session[Data] {
	return UseSessionByName[Data](c, DefaultCookieName)
}

func (s *Data) Authenticated() bool {
	return s.User != nil
}
