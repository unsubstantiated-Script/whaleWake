package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// Store provides all functions to execute database queries and transactions.
// It embeds *Queries to allow direct access to query methods and maintains a reference to the database connection.
type Store interface {
	Querier
	CreateUserWithProfileAndRoleTx(ctx context.Context, userParams CreateUserParams, profileParams CreateUserProfileParams, roleParams CreateUserRoleParams) (UserTxResult, error)
	GetUserWithProfileAndRoleTX(ctx context.Context, userID uuid.UUID) (UserTxResult, error)
	DeleteUserWithProfileAndRoleTX(ctx context.Context, userID uuid.UUID) (UserTxResult, error)
	UpdateUserWithProfileAndRoleTX(ctx context.Context, userParams UpdateUserParams, profileParams UpdateUserProfileParams, roleParams UpdateUserRoleParams) (UserTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store instance.
// Parameters:
// - db: A pointer to the database connection.
// Returns:
// - A pointer to the initialized Store.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
// Parameters:
// - ctx: The context for the transaction.
// - fn: A function that takes *Queries and performs database operations.
// Returns:
// - An error if the transaction fails or the function returns an error.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	// If there is an error, rollback the transaction
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// UserTxResult represents the result of a transactional operation involving a user.
// Fields:
// - User: The user entity.
// - UserProfile: The associated user profile entity.
// - UserRole: The associated user role entity.
type UserTxResult struct {
	User        User        `json:"user"`
	UserProfile UserProfile `json:"user_profile"`
	UserRole    UserRole    `json:"user_role"`
}

// CreateUserWithProfileAndRoleTx performs a transaction to create a user, their profile, and role.
// Parameters:
// - ctx: The context for the transaction.
// - userParams: Parameters for creating the user.
// - profileParams: Parameters for creating the user profile.
// - roleParams: Parameters for creating the user role.
// Returns:
// - A UserTxResult containing the created user, profile, and role.
// - An error if the transaction fails.
func (store *SQLStore) CreateUserWithProfileAndRoleTx(ctx context.Context, userParams CreateUserParams, profileParams CreateUserProfileParams, roleParams CreateUserRoleParams) (UserTxResult, error) {
	var result UserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, userParams)
		if err != nil {
			return err
		}

		profileParams.UserID = result.User.ID

		result.UserProfile, err = q.CreateUserProfile(ctx, profileParams)
		if err != nil {
			return err
		}

		roleParams.UserID = result.User.ID

		result.UserRole, err = q.CreateUserRole(ctx, roleParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// GetUserWithProfileAndRoleTX retrieves a user, their profile, and role in a single transaction.
// Parameters:
// - ctx: The context for the transaction.
// - userID: The UUID of the user to retrieve.
// Returns:
// - A UserTxResult containing the user, profile, and role.
// - An error if the transaction fails or the user does not exist.
func (store *SQLStore) GetUserWithProfileAndRoleTX(ctx context.Context, userID uuid.UUID) (UserTxResult, error) {
	var result UserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.GetUser(ctx, userID)
		if err != nil {
			return err
		}

		result.UserProfile, err = q.GetUserProfile(ctx, userID)
		if err != nil {
			return err
		}

		result.UserRole, err = q.GetUserRole(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// DeleteUserWithProfileAndRoleTX deletes a user, their profile, and role in a single transaction.
// Parameters:
// - ctx: The context for the transaction.
// - userID: The UUID of the user to delete.
// Returns:
// - A UserTxResult containing the deleted user, profile, and role.
// - An error if the transaction fails or the user does not exist.
func (store *SQLStore) DeleteUserWithProfileAndRoleTX(ctx context.Context, userID uuid.UUID) (UserTxResult, error) {
	var result UserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.UserRole, err = q.DeleteUserRole(ctx, userID)
		if err != nil {
			return err
		}

		result.UserProfile, err = q.DeleteUserProfile(ctx, userID)
		if err != nil {
			return err
		}

		result.User, err = q.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

// UpdateUserWithProfileAndRoleTX updates a user, their profile, and role in a single transaction.
// Parameters:
// - ctx: The context for the transaction.
// - userParams: Parameters for updating the user.
// - profileParams: Parameters for updating the user profile.
// - roleParams: Parameters for updating the user role.
// Returns:
// - A UserTxResult containing the updated user, profile, and role.
// - An error if the transaction fails or the user does not exist.
func (store *SQLStore) UpdateUserWithProfileAndRoleTX(ctx context.Context, userParams UpdateUserParams, profileParams UpdateUserProfileParams, roleParams UpdateUserRoleParams) (UserTxResult, error) {
	var result UserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.UpdateUser(ctx, userParams)
		if err != nil {
			return err
		}

		profileParams.UserID = userParams.ID

		result.UserProfile, err = q.UpdateUserProfile(ctx, profileParams)
		if err != nil {
			return err
		}

		roleParams.UserID = userParams.ID

		result.UserRole, err = q.UpdateUserRole(ctx, roleParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
