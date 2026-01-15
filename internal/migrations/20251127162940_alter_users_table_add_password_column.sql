-- +goose Up
ALTER TABLE users
    ADD COLUMN password BYTEA NOT NULL;

-- +goose Down
ALTER TABLE users
DROP
COLUMN password;
