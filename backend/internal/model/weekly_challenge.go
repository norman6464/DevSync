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
	ID            uint          `json:"id" gorm:"primaryKey"`
	UserID        uint          `json:"user_id" gorm:"not null;index"`
	Year          int           `json:"year" gorm:"not null"`
	Week          int           `json:"week" gorm:"not null"`
	ChallengeType ChallengeType `json:"challenge_type" gorm:"not null"`
	Description   string        `json:"description" gorm:"not null"`
	TargetValue   int           `json:"target_value" gorm:"not null"`
	CurrentValue  int           `json:"current_value" gorm:"default:0"`
	IsCompleted   bool          `json:"is_completed" gorm:"default:false"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
