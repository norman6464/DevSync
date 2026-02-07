package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
)

type ChatRoomHandler struct {
	roomRepo    *repository.ChatRoomRepository
	messageRepo *repository.GroupMessageRepository
	hub         *service.Hub
}

func NewChatRoomHandler(roomRepo *repository.ChatRoomRepository, messageRepo *repository.GroupMessageRepository, hub *service.Hub) *ChatRoomHandler {
	return &ChatRoomHandler{roomRepo: roomRepo, messageRepo: messageRepo, hub: hub}
}

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
	if err := h.roomRepo.Create(room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add owner as member
	if err := h.roomRepo.AddMember(room.ID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add other members
	for _, memberID := range input.MemberIDs {
		if memberID != userID {
			h.roomRepo.AddMember(room.ID, memberID)
		}
	}

	room, _ = h.roomRepo.FindByID(room.ID)
	c.JSON(http.StatusCreated, room)
}

func (h *ChatRoomHandler) GetMyRooms(c *gin.Context) {
	userID := c.GetUint("userID")
	rooms, err := h.roomRepo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func (h *ChatRoomHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isMember, err := h.roomRepo.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
		return
	}

	room, err := h.roomRepo.FindByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *ChatRoomHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.roomRepo.FindByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	if room.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can update"})
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

	if input.Name != "" {
		room.Name = input.Name
	}
	room.Description = input.Description

	if err := h.roomRepo.Update(room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (h *ChatRoomHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.roomRepo.FindByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	if room.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can delete"})
		return
	}

	if err := h.roomRepo.Delete(uint(roomID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "room deleted"})
}

func (h *ChatRoomHandler) GetMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isMember, err := h.roomRepo.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
		return
	}

	members, err := h.roomRepo.GetMembers(uint(roomID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, members)
}

func (h *ChatRoomHandler) AddMember(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isMember, err := h.roomRepo.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
		return
	}

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roomRepo.AddMember(uint(roomID), input.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member added"})
}

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

	room, err := h.roomRepo.FindByID(uint(roomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	// Owner can remove anyone, others can only remove themselves
	if room.OwnerID != userID && userID != uint(targetID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owner can remove members"})
		return
	}

	if err := h.roomRepo.RemoveMember(uint(roomID), uint(targetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

func (h *ChatRoomHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isMember, err := h.roomRepo.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, err := h.messageRepo.FindByRoomID(uint(roomID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *ChatRoomHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	roomID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isMember, err := h.roomRepo.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "not a member"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &model.GroupMessage{
		ChatRoomID: uint(roomID),
		SenderID:   userID,
		Content:    input.Content,
	}
	if err := h.messageRepo.Create(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with sender info
	msg.Sender = &model.User{}
	h.messageRepo.FindSenderByID(msg)

	// Send via WebSocket to all room members
	go func() {
		wsMsg := service.WSMessage{
			Type:       "group_message",
			SenderID:   userID,
			RoomID:     uint(roomID),
			Content:    input.Content,
			SenderName: msg.Sender.Name,
		}
		data, _ := json.Marshal(wsMsg)
		h.hub.SendToRoom(uint(roomID), userID, data)
	}()

	c.JSON(http.StatusCreated, msg)
}
