-- name: CountPostViewsReceivedByUser :one
SELECT COUNT(*) FROM post_views
JOIN posts ON posts.id = post_views.post_id
WHERE posts.user_id = $1;
