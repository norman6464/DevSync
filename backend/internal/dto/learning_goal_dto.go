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
	Title       string `json:"title" binding:"required,min=1,max=200" validate:"required,min=1,max=200"`
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

// BatchUpdateProgressRequest は目標進捗一括更新のリクエストボディ。
type BatchUpdateProgressRequest struct {
	Updates []GoalProgressUpdate `json:"updates" binding:"required,min=1,max=50"`
}

// GoalProgressUpdate は個別の目標進捗更新データ。
type GoalProgressUpdate struct {
	GoalID   uint `json:"goal_id" binding:"required"`
	Progress int  `json:"progress" binding:"min=0,max=100"`
}
