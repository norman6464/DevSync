package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
)

func TestUserDashboard_GetStats_Success(t *testing.T) {
	h, svc := setupUserDashboardHandler()
	stats := &model.UserDashboardStats{
		PostCount:     10,
		LikesReceived: 50,
		FollowerCount: 20,
	}
	svc.On("GetStats", uint(5)).Return(stats, nil)

	r := newRouter(1)
	r.GET("/users/:id/dashboard-stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/dashboard-stats", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if body["post_count"] != float64(10) {
		t.Errorf("expected post_count=10, got %v", body["post_count"])
	}
	svc.AssertExpectations(t)
}

func TestUserDashboard_GetStats_InvalidID(t *testing.T) {
	h, _ := setupUserDashboardHandler()

	r := newRouter(1)
	r.GET("/users/:id/dashboard-stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/abc/dashboard-stats", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestUserDashboard_GetStats_ServiceError(t *testing.T) {
	h, svc := setupUserDashboardHandler()
	svc.On("GetStats", uint(5)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/dashboard-stats", h.GetStats)
	w := doRequest(r, http.MethodGet, "/users/5/dashboard-stats", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
