-- name: CreateChatRoom :one
INSERT INTO chat_rooms (name, description, owner_id, created_at, updated_at)
VALUES ($1, $2, $3, now(), now())
RETURNING *;

-- name: GetChatRoomWithOwnerByID :one
-- GORMのPreload("Owner")に相当。owner_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(chat_rooms), sqlc.embed(users)
FROM chat_rooms
JOIN users ON users.id = chat_rooms.owner_id
WHERE chat_rooms.id = $1;

-- name: ListChatRoomsByUser :many
-- GORMのPreload("Owner")に相当。owner_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(chat_rooms), sqlc.embed(users)
FROM chat_rooms
JOIN chat_room_members ON chat_room_members.chat_room_id = chat_rooms.id
JOIN users ON users.id = chat_rooms.owner_id
WHERE chat_room_members.user_id = $1
ORDER BY chat_rooms.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CountChatRoomsByUser :one
-- 件数は ListChatRoomsByUser の総数取得と CountByUserID(単体メソッド) の両方から利用する。
SELECT COUNT(*) FROM chat_room_members WHERE user_id = $1;

-- name: UpdateChatRoom :one
-- GORMのSave（全カラム上書き）に相当。
UPDATE chat_rooms SET name = $2, description = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteGroupMessagesByRoom :exec
DELETE FROM group_messages WHERE chat_room_id = $1;

-- name: DeleteChatRoomMembersByRoom :exec
DELETE FROM chat_room_members WHERE chat_room_id = $1;

-- name: DeleteChatRoom :exec
DELETE FROM chat_rooms WHERE id = $1;

-- name: CreateChatRoomMember :exec
INSERT INTO chat_room_members (chat_room_id, user_id, joined_at)
VALUES ($1, $2, $3);

-- name: DeleteChatRoomMember :exec
DELETE FROM chat_room_members WHERE chat_room_id = $1 AND user_id = $2;

-- name: ListChatRoomMembersWithUser :many
-- GORMのPreload("User")に相当。user_idはNOT NULLのためINNER JOINでよい。
SELECT sqlc.embed(chat_room_members), sqlc.embed(users)
FROM chat_room_members
JOIN users ON users.id = chat_room_members.user_id
WHERE chat_room_members.chat_room_id = $1;

-- name: CountChatRoomMembership :one
SELECT COUNT(*) FROM chat_room_members WHERE chat_room_id = $1 AND user_id = $2;
