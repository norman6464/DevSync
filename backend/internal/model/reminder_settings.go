package model

import "time"

// ReminderFrequency はリマインダーの通知頻度を表す。
type ReminderFrequency string

const (
	ReminderFrequencyDaily  ReminderFrequency = "daily"  // 毎日
	ReminderFrequencyWeekly ReminderFrequency = "weekly" // 毎週
)

// ReminderSettings はユーザーの学習リマインダー設定を表す。
type ReminderSettings struct {
	ID     uint  `json:"id"`
	UserID uint  `json:"user_id"`
	User   *User `json:"user,omitempty"`

	Enabled          bool              `json:"enabled"`           // リマインダー有効/無効
	Frequency        ReminderFrequency `json:"frequency"`         // 通知頻度
	NotificationTime string            `json:"notification_time"` // 通知時間（HH:MM形式）
	InactiveDays     int               `json:"inactive_days"`     // 何日間学習がない場合に通知するか

	EnableWeb   bool `json:"enable_web"`   // Web通知を有効にする
	EnableEmail bool `json:"enable_email"` // メール通知を有効にする

	LastRemindedAt *time.Time `json:"last_reminded_at"` // 最後にリマインドした日時
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
