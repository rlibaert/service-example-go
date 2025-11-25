package domain

import (
	"context"
	"errors"
)

type Store interface {
	ContactsCreate(context.Context, *Contact) (ContactID, error)
	ContactsRead(context.Context, ContactID) (*Contact, error)
	ContactsUpdate(context.Context, ContactID, *Contact) error
	ContactsDelete(context.Context, ContactID) error

	Tx(context.Context, func(ctx context.Context, tx Store) error) error
}

type Service interface {
	ContactsCreate(context.Context, *Contact) (ContactID, error)
	ContactsRead(context.Context, ContactID) (*Contact, error)
	ContactsUpdate(context.Context, ContactID, *Contact) error
	ContactsDelete(context.Context, ContactID) error
}

var (
	ErrNotFound = errors.New("domain: not found")
	ErrInvalid  = errors.New("domain: invalid argument")
)

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
