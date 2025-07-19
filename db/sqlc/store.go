package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// UserTxResult CreateUserTxResult is the result of transfer transaction.
type UserTxResult struct {
	User        User        `json:"user"`
	UserProfile UserProfile `json:"user_profile"`
	UserRole    UserRole    `json:"user_role"`
}

// CreateUserWithProfileAndRoleTx performs a transaction to create a new user entry along with their profile and role information all in one go.
func (store *Store) CreateUserWithProfileAndRoleTx(ctx context.Context, userParams CreateUserParams, profileParams CreateUserProfileParams, roleParams CreateUserRoleParams) (UserTxResult, error) {
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

func (store *Store) GetUserWithProfileAndRoleTX(ctx context.Context, userID uuid.UUID) (UserTxResult, error) {
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

func (store *Store) DeleteUserWithProfileAndRoleTX(ctx context.Context, userID uuid.UUID) (UserTxResult, error) {
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

//TODO: Make an UpdateUserProfileRoleTx -> and tests?
