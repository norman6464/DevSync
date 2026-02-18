package dto

// CreateGoalRequest は学習目標作成のリクエストボディ。
type CreateGoalRequest struct {
	Title       string `json:"title" binding:"required" validate:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	TargetDate  string `json:"target_date"`
}

// UpdateGoalRequest は学習目標更新のリクエストボディ。
type UpdateGoalRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	TargetDate  *string `json:"target_date"`
	Progress    *int    `json:"progress"`
	Status      *string `json:"status"`
}
