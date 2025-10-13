-- +goose Up
ALTER TABLE users
    ADD CONSTRAINT users_username_key UNIQUE (username);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT users_username_key;
