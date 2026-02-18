package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateQuestionRequest は質問作成のリクエストボディ。
type CreateQuestionRequest struct {
	Title string `json:"title" binding:"required,max=500" validate:"required,max=500"`
	Body  string `json:"body" binding:"required,max=50000" validate:"required,max=50000"`
	Tags  string `json:"tags" binding:"omitempty,max=5000"`
}

// UpdateQuestionRequest は質問更新のリクエストボディ。
type UpdateQuestionRequest struct {
	Title string `json:"title" binding:"max=500" validate:"max=500"`
	Body  string `json:"body" binding:"omitempty,max=50000"`
	Tags  string `json:"tags" binding:"omitempty,max=5000"`
}

// QuestionListResponse は質問一覧レスポンス。
type QuestionListResponse struct {
	Questions []model.Question `json:"questions"`
	Total     int64            `json:"total"`
	Limit     int              `json:"limit"`
	Offset    int              `json:"offset"`
}

// QuestionDetailResponse は質問詳細レスポンス（ユーザー投票情報付き）。
type QuestionDetailResponse struct {
	Question model.Question `json:"question"`
	UserVote int            `json:"user_vote"`
}
