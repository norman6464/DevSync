package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetBookReviewStatsUseCase は指定ユーザーの書籍レビュー集計統計を取得する。
type GetBookReviewStatsUseCase struct {
	stats repository.BookReviewStatsRepository
}

// NewGetBookReviewStatsUseCase は GetBookReviewStatsUseCase を生成する。
func NewGetBookReviewStatsUseCase(stats repository.BookReviewStatsRepository) *GetBookReviewStatsUseCase {
	return &GetBookReviewStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、書籍レビュー集計統計を返す。
func (uc *GetBookReviewStatsUseCase) Execute(ctx context.Context, userID uint) (*model.BookReviewStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetBookReviewStats(ctx, userID)
}
