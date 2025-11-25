package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContactID struct{ uuid.UUID }

type Contact struct {
	ID        ContactID
	Firstname string
	Lastname  string
	Birthday  time.Time
}
