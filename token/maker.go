package token

import (
	"github.com/google/uuid"
	"time"
)

type Maker interface {
	// CreateToken generates a new token for a specific user ID and duration.
	CreateToken(userID uuid.UUID, userRole int, duration time.Duration) (string, error)
	// VerifyToken checks the validity of a token and returns the user ID if valid.
	VerifyToken(token string) (*Payload, error)
	// RefreshToken generates a new token for a user ID based on an existing token.
	RefreshToken(token string) (string, error)
}
