package domaintest

import (
	"testing"

	"github.com/rlibaert/service-example-go/domain"
)

func TestService(t *testing.T, service domain.Service) {
	id, _ := service.ContactsCreate(t.Context(), &domain.Contact{})
	_, _ = service.ContactsRead(t.Context(), id)
	_ = service.ContactsUpdate(t.Context(), id, &domain.Contact{})
	_ = service.ContactsDelete(t.Context(), id)
}
