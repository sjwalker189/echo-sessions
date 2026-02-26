package view

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Renderable interface {
	Render() templ.Component
}

func Render(ctx echo.Context, statusCode int, view templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := view.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

type Response struct {
	code    int
	content Renderable
}

// Error implements error.
func (r *Response) Error() string {
	panic("unimplemented")
}

func (r *Response) StatusCode(code int) *Response {
	r.code = code
	return r
}

func (r *Response) Send(c echo.Context) error {
	return Render(c, r.code, r.content.Render())
}

func Model(model Renderable) *Response {
	return &Response{
		code:    200,
		content: model,
	}
}
