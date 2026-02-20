package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ProjectStatsService はユーザープロジェクト活動集計統計のビジネスロジックを提供する。
type ProjectStatsService struct {
	repo repository.ProjectStatsRepositoryInterface
}

// NewProjectStatsService は新しいProjectStatsServiceインスタンスを生成する。
func NewProjectStatsService(repo repository.ProjectStatsRepositoryInterface) *ProjectStatsService {
	return &ProjectStatsService{repo: repo}
}

// GetProjectStats は指定ユーザーのプロジェクト活動集計統計を取得する。
func (s *ProjectStatsService) GetProjectStats(userID uint) (*model.ProjectStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetProjectStats(userID)
}
