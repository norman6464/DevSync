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

// mockBookmarkStatsRepo は usecase/repository.BookmarkStatsRepository のモック（ctx 付き）。
type mockBookmarkStatsRepo struct{ mock.Mock }

func (m *mockBookmarkStatsRepo) GetBookmarkStats(ctx context.Context, userID uint) (*model.BookmarkStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.BookmarkStats)
	return s, args.Error(1)
}

func setupBookmarkStatsHandler() (*BookmarkStatsHandler, *mockBookmarkStatsRepo) {
	stats := new(mockBookmarkStatsRepo)
	h := NewBookmarkStatsHandler(usecase.NewGetBookmarkStatsUseCase(stats))
	return h, stats
}

func TestBookmarkStats_GetStats_Success(t *testing.T) {
	h, stats := setupBookmarkStatsHandler()
	stats.On("GetBookmarkStats", mock.Anything, uint(5)).Return(
		&model.BookmarkStats{TotalBookmarksMade: 4, TotalBookmarksReceived: 2, BookmarksThisMonth: 3}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/bookmarks", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/bookmarks", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_bookmarks_made"].(float64) != 4 {
		t.Errorf("expected total_bookmarks_made 4, got %v", data["total_bookmarks_made"])
	}
	stats.AssertExpectations(t)
}

func TestBookmarkStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupBookmarkStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/bookmarks", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/bookmarks", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupBookmarkStatsHandler()
	stats.On("GetBookmarkStats", mock.Anything, uint(5)).Return((*model.BookmarkStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/bookmarks", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/bookmarks", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
