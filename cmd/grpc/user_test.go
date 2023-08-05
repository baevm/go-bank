package grpc

import (
	"context"
	"fmt"
	mockdb "go-bank/db/mock"
	db "go-bank/db/sqlc"
	"go-bank/internal/password"
	"go-bank/internal/testutil"
	"go-bank/pb"
	mockwrk "go-bank/internal/worker/mock"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.Users
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := password.Check(expected.password, actualArg.HashedPass)
	if err != nil {
		return false
	}

	expected.arg.HashedPass = actualArg.HashedPass
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	return true
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.Users) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

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

func Test_CreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	user, pass := RandomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.CreateResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateRequest{
				Username: user.Username,
				Password: pass,
				Email:    user.Email,
				FullName: user.FullName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						Email:    user.Email,
						FullName: user.FullName,
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(args, pass, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "ok", res.GetMessage())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			distributor := mockwrk.NewMockTaskDistributor(ctrl)

			server := NewTestServer(t, store, distributor)

			res, err := server.CreateUser(context.Background(), tc.req)

			tc.checkResponse(t, res, err)

		})

	}
}
