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

// mockUserDashboardRepo は usecase/repository.UserDashboardRepository のモック（ctx 付き）。
type mockUserDashboardRepo struct{ mock.Mock }

func (m *mockUserDashboardRepo) GetDashboardStats(ctx context.Context, userID uint) (*model.UserDashboardStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.UserDashboardStats)
	return s, args.Error(1)
}

func setupUserDashboardHandler() (*UserDashboardHandler, *mockUserDashboardRepo) {
	repo := new(mockUserDashboardRepo)
	h := NewUserDashboardHandler(usecase.NewGetUserDashboardStatsUseCase(repo))
	return h, repo
}

func TestUserDashboard_GetStats_Success(t *testing.T) {
	h, repo := setupUserDashboardHandler()
	repo.On("GetDashboardStats", mock.Anything, uint(5)).Return(
		&model.UserDashboardStats{PostCount: 10, LikesReceived: 50, FollowerCount: 20}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/dashboard-stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/dashboard-stats", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if body["post_count"] != float64(10) {
		t.Errorf("expected post_count=10, got %v", body["post_count"])
	}
	repo.AssertExpectations(t)
}

func TestUserDashboard_GetStats_InvalidID(t *testing.T) {
	h, _ := setupUserDashboardHandler()

	r := newRouter(1)
	r.GET("/users/:id/dashboard-stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/dashboard-stats", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestUserDashboard_GetStats_ServiceError(t *testing.T) {
	h, repo := setupUserDashboardHandler()
	repo.On("GetDashboardStats", mock.Anything, uint(5)).Return((*model.UserDashboardStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/dashboard-stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/dashboard-stats", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
