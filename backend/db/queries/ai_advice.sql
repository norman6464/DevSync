-- name: CreateAIAdvice :one
INSERT INTO ai_advices (user_id, type, priority, title_key, message_key, params, action_url, is_read, expires_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, false, $8, now(), now())
RETURNING *;

-- name: ListAIAdvicesByUser :many
SELECT * FROM ai_advices
WHERE user_id = $1
ORDER BY priority ASC, created_at DESC
LIMIT $2;

-- name: ListUnreadAIAdvicesByUser :many
SELECT * FROM ai_advices
WHERE user_id = $1 AND is_read = false
ORDER BY priority ASC, created_at DESC;

-- name: MarkAIAdviceAsRead :execrows
UPDATE ai_advices SET is_read = true, updated_at = now()
WHERE id = $1 AND user_id = $2;

-- name: DeleteAIAdvicesByUser :exec
DELETE FROM ai_advices WHERE user_id = $1;
