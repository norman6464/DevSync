package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetLanguageRankingUseCase は指定言語のランキングを取得する。
type GetLanguageRankingUseCase struct {
	ranking repository.RankingRepository
}

// NewGetLanguageRankingUseCase は GetLanguageRankingUseCase を生成する。
func NewGetLanguageRankingUseCase(ranking repository.RankingRepository) *GetLanguageRankingUseCase {
	return &GetLanguageRankingUseCase{ranking: ranking}
}

// Execute は期間と言語名を検証し、言語別ランキングを返す。
func (uc *GetLanguageRankingUseCase) Execute(ctx context.Context, language, period string) ([]model.RankingEntry, error) {
	if err := validateRankingPeriod(period); err != nil {
		return nil, err
	}
	if len(language) > maxLanguageNameLength {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "言語名が長すぎます", nil)
	}
	return uc.ranking.LanguageRanking(ctx, language, period)
}
