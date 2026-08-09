package handler

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockUserActivityRepo は usecase/repository.UserActivityRepository のモック（ctx 付き）。
type mockUserActivityRepo struct{ mock.Mock }

func (m *mockUserActivityRepo) FindByUserID(ctx context.Context, userID uint, activityType string, limit, offset int) ([]model.UserActivity, int64, error) {
	args := m.Called(ctx, userID, activityType, limit, offset)
	acts, _ := args.Get(0).([]model.UserActivity)
	return acts, args.Get(1).(int64), args.Error(2)
}

// setupUserActivityHandler は本物の usecase + port モックで UserActivityHandler を組む。
func setupUserActivityHandler() (*UserActivityHandler, *mockUserActivityRepo) {
	activities := new(mockUserActivityRepo)
	h := NewUserActivityHandler(usecase.NewGetActivityTimelineUseCase(activities))
	return h, activities
}

func TestUserActivityHandler_GetTimeline_Success(t *testing.T) {
	h, activities := setupUserActivityHandler()

	activities.On("FindByUserID", mock.Anything, uint(10), "", 20, 0).Return(
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
	acts := data["activities"].([]interface{})
	if len(acts) != 2 {
		t.Errorf("expected 2 activities, got %d", len(acts))
	}
	activities.AssertExpectations(t)
}

func TestUserActivityHandler_GetTimeline_WithTypeFilter(t *testing.T) {
	h, activities := setupUserActivityHandler()

	activities.On("FindByUserID", mock.Anything, uint(10), "post_created", 20, 0).Return(
		[]model.UserActivity{{ID: 1, UserID: 10, ActivityType: model.ActivityPostCreated}},
		int64(1), nil,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity?type=post_created", nil)
	assertStatus(t, w, 200)
	activities.AssertExpectations(t)
}

func TestUserActivityHandler_GetTimeline_InvalidUserID(t *testing.T) {
	h, _ := setupUserActivityHandler()

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/abc/activity", nil)
	assertStatus(t, w, 400)
}

func TestUserActivityHandler_GetTimeline_ServiceError(t *testing.T) {
	h, activities := setupUserActivityHandler()

	activities.On("FindByUserID", mock.Anything, uint(10), "", 20, 0).Return(
		[]model.UserActivity(nil), int64(0), domain.ErrNotFound,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity", nil)
	assertStatus(t, w, 404)
	activities.AssertExpectations(t)
}

func TestUserActivityHandler_GetTimeline_Empty(t *testing.T) {
	h, activities := setupUserActivityHandler()

	activities.On("FindByUserID", mock.Anything, uint(10), "", 20, 0).Return(
		[]model.UserActivity{}, int64(0), nil,
	)

	r := newRouter(1)
	r.GET("/users/:userId/activity", h.GetTimeline)

	w := doRequest(r, "GET", "/users/10/activity", nil)
	assertStatus(t, w, 200)
	data := parseJSON(t, w)
	acts := data["activities"].([]interface{})
	if len(acts) != 0 {
		t.Errorf("expected 0 activities, got %d", len(acts))
	}
	activities.AssertExpectations(t)
}
