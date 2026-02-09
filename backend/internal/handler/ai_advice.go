package handler

import (
	"errors"
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

	c.JSON(http.StatusOK, gin.H{
		"advices":              advices,
		"llm_available":        llmAvailable,
		"daily_chat_remaining": remaining,
	})
}

// MarkAsRead はアドバイスを既読にする。
func (h *AIAdviceHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid advice id"})
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.MarkAsRead(uint(id), userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "advice not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
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
		if errors.Is(err, service.ErrLLMNotConfigured) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "LLM service not configured"})
			return
		}
		if errors.Is(err, service.ErrRateLimitExceeded) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "daily chat limit reached"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process chat"})
		return
	}

	c.JSON(http.StatusOK, conv)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get conversations"})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

// GetConversation は会話詳細を取得する。
func (h *AIAdviceHandler) GetConversation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return
	}
	userID := c.GetUint("userID")

	conv, err := h.service.GetConversation(uint(id), userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, conv)
}
