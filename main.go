package main

import (
	"app/handlers"
	"app/session"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	model "app/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func newDb() *model.Queries {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic(err)
	}
	return model.New(db)
}

var db = newDb()

func main() {
	e := NewServer()
	e.Logger.Fatal(e.Start(":8888"))
}

func NewServer() *echo.Echo {
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

func routes(e *echo.Echo) {
	// Public routes
	e.GET("/", handlers.Welcome)

	// Guest routes
	guest := e.Group("", Guest("/app"))
	guest.GET("/login", handlers.Login(db))
	guest.POST("/login", handlers.LoginSubmit(db))

	// Admin routes
	admin := e.Group("", Authenticated("/login"))
	admin.POST("/logout", handlers.LogoutSubmit)
	admin.GET("/app", func(c echo.Context) error {
		csrf := session.CsrfToken(c)
		sess := session.Default(c)
		return c.HTML(http.StatusOK, fmt.Sprintf(`
			<h1>Welcome, %s</h1>
			<form action="/logout" method="POST">
				<input name="_csrf" value="%s" type="hidden" />
				<button>Logout</button>
			</form>
		`, sess.User.Name, csrf))
	})
}
