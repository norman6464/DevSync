package dto

import "github.com/norman6464/devsync/backend/internal/model"

// ChatRoomListResponse はチャットルーム一覧レスポンス（ページネーション付き）。
type ChatRoomListResponse struct {
	Rooms  []model.ChatRoom `json:"rooms"`
	Total  int64            `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

// CreateChatRoomRequest はチャットルーム作成リクエスト。
type CreateChatRoomRequest struct {
	Name        string `json:"name" binding:"required,max=200" validate:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	MemberIDs   []uint `json:"member_ids" binding:"omitempty,max=100"`
}

// UpdateChatRoomRequest はチャットルーム更新リクエスト。
type UpdateChatRoomRequest struct {
	Name        string `json:"name" binding:"omitempty,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
}

// AddChatRoomMemberRequest はメンバー追加リクエスト。
type AddChatRoomMemberRequest struct {
	UserID uint `json:"user_id" binding:"required" validate:"required"`
}

// SendChatRoomMessageRequest はチャットルームメッセージ送信リクエスト。
type SendChatRoomMessageRequest struct {
	Content string `json:"content" binding:"required,max=5000" validate:"required,max=5000"`
}
