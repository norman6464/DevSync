package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListRankingLanguagesUseCase はランキング対象の言語一覧を取得する。
type ListRankingLanguagesUseCase struct {
	ranking repository.RankingRepository
}

// NewListRankingLanguagesUseCase は ListRankingLanguagesUseCase を生成する。
func NewListRankingLanguagesUseCase(ranking repository.RankingRepository) *ListRankingLanguagesUseCase {
	return &ListRankingLanguagesUseCase{ranking: ranking}
}

// Execute は利用可能な言語一覧を返す。
func (uc *ListRankingLanguagesUseCase) Execute(ctx context.Context) ([]string, error) {
	return uc.ranking.AvailableLanguages(ctx)
}
