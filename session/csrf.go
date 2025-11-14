package session

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetCsrfToken(c echo.Context) string {
	csrfToken, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if ok {
		return csrfToken
	}
	return ""
}
