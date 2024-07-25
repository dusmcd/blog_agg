-- +goose Up
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    created_at    TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;