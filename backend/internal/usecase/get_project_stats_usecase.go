package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetProjectStatsUseCase は指定ユーザーのプロジェクト活動集計統計を取得する。
type GetProjectStatsUseCase struct {
	stats repository.ProjectStatsRepository
}

// NewGetProjectStatsUseCase は GetProjectStatsUseCase を生成する。
func NewGetProjectStatsUseCase(stats repository.ProjectStatsRepository) *GetProjectStatsUseCase {
	return &GetProjectStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、プロジェクト活動集計統計を返す。
func (uc *GetProjectStatsUseCase) Execute(ctx context.Context, userID uint) (*model.ProjectStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetProjectStats(ctx, userID)
}
