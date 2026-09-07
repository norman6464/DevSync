-- name: CountReactionsReceivedByUser :one
SELECT count(*) FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'emoji';

-- name: CountUniqueReactorsByUser :one
SELECT count(DISTINCT post_reactions.user_id) FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'emoji';

-- name: CountReactionsReceivedByUserSince :one
SELECT count(*) FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'emoji' AND post_reactions.created_at >= $2;

-- name: GetEmojiBreakdownByUser :many
SELECT post_reactions.value AS emoji, count(*) AS count
FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'emoji'
GROUP BY post_reactions.value
ORDER BY count DESC, post_reactions.value ASC;

-- name: GetTopReactedPostsByUser :many
SELECT posts.id AS id, posts.title AS title, count(*) AS reaction_count
FROM post_reactions
JOIN posts ON posts.id = post_reactions.post_id
WHERE posts.user_id = $1 AND post_reactions.kind = 'emoji'
GROUP BY posts.id, posts.title
ORDER BY reaction_count DESC, posts.id DESC
LIMIT $2;
