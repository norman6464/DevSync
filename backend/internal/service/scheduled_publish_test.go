package service

import (
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScheduledPublish_Schedule_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true}
	mockRepo.On("FindByID", uint(1)).Return(post, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *model.Post) bool {
		return p.ScheduledAt != nil && p.IsDraft == true
	})).Return(nil)

	scheduledAt := time.Now().Add(24 * time.Hour)
	result, err := svc.SchedulePublish(1, 1, scheduledAt)
	assert.NoError(t, err)
	assert.NotNil(t, result.ScheduledAt)
}

func TestScheduledPublish_Schedule_PastTime(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true}
	mockRepo.On("FindByID", uint(1)).Return(post, nil)

	pastTime := time.Now().Add(-1 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, pastTime)
	assert.Error(t, err)
}

func TestScheduledPublish_Schedule_NotDraft(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: false}
	mockRepo.On("FindByID", uint(1)).Return(post, nil)

	scheduledAt := time.Now().Add(24 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, scheduledAt)
	assert.Error(t, err)
}

func TestScheduledPublish_Schedule_Forbidden(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	post := &model.Post{ID: 1, UserID: 2, IsDraft: true}
	mockRepo.On("FindByID", uint(1)).Return(post, nil)

	scheduledAt := time.Now().Add(24 * time.Hour)
	_, err := svc.SchedulePublish(1, 1, scheduledAt)
	assert.Error(t, err)
}

func TestScheduledPublish_Cancel_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	scheduled := time.Now().Add(24 * time.Hour)
	post := &model.Post{ID: 1, UserID: 1, IsDraft: true, ScheduledAt: &scheduled}
	mockRepo.On("FindByID", uint(1)).Return(post, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *model.Post) bool {
		return p.ScheduledAt == nil && p.IsDraft == true
	})).Return(nil)

	result, err := svc.CancelSchedule(1, 1)
	assert.NoError(t, err)
	assert.Nil(t, result.ScheduledAt)
}

func TestScheduledPublish_Cancel_NotScheduled(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	post := &model.Post{ID: 1, UserID: 1, IsDraft: true, ScheduledAt: nil}
	mockRepo.On("FindByID", uint(1)).Return(post, nil)

	_, err := svc.CancelSchedule(1, 1)
	assert.Error(t, err)
}

func TestScheduledPublish_GetScheduled_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	mockNotif := new(MockNotificationService)
	svc := NewPostService(mockRepo, mockNotif)

	scheduled := time.Now().Add(24 * time.Hour)
	posts := []model.Post{
		{ID: 1, UserID: 1, Title: "予定投稿", ScheduledAt: &scheduled},
	}
	mockRepo.On("FindScheduledByUserID", uint(1)).Return(posts, nil)

	result, err := svc.GetScheduled(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}
