package model

import "time"

// ChallengeType はウィークリーチャレンジの種類を表す型。
type ChallengeType string

const (
	ChallengeDurationTotal ChallengeType = "duration_total" // 合計学習時間（分）
	ChallengeStreakDays    ChallengeType = "streak_days"    // 連続学習日数
	ChallengeCategoryCount ChallengeType = "category_count" // 異なるカテゴリ数
	ChallengeLogCount      ChallengeType = "log_count"      // 学習ログ記録回数
)

// WeeklyChallenge はウィークリーチャレンジを表す。
type WeeklyChallenge struct {
	ID            uint          `json:"id"`
	UserID        uint          `json:"user_id"`
	Year          int           `json:"year"`
	Week          int           `json:"week"`
	ChallengeType ChallengeType `json:"challenge_type"`
	Description   string        `json:"description"`
	TargetValue   int           `json:"target_value"`
	CurrentValue  int           `json:"current_value"`
	IsCompleted   bool          `json:"is_completed"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
