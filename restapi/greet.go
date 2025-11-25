package restapi

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type GreetRegisterer struct{}

func (GreetRegisterer) RegisterAPI(api huma.API) {
	type input struct {
		Who string `path:"who" maxLength:"30" example:"world" doc:"Who to greet"`
	}
	type output struct {
		Body string
	}

	handler := func(_ context.Context, i *input) (*output, error) {
		return &output{Body: fmt.Sprintf("Hello, %s!", i.Who)}, nil
	}

	huma.Get(api, "/greet/{who}", handler)
}
