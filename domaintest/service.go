package domaintest

import (
	"testing"

	"github.com/rlibaert/service-example-go/domain"
)

func TestService(t *testing.T, service domain.Service) {
	id, _ := service.ContactsCreate(t.Context(), &domain.Contact{})
	service.ContactsRead(t.Context(), id)
	service.ContactsUpdate(t.Context(), id, &domain.Contact{})
	service.ContactsDelete(t.Context(), id)
}
