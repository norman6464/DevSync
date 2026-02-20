package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// PostStatsService はユーザー投稿集計統計のビジネスロジックを提供する。
type PostStatsService struct {
	repo repository.PostStatsRepositoryInterface
}

// NewPostStatsService は新しいPostStatsServiceインスタンスを生成する。
func NewPostStatsService(repo repository.PostStatsRepositoryInterface) *PostStatsService {
	return &PostStatsService{repo: repo}
}

// GetPostStats は指定ユーザーの投稿集計統計を取得する。
func (s *PostStatsService) GetPostStats(userID uint) (*model.PostStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetPostStats(userID)
}
