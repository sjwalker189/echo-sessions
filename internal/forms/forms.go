package forms

import (
	"app/internal/session"
	"app/internal/types"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validatable interface {
	Validate() *types.HashMap[string, string]
}

type Form[T Validatable] struct {
	context echo.Context
	session *session.Session
	Errors  *types.HashMap[string, string]
	Fields  T
}

func New[T Validatable](c echo.Context) (Form[T], error) {
	sess := session.Default(c)
	form := Form[T]{
		context: c,
		session: sess,
		Errors:  sess.FormErrors,
	}

	params, err := c.FormParams()
	if err != nil {
		return form, err
	}

	err = c.Bind(&form.Fields)
	sess.FormValues = &params

	return form, err
}

func (form *Form[T]) Csrf() string {
	return session.CsrfToken(form.context)
}

func (form *Form[T]) SetFields(fields T) {
	form.Fields = fields
}

func (form *Form[T]) Validate() {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(form.Fields)
	errors := err.(validator.ValidationErrors)

	// TODO: loop and format errors and set on form.Errors hashmap
	if len(errors) > 0 {
		fmt.Println(errors)
	}

	fieldErrors := form.Fields.Validate()
	if fieldErrors != nil {
		for k, v := range *fieldErrors {
			form.Errors.Set(k, v)
		}
	}
}
