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

// mockCommentStatsRepo は usecase/repository.CommentStatsRepository のモック（ctx 付き）。
type mockCommentStatsRepo struct{ mock.Mock }

func (m *mockCommentStatsRepo) GetCommentStats(ctx context.Context, userID uint) (*model.CommentStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.CommentStats)
	return s, args.Error(1)
}

func setupCommentStatsHandler() (*CommentStatsHandler, *mockCommentStatsRepo) {
	stats := new(mockCommentStatsRepo)
	h := NewCommentStatsHandler(usecase.NewGetCommentStatsUseCase(stats))
	return h, stats
}

func TestCommentStats_GetStats_Success(t *testing.T) {
	h, stats := setupCommentStatsHandler()
	stats.On("GetCommentStats", mock.Anything, uint(5)).Return(
		&model.CommentStats{TotalComments: 3, TotalReplies: 2, CommentsReceived: 4, CommentsThisMonth: 5}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/comments", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/comments", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_comments"].(float64) != 3 {
		t.Errorf("expected total_comments 3, got %v", data["total_comments"])
	}
	stats.AssertExpectations(t)
}

func TestCommentStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupCommentStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/comments", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/comments", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupCommentStatsHandler()
	stats.On("GetCommentStats", mock.Anything, uint(5)).Return((*model.CommentStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/comments", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/comments", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
