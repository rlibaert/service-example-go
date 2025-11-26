package wrappers

import (
	"context"

	"github.com/rlibaert/service-example-go/domain"
)

// ServiceErrorHandler wraps a [domain.Service] to handle errors.
type ServiceErrorHandler struct {
	Service      domain.Service
	ErrorHandler func(context.Context, error)
}

func (service ServiceErrorHandler) handle(ctx context.Context, err error) {
	if err != nil {
		service.ErrorHandler(ctx, err)
	}
}

func (service ServiceErrorHandler) ContactsCreate(ctx context.Context, c *domain.Contact) (domain.ContactID, error) {
	id, err := service.Service.ContactsCreate(ctx, c)
	service.handle(ctx, err)
	return id, err
}

func (service ServiceErrorHandler) ContactsRead(ctx context.Context, id domain.ContactID) (*domain.Contact, error) {
	c, err := service.Service.ContactsRead(ctx, id)
	service.handle(ctx, err)
	return c, err
}

func (service ServiceErrorHandler) ContactsUpdate(ctx context.Context, id domain.ContactID, c *domain.Contact) error {
	err := service.Service.ContactsUpdate(ctx, id, c)
	service.handle(ctx, err)
	return err
}

func (service ServiceErrorHandler) ContactsDelete(ctx context.Context, id domain.ContactID) error {
	err := service.Service.ContactsDelete(ctx, id)
	service.handle(ctx, err)
	return err
}
