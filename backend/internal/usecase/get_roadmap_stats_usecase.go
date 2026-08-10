package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetRoadmapStatsUseCase は指定ユーザーのロードマップ統計を取得する。
type GetRoadmapStatsUseCase struct {
	stats repository.RoadmapStatsRepository
}

// NewGetRoadmapStatsUseCase は GetRoadmapStatsUseCase を生成する。
func NewGetRoadmapStatsUseCase(stats repository.RoadmapStatsRepository) *GetRoadmapStatsUseCase {
	return &GetRoadmapStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、ロードマップ統計を返す。
func (uc *GetRoadmapStatsUseCase) Execute(ctx context.Context, userID uint) (*model.RoadmapStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetRoadmapStats(ctx, userID)
}
