package repository

import "gorm.io/gorm"

type RankingRepository struct {
	db *gorm.DB
}

func NewRankingRepository(db *gorm.DB) *RankingRepository {
	return &RankingRepository{db: db}
}

type RankingEntry struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Score     int64  `json:"score"`
}

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

func (r *RankingRepository) AvailableLanguages() ([]string, error) {
	var languages []string
	err := r.db.Raw(`SELECT DISTINCT language FROM git_hub_language_stats ORDER BY language`).Scan(&languages).Error
	return languages, err
}
