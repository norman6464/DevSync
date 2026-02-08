package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NotificationRepository は通知データへのアクセスを提供するリポジトリ実装。
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository は新しいNotificationRepositoryインスタンスを生成する。
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create は新しい通知をデータベースに作成する。
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// CreateBatch は複数の通知を一括で作成する。
func (r *NotificationRepository) CreateBatch(notifications []*model.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.db.Create(&notifications).Error
}

// FindByUserID は指定ユーザーの通知をページネーション付きで取得する。
// notificationTypeが指定された場合、その種類でフィルタする。
func (r *NotificationRepository) FindByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	var notifications []model.Notification
	offset := (page - 1) * limit
	query := r.db.Preload("Actor").Preload("Post").Preload("Question").
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
func (r *NotificationRepository) CountByUserID(userID uint, notificationType string) (int64, error) {
	var count int64
	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}
	err := query.Count(&count).Error
	return count, err
}

// CountUnread は指定ユーザーの未読通知数を取得する。
func (r *NotificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead は指定IDの通知を既読にする。
func (r *NotificationRepository) MarkAsRead(id, userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error
}

// MarkAllAsRead は指定ユーザーの全未読通知を既読にする。
func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error
}

// Delete は指定IDの通知を削除する。
func (r *NotificationRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Notification{}).Error
}

// GetFollowerIDs は指定ユーザーのフォロワーID一覧を取得する。
// 通知の一括送信時に使用される。
func (r *NotificationRepository) GetFollowerIDs(userID uint) ([]uint, error) {
	var followerIDs []uint
	err := r.db.Model(&model.Follow{}).
		Where("followee_id = ?", userID).
		Pluck("follower_id", &followerIDs).Error
	return followerIDs, err
}
