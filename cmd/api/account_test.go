package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	mockdb "go-bank/db/mock"
	db "go-bank/db/sqlc"
	"go-bank/internal/testutil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_GetAccount(t *testing.T) {
	randomUser, _ := RandomUser(t)
	account := randomAccount(randomUser.Username)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	// build stub
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start server and send req
	server := NewTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	addAuthorization(t, req, server.tokenMaker, authorizationTypeBearer, randomUser.Username, time.Minute)

	server.router.ServeHTTP(recorder, req)

	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount(owner string) db.Accounts {
	return db.Accounts{
		ID:       testutil.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  testutil.RandomMoney(),
		Currency: testutil.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var resAccount db.Accounts
	err = json.Unmarshal(data, &resAccount)
	require.NoError(t, err)
	require.Equal(t, account, resAccount)
}
