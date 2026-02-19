package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningResourceStatsService はユーザー学習リソース活動集計統計のビジネスロジックを提供する。
type LearningResourceStatsService struct {
	repo repository.LearningResourceStatsRepositoryInterface
}

// NewLearningResourceStatsService は新しいLearningResourceStatsServiceインスタンスを生成する。
func NewLearningResourceStatsService(repo repository.LearningResourceStatsRepositoryInterface) *LearningResourceStatsService {
	return &LearningResourceStatsService{repo: repo}
}

// GetLearningResourceStats は指定ユーザーの学習リソース活動集計統計を取得する。
func (s *LearningResourceStatsService) GetLearningResourceStats(userID uint) (*model.LearningResourceStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetLearningResourceStats(userID)
}
