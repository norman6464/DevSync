package dto

// SendMessageRequest はDMメッセージ送信リクエスト。
type SendMessageRequest struct {
	Content string `json:"content" binding:"required,max=5000" validate:"required,max=5000"`
}
