-- name: ListChatRoomMemberUserIDs :many
SELECT user_id FROM chat_room_members WHERE chat_room_id = $1;
