package restapi

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type PanicRegisterer struct{}

func (PanicRegisterer) RegisterPanic(api huma.API) {
	handler := func(context.Context, *struct{}) (*struct{}, error) {
		panic("panic argument")
	}

	huma.Get(api, "/panic", handler)
}
