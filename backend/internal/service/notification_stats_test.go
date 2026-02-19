package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newNotificationStatsTestService() (*NotificationStatsService, *MockNotificationStatsRepository) {
	repo := new(MockNotificationStatsRepository)
	svc := NewNotificationStatsService(repo)
	return svc, repo
}

func TestNotificationStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newNotificationStatsTestService()
	expected := &model.NotificationStats{
		TotalNotifications:     120,
		UnreadCount:            15,
		NotificationsThisMonth: 32,
	}
	repo.On("GetNotificationStats", uint(1)).Return(expected, nil)

	result, err := svc.GetNotificationStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestNotificationStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newNotificationStatsTestService()

	_, err := svc.GetNotificationStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestNotificationStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newNotificationStatsTestService()
	repo.On("GetNotificationStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetNotificationStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestNotificationStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newNotificationStatsTestService()
	expected := &model.NotificationStats{}
	repo.On("GetNotificationStats", uint(99)).Return(expected, nil)

	result, err := svc.GetNotificationStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalNotifications)
	assert.Equal(t, int64(0), result.UnreadCount)
	assert.Equal(t, int64(0), result.NotificationsThisMonth)
	repo.AssertExpectations(t)
}
