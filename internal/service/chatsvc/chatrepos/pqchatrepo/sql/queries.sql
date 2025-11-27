-- name: CreateChat :exec
INSERT INTO chats (id, user_one_id, user_two_id, created_at)
VALUES
    ($1, $2, $3, $4);

-- name: GetChatsByUser :many
SELECT
    sqlc.embed(c), sqlc.embed(u1), sqlc.embed(u2)
FROM chats AS c
JOIN users AS u1 ON c.user_one_id = u1.id
JOIN users AS u2 ON c.user_two_id = u2.id
WHERE user_one_id = $1 OR user_two_id = $1;

-- name: GetChat :one
SELECT
    sqlc.embed(c), sqlc.embed(u1), sqlc.embed(u2)
FROM chats AS c
         JOIN users AS u1 ON c.user_one_id = u1.id
         JOIN users AS u2 ON c.user_two_id = u2.id
WHERE c.id = $1;
