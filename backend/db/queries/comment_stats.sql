-- name: CountTopLevelCommentsByUser :one
SELECT count(*) FROM comments
WHERE user_id = $1 AND parent_id IS NULL;

-- name: CountRepliesByUser :one
SELECT count(*) FROM comments
WHERE user_id = $1 AND parent_id IS NOT NULL;

-- name: CountCommentsReceivedByUser :one
SELECT count(*) FROM comments
JOIN posts ON posts.id = comments.post_id
WHERE posts.user_id = $1 AND comments.user_id != $1;

-- name: CountCommentsByUserSince :one
SELECT count(*) FROM comments
WHERE user_id = $1 AND created_at >= $2;
