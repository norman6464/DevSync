package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ChatRoomHandler はチャットルーム関連のHTTPハンドラ。
// チャットルームのCRUD・メンバー管理・メッセージ送受信を処理する。
type ChatRoomHandler struct {
	create       *usecase.CreateChatRoomUseCase
	listMine     *usecase.ListMyChatRoomsUseCase
	get          *usecase.GetChatRoomUseCase
	update       *usecase.UpdateChatRoomUseCase
	remove       *usecase.DeleteChatRoomUseCase
	listMembers  *usecase.ListChatRoomMembersUseCase
	addMember    *usecase.AddChatRoomMemberUseCase
	removeMember *usecase.RemoveChatRoomMemberUseCase
	listMessages *usecase.ListChatRoomMessagesUseCase
	sendMessage  *usecase.SendChatRoomMessageUseCase
	countMine    *usecase.CountMyChatRoomsUseCase
}

// NewChatRoomHandler は新しいChatRoomHandlerインスタンスを生成する。
func NewChatRoomHandler(
	create *usecase.CreateChatRoomUseCase,
	listMine *usecase.ListMyChatRoomsUseCase,
	get *usecase.GetChatRoomUseCase,
	update *usecase.UpdateChatRoomUseCase,
	remove *usecase.DeleteChatRoomUseCase,
	listMembers *usecase.ListChatRoomMembersUseCase,
	addMember *usecase.AddChatRoomMemberUseCase,
	removeMember *usecase.RemoveChatRoomMemberUseCase,
	listMessages *usecase.ListChatRoomMessagesUseCase,
	sendMessage *usecase.SendChatRoomMessageUseCase,
	countMine *usecase.CountMyChatRoomsUseCase,
) *ChatRoomHandler {
	return &ChatRoomHandler{
		create:       create,
		listMine:     listMine,
		get:          get,
		update:       update,
		remove:       remove,
		listMembers:  listMembers,
		addMember:    addMember,
		removeMember: removeMember,
		listMessages: listMessages,
		sendMessage:  sendMessage,
		countMine:    countMine,
	}
}

// Create は新しいチャットルームを作成する。
func (h *ChatRoomHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateChatRoomRequest](c)
	if input == nil {
		return
	}

	created, err := h.create.Execute(c.Request.Context(), usecase.CreateChatRoomInput{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
		MemberIDs:   input.MemberIDs,
	})
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
	rooms, total, err := h.listMine.Execute(c.Request.Context(), userID, limit, offset)
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
	handleGetByID(c, func(roomID, userID uint) (*model.ChatRoom, error) {
		return h.get.Execute(c.Request.Context(), roomID, userID)
	})
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

	room, err := h.update.Execute(c.Request.Context(), roomID, userID, input.Name, input.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, room)
}

// Delete は指定IDのチャットルームを削除する。
func (h *ChatRoomHandler) Delete(c *gin.Context) {
	handleDelete(c, func(roomID, userID uint) error {
		return h.remove.Execute(c.Request.Context(), roomID, userID)
	})
}

// GetMembers は指定チャットルームのメンバー一覧を取得する。
func (h *ChatRoomHandler) GetMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, ok := parseID(c, "id")
	if !ok {
		return
	}

	members, err := h.listMembers.Execute(c.Request.Context(), roomID, userID)
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

	if err := h.addMember.Execute(c.Request.Context(), roomID, userID, input.UserID); err != nil {
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

	if err := h.removeMember.Execute(c.Request.Context(), roomID, userID, targetID); err != nil {
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

	messages, err := h.listMessages.Execute(c.Request.Context(), roomID, userID, page, limit)
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

	msg, err := h.sendMessage.Execute(c.Request.Context(), roomID, userID, input.Content)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, msg)
}

// GetMyCount は認証ユーザーが参加しているチャットルーム総数を返す。
func (h *ChatRoomHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.countMine.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}
