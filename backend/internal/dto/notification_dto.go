package dto

import "github.com/norman6464/devsync/backend/internal/model"

// NotificationListResponse は通知一覧レスポンス
type NotificationListResponse struct {
	Notifications []model.Notification `json:"notifications"`
	Total         int64                `json:"total"`
	Page          int                  `json:"page"`
	Limit         int                  `json:"limit"`
}
