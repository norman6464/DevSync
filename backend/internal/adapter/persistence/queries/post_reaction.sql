-- name: CreateReaction :exec
INSERT INTO post_reactions (user_id, post_id, kind, value, created_at) VALUES ($1, $2, 'emoji', $3, now());

-- name: DeleteReaction :exec
DELETE FROM post_reactions
WHERE post_reactions.user_id = $1 AND post_reactions.post_id = $2
    AND post_reactions.kind = 'emoji' AND post_reactions.value = $3;

-- name: ListReactionCountsByPost :many
SELECT value AS emoji, COUNT(*)::bigint AS count FROM post_reactions
WHERE post_id = $1 AND kind = 'emoji'
GROUP BY value
ORDER BY count DESC;

-- name: ListUserReactionEmojisByPost :many
SELECT value AS emoji FROM post_reactions WHERE user_id = $1 AND post_id = $2 AND kind = 'emoji';

-- name: ListReactionCountsByPosts :many
SELECT post_id, value AS emoji, COUNT(*)::bigint AS count FROM post_reactions
WHERE post_id = ANY($1::bigint[]) AND kind = 'emoji'
GROUP BY post_id, value
ORDER BY post_id, count DESC;

-- name: ListUserReactionsByPosts :many
SELECT post_id, value AS emoji FROM post_reactions
WHERE user_id = $1 AND post_id = ANY($2::bigint[]) AND kind = 'emoji';

-- name: GetPostAuthorID :one
SELECT user_id FROM posts WHERE id = $1;
