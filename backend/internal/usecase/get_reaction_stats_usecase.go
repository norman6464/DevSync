package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetReactionStatsUseCase は指定ユーザーが受け取ったリアクション集計統計を取得する。
type GetReactionStatsUseCase struct {
	stats repository.ReactionStatsRepository
}

// NewGetReactionStatsUseCase は GetReactionStatsUseCase を生成する。
func NewGetReactionStatsUseCase(stats repository.ReactionStatsRepository) *GetReactionStatsUseCase {
	return &GetReactionStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、リアクション集計統計を返す。
func (uc *GetReactionStatsUseCase) Execute(ctx context.Context, userID uint) (*model.ReactionStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetReactionStats(ctx, userID)
}
