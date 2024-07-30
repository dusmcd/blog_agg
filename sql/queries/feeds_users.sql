-- name: FollowFeed :one
INSERT INTO feeds_users (id, created_at, updated_at, feed_id, user_id)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetFeedsFollowed :many
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id FROM feeds
INNER JOIN feeds_users ON feeds_users.feed_id = feeds.id
WHERE feeds_users.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feeds_users
WHERE id = $1;