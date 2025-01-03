// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID         pgtype.UUID        `json:"id"`
	UserName   string             `json:"user_name"`
	Email      string             `json:"email"`
	Password   string             `json:"password"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt  pgtype.Timestamptz `json:"updated_at"`
	VerifiedAt pgtype.Timestamptz `json:"verified_at"`
}

type UserProfile struct {
	ID            pgtype.UUID        `json:"id"`
	UserID        pgtype.UUID        `json:"user_id"`
	FirstName     pgtype.Text        `json:"first_name"`
	LastName      pgtype.Text        `json:"last_name"`
	BusinessName  pgtype.Text        `json:"business_name"`
	StreetAddress pgtype.Text        `json:"street_address"`
	City          pgtype.Text        `json:"city"`
	State         pgtype.Text        `json:"state"`
	Zip           pgtype.Text        `json:"zip"`
	CountryCode   pgtype.Int4        `json:"country_code"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
	VerifiedAt    pgtype.Timestamptz `json:"verified_at"`
}

type UserRole struct {
	ID         pgtype.UUID        `json:"id"`
	UserID     pgtype.UUID        `json:"user_id"`
	RoleID     int32              `json:"role_id"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	UpdatedAt  pgtype.Timestamptz `json:"updated_at"`
	VerifiedAt pgtype.Timestamptz `json:"verified_at"`
}
