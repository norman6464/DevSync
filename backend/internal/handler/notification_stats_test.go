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

// mockNotificationStatsRepo は usecase/repository.NotificationStatsRepository のモック（ctx 付き）。
type mockNotificationStatsRepo struct{ mock.Mock }

func (m *mockNotificationStatsRepo) GetNotificationStats(ctx context.Context, userID uint) (*model.NotificationStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.NotificationStats)
	return s, args.Error(1)
}

func setupNotificationStatsHandler() (*NotificationStatsHandler, *mockNotificationStatsRepo) {
	stats := new(mockNotificationStatsRepo)
	h := NewNotificationStatsHandler(usecase.NewGetNotificationStatsUseCase(stats))
	return h, stats
}

func TestNotificationStats_GetStats_Success(t *testing.T) {
	h, stats := setupNotificationStatsHandler()
	stats.On("GetNotificationStats", mock.Anything, uint(5)).Return(
		&model.NotificationStats{TotalNotifications: 5, UnreadCount: 2, NotificationsThisMonth: 4}, nil,
	)

	r := newRouter(1)
	r.GET("/users/:id/stats/notifications", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/notifications", nil)
	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	if data["unread_count"].(float64) != 2 {
		t.Errorf("expected unread_count 2, got %v", data["unread_count"])
	}
	stats.AssertExpectations(t)
}

func TestNotificationStats_GetStats_InvalidID(t *testing.T) {
	h, _ := setupNotificationStatsHandler()

	r := newRouter(1)
	r.GET("/users/:id/stats/notifications", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/abc/stats/notifications", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestNotificationStats_GetStats_ServiceError(t *testing.T) {
	h, stats := setupNotificationStatsHandler()
	stats.On("GetNotificationStats", mock.Anything, uint(5)).Return((*model.NotificationStats)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/stats/notifications", h.GetStats)

	w := doRequest(r, http.MethodGet, "/users/5/stats/notifications", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	stats.AssertExpectations(t)
}
