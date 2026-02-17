package dto

import "github.com/norman6464/devsync/backend/internal/model"

// AIAdviceResponse はAIアドバイス取得レスポンス。
type AIAdviceResponse struct {
	Advices            []model.AIAdvice `json:"advices"`
	LLMAvailable       bool             `json:"llm_available"`
	DailyChatRemaining int              `json:"daily_chat_remaining"`
}

// AIChatRequest はLLMチャットリクエスト。
type AIChatRequest struct {
	Message        string `json:"message" binding:"required" validate:"required"`
	ConversationID uint   `json:"conversation_id"`
}
