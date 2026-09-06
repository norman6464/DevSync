package model

import "time"

// WeeklyGoal はカテゴリ別の週間学習時間目標を表す。
type WeeklyGoal struct {
	ID            uint        `json:"id"`
	UserID        uint        `json:"user_id"`
	Category      LogCategory `json:"category"`
	TargetMinutes int         `json:"target_minutes"` // 目標学習時間（分）
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// WeeklyGoalProgress はカテゴリ別の週間学習目標の達成状況を表す。
type WeeklyGoalProgress struct {
	Category        LogCategory `json:"category"`
	TargetMinutes   int         `json:"target_minutes"`
	ActualMinutes   int         `json:"actual_minutes"`
	ProgressPercent int         `json:"progress_percent"` // 0〜100+
}
