package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// rankingRepository は [repository.RankingRepository] の sqlc(pgx) 実装。
type rankingRepository struct {
	q *sqlcgen.Queries
}

// NewRankingRepository は RankingRepository の sqlc(pgx) 実装を返す。
func NewRankingRepository(q *sqlcgen.Queries) repository.RankingRepository {
	return &rankingRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RankingRepository = (*rankingRepository)(nil)

// ContributionRanking は指定期間（weekly/monthly）のGitHubコントリビューションランキングを取得する。
// コントリビューション数の合計で降順ソートし、上位50件を返す。
func (r *rankingRepository) ContributionRanking(ctx context.Context, period string) ([]model.RankingEntry, error) {
	days := 7
	if period == "monthly" {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	rows, err := r.q.GetContributionRanking(ctx, toTimestamptzNotNull(since))
	if err != nil {
		return nil, err
	}
	entries := make([]model.RankingEntry, len(rows))
	for i, row := range rows {
		entries[i] = model.RankingEntry{
			UserID:    uint(row.UserID),
			Username:  row.Username,
			Name:      row.Name,
			AvatarURL: fromStringPtr(row.AvatarUrl),
			Score:     row.Score,
		}
	}
	return entries, nil
}

// LanguageRanking は指定プログラミング言語のバイト数ランキングを取得する。
// GitHubの言語統計データに基づき、上位50件を返す。
func (r *rankingRepository) LanguageRanking(ctx context.Context, language, period string) ([]model.RankingEntry, error) {
	rows, err := r.q.GetLanguageRanking(ctx, language)
	if err != nil {
		return nil, err
	}
	entries := make([]model.RankingEntry, len(rows))
	for i, row := range rows {
		entries[i] = model.RankingEntry{
			UserID:    uint(row.UserID),
			Username:  row.Username,
			Name:      row.Name,
			AvatarURL: fromStringPtr(row.AvatarUrl),
			Score:     row.Score,
		}
	}
	return entries, nil
}

// LevelRanking はユーザーのXP合計に基づくレベルランキングを取得する。
// 各種アクティビティから獲得したXPの合計で降順ソートし、上位50件を返す。
func (r *rankingRepository) LevelRanking(ctx context.Context) ([]model.RankingEntry, error) {
	rows, err := r.q.GetLevelRanking(ctx)
	if err != nil {
		return nil, err
	}
	entries := make([]model.RankingEntry, len(rows))
	for i, row := range rows {
		entries[i] = model.RankingEntry{
			UserID:    uint(row.UserID),
			Username:  row.Username,
			Name:      row.Name,
			AvatarURL: fromStringPtr(row.AvatarUrl),
			Score:     int64(row.Score),
		}
	}
	return entries, nil
}

// AvailableLanguages はランキング対象となるプログラミング言語の一覧を取得する。
// GitHubの言語統計データに存在する言語をアルファベット順で返す。
func (r *rankingRepository) AvailableLanguages(ctx context.Context) ([]string, error) {
	return r.q.ListAvailableLanguages(ctx)
}
