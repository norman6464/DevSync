package dto

// SendMessageRequest はDMメッセージ送信リクエスト。
type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000" validate:"required,min=1,max=5000"`
}
