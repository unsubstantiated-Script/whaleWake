package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	// ErrExpiredToken is returned when a token has expired.
	ErrExpiredToken = errors.New("token has expired")
	// ErrInvalidToken is returned when a token is invalid.
	ErrInvalidToken                      = errors.New("invalid token")
	ErrMissingPasetoEnvVariable          = errors.New("invalid paseto symmetric key")
	ErrFailedHexToSymmetricKeyConversion = errors.New("failed to convert hex string to symmetric key")
)

// Payload represents the data stored in a token.
// It includes the user ID, issued at time, expiration time, and any other relevant claims.
type Payload struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for the token
	Username  string    `json:"username"`   // The ID of the user associated with the token
	IssuedAt  time.Time `json:"issued_at"`  // The time when the token was issued in Unix timestamp format
	ExpiredAt time.Time `json:"expired_at"` // The expiration time of the token in Unix timestamp format
}

// NewPayload creates a new Payload instance with the given user ID and expiration duration.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the token payload is valid by ensuring it has not expired.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
