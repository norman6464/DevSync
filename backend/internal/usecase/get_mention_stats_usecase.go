package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetMentionStatsUseCase は指定ユーザーのメンション集計統計を取得する。
type GetMentionStatsUseCase struct {
	stats repository.MentionStatsRepository
}

// NewGetMentionStatsUseCase は GetMentionStatsUseCase を生成する。
func NewGetMentionStatsUseCase(stats repository.MentionStatsRepository) *GetMentionStatsUseCase {
	return &GetMentionStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、メンション集計統計を返す。
func (uc *GetMentionStatsUseCase) Execute(ctx context.Context, userID uint) (*model.MentionStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetMentionStats(ctx, userID)
}
