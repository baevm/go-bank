package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type Maker interface {
	Create(username string, duration time.Duration) (string, error)
	Verify(token string) (*Payload, error)
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be atleast %v characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (m *JWTMaker) Create(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) Verify(token string) (*Payload, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(m.secretKey), nil
	})

	if payload, ok := parsedToken.Claims.(*Payload); ok && parsedToken.Valid {
		return payload, nil
	} else {
		return nil, err
	}
}
