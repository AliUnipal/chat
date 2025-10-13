-- TODO: I get this big meow error, ask if the solution below is the correct one.
--  Or I should go with the three pattern thingies
-- # package
-- sql\queries.sql:7:1: relation "users" does not exist
-- sql\queries.sql:36:1: relation "users" does not exist
-- chatsvc\chatrepos\pqchatrepo\repo.go:12: running "sqlc": exit status 1

-- +goose Up
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    image_url  TEXT,
    first_name VARCHAR(30) NOT NULL,
    last_name  VARCHAR(30),
    -- TODO: Decide later on the username to be either text or phone
    username   TEXT        NOT NULL,
    created_at TIMESTAMP   NOT NULL,
    modified_at TIMESTAMP   NOT NULL
);

CREATE TABLE chats
(
    id           UUID PRIMARY KEY,
    user_one_id UUID NOT NULL REFERENCES users (id),
    user_two_id   UUID NOT NULL REFERENCES users (id),
    created_at   TIMESTAMP NOT NULL,
    modified_at   TIMESTAMP NOT NULL
);
