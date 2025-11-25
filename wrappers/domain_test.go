package wrappers_test

import (
	"context"
	"testing"

	"github.com/rlibaert/service-example-go/domain"
	"github.com/rlibaert/service-example-go/domaintest"
	"github.com/rlibaert/service-example-go/stores"
	"github.com/rlibaert/service-example-go/wrappers"
)

func TestService(t *testing.T) {
	domaintest.TestService(t, wrappers.ServiceErrorHandler{
		Service:      domain.Service(&domain.ServiceStore{Store: stores.MustNewMock()}),
		ErrorHandler: func(context.Context, error) {},
	})
}
