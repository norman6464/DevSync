package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNotificationPort は usecase/repository.NotificationReader のモック。
type mockNotificationPort struct{ mock.Mock }

func (m *mockNotificationPort) FindByUserID(ctx context.Context, userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	args := m.Called(ctx, userID, page, limit, notificationType)
	n, _ := args.Get(0).([]model.Notification)
	return n, args.Error(1)
}

func (m *mockNotificationPort) CountByUserID(ctx context.Context, userID uint, notificationType string) (int64, error) {
	args := m.Called(ctx, userID, notificationType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockNotificationPort) CountUnread(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockNotificationPort) MarkAsRead(ctx context.Context, id, userID uint) error {
	return m.Called(ctx, id, userID).Error(0)
}

func (m *mockNotificationPort) MarkAllAsRead(ctx context.Context, userID uint) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *mockNotificationPort) Delete(ctx context.Context, id, userID uint) error {
	return m.Called(ctx, id, userID).Error(0)
}

// newTestNotificationHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestNotificationHandler() (*NotificationHandler, *mockNotificationPort) {
	repo := new(mockNotificationPort)
	h := NewNotificationHandler(
		usecase.NewListNotificationsUseCase(repo),
		usecase.NewCountUnreadNotificationsUseCase(repo),
		usecase.NewMarkNotificationAsReadUseCase(repo),
		usecase.NewMarkAllNotificationsAsReadUseCase(repo),
		usecase.NewDeleteNotificationUseCase(repo),
	)
	return h, repo
}

// ---------- GetAll ----------

func TestNotificationGetAll_Success(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications", h.GetAll)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20, "").
		Return([]model.Notification{{ID: 1, UserID: 1}}, nil)
	repo.On("CountByUserID", mock.Anything, uint(1), "").Return(int64(1), nil)

	w := doRequest(r, http.MethodGet, "/notifications", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Contains(t, body, "notifications")
	assertJSONEqual(t, body, "total", float64(1))
	assertJSONEqual(t, body, "page", float64(1))
	assertJSONEqual(t, body, "limit", float64(20))
	repo.AssertExpectations(t)
}

// 通知種別のフィルタは一覧と総数の両方に効く。
func TestNotificationGetAll_WithTypeFilter(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications", h.GetAll)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20, "follow").Return([]model.Notification{}, nil)
	repo.On("CountByUserID", mock.Anything, uint(1), "follow").Return(int64(0), nil)

	w := doRequest(r, http.MethodGet, "/notifications?type=follow", nil)
	assertStatus(t, w, http.StatusOK)
	assertJSONEqual(t, parseJSON(t, w), "total", float64(0))
	repo.AssertExpectations(t)
}

func TestNotificationGetAll_WithPagination(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications", h.GetAll)

	repo.On("FindByUserID", mock.Anything, uint(1), 3, 5, "").Return([]model.Notification{}, nil)
	repo.On("CountByUserID", mock.Anything, uint(1), "").Return(int64(12), nil)

	w := doRequest(r, http.MethodGet, "/notifications?page=3&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assertJSONEqual(t, body, "page", float64(3))
	assertJSONEqual(t, body, "limit", float64(5))
	assertJSONEqual(t, body, "total", float64(12))
	repo.AssertExpectations(t)
}

func TestNotificationGetAll_ListError(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications", h.GetAll)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20, "").Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notifications", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	// 一覧取得に失敗したら総数は問い合わせない。
	repo.AssertNotCalled(t, "CountByUserID", mock.Anything, mock.Anything, mock.Anything)
}

func TestNotificationGetAll_CountError(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications", h.GetAll)

	repo.On("FindByUserID", mock.Anything, uint(1), 1, 20, "").Return([]model.Notification{}, nil)
	repo.On("CountByUserID", mock.Anything, uint(1), "").Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notifications", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ---------- GetUnreadCount ----------

func TestNotificationGetUnreadCount_Success(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications/unread-count", h.GetUnreadCount)

	repo.On("CountUnread", mock.Anything, uint(1)).Return(int64(7), nil)

	w := doRequest(r, http.MethodGet, "/notifications/unread-count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":7`)
	repo.AssertExpectations(t)
}

func TestNotificationGetUnreadCount_Error(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.GET("/notifications/unread-count", h.GetUnreadCount)

	repo.On("CountUnread", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/notifications/unread-count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ---------- MarkAsRead ----------

func TestNotificationMarkAsRead_Success(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.PUT("/notifications/:id/read", h.MarkAsRead)

	// 本人の通知だけを対象にするため userID も渡す。
	repo.On("MarkAsRead", mock.Anything, uint(5), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/notifications/5/read", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "marked as read")
	repo.AssertExpectations(t)
}

func TestNotificationMarkAsRead_InvalidID(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.PUT("/notifications/:id/read", h.MarkAsRead)

	w := doRequest(r, http.MethodPut, "/notifications/abc/read", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "MarkAsRead", mock.Anything, mock.Anything, mock.Anything)
}

func TestNotificationMarkAsRead_Error(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.PUT("/notifications/:id/read", h.MarkAsRead)

	repo.On("MarkAsRead", mock.Anything, uint(5), uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/notifications/5/read", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ---------- MarkAllAsRead ----------

func TestNotificationMarkAllAsRead_Success(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.PUT("/notifications/read-all", h.MarkAllAsRead)

	repo.On("MarkAllAsRead", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/notifications/read-all", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "all marked as read")
	repo.AssertExpectations(t)
}

func TestNotificationMarkAllAsRead_Error(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.PUT("/notifications/read-all", h.MarkAllAsRead)

	repo.On("MarkAllAsRead", mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/notifications/read-all", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ---------- Delete ----------

func TestNotificationDelete_Success(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.DELETE("/notifications/:id", h.Delete)

	repo.On("Delete", mock.Anything, uint(9), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, fmt.Sprintf("/notifications/%d", 9), nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "deleted")
	repo.AssertExpectations(t)
}

func TestNotificationDelete_InvalidID(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.DELETE("/notifications/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/notifications/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

func TestNotificationDelete_Error(t *testing.T) {
	h, repo := newTestNotificationHandler()
	r := newRouter(1)
	r.DELETE("/notifications/:id", h.Delete)

	repo.On("Delete", mock.Anything, uint(9), uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/notifications/9", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
