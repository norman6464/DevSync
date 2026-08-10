package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetBookmarkStatsUseCase は指定ユーザーのブックマーク集計統計を取得する。
type GetBookmarkStatsUseCase struct {
	stats repository.BookmarkStatsRepository
}

// NewGetBookmarkStatsUseCase は GetBookmarkStatsUseCase を生成する。
func NewGetBookmarkStatsUseCase(stats repository.BookmarkStatsRepository) *GetBookmarkStatsUseCase {
	return &GetBookmarkStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、ブックマーク集計統計を返す。
func (uc *GetBookmarkStatsUseCase) Execute(ctx context.Context, userID uint) (*model.BookmarkStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetBookmarkStats(ctx, userID)
}
