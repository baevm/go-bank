package api

import (
	"go-bank/config"
	db "go-bank/db/sqlc"
	"go-bank/internal/testutil"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *HTTPServer {
	config := config.Config{
		TOKEN_SYMMETRIC_KEY:   testutil.RandomString(32),
		ACCESS_TOKEN_DURATION: time.Minute,
	}

	server, err := NewHTTPServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
