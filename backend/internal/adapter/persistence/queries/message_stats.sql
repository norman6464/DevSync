-- name: CountMessagesSentByUser :one
SELECT count(*) FROM messages
WHERE sender_id = $1;

-- name: CountMessagesReceivedByUser :one
SELECT count(*) FROM messages
WHERE receiver_id = $1;

-- name: CountConversationsByUser :one
SELECT count(DISTINCT CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END)
FROM messages
WHERE sender_id = $1 OR receiver_id = $1;

-- name: CountMessagesSentByUserSince :one
SELECT count(*) FROM messages
WHERE sender_id = $1 AND created_at >= $2;
