// Package auth authenticates requests. M0 only declares the shape; M1
// implements session cookies backed by the database (ADR-0001, NFR-06).
package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrUnauthenticated reports a request without a valid session.
var ErrUnauthenticated = errors.New("auth: unauthenticated")

// Principal is the authenticated user of a request.
type Principal struct {
	UserID uuid.UUID
}

// Authenticator resolves a session token to the user it belongs to.
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (Principal, error)
}
