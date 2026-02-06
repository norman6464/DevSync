package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type MessageHandler struct {
	repo             *repository.MessageRepository
	notificationRepo *repository.NotificationRepository
}

func NewMessageHandler(repo *repository.MessageRepository, notificationRepo *repository.NotificationRepository) *MessageHandler {
	return &MessageHandler{repo: repo, notificationRepo: notificationRepo}
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")
	conversations, err := h.repo.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conversations)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	otherID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// Mark messages as read
	h.repo.MarkAsRead(uint(otherID), userID)

	messages, err := h.repo.GetConversation(userID, uint(otherID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	receiverID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &model.Message{
		SenderID:   userID,
		ReceiverID: uint(receiverID),
		Content:    input.Content,
	}
	if err := h.repo.Create(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create notification for message receiver
	go func(senderID, receiverID uint) {
		notification := &model.Notification{
			UserID:  receiverID,
			Type:    model.NotificationTypeMessage,
			ActorID: senderID,
		}
		h.notificationRepo.Create(notification)
	}(userID, uint(receiverID))

	c.JSON(http.StatusCreated, msg)
}
