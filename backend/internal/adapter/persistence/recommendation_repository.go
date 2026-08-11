package persistence

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// recommendationRepository は [repository.RecommendationRepository] の GORM 実装。
type recommendationRepository struct {
	db *gorm.DB
}

// NewRecommendationRepository は RecommendationRepository の GORM 実装を返す。
func NewRecommendationRepository(db *gorm.DB) repository.RecommendationRepository {
	return &recommendationRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RecommendationRepository = (*recommendationRepository)(nil)

// GetRecommendedUsers はスキルの部分一致で候補を絞り込み、共通スキル数でスコアリングして返す。
// 自分自身とフォロー済みのユーザーは候補から除く。
func (r *recommendationRepository) GetRecommendedUsers(ctx context.Context, userID uint, skills []string, limit int) ([]model.RecommendedUser, error) {
	db := r.db.WithContext(ctx)

	var followingIDs []uint
	db.Raw(`SELECT followee_id FROM follows WHERE follower_id = ?`, userID).Scan(&followingIDs)
	excludeIDs := append(followingIDs, userID)

	var conditions []string
	var args []interface{}
	for _, skill := range skills {
		conditions = append(conditions, "skills_languages LIKE ? OR skills_frameworks LIKE ?")
		pattern := "%" + escapeLikeChars(skill) + "%"
		args = append(args, pattern, pattern)
	}
	if len(conditions) == 0 {
		return []model.RecommendedUser{}, nil
	}

	args = append(args, excludeIDs)
	var candidates []model.User
	if err := db.Where("("+strings.Join(conditions, " OR ")+") AND id NOT IN ?", args...).
		Find(&candidates).Error; err != nil {
		return nil, err
	}

	skillSet := make(map[string]bool, len(skills))
	for _, s := range skills {
		skillSet[strings.TrimSpace(s)] = true
	}

	var results []model.RecommendedUser
	for _, candidate := range candidates {
		commonSkills := matchedSkills(candidate, skillSet)
		if len(commonSkills) == 0 {
			continue
		}
		results = append(results, model.RecommendedUser{
			User:         candidate,
			CommonSkills: commonSkills,
			MatchScore:   len(commonSkills),
		})
	}

	sortByMatchScore(results)
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

// matchedSkills は候補ユーザーのスキルのうち、指定集合に含まれるものを返す。
func matchedSkills(candidate model.User, skillSet map[string]bool) []string {
	candidateSkills := splitCSV(candidate.SkillsLanguages)
	candidateSkills = append(candidateSkills, splitCSV(candidate.SkillsFrameworks)...)

	var matched []string
	for _, cs := range candidateSkills {
		trimmed := strings.TrimSpace(cs)
		if trimmed != "" && skillSet[trimmed] {
			matched = append(matched, trimmed)
		}
	}
	return matched
}

// splitCSV はカンマ区切り文字列を分割する。
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// sortByMatchScore は共通スキル数の降順に並べ替える。
// 同点の並びは安定しないが、移行前の実装と同じ順序になるようアルゴリズムを変えていない。
func sortByMatchScore(results []model.RecommendedUser) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].MatchScore > results[i].MatchScore {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

// GetTrendingPosts は直近 days 日の投稿を、いいね数 + コメント数の降順で返す。
func (r *recommendationRepository) GetTrendingPosts(ctx context.Context, limit, days int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("CodeSnippets").
		Where("created_at > NOW() - INTERVAL '1 day' * ?", days).
		Order("(like_count + comment_count) DESC").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// GetTrendingResources は直近 days 日の公開リソースを、いいね数 + 保存数の降順で返す。
func (r *recommendationRepository) GetTrendingResources(ctx context.Context, limit, days int) ([]model.LearningResource, error) {
	var resources []model.LearningResource
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("is_public = ? AND created_at > NOW() - INTERVAL '1 day' * ?", true, days).
		Order("(like_count + save_count) DESC").
		Limit(limit).
		Find(&resources).Error
	return resources, err
}
