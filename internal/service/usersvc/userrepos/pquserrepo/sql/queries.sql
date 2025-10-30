-- name: CreateUser :exec
INSERT INTO users (id, image_url, first_name, last_name, username, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;
