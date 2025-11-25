package stores_test

import (
	"testing"

	"github.com/rlibaert/service-example-go/domaintest"
	"github.com/rlibaert/service-example-go/stores"
)

func TestStore(t *testing.T) {
	t.Run("mock", func(t *testing.T) {
		domaintest.TestStore(t, stores.MustNewMock())
	})
}
