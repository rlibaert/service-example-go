// Package domain provides the core application logic & entities implementing the service.
package domain

import (
	"context"
	"errors"
)

// Store is an abstraction to an underlying storage.
type Store interface {
	// ContactsCreate stores a new [Contact] object and returns its [ContactID].
	ContactsCreate(context.Context, *Contact) (ContactID, error)
	// ContactsRead retrieves a [Contact] given its [ContactID].
	ContactsRead(context.Context, ContactID) (*Contact, error)
	// ContactsUpdate updates a [Contact] given its [ContactID].
	ContactsUpdate(context.Context, ContactID, *Contact) error
	// ContactsDelete deletes a [Contact] given its [ContactID].
	ContactsDelete(context.Context, ContactID) error
}

type Service interface {
	// ContactsCreate stores a new [Contact] object and returns its [ContactID].
	ContactsCreate(context.Context, *Contact) (ContactID, error)
	// ContactsRead retrieves a [Contact] given its [ContactID].
	ContactsRead(context.Context, ContactID) (*Contact, error)
	// ContactsUpdate updates a [Contact] given its [ContactID].
	ContactsUpdate(context.Context, ContactID, *Contact) error
	// ContactsDelete deletes a [Contact] given its [ContactID].
	ContactsDelete(context.Context, ContactID) error
}

var (
	ErrNotFound = errors.New("domain: not found")
	ErrInvalid  = errors.New("domain: invalid argument")
)

// ServiceStore implements [Service] using a [Store].
type ServiceStore struct {
	Store Store
}

var _ Service = (*ServiceStore)(nil)

func (svc *ServiceStore) ContactsCreate(ctx context.Context, c *Contact) (ContactID, error) {
	return svc.Store.ContactsCreate(ctx, c)
}

func (svc *ServiceStore) ContactsRead(ctx context.Context, id ContactID) (*Contact, error) {
	return svc.Store.ContactsRead(ctx, id)
}

func (svc *ServiceStore) ContactsUpdate(ctx context.Context, id ContactID, c *Contact) error {
	return svc.Store.ContactsUpdate(ctx, id, c)
}

func (svc *ServiceStore) ContactsDelete(ctx context.Context, id ContactID) error {
	return svc.Store.ContactsDelete(ctx, id)
}
