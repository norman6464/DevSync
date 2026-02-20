package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningLogStatsService はユーザー学習ログ集計統計のビジネスロジックを提供する。
type LearningLogStatsService struct {
	repo repository.LearningLogStatsRepositoryInterface
}

// NewLearningLogStatsService は新しいLearningLogStatsServiceインスタンスを生成する。
func NewLearningLogStatsService(repo repository.LearningLogStatsRepositoryInterface) *LearningLogStatsService {
	return &LearningLogStatsService{repo: repo}
}

// GetLearningLogStats は指定ユーザーの学習ログ集計統計を取得する。
func (s *LearningLogStatsService) GetLearningLogStats(userID uint) (*model.LearningLogStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetLearningLogStats(userID)
}
