package stores

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/rlibaert/service-example-go/domain"
)

// Mock is a mock implementation of [domain.Store] for testing purposes.
type Mock struct {
	mu       sync.Mutex
	index    map[domain.ContactID]int
	contacts []*domain.Contact
}

var _ domain.Store = (*Mock)(nil)

func MustNewMock(cs ...*domain.Contact) *Mock {
	s := &Mock{index: map[domain.ContactID]int{}}
	for _, c := range cs {
		_, err := s.ContactsSet(context.Background(), c)
		if err != nil {
			panic(err)
		}
	}
	return s
}

func (s *Mock) ContactsSet(_ context.Context, c *domain.Contact) (domain.ContactID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ko := true; ko; _, ko = s.index[c.ID] {
		c.ID = domain.ContactID{UUID: uuid.New()}
	}
	s.index[c.ID] = len(s.contacts)
	s.contacts = append(s.contacts, c)
	return c.ID, nil
}

func (s *Mock) ContactsGet(_ context.Context, id domain.ContactID) (*domain.Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.index[id]
	if !ok || s.contacts[index] == nil {
		return nil, domain.ErrNotFound
	}

	return s.contacts[index], nil
}

func (s *Mock) ContactsReset(_ context.Context, id domain.ContactID, c *domain.Contact) error {
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

func (s *Mock) ContactsDel(_ context.Context, id domain.ContactID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	index, ok := s.index[id]
	if !ok || s.contacts[index] == nil {
		return domain.ErrNotFound
	}

	s.contacts[index] = nil
	return nil
}
