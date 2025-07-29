package token

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"whaleWake/util"
)

func TestPasetoMaker(t *testing.T) {

	maker, err := NewPasetoMaker()

	if err != nil {
		t.Fatalf("Failed to create PasetoMaker: %v", err)
	}

	username := util.RandomUserName()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker()
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomUserName(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, "token has expired")
	require.Nil(t, payload)
}
