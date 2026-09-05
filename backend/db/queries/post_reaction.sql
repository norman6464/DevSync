-- name: CreateReaction :exec
INSERT INTO reactions (user_id, post_id, emoji, created_at) VALUES ($1, $2, $3, now());

-- name: DeleteReaction :exec
DELETE FROM reactions WHERE user_id = $1 AND post_id = $2 AND emoji = $3;

-- name: ListReactionCountsByPost :many
SELECT emoji, COUNT(*)::bigint AS count FROM reactions
WHERE post_id = $1
GROUP BY emoji
ORDER BY count DESC;

-- name: ListUserReactionEmojisByPost :many
SELECT emoji FROM reactions WHERE user_id = $1 AND post_id = $2;

-- name: ListReactionCountsByPosts :many
SELECT post_id, emoji, COUNT(*)::bigint AS count FROM reactions
WHERE post_id = ANY($1::bigint[])
GROUP BY post_id, emoji
ORDER BY post_id, count DESC;

-- name: ListUserReactionsByPosts :many
SELECT post_id, emoji FROM reactions
WHERE user_id = $1 AND post_id = ANY($2::bigint[]);

-- name: GetPostAuthorID :one
SELECT user_id FROM posts WHERE id = $1;
