package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/rlibaert/service-example-go/domain"
	"github.com/rlibaert/service-example-go/grpc/proto"
)

type ServiceServer struct {
	proto.UnimplementedServiceServer

	Service domain.Service
}

var _ proto.ServiceServer = ServiceServer{}

func (s ServiceServer) ContactsCreate(ctx context.Context, c *proto.Contact) (*proto.ContactID, error) {
	id, err := s.Service.ContactsCreate(ctx, &domain.Contact{
		Firstname: c.Firstname,
		Lastname:  c.Firstname,
	})
	if err != nil {
		return nil, err
	}

	return &proto.ContactID{Id: id.String()}, nil
}

func (s ServiceServer) ContactsRead(ctx context.Context, id *proto.ContactID) (*proto.Contact, error) {
	c, err := s.Service.ContactsRead(ctx, domain.ContactID{UUID: uuid.MustParse(id.Id)})
	if err != nil {
		return nil, err
	}

	return &proto.Contact{
		Firstname: c.Firstname,
		Lastname:  c.Lastname,
	}, nil
}

func (s ServiceServer) ContactsDelete(ctx context.Context, id *proto.ContactID) (*proto.Empty, error) {
	err := s.Service.ContactsDelete(ctx, domain.ContactID{UUID: uuid.MustParse(id.Id)})
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}
