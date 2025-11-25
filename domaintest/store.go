package domaintest

import (
	"testing"

	"github.com/rlibaert/service-example-go/domain"
)

func TestStore(t *testing.T, store domain.Store) {
	id, _ := store.ContactsCreate(t.Context(), &domain.Contact{})
	store.ContactsRead(t.Context(), id)
	store.ContactsUpdate(t.Context(), id, &domain.Contact{})
	store.ContactsDelete(t.Context(), id)
}
