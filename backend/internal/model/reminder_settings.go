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
	ID     uint  `gorm:"primaryKey" json:"id"`
	UserID uint  `gorm:"uniqueIndex;not null" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Enabled          bool              `gorm:"default:true" json:"enabled"`              // リマインダー有効/無効
	Frequency        ReminderFrequency `gorm:"default:'daily'" json:"frequency"`         // 通知頻度
	NotificationTime string            `gorm:"default:'09:00'" json:"notification_time"` // 通知時間（HH:MM形式）
	InactiveDays     int               `gorm:"default:3" json:"inactive_days"`           // 何日間学習がない場合に通知するか

	EnableWeb   bool `gorm:"default:true" json:"enable_web"`    // Web通知を有効にする
	EnableEmail bool `gorm:"default:false" json:"enable_email"` // メール通知を有効にする

	LastRemindedAt *time.Time `json:"last_reminded_at"` // 最後にリマインドした日時
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
