package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"whaleWake/util"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		UserName: util.RandomUserName(),
		Email:    util.RandomEmail(),
		Password: util.RandomPassword(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	t.Cleanup(func() {
		_, err := testQueries.DeleteUser(context.Background(), user.ID)
		if err != nil {
			return
		}
	})
}
