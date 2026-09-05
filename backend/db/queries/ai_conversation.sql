-- name: CreateAIConversation :one
INSERT INTO ai_conversations (user_id, title, created_at, updated_at)
VALUES ($1, $2, now(), now())
RETURNING *;

-- name: ListAIConversationsByUser :many
SELECT * FROM ai_conversations
WHERE user_id = $1
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAIConversationByID :one
SELECT * FROM ai_conversations WHERE id = $1;

-- name: GetAIConversationByIDAndUser :one
SELECT * FROM ai_conversations WHERE id = $1 AND user_id = $2;

-- name: ListAIMessagesByConversationIDs :many
-- GORMのPreload("Messages")（順序未指定）に相当。まとめ取得用でid昇順（挿入順相当）とする。
SELECT * FROM ai_messages WHERE conversation_id = ANY($1::bigint[]) ORDER BY id ASC;

-- name: ListAIMessagesByConversationIDOrderedByCreatedAt :many
-- GORMのPreload("Messages", Order("created_at ASC"))に相当。
SELECT * FROM ai_messages WHERE conversation_id = $1 ORDER BY created_at ASC;

-- name: CreateAIMessage :one
INSERT INTO ai_messages (conversation_id, role, content, tokens_used, created_at)
VALUES ($1, $2, $3, $4, now())
RETURNING *;

-- name: TouchAIConversation :exec
UPDATE ai_conversations SET updated_at = $2 WHERE id = $1;

-- name: CountTodayAIMessagesByUser :one
SELECT COUNT(*) FROM ai_messages
JOIN ai_conversations ON ai_conversations.id = ai_messages.conversation_id
WHERE ai_conversations.user_id = $1 AND ai_messages.role = $2 AND ai_messages.created_at >= $3;

-- name: DeleteAIMessagesByConversationID :exec
DELETE FROM ai_messages WHERE conversation_id = $1;

-- name: DeleteAIConversation :exec
DELETE FROM ai_conversations WHERE id = $1;
