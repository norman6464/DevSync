package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// RecommendationService はレコメンド機能のビジネスロジックを提供する。
type RecommendationService struct {
	repo     repository.RecommendationRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

// NewRecommendationService は新しいRecommendationServiceインスタンスを生成する。
func NewRecommendationService(repo repository.RecommendationRepositoryInterface, userRepo repository.UserRepositoryInterface) *RecommendationService {
	return &RecommendationService{repo: repo, userRepo: userRepo}
}

// GetRecommendedUsers はスキルマッチングに基づくおすすめユーザーを返す。
func (s *RecommendationService) GetRecommendedUsers(userID uint) ([]model.RecommendedUser, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	skills := parseSkills(user.SkillsLanguages, user.SkillsFrameworks)
	if len(skills) == 0 {
		return []model.RecommendedUser{}, nil
	}

	return s.repo.GetRecommendedUsers(userID, skills, 10)
}

// GetTrendingPosts は直近7日間の人気投稿を返す。
func (s *RecommendationService) GetTrendingPosts() ([]model.Post, error) {
	return s.repo.GetTrendingPosts(10, 7)
}

// GetTrendingResources は直近30日間の人気学習リソースを返す。
func (s *RecommendationService) GetTrendingResources() ([]model.LearningResource, error) {
	return s.repo.GetTrendingResources(10, 30)
}

// parseSkills はカンマ区切りのスキル文字列をスライスに変換する。
func parseSkills(languages, frameworks string) []string {
	var skills []string
	for _, s := range splitSkills(languages) {
		if s != "" {
			skills = append(skills, s)
		}
	}
	for _, s := range splitSkills(frameworks) {
		if s != "" {
			skills = append(skills, s)
		}
	}
	return skills
}

// splitSkills はカンマ区切り文字列を分割してトリムする。
func splitSkills(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			trimmed := trimString(current)
			if trimmed != "" {
				result = append(result, trimmed)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	trimmed := trimString(current)
	if trimmed != "" {
		result = append(result, trimmed)
	}
	return result
}

// trimString は文字列の前後の空白を除去する。
func trimString(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
