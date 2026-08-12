package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// notificationRepository は [repository.NotificationReader] の GORM 実装。
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository は NotificationReader の GORM 実装を返す。
func NewNotificationRepository(db *gorm.DB) repository.NotificationReader {
	return &notificationRepository{db: db}
}

// NewNotificationCreator は NotificationCreator の GORM 実装を返す。
// 通知の作成だけを必要とする利用者（バッジ獲得など）はこちらを受け取る。
func NewNotificationCreator(db *gorm.DB) repository.NotificationCreator {
	return &notificationRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NotificationReader = (*notificationRepository)(nil)
var _ repository.NotificationCreator = (*notificationRepository)(nil)

// Create は通知を 1 件保存する。
func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// FindByUserID は指定ユーザーの通知を作成日時の降順で取得する。
// 通知の表示に必要な関連（実行者・投稿・質問）を Preload する。
func (r *notificationRepository) FindByUserID(ctx context.Context, userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	var notifications []model.Notification
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).
		Preload("Actor").Preload("Post").Preload("Question").
		Where("user_id = ?", userID)
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	err := query.Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// CountByUserID は指定ユーザーの通知総数を取得する。
func (r *notificationRepository) CountByUserID(ctx context.Context, userID uint, notificationType string) (int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID)
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// CountUnread は指定ユーザーの未読通知数を取得する。
func (r *notificationRepository) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead は指定 ID の通知を既読にする。本人の通知だけを対象にする。
func (r *notificationRepository) MarkAsRead(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error
}

// MarkAllAsRead は指定ユーザーの未読通知をすべて既読にする。
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error
}

// Delete は指定 ID の通知を削除する。本人の通知だけを対象にする。
func (r *notificationRepository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Notification{}).Error
}
