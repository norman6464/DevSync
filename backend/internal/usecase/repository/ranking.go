package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// RankingRepository はユーザーランキングの集計に対する、usecase 側が要求する契約。
type RankingRepository interface {
	ContributionRanking(ctx context.Context, period string) ([]model.RankingEntry, error)
	LanguageRanking(ctx context.Context, language, period string) ([]model.RankingEntry, error)
	LevelRanking(ctx context.Context) ([]model.RankingEntry, error)
	AvailableLanguages(ctx context.Context) ([]string, error)
}
