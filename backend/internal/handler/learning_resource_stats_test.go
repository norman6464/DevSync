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

// mockLearningResourceStatsRepo は usecase/repository.LearningResourceStatsRepository のモック（ctx 付き）。
type mockLearningResourceStatsRepo struct{ mock.Mock }

func (m *mockLearningResourceStatsRepo) GetLearningResourceStats(ctx context.Context, userID uint) (*model.LearningResourceStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.LearningResourceStats)
	return s, args.Error(1)
}

// setupLearningResourceStatsHandler は本物の usecase + port モックで組む。
func setupLearningResourceStatsHandler() (*LearningResourceStatsHandler, *mockLearningResourceStatsRepo) {
	stats := new(mockLearningResourceStatsRepo)
	h := NewLearningResourceStatsHandler(usecase.NewGetLearningResourceStatsUseCase(stats))
	return h, stats
}

func TestLearningResourceStatsHandler_GetStats_Success(t *testing.T) {
	h, stats := setupLearningResourceStatsHandler()
	stats.On("GetLearningResourceStats", mock.Anything, uint(10)).Return(
		&model.LearningResourceStats{TotalResources: 5, TotalLikes: 12, TotalSaves: 3, CategoryCount: 2}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/learning-resource-stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/10/learning-resource-stats", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_resources"].(float64) != 5 {
		t.Errorf("expected total_resources 5, got %v", data["total_resources"])
	}
	stats.AssertExpectations(t)
}

func TestLearningResourceStatsHandler_GetStats_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/learning-resource-stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/learning-resource-stats", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningResourceStatsHandler_GetStats_ServiceError(t *testing.T) {
	h, stats := setupLearningResourceStatsHandler()
	stats.On("GetLearningResourceStats", mock.Anything, uint(10)).Return((*model.LearningResourceStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/learning-resource-stats", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/10/learning-resource-stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
