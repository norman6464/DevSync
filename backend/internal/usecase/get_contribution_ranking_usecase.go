package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetContributionRankingUseCase は指定期間のコントリビューションランキングを取得する。
type GetContributionRankingUseCase struct {
	ranking repository.RankingRepository
}

// NewGetContributionRankingUseCase は GetContributionRankingUseCase を生成する。
func NewGetContributionRankingUseCase(ranking repository.RankingRepository) *GetContributionRankingUseCase {
	return &GetContributionRankingUseCase{ranking: ranking}
}

// Execute は期間を検証し、コントリビューションランキングを返す。
func (uc *GetContributionRankingUseCase) Execute(ctx context.Context, period string) ([]model.RankingEntry, error) {
	if err := validateRankingPeriod(period); err != nil {
		return nil, err
	}
	return uc.ranking.ContributionRanking(ctx, period)
}
