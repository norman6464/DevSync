package persistence

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// recommendationRepository は [repository.RecommendationRepository] の sqlc(pgx) 実装。
type recommendationRepository struct {
	q *sqlcgen.Queries
}

// NewRecommendationRepository は RecommendationRepository の sqlc(pgx) 実装を返す。
func NewRecommendationRepository(q *sqlcgen.Queries) repository.RecommendationRepository {
	return &recommendationRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RecommendationRepository = (*recommendationRepository)(nil)

// GetRecommendedUsers はスキルの部分一致で候補を絞り込み、共通スキル数でスコアリングして返す。
// 自分自身とフォロー済みのユーザーは候補から除く。
func (r *recommendationRepository) GetRecommendedUsers(ctx context.Context, userID uint, skills []string, limit int) ([]model.RecommendedUser, error) {
	if len(skills) == 0 {
		return []model.RecommendedUser{}, nil
	}

	followeeIDs, err := r.q.ListFolloweeIDs(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	excludeIDs := append(followeeIDs, int64(userID))

	skillPatterns := make([]string, len(skills))
	for i, skill := range skills {
		skillPatterns[i] = "%" + escapeLikeChars(skill) + "%"
	}

	rows, err := r.q.GetRecommendedUserCandidates(ctx, sqlcgen.GetRecommendedUserCandidatesParams{
		ExcludeIds:    excludeIDs,
		SkillPatterns: skillPatterns,
	})
	if err != nil {
		return nil, err
	}

	skillSet := make(map[string]bool, len(skills))
	for _, s := range skills {
		skillSet[strings.TrimSpace(s)] = true
	}

	var results []model.RecommendedUser
	for _, row := range rows {
		candidate := toModelUser(row)
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
	rows, err := r.q.ListTrendingPosts(ctx, sqlcgen.ListTrendingPostsParams{
		Days:  int32Param(days),
		Limit: int32Param(limit),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, len(rows))
	postIDs := make([]int64, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
		postIDs[i] = row.Post.ID
	}

	if len(postIDs) > 0 {
		snippetRows, err := r.q.ListCodeSnippetsByPostIDs(ctx, postIDs)
		if err != nil {
			return nil, err
		}
		snippetsByPostID := make(map[uint][]model.CodeSnippet)
		for _, row := range snippetRows {
			postID := uint(row.PostID)
			snippetsByPostID[postID] = append(snippetsByPostID[postID], toModelCodeSnippet(row))
		}
		for i := range posts {
			posts[i].CodeSnippets = snippetsByPostID[posts[i].ID]
		}
	}
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func toModelLearningResource(row sqlcgen.LearningResource) model.LearningResource {
	return model.LearningResource{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		URL:         fromStringPtr(row.Url),
		Category:    model.ResourceCategory(row.Category),
		Difficulty:  model.ResourceDifficulty(fromStringPtr(row.Difficulty)),
		Tags:        fromStringPtr(row.Tags),
		ImageURL:    fromStringPtr(row.ImageUrl),
		IsPublic:    row.IsPublic,
		LikeCount:   int(fromInt64PtrValue(row.LikeCount)),
		SaveCount:   int(fromInt64PtrValue(row.SaveCount)),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// GetTrendingResources は直近 days 日の公開リソースを、いいね数 + 保存数の降順で返す。
func (r *recommendationRepository) GetTrendingResources(ctx context.Context, limit, days int) ([]model.LearningResource, error) {
	rows, err := r.q.ListTrendingResources(ctx, sqlcgen.ListTrendingResourcesParams{
		Days:  int32Param(days),
		Limit: int32Param(limit),
	})
	if err != nil {
		return nil, err
	}

	resources := make([]model.LearningResource, len(rows))
	for i, row := range rows {
		resources[i] = toModelLearningResource(row.LearningResource)
		resources[i].User = toModelUser(row.User)
	}
	return resources, nil
}
