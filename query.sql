-- name: GetUser :one
SELECT * FROM users
LEFT JOIN settings ON settings.user_pk = users.public_key
WHERE public_key = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, public_key) VALUES (?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE(?, name), next_activity_change_at = COALESCE(?, next_activity_change_at)
WHERE public_key = ?
RETURNING *;

-- name: CreateSettings :one
INSERT INTO settings (
  user_pk,
  pet_species,
  pet_color,
  pet_name,
  inserted_at
) VALUES (?, ?, ?, ?, UNIXEPOCH())
RETURNING *;

-- name: GetSettings :one
SELECT * FROM settings
WHERE user_pk = ?;
