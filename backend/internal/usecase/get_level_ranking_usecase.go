package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetLevelRankingUseCase は XP 合計に基づくレベルランキングを取得する。
type GetLevelRankingUseCase struct {
	ranking repository.RankingRepository
}

// NewGetLevelRankingUseCase は GetLevelRankingUseCase を生成する。
func NewGetLevelRankingUseCase(ranking repository.RankingRepository) *GetLevelRankingUseCase {
	return &GetLevelRankingUseCase{ranking: ranking}
}

// Execute はレベルランキングを返す。
func (uc *GetLevelRankingUseCase) Execute(ctx context.Context) ([]model.RankingEntry, error) {
	return uc.ranking.LevelRanking(ctx)
}
