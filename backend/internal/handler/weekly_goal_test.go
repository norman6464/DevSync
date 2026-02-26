package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockWeeklyGoalService struct{ mock.Mock }

func (m *MockWeeklyGoalService) SetGoal(userID uint, category string, targetMinutes int) (*model.WeeklyGoal, error) {
	args := m.Called(userID, category, targetMinutes)
	if v := args.Get(0); v != nil {
		return v.(*model.WeeklyGoal), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWeeklyGoalService) GetGoals(userID uint) ([]model.WeeklyGoal, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.WeeklyGoal), args.Error(1)
}

func (m *MockWeeklyGoalService) GetProgress(userID uint) ([]model.WeeklyGoalProgress, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.WeeklyGoalProgress), args.Error(1)
}

func TestWeeklyGoal_SetGoal_Success(t *testing.T) {
	svc := new(MockWeeklyGoalService)
	h := NewWeeklyGoalHandler(svc)
	svc.On("SetGoal", uint(1), "coding", 300).Return(&model.WeeklyGoal{
		UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300,
	}, nil)
	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{
		"category": "coding", "target_minutes": 300,
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestWeeklyGoal_SetGoal_InvalidJSON(t *testing.T) {
	svc := new(MockWeeklyGoalService)
	h := NewWeeklyGoalHandler(svc)
	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestWeeklyGoal_SetGoal_ServiceError(t *testing.T) {
	svc := new(MockWeeklyGoalService)
	h := NewWeeklyGoalHandler(svc)
	svc.On("SetGoal", uint(1), "invalid", 300).Return(nil, errors.New("bad request"))
	r := newRouter(1)
	r.PUT("/weekly-goals", h.SetGoal)
	w := doRequest(r, http.MethodPut, "/weekly-goals", map[string]interface{}{
		"category": "invalid", "target_minutes": 300,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestWeeklyGoal_GetGoals_Success(t *testing.T) {
	svc := new(MockWeeklyGoalService)
	h := NewWeeklyGoalHandler(svc)
	svc.On("GetGoals", uint(1)).Return([]model.WeeklyGoal{
		{ID: 1, UserID: 1, Category: model.LogCategoryCoding, TargetMinutes: 300},
	}, nil)
	r := newRouter(1)
	r.GET("/weekly-goals", h.GetGoals)
	w := doRequest(r, http.MethodGet, "/weekly-goals", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestWeeklyGoal_GetProgress_Success(t *testing.T) {
	svc := new(MockWeeklyGoalService)
	h := NewWeeklyGoalHandler(svc)
	svc.On("GetProgress", uint(1)).Return([]model.WeeklyGoalProgress{
		{Category: model.LogCategoryCoding, TargetMinutes: 300, ActualMinutes: 150, ProgressPercent: 50},
	}, nil)
	r := newRouter(1)
	r.GET("/weekly-goals/progress", h.GetProgress)
	w := doRequest(r, http.MethodGet, "/weekly-goals/progress", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestWeeklyGoal_GetProgress_ServiceError(t *testing.T) {
	svc := new(MockWeeklyGoalService)
	h := NewWeeklyGoalHandler(svc)
	svc.On("GetProgress", uint(1)).Return([]model.WeeklyGoalProgress(nil), errors.New("db error"))
	r := newRouter(1)
	r.GET("/weekly-goals/progress", h.GetProgress)
	w := doRequest(r, http.MethodGet, "/weekly-goals/progress", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
