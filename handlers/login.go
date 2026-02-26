package handlers

import (
	"app/db"
	"app/forms"
	hash "app/lib"
	"app/session"
	"app/view"
	"app/view/pages"
	"context"
	"database/sql"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Login(model *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		form, _ := forms.New[forms.Login](c)
		return view.Render(c, 200, pages.Login(form))
	}
}

func LoginSubmit(model *db.Queries) echo.HandlerFunc {
	return func(c echo.Context) error {
		form, _ := forms.New[forms.Login](c)

		// Check validation errors
		if !form.Errors.Empty() {
			return c.Redirect(302, "/login")
		}

		// Find user to verify
		user, err := model.GetUserByEmail(c.Request().Context(), form.Fields.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				form.Errors.Set("email", "Invalid credentials")
			}
			return c.Redirect(302, "/login")
		}

		// Ensure password is valid
		if !hash.Check(form.Fields.Password, user.Password) {
			form.Errors.Set("email", "Invalid credentials")
			return c.Redirect(302, "/login")
		}

		// Authenticate the request
		sess := session.Default(c)
		sess.RegenerateID()
		sess.User = &session.AuthUser{
			ID:    strconv.FormatInt(user.ID, 10),
			Name:  user.FirstName,
			Email: user.Email,
		}

		return c.Redirect(302, "/app")
	}
}

func LogoutSubmit(c echo.Context) error {
	sess := session.Default(c)
	sess.Logout()
	sess.RegenerateID()
	return c.Redirect(302, "/login")
}
