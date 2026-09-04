-- name: CountReactionsReceivedByUser :one
SELECT count(*) FROM reactions
JOIN posts ON posts.id = reactions.post_id
WHERE posts.user_id = $1;

-- name: CountUniqueReactorsByUser :one
SELECT count(DISTINCT reactions.user_id) FROM reactions
JOIN posts ON posts.id = reactions.post_id
WHERE posts.user_id = $1;

-- name: CountReactionsReceivedByUserSince :one
SELECT count(*) FROM reactions
JOIN posts ON posts.id = reactions.post_id
WHERE posts.user_id = $1 AND reactions.created_at >= $2;

-- name: GetEmojiBreakdownByUser :many
SELECT reactions.emoji AS emoji, count(*) AS count
FROM reactions
JOIN posts ON posts.id = reactions.post_id
WHERE posts.user_id = $1
GROUP BY reactions.emoji
ORDER BY count DESC, reactions.emoji ASC;

-- name: GetTopReactedPostsByUser :many
SELECT posts.id AS id, posts.title AS title, count(*) AS reaction_count
FROM reactions
JOIN posts ON posts.id = reactions.post_id
WHERE posts.user_id = $1
GROUP BY posts.id, posts.title
ORDER BY reaction_count DESC, posts.id DESC
LIMIT $2;
