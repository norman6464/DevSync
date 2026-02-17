package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// MessageServiceInterface はMessageHandlerが依存するサービスのインターフェース。
type MessageServiceInterface interface {
	GetConversations(userID uint) ([]repository.ConversationSummary, error)
	GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error)
	SendMessage(msg *model.Message) error
}

// MessageHandler はDM（ダイレクトメッセージ）関連のHTTPハンドラ。
// 会話一覧・メッセージ取得・メッセージ送信を処理する。
type MessageHandler struct {
	service MessageServiceInterface
}

// NewMessageHandler は新しいMessageHandlerインスタンスを生成する。
func NewMessageHandler(s MessageServiceInterface) *MessageHandler {
	return &MessageHandler{service: s}
}

// GetConversations は認証ユーザーの会話一覧を返す。
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")
	conversations, err := h.service.GetConversations(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, conversations)
}

// GetMessages は指定ユーザーとの会話メッセージをページネーション付きで返す。
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	otherID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	messages, err := h.service.GetConversation(userID, otherID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, messages)
}

// SendMessage は指定ユーザーにDMを送信する。
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	receiverID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	msg := &model.Message{
		SenderID:   userID,
		ReceiverID: receiverID,
		Content:    input.Content,
	}
	if err := h.service.SendMessage(msg); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, msg)
}
