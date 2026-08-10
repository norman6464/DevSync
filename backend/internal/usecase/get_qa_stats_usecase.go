package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetQAStatsUseCase は指定ユーザーの Q&A 活動集計統計を取得する。
type GetQAStatsUseCase struct {
	stats repository.QAStatsRepository
}

// NewGetQAStatsUseCase は GetQAStatsUseCase を生成する。
func NewGetQAStatsUseCase(stats repository.QAStatsRepository) *GetQAStatsUseCase {
	return &GetQAStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、Q&A 活動集計統計を返す。
func (uc *GetQAStatsUseCase) Execute(ctx context.Context, userID uint) (*model.QAStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetQAStats(ctx, userID)
}
