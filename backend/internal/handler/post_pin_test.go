package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostPinService は PostPinServiceInterface のモック実装。
type MockPostPinService struct{ mock.Mock }

func (m *MockPostPinService) Pin(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostPinService) Unpin(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostPinService) GetByUserID(userID uint) ([]model.PostPin, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.([]model.PostPin), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostPinService) Reorder(userID uint, postIDs []uint) error {
	return m.Called(userID, postIDs).Error(0)
}
func (m *MockPostPinService) IsPinned(userID, postID uint) (bool, error) {
	args := m.Called(userID, postID)
	return args.Bool(0), args.Error(1)
}

func setupPostPinHandler() (*PostPinHandler, *MockPostPinService) {
	svc := new(MockPostPinService)
	h := NewPostPinHandler(svc)
	return h, svc
}

// ============================================================
// Pin テスト
// ============================================================

func TestPostPin_Success(t *testing.T) {
	h, svc := setupPostPinHandler()
	svc.On("Pin", uint(1), uint(5)).Return(nil)

	r := newRouter(1)
	r.POST("/posts/:postId/pin", h.Pin)

	w := doRequest(r, http.MethodPost, "/posts/5/pin", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostPin_InvalidID(t *testing.T) {
	h, _ := setupPostPinHandler()

	r := newRouter(1)
	r.POST("/posts/:postId/pin", h.Pin)

	w := doRequest(r, http.MethodPost, "/posts/abc/pin", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPin_ServiceError(t *testing.T) {
	h, svc := setupPostPinHandler()
	svc.On("Pin", uint(1), uint(5)).Return(errors.New("already pinned"))

	r := newRouter(1)
	r.POST("/posts/:postId/pin", h.Pin)

	w := doRequest(r, http.MethodPost, "/posts/5/pin", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Unpin テスト
// ============================================================

func TestPostUnpin_Success(t *testing.T) {
	h, svc := setupPostPinHandler()
	svc.On("Unpin", uint(1), uint(5)).Return(nil)

	r := newRouter(1)
	r.DELETE("/posts/:postId/pin", h.Unpin)

	w := doRequest(r, http.MethodDelete, "/posts/5/pin", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostUnpin_InvalidID(t *testing.T) {
	h, _ := setupPostPinHandler()

	r := newRouter(1)
	r.DELETE("/posts/:postId/pin", h.Unpin)

	w := doRequest(r, http.MethodDelete, "/posts/abc/pin", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestPostPinGetByUserID_Success(t *testing.T) {
	h, svc := setupPostPinHandler()
	pins := []model.PostPin{
		{UserID: 1, PostID: 10, PinOrder: 1},
		{UserID: 1, PostID: 20, PinOrder: 2},
	}
	svc.On("GetByUserID", uint(1)).Return(pins, nil)

	r := newRouter(1)
	r.GET("/users/:userId/pins", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/pins", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotNil(t, body["pins"])
	svc.AssertExpectations(t)
}

func TestPostPinGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupPostPinHandler()

	r := newRouter(1)
	r.GET("/users/:userId/pins", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/pins", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPinGetByUserID_ServiceError(t *testing.T) {
	h, svc := setupPostPinHandler()
	svc.On("GetByUserID", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/pins", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/pins", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Reorder テスト
// ============================================================

func TestPostPinReorder_Success(t *testing.T) {
	h, svc := setupPostPinHandler()
	svc.On("Reorder", uint(1), []uint{10, 20, 30}).Return(nil)

	r := newRouter(1)
	r.PUT("/pins/reorder", h.Reorder)

	w := doRequest(r, http.MethodPut, "/pins/reorder", map[string]interface{}{
		"post_ids": []uint{10, 20, 30},
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostPinReorder_InvalidBody(t *testing.T) {
	h, _ := setupPostPinHandler()

	r := newRouter(1)
	r.PUT("/pins/reorder", h.Reorder)

	w := doRequestRaw(r, http.MethodPut, "/pins/reorder", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}
