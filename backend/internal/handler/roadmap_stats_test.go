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

// mockRoadmapStatsRepo は usecase/repository.RoadmapStatsRepository のモック（ctx 付き）。
type mockRoadmapStatsRepo struct{ mock.Mock }

func (m *mockRoadmapStatsRepo) GetRoadmapStats(ctx context.Context, userID uint) (*model.RoadmapStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.RoadmapStats)
	return s, args.Error(1)
}

func setupRoadmapStatsHandler() (*RoadmapStatsHandler, *mockRoadmapStatsRepo) {
	stats := new(mockRoadmapStatsRepo)
	h := NewRoadmapStatsHandler(usecase.NewGetRoadmapStatsUseCase(stats))
	return h, stats
}

func TestRoadmapStats_GetStats_Success(t *testing.T) {
	h, stats := setupRoadmapStatsHandler()
	stats.On("GetRoadmapStats", mock.Anything, uint(5)).Return(
		&model.RoadmapStats{TotalRoadmaps: 3, ActiveRoadmaps: 2, CompletedRoadmaps: 1, TotalSteps: 12, CompletedSteps: 5}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/roadmaps", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/roadmaps", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_roadmaps"].(float64) != 3 {
		t.Errorf("expected total_roadmaps 3, got %v", data["total_roadmaps"])
	}
	stats.AssertExpectations(t)
}

func TestRoadmapStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupRoadmapStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/roadmaps", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/roadmaps", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestRoadmapStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupRoadmapStatsHandler()
	stats.On("GetRoadmapStats", mock.Anything, uint(5)).Return((*model.RoadmapStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/roadmaps", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/roadmaps", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
