-- name: CountFollowersByUser :one
SELECT count(*) FROM follows
WHERE followee_id = $1;

-- name: CountFollowingByUser :one
SELECT count(*) FROM follows
WHERE follower_id = $1;
