-- name: GetCommentByID :one
SELECT * FROM comments WHERE id = $1;
