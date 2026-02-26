package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ChatRoomServiceInterface はChatRoomHandlerが依存するサービスのインターフェース。
type ChatRoomServiceInterface interface {
	Create(room *model.ChatRoom, memberIDs []uint) (*model.ChatRoom, error)
	GetByUserID(userID uint, limit, offset int) ([]model.ChatRoom, int64, error)
	GetByID(roomID, userID uint) (*model.ChatRoom, error)
	Update(roomID, userID uint, name, description string) (*model.ChatRoom, error)
	Delete(roomID, userID uint) error
	GetMembers(roomID, userID uint) ([]model.ChatRoomMember, error)
	AddMember(roomID, userID, targetUserID uint) error
	RemoveMember(roomID, userID, targetUserID uint) error
	GetMessages(roomID, userID uint, page, limit int) ([]model.GroupMessage, error)
	SendMessage(roomID, userID uint, content string) (*model.GroupMessage, error)
	CountByUserID(userID uint) (int64, error)
}

// ChatRoomHandler はチャットルーム関連のHTTPハンドラ。
// チャットルームのCRUD・メンバー管理・メッセージ送受信を処理する。
type ChatRoomHandler struct {
	service ChatRoomServiceInterface
}

// NewChatRoomHandler は新しいChatRoomHandlerインスタンスを生成する。
func NewChatRoomHandler(s ChatRoomServiceInterface) *ChatRoomHandler {
	return &ChatRoomHandler{service: s}
}

// Create は新しいチャットルームを作成する。
func (h *ChatRoomHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateChatRoomRequest](c)
	if input == nil {
		return
	}

	room := &model.ChatRoom{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
	}

	created, err := h.service.Create(room, input.MemberIDs)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, created)
}

// GetMyRooms は現在のユーザーが参加しているチャットルーム一覧を取得する。
func (h *ChatRoomHandler) GetMyRooms(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)
	rooms, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.ChatRoomListResponse{
		Rooms:  ensureSlice(rooms),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetByID は指定IDのチャットルーム詳細を取得する。
func (h *ChatRoomHandler) GetByID(c *gin.Context) {
	handleGetByID(c, h.service.GetByID)
}

// Update は指定IDのチャットルーム情報を更新する。
func (h *ChatRoomHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.UpdateChatRoomRequest](c)
	if input == nil {
		return
	}

	room, err := h.service.Update(roomID, userID, input.Name, input.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, room)
}

// Delete は指定IDのチャットルームを削除する。
func (h *ChatRoomHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

// GetMembers は指定チャットルームのメンバー一覧を取得する。
func (h *ChatRoomHandler) GetMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}

	members, err := h.service.GetMembers(roomID, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, members)
}

// AddMember はチャットルームに新しいメンバーを追加する。
func (h *ChatRoomHandler) AddMember(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.AddChatRoomMemberRequest](c)
	if input == nil {
		return
	}

	if err := h.service.AddMember(roomID, userID, input.UserID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("member added"))
}

// RemoveMember はチャットルームからメンバーを削除する。
func (h *ChatRoomHandler) RemoveMember(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}
	targetID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	if err := h.service.RemoveMember(roomID, userID, targetID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("member removed"))
}

// GetMessages は指定チャットルームのメッセージ一覧をページネーション付きで取得する。
func (h *ChatRoomHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	messages, err := h.service.GetMessages(roomID, userID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(messages))
}

// SendMessage は指定チャットルームにメッセージを送信する。
func (h *ChatRoomHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.SendChatRoomMessageRequest](c)
	if input == nil {
		return
	}

	msg, err := h.service.SendMessage(roomID, userID, input.Content)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, msg)
}

// GetMyCount は認証ユーザーが参加しているチャットルーム総数を返す。
func (h *ChatRoomHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.service.CountByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
