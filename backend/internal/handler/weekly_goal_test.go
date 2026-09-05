package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockWeeklyGoalRepo は usecase/repository.WeeklyGoalRepository のモック（ctx 付き）。
type mockWeeklyGoalRepo struct{ mock.Mock }

func (m *mockWeeklyGoalRepo) Upsert(ctx context.Context, goal *model.WeeklyGoal) error {
	return m.Called(ctx, goal).Error(0)
}

func (m *mockWeeklyGoalRepo) GetByUserID(ctx context.Context, userID uint) ([]model.WeeklyGoal, error) {
	args := m.Called(ctx, userID)
	goals, _ := args.Get(0).([]model.WeeklyGoal)
	return goals, args.Error(1)
}

func (m *mockWeeklyGoalRepo) SumDurationByUserCategoryThisWeek(ctx context.Context, userID uint, category string) (int, error) {
	args := m.Called(ctx, userID, category)
	return args.Int(0), args.Error(1)
}

// setupWeeklyGoalHandler は本物の usecase + port モックで WeeklyGoalHandler を組む。
func setupWeeklyGoalHandler() (*WeeklyGoalHandler, *mockWeeklyGoalRepo) {
	goals := new(mockWeeklyGoalRepo)
	h := NewWeeklyGoalHandler(
		usecase.NewSetWeeklyGoalUseCase(goals),
		usecase.NewListWeeklyGoalsUseCase(goals),
		usecase.NewGetWeeklyGoalProgressUseCase(goals),
	)
	return h, goals
}

func TestWeeklyGoal_SetGoal_Success(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()
	goals.On("Upsert", mock.Anything, mock.AnythingOfType("*model.WeeklyGoal")).Return(nil)

	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{
		"category": "coding", "target_minutes": 300,
	})

	assertStatus(t, w, http.StatusOK)
	goals.AssertExpectations(t)
}

func TestWeeklyGoal_SetGoal_InvalidJSON(t *testing.T) {
	h, _ := setupWeeklyGoalHandler()

	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestWeeklyGoal_SetGoal_InvalidCategory(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()

	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{
		"category": "invalid", "target_minutes": 300,
	})

	assertStatus(t, w, http.StatusBadRequest)
	goals.AssertNotCalled(t, "Upsert")
}

func TestWeeklyGoal_SetGoal_RepoError(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()
	goals.On("Upsert", mock.Anything, mock.AnythingOfType("*model.WeeklyGoal")).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{
		"category": "coding", "target_minutes": 300,
	})

	assertStatus(t, w, http.StatusInternalServerError)
	goals.AssertExpectations(t)
}

func TestWeeklyGoal_GetGoals_Success(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()
	goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal{
		{ID: 1, UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300},
	}, nil)

	r := newRouter(1)
	r.GET("/weekly-goals", h.GetGoals)
	w := doRequest(r, http.MethodGet, "/weekly-goals", nil)

	assertStatus(t, w, http.StatusOK)
	goals.AssertExpectations(t)
}

func TestWeeklyGoal_GetGoals_RepoError(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()
	goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/weekly-goals", h.GetGoals)
	w := doRequest(r, http.MethodGet, "/weekly-goals", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	goals.AssertExpectations(t)
}

func TestWeeklyGoal_GetProgress_Success(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()
	goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal{
		{Category: model.LogCategoryCoding, TargetMinutes: 300},
	}, nil)
	goals.On("SumDurationByUserCategoryThisWeek", mock.Anything, uint(1), "coding").Return(150, nil)

	r := newRouter(1)
	r.GET("/weekly-goals/progress", h.GetProgress)
	w := doRequest(r, http.MethodGet, "/weekly-goals/progress", nil)

	assertStatus(t, w, http.StatusOK)
	goals.AssertExpectations(t)
}

func TestWeeklyGoal_GetProgress_ServiceError(t *testing.T) {
	h, goals := setupWeeklyGoalHandler()
	goals.On("GetByUserID", mock.Anything, uint(1)).Return([]model.WeeklyGoal(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/weekly-goals/progress", h.GetProgress)
	w := doRequest(r, http.MethodGet, "/weekly-goals/progress", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	goals.AssertExpectations(t)
}
