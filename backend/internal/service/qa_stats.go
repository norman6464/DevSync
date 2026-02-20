package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// QAStatsService はユーザーQ&A活動集計統計のビジネスロジックを提供する。
type QAStatsService struct {
	repo repository.QAStatsRepositoryInterface
}

// NewQAStatsService は新しいQAStatsServiceインスタンスを生成する。
func NewQAStatsService(repo repository.QAStatsRepositoryInterface) *QAStatsService {
	return &QAStatsService{repo: repo}
}

// GetQAStats は指定ユーザーのQ&A活動集計統計を取得する。
func (s *QAStatsService) GetQAStats(userID uint) (*model.QAStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetQAStats(userID)
}
