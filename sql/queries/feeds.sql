-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) VALUES
($1, $2, $3, $4, $5, (SELECT id FROM users WHERE apikey = $6))
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;