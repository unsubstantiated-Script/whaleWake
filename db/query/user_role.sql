-- name: CreateUserRole :one
INSERT INTO user_role (user_id, role_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserRole :one
SELECT *
FROM user_role
WHERE id = $1
LIMIT 1;

-- name: ListUserRoles :many
SELECT *
FROM user_role
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUserRole :one
UPDATE user_role
SET role_id    = $2,
    updated_at = statment_timestamp()
WHERE id = $1
RETURNING *;

-- name: DeleteUserRole :one
DELETE
FROM user_role
WHERE id = $1
RETURNING *;