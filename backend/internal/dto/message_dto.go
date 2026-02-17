package dto

// SendMessageRequest はDMメッセージ送信リクエスト。
type SendMessageRequest struct {
	Content string `json:"content" binding:"required" validate:"required"`
}
