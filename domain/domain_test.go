package domain_test

import (
	"testing"

	"github.com/rlibaert/service-example-go/domain"
	"github.com/rlibaert/service-example-go/domaintest"
	"github.com/rlibaert/service-example-go/stores"
)

func TestService(t *testing.T) {
	domaintest.TestService(t, &domain.ServiceStore{Store: stores.MustNewMock()})
}
