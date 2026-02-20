package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// MockLearningGoalRepository は LearningGoalRepositoryInterface のモック実装。
type MockLearningGoalRepository struct{ mock.Mock }

func (m *MockLearningGoalRepository) Create(goal *model.LearningGoal) error {
	return m.Called(goal).Error(0)
}

func (m *MockLearningGoalRepository) FindByID(id uint) (*model.LearningGoal, error) {
	args := m.Called(id)
	if goal := args.Get(0); goal != nil {
		return goal.(*model.LearningGoal), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLearningGoalRepository) GetByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningGoal), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningGoalRepository) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.LearningGoal), args.Error(1)
}

func (m *MockLearningGoalRepository) Update(goal *model.LearningGoal) error {
	return m.Called(goal).Error(0)
}

func (m *MockLearningGoalRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}

func (m *MockLearningGoalRepository) GetStats(userID uint) (*model.LearningGoalStats, error) {
	args := m.Called(userID)
	if stats := args.Get(0); stats != nil {
		return stats.(*model.LearningGoalStats), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLearningGoalRepository) GetByCategory(userID uint, category string) ([]model.LearningGoal, error) {
	args := m.Called(userID, category)
	return args.Get(0).([]model.LearningGoal), args.Error(1)
}

func (m *MockLearningGoalRepository) GetByStatus(userID uint, status string) ([]model.LearningGoal, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.LearningGoal), args.Error(1)
}

// setupLearningGoalHandler はテスト用のLearningGoalHandlerとモックを準備する。
func setupLearningGoalHandler() (*LearningGoalHandler, *MockLearningGoalRepository) {
	repo := new(MockLearningGoalRepository)
	svc := service.NewLearningGoalService(repo)
	handler := NewLearningGoalHandler(svc)
	return handler, repo
}

// ========== Create ==========

func TestLearningGoalCreate_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.POST("/goals", h.Create)

	repo.On("Create", mock.AnythingOfType("*model.LearningGoal")).Return(nil)

	w := doRequest(r, http.MethodPost, "/goals", map[string]string{
		"title": "Test Goal",
	})

	assertStatus(t, w, http.StatusCreated)
}

