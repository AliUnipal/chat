-- +goose Up
CREATE TABLE chats
(
    id           UUID PRIMARY KEY,
    current_user_id UUID NOT NULL REFERENCES users (id),
    other_user_id   UUID NOT NULL REFERENCES users (id),
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP
);