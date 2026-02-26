package forms

import "app/internal/types"

type Login struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6,max=120"`
}

func (form Login) Validate() *types.HashMap[string, string] {
	return nil
}
