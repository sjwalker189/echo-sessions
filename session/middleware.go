package session

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
)

func sessionContextKey(cookieName string) string {
	return fmt.Sprintf("session:%s", cookieName)
}

// Retreive session data for the current request with a type assertion
func UseSessionByName[T any](c echo.Context, cookieId string) *Session[T] {
	store, ok := c.Get(sessionContextKey(cookieId)).(*Session[T])
	if !ok {
		panic(fmt.Sprintf("Session \"%s\" is not present on the request. Did you forget to attach the session middleware?", cookieId))
	}
	return store
}

type SessionConfig[T any] struct {
	CookieName        string
	SaveUninitialized bool
	Store             Store[T]
	AfterResponse     func(c echo.Context, session Session[T], store Store[T])
}

func WithSessions[T any](cookieName string, store Store[T]) echo.MiddlewareFunc {
	return createSessionMiddleware(SessionConfig[T]{
		CookieName: cookieName,
		Store:      store,
	})
}

func SessionCookie[T any](config SessionConfig[T]) echo.MiddlewareFunc {
	return createSessionMiddleware(config)
}

func createSessionMiddleware[T any](config SessionConfig[T]) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isCookieNew := false

			// Read the cookie from the request or initialize a new one if not present
			cookie, err := c.Cookie(config.CookieName)
			if err != nil || cookie == nil {
				c := NewSessionCookie(config.CookieName)
				cookie = c
				isCookieNew = true
			}

			// Read the session data from the store if present
			sessId := cookie.Value
			sess, err := config.Store.Get(sessId)
			if err != nil {
				log.Println("No data present for cookie")
			}

			// Create a copy of the session to determine if it needs to be saved
			// when sending a response
			snapshot := *sess.Data

			// Expose the session on the request context so that it may
			// be accessed by handlers.
			c.Set(sessionContextKey(config.CookieName), &sess)

			// Before sending the response set the cookie header and persist
			// any changes made to the session
			c.Response().Before(func() {
				// Don't set a cookie when saveUninitialized is false and no data was
				// recorded against the session
				changed := reflect.DeepEqual(snapshot, *sess.Data)
				if config.SaveUninitialized == false && isCookieNew && !changed {
					return
				}

				// Update the cookie expires time to keep it alive
				cookie.Value = sess.ID
				cookie.Expires = time.Now().Add(24 * time.Hour)
				config.Store.Set(sess.ID, sess)
				c.SetCookie(cookie)
			})

			c.Response().After(func() {
				// Remove old session key if the ID changed
				if sess.ID != sessId {
					config.Store.Del(sessId)
				}

				if config.AfterResponse != nil {
					config.AfterResponse(c, sess, config.Store)
				}
			})

			return next(c)
		}
	}

}

