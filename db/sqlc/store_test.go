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

func TestGetUserWithProfileAndRoleTx(t *testing.T) {
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

	// Quick check on the Create Mock here.
	require.NoError(t, err)
	require.NotEmpty(t, result)

	fetched, err := store.GetUserWithProfileAndRoleTX(context.Background(), result.User.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetched)
	require.Equal(t, result.User.ID, fetched.UserProfile.UserID)
	require.Equal(t, result.User.ID, fetched.User.ID)
	require.Equal(t, result.UserProfile.ID, fetched.UserProfile.ID)
	require.Equal(t, result.UserRole.ID, fetched.UserRole.ID)

	t.Cleanup(func() {
		_, _ = testQueries.DeleteUserProfile(context.Background(), result.UserProfile.ID)
		_, _ = testQueries.DeleteUserRole(context.Background(), result.UserRole.ID)
		_, _ = testQueries.DeleteUser(context.Background(), result.User.ID)
	})

}

func TestDeleteUserWithProfileAndRoleTX(t *testing.T) {
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

	// Quick check on the Create Mock here.
	require.NoError(t, err)
	require.NotEmpty(t, result)

	_, err = store.DeleteUserWithProfileAndRoleTX(context.Background(), result.User.ID)
	require.NoError(t, err)

	_, err = store.GetUser(context.Background(), result.User.ID)
	require.Error(t, err)

	_, err = store.GetUserProfile(context.Background(), result.UserProfile.ID)
	require.Error(t, err)

	_, err = store.GetUserRole(context.Background(), result.UserRole.ID)
	require.Error(t, err)

}

func TestUpdateUserWithProfileAndRoleTX(t *testing.T) {
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

	// Quick check on the Create Mock here.
	require.NoError(t, err)
	require.NotEmpty(t, result)

	//Checking all the values of the original creation
	require.Equal(t, user.UserName, result.User.UserName)
	require.Equal(t, user.Email, result.User.Email)
	require.Equal(t, user.Password, result.User.Password)
	require.Equal(t, userProfile.FirstName, result.UserProfile.FirstName)
	require.Equal(t, userProfile.LastName, result.UserProfile.LastName)
	require.Equal(t, userProfile.BusinessName, result.UserProfile.BusinessName)
	require.Equal(t, userProfile.StreetAddress, result.UserProfile.StreetAddress)
	require.Equal(t, userProfile.City, result.UserProfile.City)
	require.Equal(t, userProfile.State, result.UserProfile.State)
	require.Equal(t, userProfile.Zip, result.UserProfile.Zip)
	require.Equal(t, userProfile.CountryCode, result.UserProfile.CountryCode)
	require.Equal(t, userRole.RoleID, result.UserRole.RoleID)

	// Adding in random updates to all the fields in the transaction
	userUpdate := User{
		UserName: util.RandomUserName(),
		Email:    util.RandomEmail(),
		Password: util.RandomPassword(),
	}

	userProfileUpdate := UserProfile{
		FirstName:     util.RandomUserName(),
		LastName:      util.RandomUserName(),
		BusinessName:  util.RandomBusinessName(),
		StreetAddress: util.RandomStreetAddress(),
		City:          util.RandomString(6),
		State:         util.RandomCountryCodeOrState(),
		Zip:           util.RandomString(5),
		CountryCode:   util.RandomCountryCodeOrState(),
	}

	userRoleUpdate := UserRole{
		RoleID: int32(util.RandomInt(1, 3)),
	}

	updatedResult, err := store.UpdateUserWithProfileAndRoleTX(context.Background(),
		UpdateUserParams{
			ID:       result.User.ID,
			UserName: userUpdate.UserName,
			Email:    userUpdate.Email,
			Password: userUpdate.Password},
		UpdateUserProfileParams{
			UserID:        result.User.ID,
			FirstName:     userProfileUpdate.FirstName,
			LastName:      userProfileUpdate.LastName,
			BusinessName:  userProfileUpdate.BusinessName,
			StreetAddress: userProfileUpdate.StreetAddress,
			City:          userProfileUpdate.City,
			State:         userProfileUpdate.State,
			Zip:           userProfileUpdate.Zip,
			CountryCode:   userProfileUpdate.CountryCode,
		},
		UpdateUserRoleParams{
			UserID: result.User.ID,
			RoleID: userRoleUpdate.RoleID,
		})

	// Quick check on the Create Mock here.
	require.NoError(t, err)
	require.NotEmpty(t, result)

	//Checking that the update took place
	require.Equal(t, userUpdate.UserName, updatedResult.User.UserName)
	require.Equal(t, userUpdate.Email, updatedResult.User.Email)
	require.Equal(t, userUpdate.Password, updatedResult.User.Password)
	require.Equal(t, userProfileUpdate.FirstName, updatedResult.UserProfile.FirstName)
	require.Equal(t, userProfileUpdate.LastName, updatedResult.UserProfile.LastName)
	require.Equal(t, userProfileUpdate.BusinessName, updatedResult.UserProfile.BusinessName)
	require.Equal(t, userProfileUpdate.StreetAddress, updatedResult.UserProfile.StreetAddress)
	require.Equal(t, userProfileUpdate.City, updatedResult.UserProfile.City)
	require.Equal(t, userProfileUpdate.State, updatedResult.UserProfile.State)
	require.Equal(t, userProfileUpdate.Zip, updatedResult.UserProfile.Zip)
	require.Equal(t, userProfileUpdate.CountryCode, updatedResult.UserProfile.CountryCode)
	require.Equal(t, userRoleUpdate.RoleID, updatedResult.UserRole.RoleID)

	//Now checking that update doesn't match the original anymore
	require.NotEqual(t, result.User.UserName, updatedResult.User.UserName)
	require.NotEqual(t, result.User.Email, updatedResult.User.Email)
	require.NotEqual(t, result.User.Password, updatedResult.User.Password)
	require.NotEqual(t, result.UserProfile.FirstName, updatedResult.UserProfile.FirstName)
	require.NotEqual(t, result.UserProfile.LastName, updatedResult.UserProfile.LastName)
	require.NotEqual(t, result.UserProfile.BusinessName, updatedResult.UserProfile.BusinessName)
	require.NotEqual(t, result.UserProfile.StreetAddress, updatedResult.UserProfile.StreetAddress)
	require.NotEqual(t, result.UserProfile.City, updatedResult.UserProfile.City)
	require.NotEqual(t, result.UserProfile.State, updatedResult.UserProfile.State)
	require.NotEqual(t, result.UserProfile.Zip, updatedResult.UserProfile.Zip)
	require.NotEqual(t, result.UserProfile.CountryCode, updatedResult.UserProfile.CountryCode)
	require.NotEqual(t, result.UserRole.RoleID, updatedResult.UserRole.RoleID)

	_, err = store.DeleteUserWithProfileAndRoleTX(context.Background(), result.User.ID)
	require.NoError(t, err)

	_, err = store.GetUser(context.Background(), result.User.ID)
	require.Error(t, err)

	_, err = store.GetUserProfile(context.Background(), result.UserProfile.ID)
	require.Error(t, err)

	_, err = store.GetUserRole(context.Background(), result.UserRole.ID)
	require.Error(t, err)
}
