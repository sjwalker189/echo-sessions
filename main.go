package main

import (
	"app/middleware"
	"app/session"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MiddlewareHandler func(http.Handler) http.Handler

type AuthUser struct {
	Name  string
	Email string
}

type CounterData struct {
	Value    int
	User     *AuthUser
	Error    string
	Messages ErrorBag
}

type ErrorBag struct {
	Value map[string]string `json:"value"`
}

func (e *ErrorBag) init() {
	if e.Value == nil {
		e.Value = make(map[string]string)
	}
}

func (e *ErrorBag) Has(key string) bool {
	e.init()
	_, ok := e.Value[key]
	return ok
}

func (e *ErrorBag) Get(key string) string {
	e.init()
	value, ok := e.Value[key]
	if !ok {
		return ""
	}
	return value
}

func (e *ErrorBag) Set(key string, value string) {
	e.init()
	e.Value[key] = value
}

func (e *ErrorBag) Clear() {
	e.Value = make(map[string]string)
}

const cookieName = "counter"

func useCounterSession(c echo.Context) *session.Session[CounterData] {
	return middleware.UseSession[CounterData](c, cookieName)
}

func Authenticated(redirectTo string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := useCounterSession(c)
			if sess.Data.User == nil {
				return c.Redirect(303, redirectTo)
			}
			return next(c)
		}
	}
}

func Guest(redirectTo string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := useCounterSession(c)
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

	counterStore := session.NewMemorySessionStore[CounterData]()

	e.Use(middleware.WithSessions(cookieName, &counterStore))

	e.GET("/", func(c echo.Context) error {
		sess := useCounterSession(c)
		sess.Data.Value += 1

		return c.JSON(http.StatusOK, map[string]any{
			"msg":   "Hello, World!",
			"count": sess.Data.Value,
		})
	})

	e.GET("/login", func(c echo.Context) error {
		sess := useCounterSession(c)
		error := sess.Data.Error

		fmt.Println(sess.Data.Messages)

		// TODO: Sessions are saved implicitly when the request is sent
		// It may be prefereable to make this explicit like the following
		// sess.Data.Field = newValue
		// sess.Save()
		sess.Data.Error = ""
		sess.Data.Messages.Clear()

		// TODO: It would be ideal to construct an API like the following:
		// errors := useFormErrors(c)
		// flash := useFlashMessages(c)
		//
		// It needs to be the case that the values are removed from the session
		// once the request ends.

		return c.HTML(http.StatusOK, fmt.Sprintf(`
			%s
			<form action="/login" method="POST">
			<input name="email"/><br>
			<input name="password"/><br>
			<button type="submit">Send</button>
			</form>
			`, error))
	}, Guest("/app"))

	e.POST("/login", func(c echo.Context) error {
		sess := useCounterSession(c)

		formData := new(LoginFormData)
		if err := c.Bind(formData); err != nil {
			return err
		}

		if formData.Email == "admin" && formData.Password == "admin" {
			sess.Data.User = &AuthUser{
				Name:  "Admin",
				Email: formData.Email,
			}
			sess.RegenerateID()
			return c.Redirect(303, "/app")
		}

		sess.Data.Messages.Set("form", "Invalid credentials")
		sess.Data.Messages.Set("email", "Email is not a valid email address")
		sess.Data.Messages.Set("password", "Password is not valid")
		sess.Data.Error = "Invalid Credentials"

		return c.Redirect(303, "/login")
	}, Guest("/app"))

	e.POST("/logout", func(c echo.Context) error {
		sess := useCounterSession(c)
		sess.Data.User = nil
		sess.RegenerateID()
		return c.Redirect(303, "/login")
	})

	e.GET("/app", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<h1>Welcome</h1>
			<form action="/logout" method="POST">
				<button>Logout</button>
			</form>
		`)
	}, Authenticated("/login"))

	e.Logger.Fatal(e.Start(":8888"))
}
