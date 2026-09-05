-- name: CreatePostLike :exec
INSERT INTO likes (user_id, post_id, created_at) VALUES ($1, $2, now());

-- name: IncrementPostLikeCount :exec
UPDATE posts SET like_count = like_count + 1 WHERE id = $1;

-- name: DeletePostLike :execrows
DELETE FROM likes WHERE user_id = $1 AND post_id = $2;

-- name: DecrementPostLikeCount :exec
UPDATE posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1;

-- name: CountPostLikeByUserAndPost :one
SELECT COUNT(*) FROM likes WHERE user_id = $1 AND post_id = $2;
