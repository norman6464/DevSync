package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// notificationStatsRepository は [repository.NotificationStatsRepository] の GORM 実装。
type notificationStatsRepository struct {
	db *gorm.DB
}

// NewNotificationStatsRepository は NotificationStatsRepository の GORM 実装を返す。
func NewNotificationStatsRepository(db *gorm.DB) repository.NotificationStatsRepository {
	return &notificationStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NotificationStatsRepository = (*notificationStatsRepository)(nil)

// GetNotificationStats は指定ユーザーの通知集計統計を返す。
func (r *notificationStatsRepository) GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.NotificationStats

	// 通知総数
	if err := db.Model(&model.Notification{}).Where("user_id = ?", userID).Count(&stats.TotalNotifications).Error; err != nil {
		return nil, err
	}

	// 未読通知数
	if err := db.Model(&model.Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&stats.UnreadCount).Error; err != nil {
		return nil, err
	}

	// 今月の通知数
	startOfMonth := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.Notification{}).Where("user_id = ? AND created_at >= ?", userID, startOfMonth).Count(&stats.NotificationsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
