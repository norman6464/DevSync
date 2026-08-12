package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// MessageHandler はDM（ダイレクトメッセージ）関連のHTTPハンドラ。
// 会話一覧・メッセージ取得・メッセージ送信を処理する。
type MessageHandler struct {
	listConversations *usecase.ListConversationsUseCase
	getConversation   *usecase.GetConversationUseCase
	send              *usecase.SendMessageUseCase
	markAsRead        *usecase.MarkMessagesAsReadUseCase
}

// NewMessageHandler は新しいMessageHandlerインスタンスを生成する。
func NewMessageHandler(
	listConversations *usecase.ListConversationsUseCase,
	getConversation *usecase.GetConversationUseCase,
	send *usecase.SendMessageUseCase,
	markAsRead *usecase.MarkMessagesAsReadUseCase,
) *MessageHandler {
	return &MessageHandler{
		listConversations: listConversations,
		getConversation:   getConversation,
		send:              send,
		markAsRead:        markAsRead,
	}
}

// GetConversations は認証ユーザーの会話一覧を返す。
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")
	conversations, err := h.listConversations.Execute(c.Request.Context(), userID)
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

	messages, err := h.getConversation.Execute(c.Request.Context(), userID, otherID, page, limit)
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

	if err := h.markAsRead.Execute(c.Request.Context(), senderID, userID); err != nil {
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
	if err := h.send.Execute(c.Request.Context(), msg); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, msg)
}
