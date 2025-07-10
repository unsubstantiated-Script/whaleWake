package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"whaleWake/util"
)

func createRandomUserRole(t *testing.T, userID uuid.UUID) UserRole {
	arg := CreateUserRoleParams{
		UserID: userID,
		RoleID: int32(util.RandomInt(1, 3)),
	}

	role, err := testQueries.CreateUserRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userID)

	require.Equal(t, arg.UserID, role.UserID)
	require.Equal(t, arg.RoleID, role.RoleID)

	require.NotZero(t, role.ID)
	require.NotZero(t, role.CreatedAt)
	require.NotZero(t, role.UpdatedAt)
	require.NotZero(t, role.VerifiedAt)

	return role
}

func TestCreateUserRole(t *testing.T) {
	user := createRandomUser(t)
	role := createRandomUserRole(t, user.ID)

	t.Cleanup(func() {
		_, err := testQueries.DeleteUserRole(context.Background(), role.ID)
		if err != nil {
			return
		}
	})
}

func TestGetUserRole(t *testing.T) {
	user := createRandomUser(t)
	role1 := createRandomUserRole(t, user.ID)
	role2, err := testQueries.GetUserRole(context.Background(), role1.ID)

	// Cleanup should be run before the require statements because if the require statements fail, the cleanup will not be run
	t.Cleanup(func() {
		_, err := testQueries.DeleteUserRole(context.Background(), role1.ID)
		if err != nil {
			return
		}
	})

	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, role1.UserID, role2.UserID)
	require.Equal(t, role1.RoleID, role2.RoleID)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
	require.Equal(t, role1.UpdatedAt, role2.UpdatedAt)
	require.Equal(t, role1.VerifiedAt, role2.VerifiedAt)
}

func TestUpdateUserRole(t *testing.T) {
	user := createRandomUser(t)
	role1 := createRandomUserRole(t, user.ID)

	arg := UpdateUserRoleParams{
		ID:     role1.ID,
		RoleID: int32(util.RandomInt(1, 3)),
	}

	role2, err := testQueries.UpdateUserRole(context.Background(), arg)

	t.Cleanup(func() {
		_, err := testQueries.DeleteUserRole(context.Background(), role2.ID)
		if err != nil {
			return
		}
	})

	require.NoError(t, err)
	require.NotEmpty(t, role2)

	require.Equal(t, role1.ID, role2.ID)
	require.Equal(t, arg.RoleID, role2.RoleID)
	require.WithinDuration(t, role1.CreatedAt, role2.CreatedAt, time.Second)
	require.NotEqual(t, role1.UpdatedAt, role2.UpdatedAt)
	require.Equal(t, role1.VerifiedAt, role2.VerifiedAt)

}

func TestDeleteUserRole(t *testing.T) {
	user := createRandomUser(t)
	role1 := createRandomUserRole(t, user.ID)

	role1, err := testQueries.DeleteUserRole(context.Background(), role1.ID)
	require.NoError(t, err)

	role2, err := testQueries.GetUserRole(context.Background(), role1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, role2)
}

func TestListUserRoles(t *testing.T) {
	user := createRandomUser(t)
	var roleSlice []UserRole

	for i := 0; i < 10; i++ {
		role := createRandomUserRole(t, user.ID)
		roleSlice = append(roleSlice, role)
	}

	arg := ListUserRolesParams{
		Limit:  5,
		Offset: 5,
	}

	roles, err := testQueries.ListUserRoles(context.Background(), arg)

	t.Cleanup(func() {
		for _, role := range roleSlice {
			_, err := testQueries.DeleteUserRole(context.Background(), role.ID)
			if err != nil {
				return
			}
		}
	})

	require.NoError(t, err)
	require.Len(t, roles, 5)

	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}
