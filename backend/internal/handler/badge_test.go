package handler

import (
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/service"
)

// ---------- GetUserBadges ----------

func TestBadgeGetUserBadges_Success(t *testing.T) {
	h, svc := setupBadgeHandler()
	badges := []service.BadgeResult{{ID: "streak_3", Name: "3 Day Streak", Earned: true}}
	svc.On("GetUserBadges", uint(1)).Return(badges, nil)

	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)
	w := doRequest(r, "GET", "/users/1/badges", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestBadgeGetUserBadges_InvalidID(t *testing.T) {
	h, _ := setupBadgeHandler()

	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)
	w := doRequest(r, "GET", "/users/abc/badges", nil)

	assertStatus(t, w, 400)
}

func TestBadgeGetUserBadges_ServiceError(t *testing.T) {
	h, svc := setupBadgeHandler()
	svc.On("GetUserBadges", uint(1)).Return([]service.BadgeResult{}, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/users/:userId/badges", h.GetUserBadges)
	w := doRequest(r, "GET", "/users/1/badges", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- NotifyBadgeEarned ----------

func TestBadgeNotifyBadgeEarned_Success(t *testing.T) {
	h, svc := setupBadgeHandler()
	svc.On("NotifyBadgeEarned", uint(1), "streak_3").Return(nil)

	r := newRouter(1)
	r.POST("/badges/notify", h.NotifyBadgeEarned)
	w := doRequest(r, "POST", "/badges/notify", map[string]interface{}{
		"badge_id": "streak_3",
	})

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestBadgeNotifyBadgeEarned_ValidationError(t *testing.T) {
	h, _ := setupBadgeHandler()

	r := newRouter(1)
	r.POST("/badges/notify", h.NotifyBadgeEarned)
	w := doRequest(r, "POST", "/badges/notify", map[string]interface{}{})

	assertStatus(t, w, 400)
}

func TestBadgeNotifyBadgeEarned_ServiceError(t *testing.T) {
	h, svc := setupBadgeHandler()
	svc.On("NotifyBadgeEarned", uint(1), "invalid").Return(fmt.Errorf("badge not found"))

	r := newRouter(1)
	r.POST("/badges/notify", h.NotifyBadgeEarned)
	w := doRequest(r, "POST", "/badges/notify", map[string]interface{}{
		"badge_id": "invalid",
	})

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}
