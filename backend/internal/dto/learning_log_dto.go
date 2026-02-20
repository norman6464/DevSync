package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateLearningLogRequest は学習ログ作成リクエスト。
type CreateLearningLogRequest struct {
	Title    string `json:"title" binding:"required,max=200" validate:"required,max=200"`
	Content  string `json:"content" binding:"required,max=50000" validate:"required,max=50000"`
	Category string `json:"category" binding:"omitempty,max=100"`
	Duration int    `json:"duration" binding:"omitempty,min=0,max=1440"`
	Source   string `json:"source" binding:"omitempty,max=500"`
}

// LearningLogListResponse は学習ログ一覧レスポンス（ページネーション付き）。
type LearningLogListResponse struct {
	Logs   []model.LearningLog `json:"logs"`
	Total  int64               `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

// UpdateLearningLogRequest は学習ログ更新リクエスト。
type UpdateLearningLogRequest struct {
	Title    *string `json:"title" binding:"omitempty,max=200"`
	Content  *string `json:"content" binding:"omitempty,max=50000"`
	Category *string `json:"category" binding:"omitempty,max=100"`
	Duration *int    `json:"duration" binding:"omitempty,min=0,max=1440"`
}
