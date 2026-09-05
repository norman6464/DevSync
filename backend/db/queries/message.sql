-- name: CreateMessage :one
INSERT INTO messages (sender_id, receiver_id, content, read, created_at)
VALUES ($1, $2, $3, false, now())
RETURNING *;

-- name: ListConversationMessages :many
-- GORMのPreload("Sender").Preload("Receiver")に相当。
-- sender_id/receiver_idともにNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(messages), sqlc.embed(sender), sqlc.embed(receiver)
FROM messages
JOIN users sender ON sender.id = messages.sender_id
JOIN users receiver ON receiver.id = messages.receiver_id
WHERE (messages.sender_id = $1 AND messages.receiver_id = $2)
    OR (messages.sender_id = $2 AND messages.receiver_id = $1)
ORDER BY messages.created_at ASC
LIMIT $3 OFFSET $4;

-- name: ListConversationSummaries :many
-- 会話相手ごとの最新メッセージと未読件数を取得する（移行前のGORM Raw SQLをそのまま踏襲）。
SELECT DISTINCT ON (other_id) other_id AS user_id, u.name, u.avatar_url,
    m.content AS last_message, m.created_at::text AS last_time,
    (SELECT COUNT(*) FROM messages unread WHERE unread.sender_id = other_id AND unread.receiver_id = $1 AND unread.read = false) AS unread_count
FROM (
    SELECT CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END AS other_id, id
    FROM messages
    WHERE sender_id = $1 OR receiver_id = $1
) sub
JOIN messages m ON m.id = sub.id
JOIN users u ON u.id = sub.other_id
ORDER BY other_id, m.created_at DESC;

-- name: MarkMessagesAsRead :exec
UPDATE messages SET read = true
WHERE sender_id = $1 AND receiver_id = $2 AND read = false;
