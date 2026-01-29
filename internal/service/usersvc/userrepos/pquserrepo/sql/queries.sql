-- name: CreateUser :exec
INSERT INTO users (id, image_url, first_name, last_name, username, password, created_at)
VALUES ($1, $2, $3, $4, $5, $6::bytea, $7);

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;
