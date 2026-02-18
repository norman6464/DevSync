package dto

// CreateLearningLogRequest は学習ログ作成リクエスト。
type CreateLearningLogRequest struct {
	Title    string `json:"title" binding:"required,max=200" validate:"required,max=200"`
	Content  string `json:"content" binding:"required,max=50000" validate:"required,max=50000"`
	Category string `json:"category" binding:"omitempty,max=100"`
	Duration int    `json:"duration" binding:"omitempty,min=0,max=1440"`
	Source   string `json:"source" binding:"omitempty,max=500"`
}

// UpdateLearningLogRequest は学習ログ更新リクエスト。
type UpdateLearningLogRequest struct {
	Title    *string `json:"title" binding:"omitempty,max=200"`
	Content  *string `json:"content" binding:"omitempty,max=50000"`
	Category *string `json:"category" binding:"omitempty,max=100"`
	Duration *int    `json:"duration" binding:"omitempty,min=0,max=1440"`
}
