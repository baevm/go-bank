package api

import (
	db "go-bank/db/sqlc"
	"go-bank/internal/password"
	"go-bank/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func RandomUser(t *testing.T) (user db.Users, hashedPass string) {
	pass := testutil.RandomString(6)
	hashedPass, err := password.Hash(pass)
	require.NoError(t, err)

	user = db.Users{
		Username:   testutil.RandomOwner(),
		Email:      testutil.RandomEmail(),
		FullName:   testutil.RandomOwner(),
		HashedPass: testutil.RandomString(20),
	}

	return
}
