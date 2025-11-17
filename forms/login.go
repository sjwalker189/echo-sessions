package forms

import (
	"app/ent"
	"app/store"
	"app/util"
	"context"
)

type Login struct {
	store  *store.UserStore
	Errors util.ErrorBag

	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=6,max=120"`
}

func NewLogin(store *store.UserStore) Login {
	errors := util.NewErrorBag()
	return Login{
		store:  store,
		Errors: errors,
	}
}

func (self Login) Validate() {
	self.Errors.ValidateStruct(self)
}

func (self Login) Submit(c context.Context) (*ent.User, error) {
	user, err := self.store.VerifyCredentials(c, self.Email, self.Password)
	return user, err
}
