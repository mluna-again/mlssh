-- name: GetUser :one
SELECT * FROM users
WHERE public_key = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, public_key) VALUES (?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(?, name), next_activity_change_at = COALESCE(?, next_activity_change_at)
RETURNING *;

