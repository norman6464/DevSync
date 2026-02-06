package model

import "time"

// ReportPeriod represents the report period type
type ReportPeriod string

const (
	ReportPeriodWeekly  ReportPeriod = "weekly"
	ReportPeriodMonthly ReportPeriod = "monthly"
)

// ActivityReport represents a user's activity report for a period
type ActivityReport struct {
	Period            ReportPeriod `json:"period"`
	StartDate         time.Time    `json:"start_date"`
	EndDate           time.Time    `json:"end_date"`
	UserID            uint         `json:"user_id"`
	TotalContributions int         `json:"total_contributions"`
	PostsCreated      int          `json:"posts_created"`
	CommentsCreated   int          `json:"comments_created"`
	LikesReceived     int          `json:"likes_received"`
	GoalsCompleted    int          `json:"goals_completed"`
	GoalsProgress     int          `json:"goals_progress"` // Average progress of active goals
	NewFollowers      int          `json:"new_followers"`
	MessagesExchanged int          `json:"messages_exchanged"`
	// Daily breakdown for charts
	DailyContributions []DailyActivity `json:"daily_contributions"`
	// Top languages used
	TopLanguages []LanguageActivity `json:"top_languages"`
}

// DailyActivity represents activity for a single day
type DailyActivity struct {
	Date          string `json:"date"`
	Contributions int    `json:"contributions"`
	Posts         int    `json:"posts"`
	Comments      int    `json:"comments"`
}

// LanguageActivity represents language usage in the period
type LanguageActivity struct {
	Language string `json:"language"`
	Bytes    int64  `json:"bytes"`
	Repos    int    `json:"repos"`
}

// ReportComparison shows comparison with previous period
type ReportComparison struct {
	ContributionsDiff int     `json:"contributions_diff"`
	PostsDiff         int     `json:"posts_diff"`
	FollowersDiff     int     `json:"followers_diff"`
	GoalsDiff         int     `json:"goals_diff"`
	TrendPercentage   float64 `json:"trend_percentage"` // Overall activity trend
}
