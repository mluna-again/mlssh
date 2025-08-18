-- name: GetUser :one
SELECT *, settings.inserted_at != NULL AS active FROM users
LEFT JOIN settings ON settings.user_pk = users.public_key
WHERE public_key = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, public_key) VALUES (?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(?, name), next_activity_change_at = COALESCE(?, next_activity_change_at)
RETURNING *;

