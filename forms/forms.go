package forms

import (
	"app/util"
	"context"
)

type Form interface {
	Submit(c context.Context) (any, error)
	Errors() util.ErrorBag
	Validate()
}
