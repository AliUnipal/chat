-- +goose Up
ALTER TABLE chats
    ADD CONSTRAINT chats_unique_users_pair UNIQUE (user_one_id, user_two_id);

-- +goose Down
ALTER TABLE chats
DROP
CONSTRAINT chats_unique_users_pair;
