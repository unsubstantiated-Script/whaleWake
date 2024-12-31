-- name: CreateUser :one
INSERT INTO users (user_name, email, password)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY id LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET user_name  = $2,
    email      = $3,
    password   = $4,
    updated_at = statment_timestamp()
WHERE id = $1 RETURNING *;

-- name: DeleteUser :one
DELETE
FROM users
WHERE id = $1 RETURNING *;