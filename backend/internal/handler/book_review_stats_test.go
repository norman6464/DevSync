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

// mockBookReviewStatsRepo は usecase/repository.BookReviewStatsRepository のモック（ctx 付き）。
type mockBookReviewStatsRepo struct{ mock.Mock }

func (m *mockBookReviewStatsRepo) GetBookReviewStats(ctx context.Context, userID uint) (*model.BookReviewStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.BookReviewStats)
	return s, args.Error(1)
}

func setupBookReviewStatsHandler() (*BookReviewStatsHandler, *mockBookReviewStatsRepo) {
	stats := new(mockBookReviewStatsRepo)
	h := NewBookReviewStatsHandler(usecase.NewGetBookReviewStatsUseCase(stats))
	return h, stats
}

func TestBookReviewStats_GetStats_Success(t *testing.T) {
	h, stats := setupBookReviewStatsHandler()
	stats.On("GetBookReviewStats", mock.Anything, uint(5)).Return(
		&model.BookReviewStats{TotalReviews: 3, AverageRating: 4, MaxRating: 5, MinRating: 3, FiveStarCount: 1}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/book-reviews", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/book-reviews", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_reviews"].(float64) != 3 {
		t.Errorf("expected total_reviews 3, got %v", data["total_reviews"])
	}
	stats.AssertExpectations(t)
}

func TestBookReviewStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupBookReviewStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/book-reviews", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/book-reviews", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupBookReviewStatsHandler()
	stats.On("GetBookReviewStats", mock.Anything, uint(5)).Return((*model.BookReviewStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/book-reviews", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/book-reviews", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
