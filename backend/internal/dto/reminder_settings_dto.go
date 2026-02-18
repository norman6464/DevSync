package dto

import "github.com/norman6464/devsync/backend/internal/model"

// UpdateReminderSettingsRequest はリマインダー設定更新のリクエストボディ。
type UpdateReminderSettingsRequest struct {
	Enabled          bool                    `json:"enabled"`
	Frequency        model.ReminderFrequency `json:"frequency" binding:"omitempty,max=50"`
	NotificationTime string                  `json:"notification_time" binding:"omitempty,max=10"`
	InactiveDays     int                     `json:"inactive_days" binding:"omitempty,min=0,max=365"`
	EnableWeb        bool                    `json:"enable_web"`
	EnableEmail      bool                    `json:"enable_email"`
}
