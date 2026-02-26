package server

import (
	"app/internal/session"

	"github.com/labstack/echo/v4"
)

func Authenticated(redirectTo string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.Default(c)
			if sess.User == nil {
				return c.Redirect(303, redirectTo)
			}

			// Do not cache authenticated route responses
			c.Response().Header().Add("Cache-Control", "no-cache")
			c.Response().Header().Add("Pragma", "no-cache")
			c.Response().Header().Add("Expires", "-1")

			return next(c)
		}
	}
}

func Guest(redirectTo string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.Default(c)
			if sess.User != nil {
				return c.Redirect(303, redirectTo)
			}
			return next(c)
		}
	}
}

func NoCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Cache-Control", "no-cache")
			c.Response().Header().Add("Pragma", "no-cache")
			c.Response().Header().Add("Expires", "-1")
			return next(c)
		}
	}
}
