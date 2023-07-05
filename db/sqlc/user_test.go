package db

import (
	"context"
	"go-bank/internal/password"
	"go-bank/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	hashedPass, err := password.Hash(testutil.RandomString(5))
	require.NoError(t, err)

	args := CreateUserParams{
		Email:      testutil.RandomEmail(),
		HashedPass: hashedPass,
		Username:   testutil.RandomOwner(),
		FullName:   testutil.RandomOwner(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Email, user.Email)
	require.Equal(t, args.HashedPass, user.HashedPass)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.FullName, user.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func Test_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func Test_GetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPass, user2.HashedPass)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
