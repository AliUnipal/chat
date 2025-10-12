-- name: CreateChat :exec
INSERT INTO chats (id, current_user_id, other_user_id, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, $5),
    ($6, $3, $2, $4, $5);

-- name: GetChatsByUser :many
SELECT *
FROM chats
WHERE current_user_id = $1;

-- name: GetChat :one
SELECT *
FROM chats
WHERE id = $1;
