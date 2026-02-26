package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// MessageServiceInterface はMessageHandlerが依存するサービスのインターフェース。
type MessageServiceInterface interface {
	GetConversations(userID uint) ([]model.ConversationSummary, error)
	GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error)
	SendMessage(msg *model.Message) error
	MarkAsRead(senderID, receiverID uint) error
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
	respondOK(c, ensureSlice(conversations))
}

// GetMessages は指定ユーザーとの会話メッセージをページネーション付きで返す。
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	otherID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	page, limit := parsePagination(c)

	messages, err := h.service.GetConversation(userID, otherID, page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(messages))
}

// MarkAsRead は指定ユーザーからのメッセージを既読にマークする。
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID")
	senderID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	if err := h.service.MarkAsRead(senderID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"message": "既読にしました"})
}

// SendMessage は指定ユーザーにDMを送信する。
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	receiverID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	input := bindJSON[dto.SendMessageRequest](c)
	if input == nil {
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
