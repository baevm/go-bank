package token

import (
	"go-bank/internal/testutil"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func Test_JWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(testutil.RandomString(32))
	require.NoError(t, err)

	username := testutil.RandomString(10)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := maker.Create(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Verify(token)

	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt.Time, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiresAt.Time, expiredAt, time.Second)
}

func Test_ExpiredJWT(t *testing.T) {
	maker, err := NewJWTMaker(testutil.RandomString(32))
	require.NoError(t, err)

	username := testutil.RandomString(10)

	token, err := maker.Create(username, -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Verify(token)
	require.Error(t, err)
	require.Nil(t, payload)
}

func Test_InvalidSignature(t *testing.T) {
	payload, err := NewPayload(testutil.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	require.NoError(t, err)

	maker, err := NewJWTMaker(testutil.RandomString(32))
	require.NoError(t, err)

	verifiedPayload, err := maker.Verify(token)

	require.Error(t, err)
	require.Nil(t, verifiedPayload)
}
