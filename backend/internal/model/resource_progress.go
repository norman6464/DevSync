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
	ID                uint                   `json:"id" gorm:"primaryKey"`
	UserID            uint                   `json:"user_id" gorm:"not null;uniqueIndex:idx_resource_progress"`
	ResourceID        uint                   `json:"resource_id" gorm:"not null;uniqueIndex:idx_resource_progress"`
	Resource          *LearningResource      `json:"resource,omitempty" gorm:"foreignKey:ResourceID"`
	Status            ResourceProgressStatus `json:"status" gorm:"size:20;default:'not_started'"`
	CompletionPercent int                    `json:"completion_percent" gorm:"default:0"`
	Note              string                 `json:"note" gorm:"type:text"`
	StartedAt         *time.Time             `json:"started_at"`
	CompletedAt       *time.Time             `json:"completed_at"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}
