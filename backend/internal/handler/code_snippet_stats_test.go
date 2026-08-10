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

// mockCodeSnippetStatsRepo は usecase/repository.CodeSnippetStatsRepository のモック（ctx 付き）。
type mockCodeSnippetStatsRepo struct{ mock.Mock }

func (m *mockCodeSnippetStatsRepo) GetCodeSnippetStats(ctx context.Context, userID uint) (*model.CodeSnippetStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.CodeSnippetStats)
	return s, args.Error(1)
}

func setupCodeSnippetStatsHandler() (*CodeSnippetStatsHandler, *mockCodeSnippetStatsRepo) {
	stats := new(mockCodeSnippetStatsRepo)
	h := NewCodeSnippetStatsHandler(usecase.NewGetCodeSnippetStatsUseCase(stats))
	return h, stats
}

func TestCodeSnippetStats_GetStats_Success(t *testing.T) {
	h, stats := setupCodeSnippetStatsHandler()
	stats.On("GetCodeSnippetStats", mock.Anything, uint(5)).Return(
		&model.CodeSnippetStats{TotalSnippets: 4, TotalComments: 9, LanguageCount: 3}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/snippets", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/snippets", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_snippets"].(float64) != 4 {
		t.Errorf("expected total_snippets 4, got %v", data["total_snippets"])
	}
	stats.AssertExpectations(t)
}

func TestCodeSnippetStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupCodeSnippetStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/snippets", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/snippets", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCodeSnippetStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupCodeSnippetStatsHandler()
	stats.On("GetCodeSnippetStats", mock.Anything, uint(5)).Return((*model.CodeSnippetStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/snippets", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/snippets", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
