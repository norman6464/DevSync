package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// FollowStatsService はユーザーフォロー関係集計統計のビジネスロジックを提供する。
type FollowStatsService struct {
	repo repository.FollowStatsRepositoryInterface
}

// NewFollowStatsService は新しいFollowStatsServiceインスタンスを生成する。
func NewFollowStatsService(repo repository.FollowStatsRepositoryInterface) *FollowStatsService {
	return &FollowStatsService{repo: repo}
}

// GetFollowStats は指定ユーザーのフォロー関係集計統計を取得する。
func (s *FollowStatsService) GetFollowStats(userID uint) (*model.FollowStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetFollowStats(userID)
}
