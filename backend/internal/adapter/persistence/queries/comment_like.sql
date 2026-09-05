-- name: CreateCommentLike :exec
INSERT INTO comment_likes (user_id, comment_id, created_at)
VALUES ($1, $2, now());

-- name: IncrementCommentLikeCount :exec
UPDATE comments
SET like_count = like_count + 1
WHERE id = $1;

-- name: DeleteCommentLike :exec
DELETE FROM comment_likes
WHERE user_id = $1 AND comment_id = $2;

-- name: DecrementCommentLikeCount :exec
UPDATE comments
SET like_count = like_count - 1
WHERE id = $1 AND like_count > 0;

-- name: CountCommentLikesByUserAndComment :one
SELECT count(*) FROM comment_likes
WHERE user_id = $1 AND comment_id = $2;

-- name: CountCommentLikesByComment :one
SELECT count(*) FROM comment_likes
WHERE comment_id = $1;
