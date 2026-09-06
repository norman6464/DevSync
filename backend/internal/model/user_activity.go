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
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	ActivityType ActivityType `json:"activity_type"`
	TargetType   string       `json:"target_type"`
	TargetID     uint         `json:"target_id"`
	Metadata     string       `json:"metadata"`
	CreatedAt    time.Time    `json:"created_at"`
}
