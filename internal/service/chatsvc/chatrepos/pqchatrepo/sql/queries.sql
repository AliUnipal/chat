-- name: CreateChat :exec
INSERT INTO chats (id, user_one_id, user_two_id, created_at, modified_at)
VALUES
    ($1, $2, $3, $4, $5);

-- name: GetChatsByUser :many
SELECT
    c.id AS chat_id,
    c.created_at AS chat_created_at,
    c.modified_at AS chat_modified_at,

    u1.id AS user_one_id,
    u1.first_name AS user_one_first_name,
    u1.last_name AS user_one_last_name,
    u1.username AS user_one_username,
    u1.image_url AS user_one_image_url,
    u1.created_at AS user_one_created_at,
    u1.modified_at AS user_one_modified_At,


    u2.id AS user_two_id,
    u2.first_name AS user_two_first_name,
    u2.last_name AS user_two_last_name,
    u2.username AS user_two_username,
    u2.image_url AS user_two_image_url,
    u2.created_at AS user_two_created_at,
    u2.modified_at AS user_two_modified_At
FROM chats AS c
JOIN users AS u1 ON c.user_one_id = u1.id
JOIN users AS u2 ON c.user_two_id = u2.id
WHERE user_one_id = $1 OR user_two_id = $1;

-- name: GetChat :one
SELECT
    c.id AS chat_id,
    c.created_at AS chat_created_at,
    c.modified_at AS chat_modified_at,

    u1.id AS user_one_id,
    u1.first_name AS user_one_first_name,
    u1.last_name AS user_one_last_name,
    u1.username AS user_one_username,
    u1.image_url AS user_one_image_url,
    u1.created_at AS user_one_created_at,
    u1.modified_at AS user_one_modified_At,


    u2.id AS user_two_id,
    u2.first_name AS user_two_first_name,
    u2.last_name AS user_two_last_name,
    u2.username AS user_two_username,
    u2.image_url AS user_two_image_url,
    u2.created_at AS user_two_created_at,
    u2.modified_at AS user_two_modified_At
FROM chats AS c
         JOIN users AS u1 ON c.user_one_id = u1.id
         JOIN users AS u2 ON c.user_two_id = u2.id
WHERE c.id = $1;
