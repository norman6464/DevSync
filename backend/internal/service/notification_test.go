package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestNotificationService() (*NotificationService, *MockNotificationRepository) {
	repo := new(MockNotificationRepository)
	svc := NewNotificationService(repo)
	return svc, repo
}

// ============================================================
// NotifyFollowers
// ============================================================

func TestNotifyFollowers_Success(t *testing.T) {
	svc, repo := newTestNotificationService()

	repo.On("GetFollowerIDs", uint(1)).Return([]uint{2, 3, 4}, nil)
	repo.On("CreateBatch", mock.MatchedBy(func(notifications []*model.Notification) bool {
		return len(notifications) == 3
	})).Return(nil)

	svc.NotifyFollowers(1, 10, model.NotificationTypePost)
	repo.AssertExpectations(t)
}

func TestNotifyFollowers_NoFollowers(t *testing.T) {
	svc, repo := newTestNotificationService()

	repo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)

	// フォロワー0人の場合、CreateBatchは呼ばれない
	svc.NotifyFollowers(1, 10, model.NotificationTypePost)
	repo.AssertNotCalled(t, "CreateBatch")
}

// ============================================================
// CountUnread
// ============================================================

func TestNotificationCountUnread_Success(t *testing.T) {
	svc, repo := newTestNotificationService()

	repo.On("CountUnread", uint(1)).Return(int64(5), nil)

	count, err := svc.CountUnread(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	repo.AssertExpectations(t)
}

// ============================================================
// MarkAllAsRead
// ============================================================

func TestNotificationMarkAllAsRead_Success(t *testing.T) {
	svc, repo := newTestNotificationService()

	repo.On("MarkAllAsRead", uint(1)).Return(nil)

	err := svc.MarkAllAsRead(1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// CreateNotification
// ============================================================

func TestCreateNotification_Success(t *testing.T) {
	svc, repo := newTestNotificationService()

	notification := &model.Notification{UserID: 1, Type: model.NotificationTypePost, ActorID: 2}
	repo.On("Create", notification).Return(nil)

	err := svc.CreateNotification(notification)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
