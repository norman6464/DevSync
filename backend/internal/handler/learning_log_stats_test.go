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

// mockLearningLogStatsRepo は usecase/repository.LearningLogStatsRepository のモック（ctx 付き）。
type mockLearningLogStatsRepo struct{ mock.Mock }

func (m *mockLearningLogStatsRepo) GetLearningLogStats(ctx context.Context, userID uint) (*model.LearningLogStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.LearningLogStats)
	return s, args.Error(1)
}

func setupLearningLogStatsHandler() (*LearningLogStatsHandler, *mockLearningLogStatsRepo) {
	stats := new(mockLearningLogStatsRepo)
	h := NewLearningLogStatsHandler(usecase.NewGetLearningLogStatsUseCase(stats))
	return h, stats
}

func TestLearningLogStats_GetStats_Success(t *testing.T) {
	h, stats := setupLearningLogStatsHandler()
	stats.On("GetLearningLogStats", mock.Anything, uint(5)).Return(
		&model.LearningLogStats{TotalLogs: 4, TotalDuration: 300, CategoryCount: 2, LogsThisMonth: 3}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/learning-logs", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/learning-logs", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_logs"].(float64) != 4 {
		t.Errorf("expected total_logs 4, got %v", data["total_logs"])
	}
	stats.AssertExpectations(t)
}

func TestLearningLogStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupLearningLogStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/learning-logs", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/learning-logs", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLogStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupLearningLogStatsHandler()
	stats.On("GetLearningLogStats", mock.Anything, uint(5)).Return((*model.LearningLogStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/learning-logs", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/learning-logs", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
