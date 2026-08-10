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

// mockProjectStatsRepo は usecase/repository.ProjectStatsRepository のモック（ctx 付き）。
type mockProjectStatsRepo struct{ mock.Mock }

func (m *mockProjectStatsRepo) GetProjectStats(ctx context.Context, userID uint) (*model.ProjectStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ProjectStats)
	return s, args.Error(1)
}

func setupProjectStatsHandler() (*ProjectStatsHandler, *mockProjectStatsRepo) {
	stats := new(mockProjectStatsRepo)
	h := NewProjectStatsHandler(usecase.NewGetProjectStatsUseCase(stats))
	return h, stats
}

func TestProjectStats_GetStats_Success(t *testing.T) {
	h, stats := setupProjectStatsHandler()
	stats.On("GetProjectStats", mock.Anything, uint(5)).Return(
		&model.ProjectStats{TotalProjects: 4, FeaturedProjects: 1, OngoingProjects: 3, CompletedProjects: 1}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/projects", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/projects", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_projects"].(float64) != 4 {
		t.Errorf("expected total_projects 4, got %v", data["total_projects"])
	}
	stats.AssertExpectations(t)
}

func TestProjectStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupProjectStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/projects", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/projects", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestProjectStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupProjectStatsHandler()
	stats.On("GetProjectStats", mock.Anything, uint(5)).Return((*model.ProjectStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/projects", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/projects", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
