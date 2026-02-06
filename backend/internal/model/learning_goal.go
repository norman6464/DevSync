package model

import "time"

// GoalStatus represents the status of a learning goal
type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusPaused    GoalStatus = "paused"
)

// GoalCategory represents the category of a learning goal
type GoalCategory string

const (
	GoalCategoryLanguage  GoalCategory = "language"
	GoalCategoryFramework GoalCategory = "framework"
	GoalCategorySkill     GoalCategory = "skill"
	GoalCategoryProject   GoalCategory = "project"
	GoalCategoryOther     GoalCategory = "other"
)

// LearningGoal represents a user's learning goal
type LearningGoal struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UserID      uint         `json:"user_id" gorm:"not null;index"`
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description"`
	Category    GoalCategory `json:"category" gorm:"default:'other'"`
	TargetDate  *time.Time   `json:"target_date"`
	Progress    int          `json:"progress" gorm:"default:0"` // 0-100
	Status      GoalStatus   `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CompletedAt *time.Time   `json:"completed_at"`
}

// LearningGoalStats represents aggregated learning goal statistics for a user
type LearningGoalStats struct {
	TotalGoals     int `json:"total_goals"`
	ActiveGoals    int `json:"active_goals"`
	CompletedGoals int `json:"completed_goals"`
	AverageProgress int `json:"average_progress"`
}
