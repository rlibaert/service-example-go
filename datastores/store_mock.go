package datastores

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/rlibaert/service-example-go/domain"
)

// StoreMock is a mock implementation of [domain.Store] for testing purposes.
type StoreMock struct {
	mu       sync.Mutex
	index    map[domain.ContactID]int
	contacts []*domain.Contact
}

var _ domain.Store = (*StoreMock)(nil)

func MustNewStoreMock(cs ...*domain.Contact) *StoreMock {
	s := &StoreMock{index: map[domain.ContactID]int{}}
	for _, c := range cs {
		_, err := s.ContactsCreate(context.Background(), c)
		if err != nil {
			panic(err)
		}
	}
	return s
}

func (s *StoreMock) Tx(ctx context.Context, f func(context.Context, domain.Store) error) error {
	err := f(ctx, s)
	if err != nil {
		panic("store: cannot rollback")
	}
	return nil
}

func (s *StoreMock) ContactsCreate(_ context.Context, c *domain.Contact) (domain.ContactID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ko := true; ko; _, ko = s.index[c.ID] {
		c.ID = domain.ContactID{UUID: uuid.New()}
	}
	s.index[c.ID] = len(s.contacts)
	s.contacts = append(s.contacts, c)
	return c.ID, nil
}

func (s *StoreMock) ContactsRead(_ context.Context, id domain.ContactID) (*domain.Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.index[id]
	if !ok || s.contacts[index] == nil {
		return nil, domain.ErrNotFound
	}

	return s.contacts[index], nil
}

func (s *StoreMock) ContactsUpdate(_ context.Context, id domain.ContactID, c *domain.Contact) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.index[id]
	if !ok || s.contacts[index] == nil {
		return domain.ErrNotFound
	}

	c.ID = id
	s.contacts[index] = c
	return nil
}

func (s *StoreMock) ContactsDelete(_ context.Context, id domain.ContactID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.index[id]
	if !ok || s.contacts[index] == nil {
		return domain.ErrNotFound
	}

	s.contacts[index] = nil
	return nil
}
