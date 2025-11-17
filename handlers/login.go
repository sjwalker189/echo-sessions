package handlers

import (
	"app/forms"
	"app/session"
	"app/views"
	"app/views/pages"
	"errors"

	"github.com/labstack/echo/v4"
)

func LoginGet(store *store.UserStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess := session.Default(c)
		form := forms.NewLogin(store)
		form.Errors = sess.Errors()
		form.Email = "admin@rjs.co.nz"
		form.Password = "admin123"
		return views.Render(c, 200, pages.Login(form))
	}
}

func LoginPost(store *store.UserStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		form := forms.NewLogin(store)
		if err := c.Bind(&form); err != nil {
			return errors.New("decode error")
		}

		form.Validate()
		if form.Errors.Any() {
			sess.AddErrors(form.Errors)
			return c.Redirect(302, RouteLogin)
		}

		ctx := c.Request().Context()
		user, err := form.Submit(ctx)
		if err != nil {
			form.Errors.Set("email", "Server error. Please try again soon.")
			sess.AddErrors(form.Errors)
			return c.Redirect(302, RouteLogin)
		}

		sess.Authenticate(user.ID)

		return c.Redirect(302, RouteAdmin)
	}
}

func LogoutPost() echo.HandlerFunc {
	return func(c echo.Context) error {
		session.Default(c).Logout()
		return c.Redirect(302, RouteLogin)
	}
}
