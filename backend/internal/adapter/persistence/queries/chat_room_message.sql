-- name: CreateGroupMessage :one
INSERT INTO group_messages (chat_room_id, sender_id, content, created_at)
VALUES ($1, $2, $3, now())
RETURNING *;

-- name: ListGroupMessagesByRoom :many
-- GORMのPreload("Sender")に相当。sender_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(group_messages), sqlc.embed(users)
FROM group_messages
JOIN users ON users.id = group_messages.sender_id
WHERE group_messages.chat_room_id = $1
ORDER BY group_messages.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetUserByIDForChatSender :one
SELECT * FROM users WHERE id = $1;
