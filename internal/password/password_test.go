package password

import (
	"fmt"
	"go-bank/internal/testutil"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createHash(t *testing.T, password string) string {
	hash, err := Hash(password)

	require.NoError(t, err)
	require.NotEmpty(t, hash)

	return hash
}

func Test_Hash(t *testing.T) {
	password := testutil.RandomString(10)

	createHash(t, password)
}

func Test_Check(t *testing.T) {
	password := testutil.RandomString(10)
	hashed := createHash(t, password)

	err := Check(hashed, password)
	require.NoError(t, err)
}

func Test_Hash_Error(t *testing.T) {
	password := testutil.RandomString(80)

	hash, err := Hash(password)

	fmt.Println(err)

	require.EqualError(t, err, bcrypt.ErrPasswordTooLong.Error())
	require.Empty(t, hash)
}

func Test_Check_Error(t *testing.T) {
	password := testutil.RandomString(10)
	hashed := createHash(t, password)

	err := Check(hashed, testutil.RandomString(10))
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

func Test_SamePasswordHash(t *testing.T) {
	password := testutil.RandomString(10)
	hashed1 := createHash(t, password)
	hashed2 := createHash(t, password)

	require.NotEqual(t, hashed1, hashed2)
}
