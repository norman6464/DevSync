package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NotificationStatsService はユーザー通知集計統計のビジネスロジックを提供する。
type NotificationStatsService struct {
	repo repository.NotificationStatsRepositoryInterface
}

// NewNotificationStatsService は新しいNotificationStatsServiceインスタンスを生成する。
func NewNotificationStatsService(repo repository.NotificationStatsRepositoryInterface) *NotificationStatsService {
	return &NotificationStatsService{repo: repo}
}

// GetNotificationStats は指定ユーザーの通知集計統計を取得する。
func (s *NotificationStatsService) GetNotificationStats(userID uint) (*model.NotificationStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetNotificationStats(userID)
}
