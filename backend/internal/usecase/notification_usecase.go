package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListNotificationsUseCase は通知一覧と総数を取得する。
type ListNotificationsUseCase struct {
	notifications repository.NotificationReader
}

// NewListNotificationsUseCase は ListNotificationsUseCase を生成する。
func NewListNotificationsUseCase(notifications repository.NotificationReader) *ListNotificationsUseCase {
	return &ListNotificationsUseCase{notifications: notifications}
}

// Execute は通知を新しい順に返し、あわせて絞り込み後の総数を返す。
func (uc *ListNotificationsUseCase) Execute(ctx context.Context, userID uint, page, limit int, notificationType string) ([]model.Notification, int64, error) {
	notifications, err := uc.notifications.FindByUserID(ctx, userID, page, limit, notificationType)
	if err != nil {
		return nil, 0, err
	}

	total, err := uc.notifications.CountByUserID(ctx, userID, notificationType)
	if err != nil {
		return nil, 0, err
	}
	return notifications, total, nil
}

// CountUnreadNotificationsUseCase は未読通知の数を取得する。
type CountUnreadNotificationsUseCase struct {
	notifications repository.NotificationReader
}

// NewCountUnreadNotificationsUseCase は CountUnreadNotificationsUseCase を生成する。
func NewCountUnreadNotificationsUseCase(notifications repository.NotificationReader) *CountUnreadNotificationsUseCase {
	return &CountUnreadNotificationsUseCase{notifications: notifications}
}

// Execute は未読通知の数を返す。
func (uc *CountUnreadNotificationsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.notifications.CountUnread(ctx, userID)
}

// MarkNotificationAsReadUseCase は通知 1 件を既読にする。
type MarkNotificationAsReadUseCase struct {
	notifications repository.NotificationReader
}

// NewMarkNotificationAsReadUseCase は MarkNotificationAsReadUseCase を生成する。
func NewMarkNotificationAsReadUseCase(notifications repository.NotificationReader) *MarkNotificationAsReadUseCase {
	return &MarkNotificationAsReadUseCase{notifications: notifications}
}

// Execute は本人の通知だけを既読にする。
func (uc *MarkNotificationAsReadUseCase) Execute(ctx context.Context, id, userID uint) error {
	return uc.notifications.MarkAsRead(ctx, id, userID)
}

// MarkAllNotificationsAsReadUseCase は未読通知をすべて既読にする。
type MarkAllNotificationsAsReadUseCase struct {
	notifications repository.NotificationReader
}

// NewMarkAllNotificationsAsReadUseCase は MarkAllNotificationsAsReadUseCase を生成する。
func NewMarkAllNotificationsAsReadUseCase(notifications repository.NotificationReader) *MarkAllNotificationsAsReadUseCase {
	return &MarkAllNotificationsAsReadUseCase{notifications: notifications}
}

// Execute は本人の未読通知をすべて既読にする。
func (uc *MarkAllNotificationsAsReadUseCase) Execute(ctx context.Context, userID uint) error {
	return uc.notifications.MarkAllAsRead(ctx, userID)
}

// DeleteNotificationUseCase は通知 1 件を削除する。
type DeleteNotificationUseCase struct {
	notifications repository.NotificationReader
}

// NewDeleteNotificationUseCase は DeleteNotificationUseCase を生成する。
func NewDeleteNotificationUseCase(notifications repository.NotificationReader) *DeleteNotificationUseCase {
	return &DeleteNotificationUseCase{notifications: notifications}
}

// Execute は本人の通知だけを削除する。
func (uc *DeleteNotificationUseCase) Execute(ctx context.Context, id, userID uint) error {
	return uc.notifications.Delete(ctx, id, userID)
}
