package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostViewService は PostViewServiceInterface のモック実装。
type MockPostViewService struct{ mock.Mock }

func (m *MockPostViewService) RecordView(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostViewService) GetViewCount(postID uint) (int64, error) {
	args := m.Called(postID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockPostViewService) HasViewed(userID, postID uint) (bool, error) {
	args := m.Called(userID, postID)
	return args.Bool(0), args.Error(1)
}
func (m *MockPostViewService) GetMostViewed(limit int) ([]model.ViewCount, error) {
	args := m.Called(limit)
	if v := args.Get(0); v != nil {
		return v.([]model.ViewCount), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupPostViewHandler() (*PostViewHandler, *MockPostViewService) {
	svc := new(MockPostViewService)
	h := NewPostViewHandler(svc)
	return h, svc
}

// ============================================================
// RecordView テスト
// ============================================================

func TestPostViewRecordView_Success(t *testing.T) {
	h, svc := setupPostViewHandler()
	svc.On("RecordView", uint(1), uint(5)).Return(nil)

	r := newRouter(1)
	r.POST("/posts/:postId/views", h.RecordView)

	w := doRequest(r, http.MethodPost, "/posts/5/views", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, "記録しました", body["message"])
	svc.AssertExpectations(t)
}

func TestPostViewRecordView_InvalidID(t *testing.T) {
	h, _ := setupPostViewHandler()

	r := newRouter(1)
	r.POST("/posts/:postId/views", h.RecordView)

	w := doRequest(r, http.MethodPost, "/posts/abc/views", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostViewRecordView_ServiceError(t *testing.T) {
	h, svc := setupPostViewHandler()
	svc.On("RecordView", uint(1), uint(5)).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/posts/:postId/views", h.RecordView)

	w := doRequest(r, http.MethodPost, "/posts/5/views", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetViewCount テスト
// ============================================================

func TestPostViewGetViewCount_Success(t *testing.T) {
	h, svc := setupPostViewHandler()
	svc.On("GetViewCount", uint(5)).Return(int64(42), nil)

	r := newRouter(1)
	r.GET("/posts/:postId/view-count", h.GetViewCount)

	w := doRequest(r, http.MethodGet, "/posts/5/view-count", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(42), body["view_count"])
	svc.AssertExpectations(t)
}

func TestPostViewGetViewCount_InvalidID(t *testing.T) {
	h, _ := setupPostViewHandler()

	r := newRouter(1)
	r.GET("/posts/:postId/view-count", h.GetViewCount)

	w := doRequest(r, http.MethodGet, "/posts/abc/view-count", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostViewGetViewCount_ServiceError(t *testing.T) {
	h, svc := setupPostViewHandler()
	svc.On("GetViewCount", uint(5)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/:postId/view-count", h.GetViewCount)

	w := doRequest(r, http.MethodGet, "/posts/5/view-count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetMostViewed テスト
// ============================================================

func TestPostViewGetMostViewed_Success(t *testing.T) {
	h, svc := setupPostViewHandler()
	viewCounts := []model.ViewCount{
		{PostID: 1, Count: 100},
		{PostID: 2, Count: 50},
	}
	svc.On("GetMostViewed", 20).Return(viewCounts, nil)

	r := newRouter(1)
	r.GET("/posts/trending/most-viewed", h.GetMostViewed)

	w := doRequest(r, http.MethodGet, "/posts/trending/most-viewed", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostViewGetMostViewed_ServiceError(t *testing.T) {
	h, svc := setupPostViewHandler()
	svc.On("GetMostViewed", 20).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/trending/most-viewed", h.GetMostViewed)

	w := doRequest(r, http.MethodGet, "/posts/trending/most-viewed", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
