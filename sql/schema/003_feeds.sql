-- +goose Up
CREATE TABLE feeds (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id VARCHAR(255) NOT NULL REFERENCES users ON DELETE CASCADE,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE feeds;