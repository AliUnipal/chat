-- name: CreateMessage :exec
INSERT INTO messages (id, sender_id, chat_id, content, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetMessages :many
SELECT *
FROM messages
WHERE chat_id = $1;
