package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
	"whaleWake/util"
)

func createRandomUserProfile(t *testing.T, userID uuid.UUID) UserProfile {
	arg := CreateUserProfileParams{
		UserID:        userID,
		FirstName:     util.RandomUserName(),
		LastName:      util.RandomUserName(),
		BusinessName:  util.RandomBusinessName(),
		StreetAddress: util.RandomStreetAddress(),
		City:          util.RandomString(6),
		State:         util.RandomCountryCodeOrState(),
		Zip:           strconv.FormatInt(util.RandomInt(11111, 99999), 10),
		CountryCode:   util.RandomCountryCodeOrState(),
	}

	profile, err := testQueries.CreateUserProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, profile)

	return profile
}

func TestCreateUserProfile(t *testing.T) {
	user := createRandomUser(t)
	profile := createRandomUserProfile(t, user.ID)

	t.Cleanup(func() {
		_, _ = testQueries.DeleteUserProfile(context.Background(), profile.UserID)
		_, _ = testQueries.DeleteUser(context.Background(), user.ID)
	})
}

func TestGetUserProfile(t *testing.T) {
	user := createRandomUser(t)
	profile1 := createRandomUserProfile(t, user.ID)
	profile2, err := testQueries.GetUserProfile(context.Background(), profile1.UserID)

	// Cleanup should be run before the require statements because if the require statements fail, the cleanup will not be run
	t.Cleanup(func() {
		_, _ = testQueries.DeleteUserProfile(context.Background(), profile1.UserID)
		_, _ = testQueries.DeleteUser(context.Background(), user.ID)
	})

	require.NoError(t, err)
	require.NotEmpty(t, profile2)

	require.Equal(t, profile1.ID, profile2.ID)
	require.Equal(t, profile1.UserID, profile2.UserID)
	require.Equal(t, profile1.FirstName, profile2.FirstName)
	require.Equal(t, profile1.LastName, profile2.LastName)
	require.Equal(t, profile1.BusinessName, profile2.BusinessName)
	require.Equal(t, profile1.StreetAddress, profile2.StreetAddress)
	require.Equal(t, profile1.City, profile2.City)
	require.Equal(t, profile1.State, profile2.State)
	require.Equal(t, profile1.Zip, profile2.Zip)
	require.Equal(t, profile1.CountryCode, profile2.CountryCode)
	require.WithinDuration(t, profile1.CreatedAt, profile2.CreatedAt, time.Second)
	require.Equal(t, profile1.UpdatedAt, profile2.UpdatedAt)
	require.Equal(t, profile1.VerifiedAt, profile2.VerifiedAt)
}

func TestUpdateUserProfile(t *testing.T) {
	user := createRandomUser(t)
	profile1 := createRandomUserProfile(t, user.ID)

	arg := UpdateUserProfileParams{
		UserID:        profile1.UserID,
		FirstName:     util.RandomUserName(),
		LastName:      util.RandomUserName(),
		BusinessName:  util.RandomBusinessName(),
		StreetAddress: util.RandomStreetAddress(),
		City:          util.RandomString(6),
		State:         util.RandomCountryCodeOrState(),
		Zip:           strconv.FormatInt(util.RandomInt(11111, 99999), 10),
		CountryCode:   util.RandomCountryCodeOrState(),
	}

	profile2, err := testQueries.UpdateUserProfile(context.Background(), arg)

	t.Cleanup(func() {
		_, _ = testQueries.DeleteUserProfile(context.Background(), profile2.UserID)
		_, _ = testQueries.DeleteUser(context.Background(), user.ID)
	})

	require.NoError(t, err)
	require.NotEmpty(t, profile2)

	require.Equal(t, profile1.ID, profile2.ID)
	require.Equal(t, arg.FirstName, profile2.FirstName)
	require.Equal(t, arg.LastName, profile2.LastName)
	require.Equal(t, arg.BusinessName, profile2.BusinessName)
	require.Equal(t, arg.StreetAddress, profile2.StreetAddress)
	require.Equal(t, arg.City, profile2.City)
	require.Equal(t, arg.State, profile2.State)
	require.Equal(t, arg.Zip, profile2.Zip)
	require.Equal(t, arg.CountryCode, profile2.CountryCode)
	require.WithinDuration(t, profile1.CreatedAt, profile2.CreatedAt, time.Second)
	require.NotEqual(t, profile1.UpdatedAt, profile2.UpdatedAt)
	require.Equal(t, profile1.VerifiedAt, profile2.VerifiedAt)
}

func TestDeleteUserProfile(t *testing.T) {
	user := createRandomUser(t)
	profile1 := createRandomUserProfile(t, user.ID)

	profile1, err := testQueries.DeleteUserProfile(context.Background(), profile1.UserID)
	require.NoError(t, err)

	profile2, err := testQueries.GetUserProfile(context.Background(), profile1.UserID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, profile2)

	t.Cleanup(func() {
		_, _ = testQueries.DeleteUserProfile(context.Background(), profile1.UserID)
		_, _ = testQueries.DeleteUser(context.Background(), user.ID)
	})
}

func TestListUserProfiles(t *testing.T) {
	user := createRandomUser(t)

	var profileSlice []UserProfile

	for i := 0; i < 10; i++ {
		profile := createRandomUserProfile(t, user.ID)
		profileSlice = append(profileSlice, profile)
	}

	arg := ListUserProfilesParams{
		Limit:  5,
		Offset: 5,
	}

	profiles, err := testQueries.ListUserProfiles(context.Background(), arg)

	t.Cleanup(func() {
		for _, profile := range profileSlice {
			_, _ = testQueries.DeleteUserProfile(context.Background(), profile.UserID)
			_, _ = testQueries.DeleteUser(context.Background(), user.ID)
		}
	})

	require.NoError(t, err)
	require.Len(t, profiles, 5)

	for _, profile := range profiles {
		require.NotEmpty(t, profile)
	}

}
