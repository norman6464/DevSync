package dto

import "github.com/norman6464/devsync/backend/internal/model"

// UpdateReminderSettingsRequest はリマインダー設定更新のリクエストボディ。
type UpdateReminderSettingsRequest struct {
	Enabled          bool                    `json:"enabled"`
	Frequency        model.ReminderFrequency `json:"frequency"`
	NotificationTime string                  `json:"notification_time"`
	InactiveDays     int                     `json:"inactive_days"`
	EnableWeb        bool                    `json:"enable_web"`
	EnableEmail      bool                    `json:"enable_email"`
}
