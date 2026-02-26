package server

import (
	"app/internal/handlers"
	"app/internal/view"
	"app/internal/view/pages"

	"github.com/labstack/echo/v4"
)

func routes(e *echo.Echo) {
	e.GET("/", handlers.Home)
	e.GET("/login", func (c echo.Context) error {
		return view.Model(pages.Login{}).Send(c)
	})
	// Guest routes
	// guest := e.Group("", Guest("/app"))
	// guest.GET("/login", handlers.Login())

	// Admin routes
	// admin := e.Group("", Authenticated("/login"))
	// admin.POST("/logout", handlers.LogoutSubmit)
	// admin.GET("/app", func(c echo.Context) error {
	// 	csrf := session.CsrfToken(c)
	// 	sess := session.Default(c)
	// 	return c.HTML(http.StatusOK, fmt.Sprintf(`
	// 		<h1>Welcome, %s</h1>
	// 		<form action="/logout" method="POST">
	// 			<input name="_csrf" value="%s" type="hidden" />
	// 			<button>Logout</button>
	// 		</form>
	// 	`, sess.User.Name, csrf))
	// })
}
