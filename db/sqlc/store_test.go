package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"whaleWake/util"
)

func TestCreateUserWithProfileAndRoleTx(t *testing.T) {
	store := NewStore(testDB)

	user := User{
		UserName: util.RandomUserName(),
		Email:    util.RandomEmail(),
		Password: util.RandomPassword(),
	}

	userProfile := UserProfile{
		FirstName:     util.RandomUserName(),
		LastName:      util.RandomUserName(),
		BusinessName:  util.RandomBusinessName(),
		StreetAddress: util.RandomStreetAddress(),
		City:          util.RandomString(6),
		State:         util.RandomCountryCodeOrState(),
		Zip:           util.RandomString(5),
		CountryCode:   util.RandomCountryCodeOrState(),
	}

	userRole := UserRole{
		RoleID: int32(util.RandomInt(1, 3)),
	}

	result, err := store.CreateUserWithProfileAndRoleTx(context.Background(),
		CreateUserParams{
			UserName: user.UserName,
			Email:    user.Email,
			Password: user.Password},
		CreateUserProfileParams{
			FirstName:     userProfile.FirstName,
			LastName:      userProfile.LastName,
			BusinessName:  userProfile.BusinessName,
			StreetAddress: userProfile.StreetAddress,
			City:          userProfile.City,
			State:         userProfile.State,
			Zip:           userProfile.Zip,
			CountryCode:   userProfile.CountryCode,
		},
		CreateUserRoleParams{
			RoleID: userRole.RoleID,
		})

	// Check for errors
	require.NoError(t, err)
	require.NotEmpty(t, result)

	//Check if the user, user profile, and user role were created successfully
	userResult := result.User
	require.NotEmpty(t, result.User)
	require.Equal(t, user.UserName, userResult.UserName)
	require.Equal(t, user.Email, userResult.Email)
	require.Equal(t, user.Password, userResult.Password)

	_, err = store.GetUser(context.Background(), userResult.ID)
	require.NoError(t, err)

	userProfileResult := result.UserProfile
	require.NotEmpty(t, userProfileResult)
	//While this does peel off of the transaction result, it's still viable for testing the ID in the UserProfile table
	require.Equal(t, userResult.ID, userProfileResult.UserID)
	require.Equal(t, userProfile.FirstName, userProfileResult.FirstName)
	require.Equal(t, userProfile.LastName, userProfileResult.LastName)
	require.Equal(t, userProfile.BusinessName, userProfileResult.BusinessName)
	require.Equal(t, userProfile.StreetAddress, userProfileResult.StreetAddress)
	require.Equal(t, userProfile.City, userProfileResult.City)
	require.Equal(t, userProfile.State, userProfileResult.State)
	require.Equal(t, userProfile.Zip, userProfileResult.Zip)
	require.Equal(t, userProfile.CountryCode, userProfileResult.CountryCode)

	_, err = store.GetUserProfile(context.Background(), userProfileResult.ID)
	require.NoError(t, err)

	userRoleResult := result.UserRole
	require.NotEmpty(t, userRoleResult)
	//While this does peel off of the transaction result, it's still viable for testing the ID in the UserRole table
	require.Equal(t, userResult.ID, userRoleResult.UserID)
	require.Equal(t, userRole.RoleID, userRoleResult.RoleID)

	_, err = store.GetUserRole(context.Background(), userRoleResult.ID)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = testQueries.DeleteUserProfile(context.Background(), result.UserProfile.ID)
		if err != nil {
			return
		}

		_, err = testQueries.DeleteUserRole(context.Background(), result.UserRole.ID)
		if err != nil {
			return
		}

		_, err = testQueries.DeleteUser(context.Background(), result.User.ID)
		if err != nil {
			return
		}
	})

}
