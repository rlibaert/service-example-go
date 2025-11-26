package domain

import (
	"time"

	"github.com/google/uuid"
)

// ContactID represents a [Contact] ID.
type ContactID struct{ uuid.UUID }

// Contact contains a contact's personal informations.
type Contact struct {
	ID        ContactID
	Firstname string
	Lastname  string
	Birthday  time.Time
}
