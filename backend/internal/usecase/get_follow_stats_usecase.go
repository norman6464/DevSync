package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetFollowStatsUseCase は指定ユーザーのフォロー関係集計統計を取得する。
type GetFollowStatsUseCase struct {
	stats repository.FollowStatsRepository
}

// NewGetFollowStatsUseCase は GetFollowStatsUseCase を生成する。
func NewGetFollowStatsUseCase(stats repository.FollowStatsRepository) *GetFollowStatsUseCase {
	return &GetFollowStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、フォロー関係集計統計を返す。
func (uc *GetFollowStatsUseCase) Execute(ctx context.Context, userID uint) (*model.FollowStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetFollowStats(ctx, userID)
}
