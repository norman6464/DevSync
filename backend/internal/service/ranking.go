package service

import (
	"github.com/norman6464/devsync/backend/internal/repository"
)

// RankingService はランキングのビジネスロジックを提供する。
// リポジトリ層に処理を委譲する薄いラッパー。
type RankingService struct {
	repo repository.RankingRepositoryInterface
}

// NewRankingService は新しいRankingServiceインスタンスを生成する。
func NewRankingService(repo repository.RankingRepositoryInterface) *RankingService {
	return &RankingService{repo: repo}
}

// ContributionRanking は指定期間のコントリビューションランキングを返す。
func (s *RankingService) ContributionRanking(period string) ([]repository.RankingEntry, error) {
	return s.repo.ContributionRanking(period)
}

// LanguageRanking は指定言語・期間の言語別ランキングを返す。
func (s *RankingService) LanguageRanking(language, period string) ([]repository.RankingEntry, error) {
	return s.repo.LanguageRanking(language, period)
}

// AvailableLanguages はランキング対象の言語一覧を返す。
func (s *RankingService) AvailableLanguages() ([]string, error) {
	return s.repo.AvailableLanguages()
}
