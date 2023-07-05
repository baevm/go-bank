package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// The RegisteredClaims is embedded in the custom type to allow for easy encoding, parsing and validation of registered claims.
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`

	jwt.RegisteredClaims
}

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	expTime := time.Now().Add(duration)

	payload := &Payload{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt.Time) {
		return ErrTokenExpired
	}

	return nil
}
