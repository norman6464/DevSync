package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestUserActivityService() (*UserActivityService, *MockUserActivityRepository) {
	repo := new(MockUserActivityRepository)
	svc := NewUserActivityService(repo)
	return svc, repo
}

// --- RecordActivity ---

func TestUserActivityService_RecordActivity_Success(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("Create", mock.MatchedBy(func(a *model.UserActivity) bool {
		return a.UserID == 10 && a.ActivityType == model.ActivityPostCreated && a.TargetType == "post" && a.TargetID == 5
	})).Return(nil)

	err := svc.RecordActivity(10, model.ActivityPostCreated, "post", 5, "")
	assert.NoError(t, err)
}

func TestUserActivityService_RecordActivity_WithMetadata(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("Create", mock.MatchedBy(func(a *model.UserActivity) bool {
		return a.Metadata == `{"title":"Hello"}`
	})).Return(nil)

	err := svc.RecordActivity(10, model.ActivityPostCreated, "post", 5, `{"title":"Hello"}`)
	assert.NoError(t, err)
}

func TestUserActivityService_RecordActivity_RepoError(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("Create", mock.Anything).Return(errors.New("db error"))

	err := svc.RecordActivity(10, model.ActivityPostCreated, "post", 5, "")
	assert.Error(t, err)
}

// --- GetTimeline ---

func TestUserActivityService_GetTimeline_NoFilter(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("FindByUserID", uint(10), "", 20, 0).Return(
		[]model.UserActivity{
			{ID: 1, UserID: 10, ActivityType: model.ActivityPostCreated},
			{ID: 2, UserID: 10, ActivityType: model.ActivityCommentCreated},
		},
		int64(2), nil,
	)

	list, total, err := svc.GetTimeline(10, "", 20, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, int64(2), total)
}

func TestUserActivityService_GetTimeline_WithTypeFilter(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("FindByUserID", uint(10), "post_created", 20, 0).Return(
		[]model.UserActivity{{ID: 1, UserID: 10, ActivityType: model.ActivityPostCreated}},
		int64(1), nil,
	)

	list, total, err := svc.GetTimeline(10, "post_created", 20, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), total)
}

func TestUserActivityService_GetTimeline_RepoError(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("FindByUserID", uint(10), "", 20, 0).Return(
		[]model.UserActivity{}, int64(0), errors.New("db error"),
	)

	_, _, err := svc.GetTimeline(10, "", 20, 0)
	assert.Error(t, err)
}

func TestUserActivityService_GetTimeline_Empty(t *testing.T) {
	svc, repo := newTestUserActivityService()

	repo.On("FindByUserID", uint(10), "", 20, 0).Return(
		[]model.UserActivity{}, int64(0), nil,
	)

	list, total, err := svc.GetTimeline(10, "", 20, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
	assert.Equal(t, int64(0), total)
}
