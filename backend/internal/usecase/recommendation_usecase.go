package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

const (
	// recommendedUserLimit はおすすめユーザーの最大件数。
	recommendedUserLimit = 10
	// trendingPostLimit / trendingPostDays は人気投稿の件数と対象期間。
	trendingPostLimit = 10
	trendingPostDays  = 7
	// trendingResourceLimit / trendingResourceDays は人気学習リソースの件数と対象期間。
	trendingResourceLimit = 10
	trendingResourceDays  = 30
)

// errRecommendationUserNotFound は対象ユーザーが見つからないときに返すエラー。
// DomainError ではないため handler では 500 になり、不在を素の DB エラーとして扱っていた
// 移行前の挙動と一致する。
var errRecommendationUserNotFound = errors.New("ユーザーが見つかりません")

// GetRecommendedUsersUseCase はスキルマッチングでおすすめユーザーを取得する。
type GetRecommendedUsersUseCase struct {
	recommendations repository.RecommendationRepository
	users           repository.UserSkillsReader
}

// NewGetRecommendedUsersUseCase は GetRecommendedUsersUseCase を生成する。
func NewGetRecommendedUsersUseCase(
	recommendations repository.RecommendationRepository,
	users repository.UserSkillsReader,
) *GetRecommendedUsersUseCase {
	return &GetRecommendedUsersUseCase{recommendations: recommendations, users: users}
}

// Execute はプロフィールのスキルを分解し、共通点を持つユーザーを返す。
// スキルが 1 つも無ければ空の一覧を返す。
func (uc *GetRecommendedUsersUseCase) Execute(ctx context.Context, userID uint) ([]model.RecommendedUser, error) {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errRecommendationUserNotFound
	}

	skills := ParseSkills(user.SkillsLanguages, user.SkillsFrameworks)
	if len(skills) == 0 {
		return []model.RecommendedUser{}, nil
	}
	return uc.recommendations.GetRecommendedUsers(ctx, userID, skills, recommendedUserLimit)
}

// ParseSkills はカンマ区切りの言語・フレームワークを 1 つのスキル一覧へ分解する純粋関数。
// 前後の空白（スペースとタブ）を落とし、空の要素は捨てる。
func ParseSkills(languages, frameworks string) []string {
	var skills []string
	for _, source := range []string{languages, frameworks} {
		for _, s := range strings.Split(source, ",") {
			if trimmed := strings.Trim(s, " \t"); trimmed != "" {
				skills = append(skills, trimmed)
			}
		}
	}
	return skills
}

// GetTrendingPostsUseCase は人気投稿を取得する。
type GetTrendingPostsUseCase struct {
	recommendations repository.RecommendationRepository
}

// NewGetTrendingPostsUseCase は GetTrendingPostsUseCase を生成する。
func NewGetTrendingPostsUseCase(recommendations repository.RecommendationRepository) *GetTrendingPostsUseCase {
	return &GetTrendingPostsUseCase{recommendations: recommendations}
}

// Execute は直近 7 日の人気投稿を 10 件返す。
func (uc *GetTrendingPostsUseCase) Execute(ctx context.Context) ([]model.Post, error) {
	return uc.recommendations.GetTrendingPosts(ctx, trendingPostLimit, trendingPostDays)
}

// GetTrendingResourcesUseCase は人気の学習リソースを取得する。
type GetTrendingResourcesUseCase struct {
	recommendations repository.RecommendationRepository
}

// NewGetTrendingResourcesUseCase は GetTrendingResourcesUseCase を生成する。
func NewGetTrendingResourcesUseCase(recommendations repository.RecommendationRepository) *GetTrendingResourcesUseCase {
	return &GetTrendingResourcesUseCase{recommendations: recommendations}
}

// Execute は直近 30 日の人気学習リソースを 10 件返す。
func (uc *GetTrendingResourcesUseCase) Execute(ctx context.Context) ([]model.LearningResource, error) {
	return uc.recommendations.GetTrendingResources(ctx, trendingResourceLimit, trendingResourceDays)
}
