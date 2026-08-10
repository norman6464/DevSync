package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetCommentStatsUseCase は指定ユーザーのコメント活動集計統計を取得する。
type GetCommentStatsUseCase struct {
	stats repository.CommentStatsRepository
}

// NewGetCommentStatsUseCase は GetCommentStatsUseCase を生成する。
func NewGetCommentStatsUseCase(stats repository.CommentStatsRepository) *GetCommentStatsUseCase {
	return &GetCommentStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、コメント活動集計統計を返す。
func (uc *GetCommentStatsUseCase) Execute(ctx context.Context, userID uint) (*model.CommentStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetCommentStats(ctx, userID)
}
