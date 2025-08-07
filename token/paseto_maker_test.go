package token

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"whaleWake/util"
)

func TestPasetoMaker(t *testing.T) {

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	maker, err := NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		t.Fatalf("Failed to create PasetoMaker: %v", err)
	}

	userID := util.RandomUUID()
	roleID := 1
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(userID, roleID, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	maker, err := NewPasetoMaker(config.TokenSymmetricKey)
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomUUID(), 1, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, "token has expired")
	require.Nil(t, payload)
}
