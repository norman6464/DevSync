package model

// NotificationStats はユーザーの通知に関する集計統計を表す。
type NotificationStats struct {
	TotalNotifications    int64 `json:"total_notifications"`
	UnreadCount           int64 `json:"unread_count"`
	NotificationsThisMonth int64 `json:"notifications_this_month"`
}
