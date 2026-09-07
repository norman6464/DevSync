-- name: CountPostViewsReceivedByUser :one
SELECT COUNT(*) FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'view';
