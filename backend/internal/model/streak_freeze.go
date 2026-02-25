package model

import "time"

// StreakFreeze は学習ストリークのフリーズ（保護）を表す。
// ユーザーが忙しい日にストリークを守るために使用する。
type StreakFreeze struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	UsedDate  string    `json:"used_date" gorm:"not null;type:varchar(10)"` // "YYYY-MM-DD" 形式
	Month     int       `json:"month" gorm:"not null"`                      // 使用月（1-12）
	Year      int       `json:"year" gorm:"not null"`                       // 使用年
	CreatedAt time.Time `json:"created_at"`
}

// StreakFreezeStatus は今月のフリーズ使用状況を表す。
type StreakFreezeStatus struct {
	MaxFreezes   int            `json:"max_freezes"`   // 月あたりの最大フリーズ回数
	UsedFreezes  int            `json:"used_freezes"`  // 今月の使用済み回数
	Remaining    int            `json:"remaining"`     // 残り回数
	UsedDates    []string       `json:"used_dates"`    // 使用済み日付リスト
	TodayUsed    bool           `json:"today_used"`    // 今日フリーズを使ったか
	CanUseToday  bool           `json:"can_use_today"` // 今日使用可能か
}

// MaxFreezesPerMonth は月あたりの最大フリーズ回数。
const MaxFreezesPerMonth = 2
