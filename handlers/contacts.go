package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	ds "github.com/rlibaert/service-example-go/datastores"
)

type Contacts struct {
	Store        ds.ContactsStore
	ErrorHandler func(context.Context, error)
}

type ContactIDModel struct {
	ID ds.ContactID `json:"id" readOnly:"true"`
}

type ContactModel struct {
	ContactIDModel

	Firstname string `json:"firstname" example:"john"`
	Lastname  string `json:"lastname"  example:"smith"`
	Birthday  string `json:"birthday"  example:"1999-12-31" format:"date"`
}

func (h *Contacts) RegisterCreate(api huma.API) { // called by [huma.AutoRegister]
	huma.Put(api, "/",
		handlerWithErrorHandler(h.create, h.ErrorHandler),
		opErrors(http.StatusUnprocessableEntity, http.StatusInternalServerError),
	)
}

type ContactsCreateOutput struct {
	Body ContactIDModel
}

func (h *Contacts) create(ctx context.Context, input *struct {
	Body ContactModel
}) (*ContactsCreateOutput, error) {
	birthday, err := time.Parse(time.DateOnly, input.Body.Birthday)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("invalid format for birthday", err)
	}

	id, err := h.Store.Create(ctx, &ds.Contact{
		Firstname: input.Body.Firstname,
		Lastname:  input.Body.Lastname,
		Birthday:  birthday,
	})
	if err != nil {
		return nil, err
	}

	return &ContactsCreateOutput{Body: ContactIDModel{id}}, nil
}

func (h *Contacts) RegisterList(api huma.API) { // called by [huma.AutoRegister]
	huma.Get(api, "/",
		handlerWithErrorHandler(h.list, h.ErrorHandler),
		opErrors(http.StatusInternalServerError),
	)
}

type ContactsListOutput struct {
	Body []ContactModel
}

func (h *Contacts) list(ctx context.Context, input *struct {
	Page int `query:"page" default:"1"  minimum:"1"`
	Size int `query:"size" default:"10" minimum:"1" maximum:"100"`
}) (*ContactsListOutput, error) {
	contacts, err := h.Store.List(ctx, (input.Page-1)*input.Size, input.Size)
	if err != nil {
		return nil, err
	}

	body := make([]ContactModel, 0, len(contacts))
	for _, contact := range contacts {
		body = append(body, ContactModel{
			ContactIDModel: ContactIDModel{contact.ID},
			Firstname:      contact.Firstname,
			Lastname:       contact.Lastname,
			Birthday:       contact.Birthday.Format(time.DateOnly),
		})
	}

	return &ContactsListOutput{Body: body}, nil
}

func (h *Contacts) RegisterGet(api huma.API) { // called by [huma.AutoRegister]
	huma.Get(api, "/{id}",
		handlerWithErrorHandler(h.get, h.ErrorHandler),
		opErrors(http.StatusNotFound, http.StatusInternalServerError),
	)
}

type ContactsGetOutput struct {
	Body ContactModel
}

func (h *Contacts) get(ctx context.Context, input *struct {
	ID ds.ContactID `path:"id" doc:"ID of the contact to get"`
}) (*ContactsGetOutput, error) {
	contact, err := h.Store.Get(ctx, input.ID)
	switch {
	case err == nil:
		return &ContactsGetOutput{Body: ContactModel{
			ContactIDModel: ContactIDModel{contact.ID},
			Firstname:      contact.Firstname,
			Lastname:       contact.Lastname,
			Birthday:       contact.Birthday.Format(time.DateOnly),
		}}, nil

	case errors.Is(err, ds.ErrObjectNotFound):
		return nil, huma.Error404NotFound("id not found", err)

	default:
		return nil, err
	}
}

func (h *Contacts) RegisterDelete(api huma.API) { // called by [huma.AutoRegister]
	huma.Delete(api, "/{id}",
		handlerWithErrorHandler(h.delete, h.ErrorHandler),
		opErrors(http.StatusInternalServerError),
	)
}

func (h *Contacts) delete(ctx context.Context, input *struct {
	ID ds.ContactID `path:"id" doc:"ID of the contact to delete"`
}) (*struct{}, error) {
	return nil, h.Store.Delete(ctx, input.ID)
}
