package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// MessageHandler はDM（ダイレクトメッセージ）関連のHTTPハンドラ。
// 会話一覧・メッセージ取得・メッセージ送信を処理する。
type MessageHandler struct {
	service *service.MessageService
}

// NewMessageHandler は新しいMessageHandlerインスタンスを生成する。
func NewMessageHandler(s *service.MessageService) *MessageHandler {
	return &MessageHandler{service: s}
}

// GetConversations は認証ユーザーの会話一覧を返す。
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")
	conversations, err := h.service.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conversations)
}

// GetMessages は指定ユーザーとの会話メッセージをページネーション付きで返す。
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	otherID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, err := h.service.GetConversation(userID, uint(otherID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// SendMessage は指定ユーザーにDMを送信する。
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
	if err := h.service.SendMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}
