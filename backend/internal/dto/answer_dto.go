package dto

// CreateAnswerRequest は回答作成のリクエストボディ。
type CreateAnswerRequest struct {
	Body string `json:"body" binding:"required"`
}

// UpdateAnswerRequest は回答更新のリクエストボディ。
type UpdateAnswerRequest struct {
	Body string `json:"body" binding:"required"`
}
