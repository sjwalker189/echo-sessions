package main

import (
	"app/session"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//
// type ErrorBag struct {
// 	Value map[string]string `json:"value"`
// }
//
// func (e *ErrorBag) init() {
// 	if e.Value == nil {
// 		e.Value = make(map[string]string)
// 	}
// }
//
// func (e *ErrorBag) Has(key string) bool {
// 	e.init()
// 	_, ok := e.Value[key]
// 	return ok
// }
//
// func (e *ErrorBag) Get(key string) string {
// 	e.init()
// 	value, ok := e.Value[key]
// 	if !ok {
// 		return ""
// 	}
// 	return value
// }

// func (e *ErrorBag) Set(key string, value string) {
// 	e.init()
// 	e.Value[key] = value
// }
//
// func (e *ErrorBag) Clear() {
// 	e.Value = make(map[string]string)
// }

func Authenticated(redirectTo string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := session.DefaultSession(c)
			if sess.Data.User == nil {
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
			sess := session.DefaultSession(c)
			if sess.Data.User != nil {
				return c.Redirect(303, redirectTo)
			}
			return next(c)
		}
	}
}

type LoginFormData struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.ContextTimeout(5 * time.Second))
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_csrf",
	}))

	e.Use(session.SessionCookie(session.SessionConfig[session.SessionData]{
		CookieName: "session",
		Store:      session.NewMemorySessionStore[session.SessionData](),
		AfterResponse: func(c echo.Context, s session.Session[session.SessionData], store session.Store[session.SessionData]) {
			// Clear any values from the session that should only live for one request
			// Update the session store if nesssary
			if len(s.Data.Flashes) > 0 || len(s.Data.Errors) > 0 {
				s.Data.Flashes = nil
				s.Data.Errors = nil
				store.Set(s.ID, s)
			}
		},
	}))

	// Auth sessions

	e.GET("/", func(c echo.Context) error {
		sess := session.DefaultSession(c)
		sess.Data.Value += 1

		return c.JSON(http.StatusOK, map[string]any{
			"msg":   "Hello, World!",
			"count": sess.Data.Value,
		})
	})

	e.GET("/login", func(c echo.Context) error {
		csrf := session.GetCsrfToken(c)
		sess := session.DefaultSession(c)

		// Accumulate all flash messages
		var sb strings.Builder
		for _, m := range sess.Data.Flashes {
			sb.WriteString(fmt.Sprintf("%s<br>", m.Message))
		}

		return c.HTML(http.StatusOK, fmt.Sprintf(`
			%s
			<form action="/login" method="POST">
			<input name="_csrf" value="%s" type="hidden" />
			<input name="email"/><br>
			<input name="password"/><br>
			<button type="submit">Send</button>
			</form>
		`, sb.String(), csrf))
	}, Guest("/app"))

	e.POST("/login", func(c echo.Context) error {
		sess := session.DefaultSession(c)

		formData := new(LoginFormData)
		if err := c.Bind(formData); err != nil {
			return err
		}

		if formData.Email == "admin" && formData.Password == "admin" {
			sess.Data.User = &session.AuthUser{
				Name:  "Admin",
				Email: formData.Email,
			}
			sess.RegenerateID()
			return c.Redirect(303, "/app")
		}

		sess.Data.Flash(session.FlashMessage{
			Message: "Invalid credentials",
			Type:    "Notice",
		})

		return c.Redirect(303, "/login")
	}, Guest("/app"))

	e.POST("/logout", func(c echo.Context) error {
		sess := session.DefaultSession(c)
		sess.Data.User = nil
		sess.RegenerateID()
		return c.Redirect(303, "/login")
	})

	e.GET("/app", func(c echo.Context) error {
		csrf := session.GetCsrfToken(c)
		return c.HTML(http.StatusOK, fmt.Sprintf(`
			<h1>Welcome</h1>
			<form action="/logout" method="POST">
				<input name="_csrf" value="%s" type="hidden" />
				<button>Logout</button>
			</form>
		`, csrf))
	}, Authenticated("/login"))

	e.Logger.Fatal(e.Start(":8888"))
}
