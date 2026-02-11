package service

import "github.com/norman6464/devsync/backend/internal/model"

// NotificationServiceInterface は通知サービスの契約を定義する。
// サービス間依存でインターフェースとして注入される。
type NotificationServiceInterface interface {
	CreateNotification(notification *model.Notification) error
	NotifyFollowers(actorID uint, postID uint, notificationType model.NotificationType)
}
