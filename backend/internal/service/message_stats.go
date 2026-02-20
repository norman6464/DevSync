package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// MessageStatsService はユーザーメッセージ集計統計のビジネスロジックを提供する。
type MessageStatsService struct {
	repo repository.MessageStatsRepositoryInterface
}

// NewMessageStatsService は新しいMessageStatsServiceインスタンスを生成する。
func NewMessageStatsService(repo repository.MessageStatsRepositoryInterface) *MessageStatsService {
	return &MessageStatsService{repo: repo}
}

// GetMessageStats は指定ユーザーのメッセージ集計統計を取得する。
func (s *MessageStatsService) GetMessageStats(userID uint) (*model.MessageStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetMessageStats(userID)
}
