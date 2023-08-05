package grpc

import (
	"go-bank/config"
	db "go-bank/db/sqlc"
	"go-bank/internal/testutil"
	"go-bank/internal/worker"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store, distributor worker.TaskDistributor) *GrpcServer {
	config := config.Config{
		TOKEN_SYMMETRIC_KEY:   testutil.RandomString(32),
		ACCESS_TOKEN_DURATION: time.Minute,
	}

	server, err := NewGrpcServer(config, store, distributor)
	require.NoError(t, err)

	return server
}
