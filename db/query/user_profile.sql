-- name: CreateUserProfile :one
INSERT INTO user_profile (user_id,
                          first_name,
                          last_name,
                          business_name,
                          street_address,
                          city,
                          state,
                          zip,
                          country_code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: GetUserProfile :one
SELECT *
FROM user_profile
WHERE user_id = $1 LIMIT 1;

-- name: ListUserProfiles :many
SELECT *
FROM user_profile
ORDER BY id LIMIT $1
OFFSET $2;

-- name: UpdateUserProfile :one
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
WHERE user_id = $1 RETURNING *;

-- name: DeleteUserProfile :one
DELETE
FROM user_profile
WHERE user_id = $1 RETURNING *;