package model

// LearningLogStats はユーザーの学習ログ集計統計を表す。
type LearningLogStats struct {
	TotalLogs     int64 `json:"total_logs"`
	TotalDuration int64 `json:"total_duration"`
	CategoryCount int64 `json:"category_count"`
	LogsThisMonth int64 `json:"logs_this_month"`
}
