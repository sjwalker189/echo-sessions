package handlers

import (
	"app/web/views"
	"app/web/views/pages"

	"github.com/labstack/echo/v4"
)

func handleAdminGet(c echo.Context) error {
	return c.JSON(200, "control panel")
	return views.Render(c, 200, pages.Dashboard())
}
