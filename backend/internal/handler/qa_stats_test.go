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

// mockQAStatsRepo は usecase/repository.QAStatsRepository のモック（ctx 付き）。
type mockQAStatsRepo struct{ mock.Mock }

func (m *mockQAStatsRepo) GetQAStats(ctx context.Context, userID uint) (*model.QAStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.QAStats)
	return s, args.Error(1)
}

func setupQAStatsHandler() (*QAStatsHandler, *mockQAStatsRepo) {
	stats := new(mockQAStatsRepo)
	h := NewQAStatsHandler(usecase.NewGetQAStatsUseCase(stats))
	return h, stats
}

func TestQAStats_GetStats_Success(t *testing.T) {
	h, stats := setupQAStatsHandler()
	stats.On("GetQAStats", mock.Anything, uint(5)).Return(
		&model.QAStats{TotalQuestions: 2, TotalAnswers: 3, BestAnswerCount: 1, TotalVotesReceived: 12}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/qa", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/qa", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_votes_received"].(float64) != 12 {
		t.Errorf("expected total_votes_received 12, got %v", data["total_votes_received"])
	}
	stats.AssertExpectations(t)
}

func TestQAStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupQAStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/qa", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/qa", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestQAStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupQAStatsHandler()
	stats.On("GetQAStats", mock.Anything, uint(5)).Return((*model.QAStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/qa", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/qa", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
