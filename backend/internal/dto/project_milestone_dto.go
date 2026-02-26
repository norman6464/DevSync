package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateMilestoneRequest はマイルストーン作成のリクエストボディ。
type CreateMilestoneRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"omitempty,max=5000"`
	DueDate     string `json:"due_date" binding:"omitempty,max=20"`
}

// UpdateMilestoneRequest はマイルストーン更新のリクエストボディ。
type UpdateMilestoneRequest struct {
	Title       string `json:"title" binding:"omitempty,max=200"`
	Description string `json:"description" binding:"omitempty,max=5000"`
	DueDate     string `json:"due_date" binding:"omitempty,max=20"`
	Status      string `json:"status" binding:"omitempty"`
}

// MilestoneListResponse はマイルストーン一覧レスポンス。
type MilestoneListResponse struct {
	Milestones []model.ProjectMilestone `json:"milestones"`
}
