package model

import "time"

// ActivityType はアクティビティの種別を表す型。
type ActivityType string

const (
	ActivityPostCreated    ActivityType = "post_created"
	ActivityCommentCreated ActivityType = "comment_created"
	ActivityResourceShared ActivityType = "resource_shared"
	ActivityGoalCompleted  ActivityType = "goal_completed"
	ActivityProjectCreated ActivityType = "project_created"
	ActivitySnippetCreated ActivityType = "snippet_created"
)

// UserActivity はユーザーのアクティビティログを表す。
type UserActivity struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	UserID       uint         `json:"user_id" gorm:"not null;index"`
	ActivityType ActivityType `json:"activity_type" gorm:"size:50;not null;index"`
	TargetType   string       `json:"target_type" gorm:"size:50;not null"`
	TargetID     uint         `json:"target_id" gorm:"not null"`
	Metadata     string       `json:"metadata" gorm:"type:text"`
	CreatedAt    time.Time    `json:"created_at" gorm:"index"`
}
