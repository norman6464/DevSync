package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetCodeSnippetStatsUseCase は指定ユーザーのコードスニペット活動集計統計を取得する。
type GetCodeSnippetStatsUseCase struct {
	stats repository.CodeSnippetStatsRepository
}

// NewGetCodeSnippetStatsUseCase は GetCodeSnippetStatsUseCase を生成する。
func NewGetCodeSnippetStatsUseCase(stats repository.CodeSnippetStatsRepository) *GetCodeSnippetStatsUseCase {
	return &GetCodeSnippetStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、コードスニペット活動集計統計を返す。
func (uc *GetCodeSnippetStatsUseCase) Execute(ctx context.Context, userID uint) (*model.CodeSnippetStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetCodeSnippetStats(ctx, userID)
}
