package handlers

import (
	"app/internal/view"
	"app/internal/view/pages"

	"github.com/labstack/echo/v4"
)

type State struct {}

func Home(c echo.Context) error {
	return view.Model(pages.Home{
		DocTitle: "Accelerator",
	}).Send(c)
}
