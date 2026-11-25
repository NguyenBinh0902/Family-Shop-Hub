package kernel

import "github.com/google/uuid"

// NewID returns a new UUID version 7. Version 7 IDs start with a timestamp, so
// they sort by creation time and keep B-tree indexes compact.
func NewID() uuid.UUID {
	// NewV7 only fails if crypto/rand fails, which since Go 1.24 crashes the
	// program instead of returning an error, so Must never panics in practice.
	return uuid.Must(uuid.NewV7())
}
