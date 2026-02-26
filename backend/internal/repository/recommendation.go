package repository

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// RecommendationRepository はレコメンド機能のデータ取得を提供するリポジトリ実装。
type RecommendationRepository struct {
	db *gorm.DB
}

// NewRecommendationRepository は新しいRecommendationRepositoryインスタンスを生成する。
func NewRecommendationRepository(db *gorm.DB) *RecommendationRepository {
	return &RecommendationRepository{db: db}
}

// GetRecommendedUsers はスキルマッチングに基づくおすすめユーザーを返す。
// 自分自身とフォロー済みユーザーを除外し、共通スキル数でスコアリングする。
func (r *RecommendationRepository) GetRecommendedUsers(userID uint, skills []string, limit int) ([]model.RecommendedUser, error) {
	// フォロー済みユーザーIDを取得
	var followingIDs []uint
	r.db.Raw(`SELECT followee_id FROM follows WHERE follower_id = ?`, userID).Scan(&followingIDs)

	// 除外IDリスト（自分 + フォロー済み）
	excludeIDs := append(followingIDs, userID)

	// スキルの部分一致でフィルタリング用のOR条件を構築
	var conditions []string
	var args []interface{}
	for _, skill := range skills {
		conditions = append(conditions, "skills_languages LIKE ? OR skills_frameworks LIKE ?")
		pattern := "%" + EscapeLikeChars(skill) + "%"
		args = append(args, pattern, pattern)
	}

	if len(conditions) == 0 {
		return []model.RecommendedUser{}, nil
	}

	whereClause := "(" + strings.Join(conditions, " OR ") + ") AND id NOT IN ?"
	args = append(args, excludeIDs)

	var candidates []model.User
	if err := r.db.Where(whereClause, args...).Find(&candidates).Error; err != nil {
		return nil, err
	}

	// Go側で正確なマッチングとスコア計算
	skillSet := make(map[string]bool)
	for _, s := range skills {
		skillSet[strings.TrimSpace(s)] = true
	}

	var results []model.RecommendedUser
	for _, candidate := range candidates {
		var commonSkills []string
		candidateSkills := parseCSV(candidate.SkillsLanguages)
		candidateSkills = append(candidateSkills, parseCSV(candidate.SkillsFrameworks)...)

		for _, cs := range candidateSkills {
			trimmed := strings.TrimSpace(cs)
			if trimmed != "" && skillSet[trimmed] {
				commonSkills = append(commonSkills, trimmed)
			}
		}

		if len(commonSkills) > 0 {
			results = append(results, model.RecommendedUser{
				User:         candidate,
				CommonSkills: commonSkills,
				MatchScore:   len(commonSkills),
			})
		}
	}

	// スコアの降順でソート
	sortByScore(results)

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetTrendingPosts は直近指定日数内の人気投稿を返す。
func (r *RecommendationRepository) GetTrendingPosts(limit int, days int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.
		Preload("User").
		Preload("CodeSnippets").
		Where("created_at > NOW() - INTERVAL '1 day' * ?", days).
		Order("(like_count + comment_count) DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// GetTrendingResources は直近指定日数内の人気学習リソースを返す。
func (r *RecommendationRepository) GetTrendingResources(limit int, days int) ([]model.LearningResource, error) {
	var resources []model.LearningResource
	err := r.db.
		Preload("User").
		Where("is_public = ? AND created_at > NOW() - INTERVAL '1 day' * ?", true, days).
		Order("(like_count + save_count) DESC").
		Limit(limit).
		Find(&resources).Error
	return resources, err
}

// parseCSV はカンマ区切り文字列をスライスに変換する。
func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// sortByScore はRecommendedUserをスコアの降順でソートする。
func sortByScore(results []model.RecommendedUser) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].MatchScore > results[i].MatchScore {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}
