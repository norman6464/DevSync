package dto

// CreateChatRoomRequest はチャットルーム作成リクエスト。
type CreateChatRoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MemberIDs   []uint `json:"member_ids"`
}

// UpdateChatRoomRequest はチャットルーム更新リクエスト。
type UpdateChatRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AddChatRoomMemberRequest はメンバー追加リクエスト。
type AddChatRoomMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// SendChatRoomMessageRequest はチャットルームメッセージ送信リクエスト。
type SendChatRoomMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
