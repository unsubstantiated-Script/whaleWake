package api

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	db "whaleWake/db/sqlc"
	"whaleWake/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomSymmetricKey(),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

//func TestMain(m *testing.M) {
//	gin.SetMode(gin.TestMode)
//
//	os.Exit(m.Run())
//}
