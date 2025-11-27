-- +goose Up
CREATE TABLE messages
(
    id          UUID PRIMARY KEY,
    sender_id   UUID      NOT NULL REFERENCES users (id),
    chat_id     UUID      NOT NULL REFERENCES chats (id) ON DELETE CASCADE,
    content     BYTEA     NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    modified_At TIMESTAMP
);
