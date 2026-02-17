package dto

// AIChatRequest はLLMチャットリクエスト。
type AIChatRequest struct {
	Message        string `json:"message" binding:"required" validate:"required"`
	ConversationID uint   `json:"conversation_id"`
}
