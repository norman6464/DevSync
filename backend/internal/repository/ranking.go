package repository

import "gorm.io/gorm"

// RankingRepository はユーザーランキングデータの集計・取得を提供するリポジトリ実装。
type RankingRepository struct {
	db *gorm.DB
}

// NewRankingRepository は新しいRankingRepositoryインスタンスを生成する。
func NewRankingRepository(db *gorm.DB) *RankingRepository {
	return &RankingRepository{db: db}
}

// RankingEntry はランキング表示用の1エントリを表す。
type RankingEntry struct {
	UserID    uint   `json:"user_id"`    // ユーザーID
	Name      string `json:"name"`       // ユーザー名
	AvatarURL string `json:"avatar_url"` // アバター画像URL
	Score     int64  `json:"score"`      // ランキングスコア
}

// ContributionRanking は指定期間（weekly/monthly）のGitHubコントリビューションランキングを取得する。
// コントリビューション数の合計で降順ソートし、上位50件を返す。
func (r *RankingRepository) ContributionRanking(period string) ([]RankingEntry, error) {
	interval := "7 days"
	if period == "monthly" {
		interval = "30 days"
	}

	var entries []RankingEntry
	err := r.db.Raw(`
		SELECT u.id as user_id, u.name, u.avatar_url, COALESCE(SUM(gc.count), 0) as score
		FROM users u
		JOIN git_hub_contributions gc ON gc.user_id = u.id
		WHERE gc.date >= NOW() - INTERVAL '`+interval+`'
		GROUP BY u.id
		HAVING SUM(gc.count) > 0
		ORDER BY score DESC
		LIMIT 50
	`).Scan(&entries).Error
	return entries, err
}

// LanguageRanking は指定プログラミング言語のバイト数ランキングを取得する。
// GitHubの言語統計データに基づき、上位50件を返す。
func (r *RankingRepository) LanguageRanking(language, period string) ([]RankingEntry, error) {
	var entries []RankingEntry
	err := r.db.Raw(`
		SELECT u.id as user_id, u.name, u.avatar_url, gls.bytes as score
		FROM users u
		JOIN git_hub_language_stats gls ON gls.user_id = u.id
		WHERE gls.language = ?
		ORDER BY gls.bytes DESC
		LIMIT 50
	`, language).Scan(&entries).Error
	return entries, err
}

// AvailableLanguages はランキング対象となるプログラミング言語の一覧を取得する。
// GitHubの言語統計データに存在する言語をアルファベット順で返す。
func (r *RankingRepository) AvailableLanguages() ([]string, error) {
	var languages []string
	err := r.db.Raw(`SELECT DISTINCT language FROM git_hub_language_stats ORDER BY language`).Scan(&languages).Error
	return languages, err
}
