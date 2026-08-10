package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetNotificationStatsUseCase は指定ユーザーの通知集計統計を取得する。
type GetNotificationStatsUseCase struct {
	stats repository.NotificationStatsRepository
}

// NewGetNotificationStatsUseCase は GetNotificationStatsUseCase を生成する。
func NewGetNotificationStatsUseCase(stats repository.NotificationStatsRepository) *GetNotificationStatsUseCase {
	return &GetNotificationStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、通知集計統計を返す。
func (uc *GetNotificationStatsUseCase) Execute(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetNotificationStats(ctx, userID)
}
