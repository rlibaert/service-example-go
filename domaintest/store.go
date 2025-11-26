// Package domaintest provides utilities for [domain] testing.
package domaintest

import (
	"testing"

	"github.com/rlibaert/service-example-go/domain"
)

// TestStore tests a [domain.Store] implementation.
func TestStore(t *testing.T, store domain.Store) {
	// extremely basic test
	id, _ := store.ContactsSet(t.Context(), &domain.Contact{})
	_, _ = store.ContactsGet(t.Context(), id)
	_ = store.ContactsReset(t.Context(), id, &domain.Contact{})
	_ = store.ContactsDel(t.Context(), id)
}

// TestService tests a [domain.Service] implementation.
func TestService(t *testing.T, service domain.Service) {
	// extremely basic test
	id, _ := service.ContactsCreate(t.Context(), &domain.Contact{})
	_, _ = service.ContactsRead(t.Context(), id)
	_ = service.ContactsUpdate(t.Context(), id, &domain.Contact{})
	_ = service.ContactsDelete(t.Context(), id)
}
