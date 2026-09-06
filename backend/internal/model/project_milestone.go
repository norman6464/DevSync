package model

import "time"

// MilestoneStatus はマイルストーンのステータスを表す型。
type MilestoneStatus string

const (
	MilestoneNotStarted MilestoneStatus = "not_started"
	MilestoneInProgress MilestoneStatus = "in_progress"
	MilestoneCompleted  MilestoneStatus = "completed"
)

// ProjectMilestone はプロジェクトのマイルストーンを管理する。
type ProjectMilestone struct {
	ID          uint            `json:"id"`
	ProjectID   uint            `json:"project_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      MilestoneStatus `json:"status"`
	DueDate     *time.Time      `json:"due_date"`
	CompletedAt *time.Time      `json:"completed_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
