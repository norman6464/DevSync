package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationStatsRepository はユーザー通知集計統計の取得に対する、usecase 側が要求する契約。
type NotificationStatsRepository interface {
	GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error)
}
