package handler

import (
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// ---------- GetAll ----------

func TestNotificationGetAll_Success(t *testing.T) {
	h, svc := setupNotificationHandler()
	notifs := []model.Notification{{ID: 1, UserID: 1}}
	svc.On("GetByUserID", uint(1), 1, 20, "").Return(notifs, nil)
	svc.On("CountByUserID", uint(1), "").Return(int64(1), nil)

	r := newRouter(1)
	r.GET("/notifications", h.GetAll)
	w := doRequest(r, "GET", "/notifications", nil)
	assertStatus(t, w, 200)

	body := parseJSON(t, w)
	_, hasNotifs := body["notifications"]
	if !hasNotifs {
		t.Error("レスポンスに notifications フィールドがない")
	}
	svc.AssertExpectations(t)
}

func TestNotificationGetAll_WithTypeFilter(t *testing.T) {
	h, svc := setupNotificationHandler()
	notifs := []model.Notification{}
	svc.On("GetByUserID", uint(1), 1, 20, "follow").Return(notifs, nil)
	svc.On("CountByUserID", uint(1), "follow").Return(int64(0), nil)

	r := newRouter(1)
	r.GET("/notifications", h.GetAll)
	w := doRequest(r, "GET", "/notifications?type=follow", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestNotificationGetAll_ServiceError(t *testing.T) {
	h, svc := setupNotificationHandler()
	notifs := []model.Notification{}
	svc.On("GetByUserID", uint(1), 1, 20, "").Return(notifs, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/notifications", h.GetAll)
	w := doRequest(r, "GET", "/notifications", nil)
	assertStatus(t, w, 500)
}

// ---------- GetUnreadCount ----------

func TestNotificationGetUnreadCount_Success(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("CountUnread", uint(1)).Return(int64(5), nil)

	r := newRouter(1)
	r.GET("/notifications/unread", h.GetUnreadCount)
	w := doRequest(r, "GET", "/notifications/unread", nil)
	assertStatus(t, w, 200)

	body := parseJSON(t, w)
	count, ok := body["count"].(float64)
	if !ok || count != 5 {
		t.Errorf("count should be 5, got %v", body["count"])
	}
}

func TestNotificationGetUnreadCount_Error(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("CountUnread", uint(1)).Return(int64(0), fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/notifications/unread", h.GetUnreadCount)
	w := doRequest(r, "GET", "/notifications/unread", nil)
	assertStatus(t, w, 500)
}

// ---------- MarkAsRead ----------

func TestNotificationMarkAsRead_Success(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("MarkAsRead", uint(10), uint(1)).Return(nil)

	r := newRouter(1)
	r.PATCH("/notifications/:id/read", h.MarkAsRead)
	w := doRequest(r, "PATCH", "/notifications/10/read", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestNotificationMarkAsRead_InvalidID(t *testing.T) {
	h, _ := setupNotificationHandler()

	r := newRouter(1)
	r.PATCH("/notifications/:id/read", h.MarkAsRead)
	w := doRequest(r, "PATCH", "/notifications/abc/read", nil)
	assertStatus(t, w, 400)
}

func TestNotificationMarkAsRead_NotFound(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("MarkAsRead", uint(999), uint(1)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.PATCH("/notifications/:id/read", h.MarkAsRead)
	w := doRequest(r, "PATCH", "/notifications/999/read", nil)
	assertStatus(t, w, 404)
}

// ---------- MarkAllAsRead ----------

func TestNotificationMarkAllAsRead_Success(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("MarkAllAsRead", uint(1)).Return(nil)

	r := newRouter(1)
	r.PATCH("/notifications/read-all", h.MarkAllAsRead)
	w := doRequest(r, "PATCH", "/notifications/read-all", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestNotificationMarkAllAsRead_Error(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("MarkAllAsRead", uint(1)).Return(fmt.Errorf("internal error"))

	r := newRouter(1)
	r.PATCH("/notifications/read-all", h.MarkAllAsRead)
	w := doRequest(r, "PATCH", "/notifications/read-all", nil)
	assertStatus(t, w, 500)
}

// ---------- Delete ----------

func TestNotificationDelete_Success(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("Delete", uint(10), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/notifications/:id", h.Delete)
	w := doRequest(r, "DELETE", "/notifications/10", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestNotificationDelete_NotFound(t *testing.T) {
	h, svc := setupNotificationHandler()
	svc.On("Delete", uint(999), uint(1)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.DELETE("/notifications/:id", h.Delete)
	w := doRequest(r, "DELETE", "/notifications/999", nil)
	assertStatus(t, w, 404)
}

func TestNotificationDelete_InvalidID(t *testing.T) {
	h, _ := setupNotificationHandler()

	r := newRouter(1)
	r.DELETE("/notifications/:id", h.Delete)
	w := doRequest(r, "DELETE", "/notifications/abc", nil)
	assertStatus(t, w, 400)
}
