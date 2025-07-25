// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user_profile.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUserProfile = `-- name: CreateUserProfile :one
INSERT INTO user_profile (user_id,
                          first_name,
                          last_name,
                          business_name,
                          street_address,
                          city,
                          state,
                          zip,
                          country_code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, user_id, first_name, last_name, business_name, street_address, city, state, zip, country_code, created_at, updated_at, verified_at
`

type CreateUserProfileParams struct {
	UserID        uuid.UUID `json:"user_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BusinessName  string    `json:"business_name"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Zip           string    `json:"zip"`
	CountryCode   string    `json:"country_code"`
}

func (q *Queries) CreateUserProfile(ctx context.Context, arg CreateUserProfileParams) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, createUserProfile,
		arg.UserID,
		arg.FirstName,
		arg.LastName,
		arg.BusinessName,
		arg.StreetAddress,
		arg.City,
		arg.State,
		arg.Zip,
		arg.CountryCode,
	)
	var i UserProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.BusinessName,
		&i.StreetAddress,
		&i.City,
		&i.State,
		&i.Zip,
		&i.CountryCode,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VerifiedAt,
	)
	return i, err
}

const deleteUserProfile = `-- name: DeleteUserProfile :one
DELETE
FROM user_profile
WHERE user_id = $1 RETURNING id, user_id, first_name, last_name, business_name, street_address, city, state, zip, country_code, created_at, updated_at, verified_at
`

func (q *Queries) DeleteUserProfile(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, deleteUserProfile, userID)
	var i UserProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.BusinessName,
		&i.StreetAddress,
		&i.City,
		&i.State,
		&i.Zip,
		&i.CountryCode,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VerifiedAt,
	)
	return i, err
}

const getUserProfile = `-- name: GetUserProfile :one
SELECT id, user_id, first_name, last_name, business_name, street_address, city, state, zip, country_code, created_at, updated_at, verified_at
FROM user_profile
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUserProfile(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, getUserProfile, userID)
	var i UserProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.BusinessName,
		&i.StreetAddress,
		&i.City,
		&i.State,
		&i.Zip,
		&i.CountryCode,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VerifiedAt,
	)
	return i, err
}

const listUserProfiles = `-- name: ListUserProfiles :many
SELECT id, user_id, first_name, last_name, business_name, street_address, city, state, zip, country_code, created_at, updated_at, verified_at
FROM user_profile
ORDER BY id LIMIT $1
OFFSET $2
`

type ListUserProfilesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUserProfiles(ctx context.Context, arg ListUserProfilesParams) ([]UserProfile, error) {
	rows, err := q.db.QueryContext(ctx, listUserProfiles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserProfile{}
	for rows.Next() {
		var i UserProfile
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.BusinessName,
			&i.StreetAddress,
			&i.City,
			&i.State,
			&i.Zip,
			&i.CountryCode,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.VerifiedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserProfile = `-- name: UpdateUserProfile :one
UPDATE user_profile
SET first_name = $2,
    last_name = $3,
    business_name = $4,
    street_address = $5,
    city = $6,
    state = $7,
    zip = $8,
    country_code = $9,
    updated_at = STATEMENT_TIMESTAMP()
WHERE user_id = $1 RETURNING id, user_id, first_name, last_name, business_name, street_address, city, state, zip, country_code, created_at, updated_at, verified_at
`

type UpdateUserProfileParams struct {
	UserID        uuid.UUID `json:"user_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BusinessName  string    `json:"business_name"`
	StreetAddress string    `json:"street_address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Zip           string    `json:"zip"`
	CountryCode   string    `json:"country_code"`
}

func (q *Queries) UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (UserProfile, error) {
	row := q.db.QueryRowContext(ctx, updateUserProfile,
		arg.UserID,
		arg.FirstName,
		arg.LastName,
		arg.BusinessName,
		arg.StreetAddress,
		arg.City,
		arg.State,
		arg.Zip,
		arg.CountryCode,
	)
	var i UserProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.BusinessName,
		&i.StreetAddress,
		&i.City,
		&i.State,
		&i.Zip,
		&i.CountryCode,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VerifiedAt,
	)
	return i, err
}
