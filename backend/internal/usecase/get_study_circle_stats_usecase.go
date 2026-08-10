package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetStudyCircleStatsUseCase は指定サークルの集計統計を取得する。
type GetStudyCircleStatsUseCase struct {
	stats repository.StudyCircleStatsRepository
}

// NewGetStudyCircleStatsUseCase は GetStudyCircleStatsUseCase を生成する。
func NewGetStudyCircleStatsUseCase(stats repository.StudyCircleStatsRepository) *GetStudyCircleStatsUseCase {
	return &GetStudyCircleStatsUseCase{stats: stats}
}

// Execute はサークル ID を検証し、サークルの集計統計を返す。
func (uc *GetStudyCircleStatsUseCase) Execute(ctx context.Context, circleID uint) (*model.StudyCircleStats, error) {
	if err := domain.ValidateRequiredID(circleID, "circleID"); err != nil {
		return nil, err
	}
	return uc.stats.GetCircleStats(ctx, circleID)
}