func TestLearningGoalCreate_ValidationError(t *testing.T) {
	h, _ := setupLearningGoalHandler()
	r := newRouter(1)
	r.POST("/goals", h.Create)

	// title is required
	w := doRequest(r, http.MethodPost, "/goals", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningGoalCreate_InvalidJSON(t *testing.T) {
	h, _ := setupLearningGoalHandler()
	r := newRouter(1)
	r.POST("/goals", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/goals", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
}

// ========== Update ==========

func TestLearningGoalUpdate_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.PUT("/goals/:id", h.Update)

	goal := &model.LearningGoal{}
	goal.ID = 10
	goal.UserID = 1
	goal.Title = "Old Title"

	repo.On("FindByID", uint(10)).Return(goal, nil)
	repo.On("Update", mock.AnythingOfType("*model.LearningGoal")).Return(nil)

	newTitle := "Updated Title"
	w := doRequest(r, http.MethodPut, "/goals/10", map[string]interface{}{
		"title": &newTitle,
	})

	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalUpdate_Forbidden(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.PUT("/goals/:id", h.Update)

	goal := &model.LearningGoal{}
	goal.ID = 10
	goal.UserID = 999 // 別のユーザー
	goal.Title = "Goal"

	repo.On("FindByID", uint(10)).Return(goal, nil)

	newTitle := "Updated"
	w := doRequest(r, http.MethodPut, "/goals/10", map[string]interface{}{
		"title": &newTitle,
	})

	assertStatus(t, w, http.StatusForbidden)
}

func TestLearningGoalUpdate_NotFound(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.PUT("/goals/:id", h.Update)

	repo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	newTitle := "Updated"
	w := doRequest(r, http.MethodPut, "/goals/999", map[string]interface{}{
		"title": &newTitle,
	})

	assertStatus(t, w, http.StatusNotFound)
}

// ========== Delete ==========

func TestLearningGoalDelete_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.DELETE("/goals/:id", h.Delete)

	goal := &model.LearningGoal{}
	goal.ID = 10
	goal.UserID = 1

	repo.On("FindByID", uint(10)).Return(goal, nil)
	repo.On("Delete", uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/goals/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalDelete_Forbidden(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.DELETE("/goals/:id", h.Delete)

	goal := &model.LearningGoal{}
	goal.ID = 10
	goal.UserID = 999 // 別のユーザー

	repo.On("FindByID", uint(10)).Return(goal, nil)

	w := doRequest(r, http.MethodDelete, "/goals/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestLearningGoalDelete_NotFound(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.DELETE("/goals/:id", h.Delete)

	repo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/goals/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ========== GetByID ==========

func TestLearningGoalGetByID_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/:id", h.GetByID)

	goal := &model.LearningGoal{}
	goal.ID = 10
	goal.Title = "Test Goal"

	repo.On("FindByID", uint(10)).Return(goal, nil)

	w := doRequest(r, http.MethodGet, "/goals/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetByID_NotFound(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/:id", h.GetByID)

	repo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/goals/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ========== GetMyGoals ==========

func TestLearningGoalGetMyGoals_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/my", h.GetMyGoals)

	goals := []model.LearningGoal{
		{Title: "Goal 1"},
		{Title: "Goal 2"},
	}

	repo.On("GetByUserID", uint(1), 20, 0).Return(goals, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/goals/my", nil)
	assertStatus(t, w, http.StatusOK)
}

// ========== GetByUserID ==========

func TestLearningGoalGetByUserID_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/users/:userId/goals", h.GetByUserID)

	goals := []model.LearningGoal{{Title: "Goal 1"}, {Title: "Goal 2"}}
	repo.On("GetByUserID", uint(5), 20, 0).Return(goals, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/users/5/goals", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if body["total"] != float64(2) {
		t.Errorf("expected total=2, got %v", body["total"])
	}
}

func TestLearningGoalGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/users/:userId/goals", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/goals", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningGoalGetByUserID_ServiceError(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/users/:userId/goals", h.GetByUserID)

	repo.On("GetByUserID", uint(5), 20, 0).Return([]model.LearningGoal{}, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/5/goals", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ========== GetDeadlineAlerts ==========

func TestLearningGoalGetDeadlineAlerts_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/deadline-alerts", h.GetDeadlineAlerts)

	repo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, nil)

	w := doRequest(r, http.MethodGet, "/goals/deadline-alerts", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetDeadlineAlerts_ServiceError(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/deadline-alerts", h.GetDeadlineAlerts)

	repo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/goals/deadline-alerts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ========== GetByCategory ==========

func TestLearningGoalGetByCategory_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/category/:category", h.GetByCategory)

	goals := []model.LearningGoal{{Title: "Go学習", Category: model.GoalCategoryLanguage}}
	repo.On("GetByCategory", uint(1), "language").Return(goals, nil)

	w := doRequest(r, http.MethodGet, "/goals/category/language", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetByCategory_NilResult(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/category/:category", h.GetByCategory)

	var nilGoals []model.LearningGoal
	repo.On("GetByCategory", uint(1), "framework").Return(nilGoals, nil)

	w := doRequest(r, http.MethodGet, "/goals/category/framework", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetByCategory_InvalidCategory(t *testing.T) {
	h, _ := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/category/:category", h.GetByCategory)

	w := doRequest(r, http.MethodGet, "/goals/category/invalid", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ========== GetByStatus ==========

func TestLearningGoalGetByStatus_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/status/:status", h.GetByStatus)

	goals := []model.LearningGoal{{Title: "完了目標", Status: model.GoalStatusCompleted}}
	repo.On("GetByStatus", uint(1), "completed").Return(goals, nil)

	w := doRequest(r, http.MethodGet, "/goals/status/completed", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetByStatus_NilResult(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/status/:status", h.GetByStatus)

	var nilGoals []model.LearningGoal
	repo.On("GetByStatus", uint(1), "paused").Return(nilGoals, nil)

	w := doRequest(r, http.MethodGet, "/goals/status/paused", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetByStatus_InvalidStatus(t *testing.T) {
	h, _ := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/goals/status/:status", h.GetByStatus)

	w := doRequest(r, http.MethodGet, "/goals/status/invalid", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ========== GetStats ==========

func TestLearningGoalGetStats_Success(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/users/:userId/goals/stats", h.GetStats)

	stats := &model.LearningGoalStats{
		TotalGoals:      10,
		ActiveGoals:     3,
		CompletedGoals:  5,
		AverageProgress: 60,
	}

	repo.On("GetStats", uint(1)).Return(stats, nil)

	w := doRequest(r, http.MethodGet, "/users/1/goals/stats", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestLearningGoalGetStats_ServiceError(t *testing.T) {
	h, repo := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/users/:userId/goals/stats", h.GetStats)

	repo.On("GetStats", uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/1/goals/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestLearningGoalGetStats_InvalidID(t *testing.T) {
	h, _ := setupLearningGoalHandler()
	r := newRouter(1)
	r.GET("/users/:userId/goals/stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/goals/stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
}
