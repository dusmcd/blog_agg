-- +goose Up
CREATE TABLE feeds_users (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    feed_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    UNIQUE(feed_id, user_id)
);

-- +goose Down
DROP TABLE feeds_users;