-- name: CreateFollow :exec
INSERT INTO follows (follower_id, followee_id, created_at)
VALUES ($1, $2, now());

-- name: DeleteFollow :exec
DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2;

-- name: ListFollowers :many
-- 件数は follow_stats.sql の CountFollowersByUser を再利用する。
SELECT u.* FROM users u
JOIN follows f ON f.follower_id = u.id
WHERE f.followee_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListFollowing :many
-- 件数は follow_stats.sql の CountFollowingByUser を再利用する。
SELECT u.* FROM users u
JOIN follows f ON f.followee_id = u.id
WHERE f.follower_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;
