package model

import "time"

// ResourceProgressStatus はリソース進捗のステータスを表す型。
type ResourceProgressStatus string

const (
	ResourceProgressNotStarted ResourceProgressStatus = "not_started"
	ResourceProgressInProgress ResourceProgressStatus = "in_progress"
	ResourceProgressCompleted  ResourceProgressStatus = "completed"
)

// ResourceProgress は学習リソースに対するユーザーの進捗を管理する。
type ResourceProgress struct {
	ID                uint                   `json:"id"`
	UserID            uint                   `json:"user_id"`
	ResourceID        uint                   `json:"resource_id"`
	Resource          *LearningResource      `json:"resource,omitempty"`
	Status            ResourceProgressStatus `json:"status"`
	CompletionPercent int                    `json:"completion_percent"`
	Note              string                 `json:"note"`
	StartedAt         *time.Time             `json:"started_at"`
	CompletedAt       *time.Time             `json:"completed_at"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}
