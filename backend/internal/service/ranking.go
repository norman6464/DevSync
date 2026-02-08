package service

import (
	"github.com/norman6464/devsync/backend/internal/repository"
)

// RankingService handles ranking business logic.
type RankingService struct {
	repo repository.RankingRepositoryInterface
}

// NewRankingService creates a new RankingService.
func NewRankingService(repo repository.RankingRepositoryInterface) *RankingService {
	return &RankingService{repo: repo}
}

// ContributionRanking returns the contribution ranking for a period.
func (s *RankingService) ContributionRanking(period string) ([]repository.RankingEntry, error) {
	return s.repo.ContributionRanking(period)
}

// LanguageRanking returns the language ranking for a language and period.
func (s *RankingService) LanguageRanking(language, period string) ([]repository.RankingEntry, error) {
	return s.repo.LanguageRanking(language, period)
}

// AvailableLanguages returns the list of available languages.
func (s *RankingService) AvailableLanguages() ([]string, error) {
	return s.repo.AvailableLanguages()
}
