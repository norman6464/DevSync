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

// mockPostStatsRepo は usecase/repository.PostStatsRepository のモック（ctx 付き）。
type mockPostStatsRepo struct{ mock.Mock }

func (m *mockPostStatsRepo) GetPostStats(ctx context.Context, userID uint) (*model.PostStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.PostStats)
	return s, args.Error(1)
}

func setupPostStatsHandler() (*PostStatsHandler, *mockPostStatsRepo) {
	stats := new(mockPostStatsRepo)
	h := NewPostStatsHandler(usecase.NewGetPostStatsUseCase(stats))
	return h, stats
}

func TestPostStats_GetStats_Success(t *testing.T) {
	h, stats := setupPostStatsHandler()
	stats.On("GetPostStats", mock.Anything, uint(5)).Return(
		&model.PostStats{TotalPosts: 5, PublishedPosts: 4, DraftPosts: 1, TotalLikesReceived: 12, TotalComments: 7, PostsThisMonth: 2}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/posts", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/posts", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_posts"].(float64) != 5 {
		t.Errorf("expected total_posts 5, got %v", data["total_posts"])
	}
	stats.AssertExpectations(t)
}

func TestPostStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupPostStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/posts", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupPostStatsHandler()
	stats.On("GetPostStats", mock.Anything, uint(5)).Return((*model.PostStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/posts", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
