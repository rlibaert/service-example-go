// Package restapi provide primitives to expose [domain] with a REST interface.
package restapi

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/rlibaert/service-example-go/domain"
)

// ServiceRegisterer registers endpoints in a [huma.API] to expose a [domain.Service] with a REST interface.
type ServiceRegisterer struct {
	Service domain.Service
}

type ContactIDModel struct {
	ID domain.ContactID `json:"id" readOnly:"true"`
}

type ContactModel struct {
	ContactIDModel

	Firstname string `json:"firstname" example:"john"`
	Lastname  string `json:"lastname"  example:"smith"`
	Birthday  string `json:"birthday"  example:"1999-12-31" format:"date"`
}

func (reg ServiceRegisterer) RegisterContactsCreate(api huma.API) {
	type input struct {
		Body ContactModel
	}
	type output struct {
		Body ContactIDModel
	}

	handler := func(ctx context.Context, i *input) (*output, error) {
		birthday, err := time.Parse(time.DateOnly, i.Body.Birthday)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid format for birthday", err)
		}

		id, err := reg.Service.ContactsCreate(ctx, &domain.Contact{
			Firstname: i.Body.Firstname,
			Lastname:  i.Body.Lastname,
			Birthday:  birthday,
		})
		if err != nil {
			return nil, err
		}

		return &output{Body: ContactIDModel{id}}, nil
	}

	huma.Post(api, "/contacts", handler)
}

func (reg ServiceRegisterer) RegisterContactsRead(api huma.API) {
	type input struct {
		ContactID domain.ContactID `path:"id"`
	}
	type output struct {
		Body ContactModel
	}

	handler := func(ctx context.Context, input *input) (*output, error) {
		c, err := reg.Service.ContactsRead(ctx, input.ContactID)
		if err != nil {
			return nil, err
		}

		return &output{Body: ContactModel{
			ContactIDModel: ContactIDModel{c.ID},
			Firstname:      c.Firstname,
			Lastname:       c.Lastname,
			Birthday:       c.Birthday.Format(time.DateOnly),
		}}, nil
	}

	huma.Get(api, "/contacts/{id}", handler)
}

func (reg ServiceRegisterer) RegisterContactsUpdate(api huma.API) {
	type input struct {
		ContactID domain.ContactID `path:"id"`
		Body      ContactModel
	}
	type output struct{}

	handler := func(ctx context.Context, i *input) (*output, error) {
		birthday, err := time.Parse(time.DateOnly, i.Body.Birthday)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid format for birthday", err)
		}

		return nil, reg.Service.ContactsUpdate(ctx, i.ContactID, &domain.Contact{
			Firstname: i.Body.Firstname,
			Lastname:  i.Body.Lastname,
			Birthday:  birthday,
		})
	}

	huma.Put(api, "/contacts/{id}", handler)
}

func (reg ServiceRegisterer) RegisterContactsDelete(api huma.API) {
	type input struct {
		ContactID domain.ContactID `path:"id"`
	}
	type output struct{}

	handler := func(ctx context.Context, input *input) (*output, error) {
		return nil, reg.Service.ContactsDelete(ctx, input.ContactID)
	}

	huma.Delete(api, "/contacts/{id}", handler)
}
