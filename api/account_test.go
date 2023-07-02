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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_GetAccount(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	// build stub
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start server and send req
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	server.router.ServeHTTP(recorder, req)

	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       testutil.RandomInt(1, 1000),
		Owner:    testutil.RandomOwner(),
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
