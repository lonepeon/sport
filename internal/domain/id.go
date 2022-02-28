package domain

import "github.com/google/uuid"

// ID represents an internal identifier
type ID uuid.UUID

// NewID generates a new internal identifier
func NewID() ID {
	return ID(uuid.New())
}

// String returns the string representation of the identifier
func (id ID) String() string {
	return uuid.UUID(id).String()
}
