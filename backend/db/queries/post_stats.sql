-- name: CountPublishedPostsByUser :one
SELECT COUNT(*) FROM posts WHERE user_id = $1 AND is_draft = false;

-- name: CountDraftPostsByUser :one
SELECT COUNT(*) FROM posts WHERE user_id = $1 AND is_draft = true;

-- name: SumPostCommentsReceivedByUser :one
SELECT COALESCE(SUM(comment_count), 0)::bigint FROM posts WHERE user_id = $1;

-- name: CountPostsByUserSince :one
SELECT COUNT(*) FROM posts WHERE user_id = $1 AND created_at >= $2;
