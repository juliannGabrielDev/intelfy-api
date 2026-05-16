-- name: CreateUser :one
INSERT INTO users (
	id, username, email, password_hash, role
) VALUES (
	$1, $2, $3, $4, COALESCE(NULLIF(sqlc.arg(role)::varchar, ''), 'user')
)
RETURNING id, username, email, role, created_at;

-- name: GetUsers :many
SELECT id, username, email, role, created_at FROM users
ORDER BY username ASC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: GetUserByEmail :one
SELECT id, username, email, role, created_at, password_hash FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, email, role, created_at, password_hash FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, password_hash = $4, role = $5
WHERE id = $1
RETURNING id, username, email, role, created_at;