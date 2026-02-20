package dto

import "github.com/norman6464/devsync/backend/internal/model"

// GoalListResponse は学習目標一覧レスポンス。
type GoalListResponse struct {
	Goals  []model.LearningGoal `json:"goals"`
	Total  int64                `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

// CreateGoalRequest は学習目標作成のリクエストボディ。
type CreateGoalRequest struct {
	Title       string `json:"title" binding:"required,max=200" validate:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Category    string `json:"category" binding:"omitempty,max=100"`
	TargetDate  string `json:"target_date" binding:"omitempty,max=20"`
}

// UpdateGoalRequest は学習目標更新のリクエストボディ。
type UpdateGoalRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	Category    *string `json:"category" binding:"omitempty,max=100"`
	TargetDate  *string `json:"target_date" binding:"omitempty,max=20"`
	Progress    *int    `json:"progress" binding:"omitempty,min=0,max=100"`
	Status      *string `json:"status" binding:"omitempty,max=50"`
}
