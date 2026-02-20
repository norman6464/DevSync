package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// RoadmapStatsService はユーザーロードマップ統計のビジネスロジックを提供する。
type RoadmapStatsService struct {
	repo repository.RoadmapStatsRepositoryInterface
}

// NewRoadmapStatsService は新しいRoadmapStatsServiceインスタンスを生成する。
func NewRoadmapStatsService(repo repository.RoadmapStatsRepositoryInterface) *RoadmapStatsService {
	return &RoadmapStatsService{repo: repo}
}

// GetRoadmapStats は指定ユーザーのロードマップ統計を取得する。
func (s *RoadmapStatsService) GetRoadmapStats(userID uint) (*model.RoadmapStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetRoadmapStats(userID)
}
