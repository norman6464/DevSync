package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetMessageStatsUseCase は指定ユーザーのメッセージ集計統計を取得する。
type GetMessageStatsUseCase struct {
	stats repository.MessageStatsRepository
}

// NewGetMessageStatsUseCase は GetMessageStatsUseCase を生成する。
func NewGetMessageStatsUseCase(stats repository.MessageStatsRepository) *GetMessageStatsUseCase {
	return &GetMessageStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、メッセージ集計統計を返す。
func (uc *GetMessageStatsUseCase) Execute(ctx context.Context, userID uint) (*model.MessageStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetMessageStats(ctx, userID)
}
