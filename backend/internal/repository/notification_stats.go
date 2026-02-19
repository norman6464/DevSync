package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NotificationStatsRepository はユーザー通知集計統計の取得を担当するリポジトリ実装。
type NotificationStatsRepository struct {
	db *gorm.DB
}

// NewNotificationStatsRepository は新しいNotificationStatsRepositoryインスタンスを生成する。
func NewNotificationStatsRepository(db *gorm.DB) *NotificationStatsRepository {
	return &NotificationStatsRepository{db: db}
}

// GetNotificationStats は指定ユーザーの通知集計統計を返す。
func (r *NotificationStatsRepository) GetNotificationStats(userID uint) (*model.NotificationStats, error) {
	var stats model.NotificationStats

	// 通知総数
	if err := r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Count(&stats.TotalNotifications).Error; err != nil {
		return nil, err
	}

	// 未読通知数
	if err := r.db.Model(&model.Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&stats.UnreadCount).Error; err != nil {
		return nil, err
	}

	// 今月の通知数
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := r.db.Model(&model.Notification{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.NotificationsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
