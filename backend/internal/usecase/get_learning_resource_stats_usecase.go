package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetLearningResourceStatsUseCase は指定ユーザーの学習リソース活動集計統計を取得する。
type GetLearningResourceStatsUseCase struct {
	stats repository.LearningResourceStatsRepository
}

// NewGetLearningResourceStatsUseCase は GetLearningResourceStatsUseCase を生成する。
func NewGetLearningResourceStatsUseCase(stats repository.LearningResourceStatsRepository) *GetLearningResourceStatsUseCase {
	return &GetLearningResourceStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、学習リソース活動集計統計を返す。
func (uc *GetLearningResourceStatsUseCase) Execute(ctx context.Context, userID uint) (*model.LearningResourceStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetLearningResourceStats(ctx, userID)
}
