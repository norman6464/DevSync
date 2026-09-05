-- name: CountMentionsReceivedByUser :one
SELECT count(*) FROM mentions
WHERE user_id = $1;

-- name: CountMentionsMadeByUser :one
SELECT count(*) FROM mentions
WHERE actor_id = $1;

-- name: CountMentionsReceivedByUserSince :one
SELECT count(*) FROM mentions
WHERE user_id = $1 AND created_at >= $2;
