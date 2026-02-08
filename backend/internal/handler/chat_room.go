package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// ChatRoomHandler はチャットルーム関連のHTTPハンドラ。
// チャットルームのCRUD・メンバー管理・メッセージ送受信を処理する。
type ChatRoomHandler struct {
	service *service.ChatRoomService
}

// NewChatRoomHandler は新しいChatRoomHandlerインスタンスを生成する。
func NewChatRoomHandler(s *service.ChatRoomService) *ChatRoomHandler {
	return &ChatRoomHandler{service: s}
}

// Create は新しいチャットルームを作成する。
func (h *ChatRoomHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		MemberIDs   []uint `json:"member_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := &model.ChatRoom{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     userID,
	}

	created, err := h.service.Create(room, input.MemberIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetMyRooms は現在のユーザーが参加しているチャットルーム一覧を取得する。
func (h *ChatRoomHandler) GetMyRooms(c *gin.Context) {
	userID := c.GetUint("userID")
	rooms, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// GetByID は指定IDのチャットルーム詳細を取得する。
func (h *ChatRoomHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.service.GetByID(uint(roomID), userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Update は指定IDのチャットルーム情報を更新する。
func (h *ChatRoomHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.service.Update(uint(roomID), userID, input.Name, input.Description)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only owner can update"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Delete は指定IDのチャットルームを削除する。
func (h *ChatRoomHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	if err := h.service.Delete(uint(roomID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}

// GetMembers は指定チャットルームのメンバー一覧を取得する。
func (h *ChatRoomHandler) GetMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	members, err := h.service.GetMembers(uint(roomID), userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

// AddMember はチャットルームに新しいメンバーを追加する。
func (h *ChatRoomHandler) AddMember(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddMember(uint(roomID), userID, input.UserID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

// RemoveMember はチャットルームからメンバーを削除する。
func (h *ChatRoomHandler) RemoveMember(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}
	targetID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.RemoveMember(uint(roomID), userID, uint(targetID)); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only owner can remove members"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

// GetMessages は指定チャットルームのメッセージ一覧をページネーション付きで取得する。
func (h *ChatRoomHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, err := h.service.GetMessages(uint(roomID), userID, page, limit)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// SendMessage は指定チャットルームにメッセージを送信する。
func (h *ChatRoomHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.service.SendMessage(uint(roomID), userID, input.Content)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}
