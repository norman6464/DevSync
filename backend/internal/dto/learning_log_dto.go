package dto

// CreateLearningLogRequest は学習ログ作成リクエスト。
type CreateLearningLogRequest struct {
	Title    string `json:"title" binding:"required" validate:"required"`
	Content  string `json:"content" binding:"required" validate:"required"`
	Category string `json:"category"`
	Duration int    `json:"duration"`
	Source   string `json:"source"`
}

// UpdateLearningLogRequest は学習ログ更新リクエスト。
type UpdateLearningLogRequest struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Category *string `json:"category"`
	Duration *int    `json:"duration"`
}
