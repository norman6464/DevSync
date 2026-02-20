package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// MentionStatsService はユーザーメンション集計統計のビジネスロジックを提供する。
type MentionStatsService struct {
	repo repository.MentionStatsRepositoryInterface
}

// NewMentionStatsService は新しいMentionStatsServiceインスタンスを生成する。
func NewMentionStatsService(repo repository.MentionStatsRepositoryInterface) *MentionStatsService {
	return &MentionStatsService{repo: repo}
}

// GetMentionStats は指定ユーザーのメンション集計統計を取得する。
func (s *MentionStatsService) GetMentionStats(userID uint) (*model.MentionStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetMentionStats(userID)
}
