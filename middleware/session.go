package middleware

import (
	"app/session"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

const sessionContextKey = "session"

// Retreive session data for the current request with a type assertion
func UseSession[T any](c echo.Context, cookieId string) *session.Session[T] {
	store, ok := c.Get(fmt.Sprintf("session:%s", cookieId)).(*session.Session[T])
	if !ok {
		panic(fmt.Sprintf("Session \"%s\" is not present on the request. Did you forget to attach the session middleware?", cookieId))
	}
	return store
}

func WithSessions[T any](cookieId string, store session.Store[T]) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Read the cookie from the request or initialize a new one if not present
			cookie, err := c.Cookie(cookieId)
			if err != nil || cookie == nil {
				c := session.NewSessionCookie(cookieId)
				cookie = c
			}

			// Read the session data from the store if present
			sessId := cookie.Value
			sess, err := store.Get(sessId)
			if err != nil {
				log.Println("No data present for cookie")
			}

			// Expose the session on the request context so that it may
			// be accessed by handlers.
			c.Set(fmt.Sprintf("session:%s", cookieId), &sess)

			// Before sending the response set the cookie header and persist
			// any changes made to the session
			c.Response().Before(func() {
				cookie.Value = sess.ID
				cookie.Expires = time.Now().Add(24 * time.Hour)

				// TODO: if SaveUninitialized is false and cookie is not modified and cookie is new then
				// do not save or set the cookie
				store.Set(sess.ID, sess)

				// If the session id changed, delete the old entry
				if sess.ID != sessId {
					store.Del(sessId)
				}

				c.SetCookie(cookie)
			})

			return next(c)
		}
	}
}
