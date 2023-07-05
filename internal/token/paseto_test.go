package token

import (
	"go-bank/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_PasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(testutil.RandomString(32))
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

func Test_ExpiredPaseto(t *testing.T) {
	maker, err := NewPasetoMaker(testutil.RandomString(32))
	require.NoError(t, err)

	username := testutil.RandomString(10)

	token, err := maker.Create(username, -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Verify(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
