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