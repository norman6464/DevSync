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

// mockFollowStatsRepo は usecase/repository.FollowStatsRepository のモック（ctx 付き）。
type mockFollowStatsRepo struct{ mock.Mock }

func (m *mockFollowStatsRepo) GetFollowStats(ctx context.Context, userID uint) (*model.FollowStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.FollowStats)
	return s, args.Error(1)
}

func setupFollowStatsHandler() (*FollowStatsHandler, *mockFollowStatsRepo) {
	stats := new(mockFollowStatsRepo)
	h := NewFollowStatsHandler(usecase.NewGetFollowStatsUseCase(stats))
	return h, stats
}

func TestFollowStats_GetStats_Success(t *testing.T) {
	h, stats := setupFollowStatsHandler()
	stats.On("GetFollowStats", mock.Anything, uint(5)).Return(
		&model.FollowStats{FollowerCount: 3, FollowingCount: 2}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/follows", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/follows", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["follower_count"].(float64) != 3 {
		t.Errorf("expected follower_count 3, got %v", data["follower_count"])
	}
	stats.AssertExpectations(t)
}

func TestFollowStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupFollowStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/follows", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/follows", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestFollowStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupFollowStatsHandler()
	stats.On("GetFollowStats", mock.Anything, uint(5)).Return((*model.FollowStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/follows", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/follows", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
