package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// notificationStatsRepository は [repository.NotificationStatsRepository] の sqlc(pgx) 実装。
type notificationStatsRepository struct {
	q *sqlcgen.Queries
}

// NewNotificationStatsRepository は NotificationStatsRepository の sqlc(pgx) 実装を返す。
func NewNotificationStatsRepository(q *sqlcgen.Queries) repository.NotificationStatsRepository {
	return &notificationStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NotificationStatsRepository = (*notificationStatsRepository)(nil)

// GetNotificationStats は指定ユーザーの通知集計統計を返す。
func (r *notificationStatsRepository) GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	total, err := r.q.CountNotificationsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	unread, err := r.q.CountUnreadNotificationsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	startOfMonth := domain.StartOfMonth(time.Now())
	thisMonth, err := r.q.CountNotificationsByUserSince(ctx, sqlcgen.CountNotificationsByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startOfMonth),
	})
	if err != nil {
		return nil, err
	}

	return &model.NotificationStats{
		TotalNotifications:     total,
		UnreadCount:            unread,
		NotificationsThisMonth: thisMonth,
	}, nil
}
