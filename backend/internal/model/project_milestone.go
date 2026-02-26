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
	ID          uint            `json:"id" gorm:"primaryKey"`
	ProjectID   uint            `json:"project_id" gorm:"not null;index"`
	Title       string          `json:"title" gorm:"size:200;not null"`
	Description string          `json:"description" gorm:"type:text"`
	Status      MilestoneStatus `json:"status" gorm:"size:20;default:'not_started'"`
	DueDate     *time.Time      `json:"due_date"`
	CompletedAt *time.Time      `json:"completed_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
