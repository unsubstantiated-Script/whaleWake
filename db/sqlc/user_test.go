package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"whaleWake/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		UserName: util.RandomUserName(),
		Email:    util.RandomEmail(),
		Password: hashedPassword,
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
	require.NotZero(t, user.VerifiedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)

	t.Cleanup(func() {
		_, err := testQueries.DeleteUser(context.Background(), user.ID)
		if err != nil {
			return
		}
	})
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)

	// Cleanup should be run before the require statements because if the require statements fail, the cleanup will not be run
	t.Cleanup(func() {
		_, err := testQueries.DeleteUser(context.Background(), user1.ID)
		if err != nil {
			return
		}
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.Equal(t, user1.UpdatedAt, user2.UpdatedAt)
	require.Equal(t, user1.VerifiedAt, user2.VerifiedAt)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		ID:       user1.ID,
		UserName: user1.UserName,
		Email:    util.RandomEmail(),
		Password: user1.Password,
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	// Cleanup should be run before the require statements because if the require statements fail, the cleanup will not be run
	t.Cleanup(func() {
		_, err := testQueries.DeleteUser(context.Background(), user1.ID)
		if err != nil {
			return
		}
	})

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	// Because this is the one item getting updated, it relies upon the arg struct vs. the user1 struct
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.NotEqual(t, user1.UpdatedAt, user2.UpdatedAt)
	require.Equal(t, user1.VerifiedAt, user2.VerifiedAt)

}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	user1, err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	var userSlice []User

	for i := 0; i < 10; i++ {
		userSlice = append(userSlice, createRandomUser(t))
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	t.Cleanup(func() {
		for _, user := range userSlice {
			_, err := testQueries.DeleteUser(context.Background(), user.ID)
			if err != nil {
				return
			}
		}
	})

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

}
