package domaintest

import (
	"testing"

	"github.com/rlibaert/service-example-go/domain"
)

func TestStore(t *testing.T, store domain.Store) {
	id, _ := store.ContactsCreate(t.Context(), &domain.Contact{})
	_, _ = store.ContactsRead(t.Context(), id)
	_ = store.ContactsUpdate(t.Context(), id, &domain.Contact{})
	_ = store.ContactsDelete(t.Context(), id)
}
