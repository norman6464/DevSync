package handler

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// --- Mock ---

type MockUserActivityService struct{ mock.Mock }

func (m *MockUserActivityService) GetTimeline(userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	args := m.Called(userID, activityType, limit, offset)
	return args.Get(0).([]model.UserActivity), args.Get(1).(int64), args.Error(2)
}

func setupUserActivityHandler() (*UserActivityHandler, *MockUserActivityService) {
	svc := new(MockUserActivityService)
	h := NewUserActivityHandler(svc)
	return h, svc
}

// --- GetTimeline ---

func TestUserActivityHandler_GetTimeline_Success(t *testing.T) {
	h, svc := setupUserActivityHandler()

	svc.On("GetTimeline", uint(10), "", 20, 0).Return(
		[]model.UserActivity{
			{ID: 1, UserID: 10, ActivityType: model.ActivityPostCreated, TargetType: "post", TargetID: 5},
			{ID: 2, UserID: 10, ActivityType: model.ActivityCommentCreated, TargetType: "comment", TargetID: 3},
		},
		int64(2), nil,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity", nil)
	assertStatus(t, w, 200)
	data := parseJSON(t, w)
	activities := data["activities"].([]interface{})
	if len(activities) != 2 {
		t.Errorf("expected 2 activities, got %d", len(activities))
	}
}

func TestUserActivityHandler_GetTimeline_WithTypeFilter(t *testing.T) {
	h, svc := setupUserActivityHandler()

	svc.On("GetTimeline", uint(10), "post_created", 20, 0).Return(
		[]model.UserActivity{{ID: 1, UserID: 10, ActivityType: model.ActivityPostCreated}},
		int64(1), nil,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity?type=post_created", nil)
	assertStatus(t, w, 200)
}

func TestUserActivityHandler_GetTimeline_InvalidUserID(t *testing.T) {
	h, _ := setupUserActivityHandler()

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/abc/activity", nil)
	assertStatus(t, w, 400)
}

func TestUserActivityHandler_GetTimeline_ServiceError(t *testing.T) {
	h, svc := setupUserActivityHandler()

	svc.On("GetTimeline", uint(10), "", 20, 0).Return(
		[]model.UserActivity{}, int64(0), service.ErrNotFound,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity", nil)
	assertStatus(t, w, 404)
}

func TestUserActivityHandler_GetTimeline_Empty(t *testing.T) {
	h, svc := setupUserActivityHandler()

	svc.On("GetTimeline", uint(10), "", 20, 0).Return(
		[]model.UserActivity{}, int64(0), nil,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity", nil)
	assertStatus(t, w, 200)
	data := parseJSON(t, w)
	activities := data["activities"].([]interface{})
	if len(activities) != 0 {
		t.Errorf("expected 0 activities, got %d", len(activities))
	}
}
