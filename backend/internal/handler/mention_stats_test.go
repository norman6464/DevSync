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

// mockMentionStatsRepo は usecase/repository.MentionStatsRepository のモック（ctx 付き）。
type mockMentionStatsRepo struct{ mock.Mock }

func (m *mockMentionStatsRepo) GetMentionStats(ctx context.Context, userID uint) (*model.MentionStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.MentionStats)
	return s, args.Error(1)
}

func setupMentionStatsHandler() (*MentionStatsHandler, *mockMentionStatsRepo) {
	stats := new(mockMentionStatsRepo)
	h := NewMentionStatsHandler(usecase.NewGetMentionStatsUseCase(stats))
	return h, stats
}

func TestMentionStats_GetStats_Success(t *testing.T) {
	h, stats := setupMentionStatsHandler()
	stats.On("GetMentionStats", mock.Anything, uint(5)).Return(
		&model.MentionStats{MentionsReceived: 3, MentionsMade: 2, MentionsThisMonth: 2}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/mentions", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/mentions", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["mentions_received"].(float64) != 3 {
		t.Errorf("expected mentions_received 3, got %v", data["mentions_received"])
	}
	stats.AssertExpectations(t)
}

func TestMentionStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupMentionStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/mentions", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/mentions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestMentionStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupMentionStatsHandler()
	stats.On("GetMentionStats", mock.Anything, uint(5)).Return((*model.MentionStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/mentions", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/mentions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
