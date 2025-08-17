-- name: GetUser :one
SELECT * FROM users
WHERE public_key = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, public_key) VALUES (?, ?)
RETURNING *;
