package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// rankingRepository は [repository.RankingRepository] の GORM 実装。
type rankingRepository struct {
	db *gorm.DB
}

// NewRankingRepository は RankingRepository の GORM 実装を返す。
func NewRankingRepository(db *gorm.DB) repository.RankingRepository {
	return &rankingRepository{db: db}
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

	var entries []model.RankingEntry
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.id as user_id, u.username, u.name, u.avatar_url, COALESCE(SUM(gc.count), 0) as score
		FROM users u
		JOIN git_hub_contributions gc ON gc.user_id = u.id
		WHERE gc.date >= ?
		GROUP BY u.id
		HAVING SUM(gc.count) > 0
		ORDER BY score DESC
		LIMIT 50
	`, since).Scan(&entries).Error
	return entries, err
}

// LanguageRanking は指定プログラミング言語のバイト数ランキングを取得する。
// GitHubの言語統計データに基づき、上位50件を返す。
func (r *rankingRepository) LanguageRanking(ctx context.Context, language, period string) ([]model.RankingEntry, error) {
	var entries []model.RankingEntry
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.id as user_id, u.username, u.name, u.avatar_url, gls.bytes as score
		FROM users u
		JOIN git_hub_language_stats gls ON gls.user_id = u.id
		WHERE gls.language = ?
		ORDER BY gls.bytes DESC
		LIMIT 50
	`, language).Scan(&entries).Error
	return entries, err
}

// LevelRanking はユーザーのXP合計に基づくレベルランキングを取得する。
// 各種アクティビティから獲得したXPの合計で降順ソートし、上位50件を返す。
func (r *rankingRepository) LevelRanking(ctx context.Context) ([]model.RankingEntry, error) {
	var entries []model.RankingEntry
	err := r.db.WithContext(ctx).Raw(`
		SELECT u.id as user_id, u.username, u.name, u.avatar_url,
			COALESCE(ll.xp, 0) + COALESCE(p.xp, 0) + COALESCE(gh.xp, 0) +
			COALESCE(g.xp, 0) + COALESCE(c.xp, 0) + COALESCE(lk.xp, 0) as score
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) * 10 + COALESCE(SUM(duration), 0) / 2 as xp
			FROM learning_logs GROUP BY user_id
		) ll ON ll.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) * 30 as xp
			FROM posts GROUP BY user_id
		) p ON p.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(DISTINCT date) * 5 as xp
			FROM git_hub_contributions WHERE count > 0 GROUP BY user_id
		) gh ON gh.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) * 50 as xp
			FROM learning_goals WHERE status = 'completed' GROUP BY user_id
		) g ON g.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COUNT(*) * 5 as xp
			FROM comments GROUP BY user_id
		) c ON c.user_id = u.id
		LEFT JOIN (
			SELECT user_id, COALESCE(SUM(like_count), 0) * 3 as xp
			FROM posts GROUP BY user_id
		) lk ON lk.user_id = u.id
		HAVING score > 0
		ORDER BY score DESC
		LIMIT 50
	`).Scan(&entries).Error
	return entries, err
}

// AvailableLanguages はランキング対象となるプログラミング言語の一覧を取得する。
// GitHubの言語統計データに存在する言語をアルファベット順で返す。
func (r *rankingRepository) AvailableLanguages(ctx context.Context) ([]string, error) {
	var languages []string
	err := r.db.WithContext(ctx).Raw(`SELECT DISTINCT language FROM git_hub_language_stats ORDER BY language`).Scan(&languages).Error
	return languages, err
}
