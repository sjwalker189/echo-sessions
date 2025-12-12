package session

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
)

func sessionContextKey(cookieName string) string {
	return fmt.Sprintf("session:%s", cookieName)
}

type Config struct {
	CookieName        string
	SaveUninitialized bool
	Store             Store
	AfterResponse     func(c echo.Context, session Session, store Store)
}

func WithSessions(cookieName string, store Store) echo.MiddlewareFunc {
	return createSessionMiddleware(Config{
		CookieName: cookieName,
		Store:      store,
	})
}

func SessionCookie(config Config) echo.MiddlewareFunc {
	return createSessionMiddleware(config)
}

func createSessionMiddleware(config Config) echo.MiddlewareFunc {
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
			snapshot := sess
			ref := &sess

			// Expose the session on the request context so that it may
			// be accessed by handlers.
			c.Set(sessionContextKey(config.CookieName), ref)

			// Restore previous form submission values if present so that these can be rebound
			req := c.Request()
			if req.Method == http.MethodGet && sess.FormValues != nil {
				req := c.Request()
				req.ContentLength = -1                                           // unknown length. ContentLength cannot be zero for c.Bind() to work
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm) // TODO: Should store the content type when setting sess.FormValues
				req.PostForm = *sess.FormValues
				req.Form = *sess.FormValues
			}

			// Before sending the response set the cookie header and persist
			// any changes made to the session
			c.Response().Before(func() {
				// Don't set a cookie when saveUninitialized is false and no data was
				// recorded against the session
				changed := reflect.DeepEqual(snapshot, *ref)
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

				// TODO: I don't super love this because the session is stored twice
				// for every request. There is also a (very unlikely) race condition
				// where a subsequent request reuses the session data before it's
				// cleared from the store.
				if len(sess.Flashes) > 0 || !sess.FormErrors.Empty() || sess.FormValues != nil {
					fmt.Println("Clearing flash/errors from session")
					sess.Flush()
					config.Store.Set(sess.ID, sess)
				}
			})

			return next(c)
		}
	}

}
