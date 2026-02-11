package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AIAdviceHandler はAIアドバイス関連のHTTPハンドラ。
type AIAdviceHandler struct {
	service *service.AIAdviceService
}

// NewAIAdviceHandler は新しいAIAdviceHandlerインスタンスを生成する。
func NewAIAdviceHandler(s *service.AIAdviceService) *AIAdviceHandler {
	return &AIAdviceHandler{service: s}
}

// GetAdvice はルールベースアドバイスを取得する。
// LLM利用可否と本日の残りチャット回数も返す。
func (h *AIAdviceHandler) GetAdvice(c *gin.Context) {
	userID := c.GetUint("userID")

	// ルールエンジンを実行してリアルタイムにアドバイス生成
	advices := h.service.GenerateAdvice(userID)

	// LLM利用可否
	llmAvailable := h.service.IsLLMAvailable()

	// 本日の残りチャット回数
	remaining := 0
	if llmAvailable {
		var err error
		remaining, err = h.service.GetDailyChatRemaining(userID)
		if err != nil {
			remaining = 0
		}
	}

	respondOK(c, gin.H{
		"advices":              advices,
		"llm_available":        llmAvailable,
		"daily_chat_remaining": remaining,
	})
}

// MarkAsRead はアドバイスを既読にする。
func (h *AIAdviceHandler) MarkAsRead(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.MarkAsRead(id, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "advice not found"})
		return
	}

	respondOK(c, gin.H{"message": "marked as read"})
}

// Chat はLLMとのチャットメッセージを送信する。
func (h *AIAdviceHandler) Chat(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Message        string `json:"message" binding:"required"`
		ConversationID uint   `json:"conversation_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.service.Chat(userID, input.Message, input.ConversationID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, conv)
}

// DeleteConversation は会話を削除する。
func (h *AIAdviceHandler) DeleteConversation(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.DeleteConversation(id, userID); err != nil {
		log.Printf("会話削除エラー: %v", err)
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "conversation deleted"})
}

// GetConversations は会話履歴一覧を取得する。
func (h *AIAdviceHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("userID")

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	conversations, err := h.service.GetConversations(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, conversations)
}

// GetConversation は会話詳細を取得する。
func (h *AIAdviceHandler) GetConversation(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	conv, err := h.service.GetConversation(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, conv)
}
