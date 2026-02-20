package handler

import (
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
