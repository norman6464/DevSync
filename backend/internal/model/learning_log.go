package model

import "time"

// LogCategory represents the category of a learning log
type LogCategory string

const (
	LogCategoryCoding  LogCategory = "coding"
	LogCategoryReading LogCategory = "reading"
	LogCategoryCourse  LogCategory = "course"
	LogCategoryMeetup  LogCategory = "meetup"
	LogCategoryOther   LogCategory = "other"
)

// LearningLog represents a daily learning diary entry
type LearningLog struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    uint        `json:"user_id" gorm:"not null;index"`
	Title     string      `json:"title" gorm:"not null"`
	Content   string      `json:"content" gorm:"type:text;not null"`
	Category  LogCategory `json:"category" gorm:"default:'other'"`
	Duration  int         `json:"duration" gorm:"default:0"` // minutes
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// CalendarEntry represents a single day's log count for calendar view
type CalendarEntry struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// StreakInfo represents streak data calculated from learning logs
type StreakInfo struct {
	CurrentStreak int    `json:"current_streak"`
	LongestStreak int    `json:"longest_streak"`
	TotalDays     int    `json:"total_days"`
	LastLogDate   string `json:"last_log_date"`
}
