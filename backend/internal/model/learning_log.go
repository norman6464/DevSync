package model

import "time"

// LogCategory は学習ログのカテゴリを表す型。
type LogCategory string

// 学習ログのカテゴリ定数群。
const (
	LogCategoryCoding  LogCategory = "coding"  // コーディング
	LogCategoryReading LogCategory = "reading" // 読書
	LogCategoryCourse  LogCategory = "course"  // 講座・コース
	LogCategoryMeetup  LogCategory = "meetup"  // 勉強会・ミートアップ
	LogCategoryOther   LogCategory = "other"   // その他
)

// ValidCategories は有効なカテゴリ値のマップ（バリデーション用）。
var ValidCategories = map[LogCategory]bool{
	LogCategoryCoding:  true,
	LogCategoryReading: true,
	LogCategoryCourse:  true,
	LogCategoryMeetup:  true,
	LogCategoryOther:   true,
}

// LogSource は学習ログの記録元を表す型。
type LogSource string

// 学習ログのソース定数群。
const (
	LogSourceManual   LogSource = "manual"   // 手動入力
	LogSourcePomodoro LogSource = "pomodoro" // ポモドーロタイマー
)

// ValidSources は有効なソース値のマップ（バリデーション用）。
var ValidSources = map[LogSource]bool{
	LogSourceManual:   true,
	LogSourcePomodoro: true,
}

// LearningLog は日々の学習記録を表す。
type LearningLog struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserID    uint        `json:"user_id" gorm:"not null;index"`
	Title     string      `json:"title" gorm:"not null"`
	Content   string      `json:"content" gorm:"type:text;not null"`
	Category  LogCategory `json:"category" gorm:"default:'other'"`
	Duration  int         `json:"duration" gorm:"default:0"`      // 学習時間（分単位）
	Source    LogSource   `json:"source" gorm:"default:'manual'"` // 記録元（manual/pomodoro）
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// CalendarEntry はカレンダービュー用の日別ログ件数を表す。
type CalendarEntry struct {
	Date  string `json:"date"`  // "YYYY-MM-DD" 形式
	Count int    `json:"count"` // その日のログ件数
}

// StreakInfo は学習ログから算出した連続学習情報を表す。
type StreakInfo struct {
	CurrentStreak int    `json:"current_streak"` // 現在の連続日数
	LongestStreak int    `json:"longest_streak"` // 最長連続日数
	TotalDays     int    `json:"total_days"`     // 合計学習日数
	LastLogDate   string `json:"last_log_date"`  // 最後にログを記録した日
}
