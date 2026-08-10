package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetLearningLogStatsUseCase は指定ユーザーの学習ログ集計統計を取得する。
type GetLearningLogStatsUseCase struct {
	stats repository.LearningLogStatsRepository
}

// NewGetLearningLogStatsUseCase は GetLearningLogStatsUseCase を生成する。
func NewGetLearningLogStatsUseCase(stats repository.LearningLogStatsRepository) *GetLearningLogStatsUseCase {
	return &GetLearningLogStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、学習ログ集計統計を返す。
func (uc *GetLearningLogStatsUseCase) Execute(ctx context.Context, userID uint) (*model.LearningLogStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetLearningLogStats(ctx, userID)
}
