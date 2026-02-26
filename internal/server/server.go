package server

import (
	"app/internal/session"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.ContextTimeout(5 * time.Second))
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_csrf",
	}))

	e.Use(session.SessionCookie(session.Config{
		CookieName: "session",
		Store:      session.NewMemorySessionStore(),
	}))

	e.Use(NoCache()) // Dev only?
	e.Static("/static", "public")

	routes(e)

	return e
}
