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

// mockMessageStatsRepo は usecase/repository.MessageStatsRepository のモック（ctx 付き）。
type mockMessageStatsRepo struct{ mock.Mock }

func (m *mockMessageStatsRepo) GetMessageStats(ctx context.Context, userID uint) (*model.MessageStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.MessageStats)
	return s, args.Error(1)
}

func setupMessageStatsHandler() (*MessageStatsHandler, *mockMessageStatsRepo) {
	stats := new(mockMessageStatsRepo)
	h := NewMessageStatsHandler(usecase.NewGetMessageStatsUseCase(stats))
	return h, stats
}

func TestMessageStats_GetStats_Success(t *testing.T) {
	h, stats := setupMessageStatsHandler()
	stats.On("GetMessageStats", mock.Anything, uint(5)).Return(
		&model.MessageStats{TotalSent: 5, TotalReceived: 3, ConversationsCount: 2, MessagesThisMonth: 4}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/messages", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/messages", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["total_sent"].(float64) != 5 {
		t.Errorf("expected total_sent 5, got %v", data["total_sent"])
	}
	stats.AssertExpectations(t)
}

func TestMessageStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupMessageStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/messages", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/messages", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestMessageStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupMessageStatsHandler()
	stats.On("GetMessageStats", mock.Anything, uint(5)).Return((*model.MessageStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/messages", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/messages", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
