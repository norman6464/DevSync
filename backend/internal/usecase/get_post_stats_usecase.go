package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetPostStatsUseCase は指定ユーザーの投稿集計統計を取得する。
type GetPostStatsUseCase struct {
	stats repository.PostStatsRepository
}

// NewGetPostStatsUseCase は GetPostStatsUseCase を生成する。
func NewGetPostStatsUseCase(stats repository.PostStatsRepository) *GetPostStatsUseCase {
	return &GetPostStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、投稿集計統計を返す。
func (uc *GetPostStatsUseCase) Execute(ctx context.Context, userID uint) (*model.PostStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetPostStats(ctx, userID)
}
