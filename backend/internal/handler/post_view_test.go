package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPostViewRepo は usecase/repository.PostViewRepository のモック（ctx 付き）。
type mockPostViewRepo struct{ mock.Mock }

func (m *mockPostViewRepo) RecordView(ctx context.Context, view *model.PostView) error {
	return m.Called(ctx, view).Error(0)
}
func (m *mockPostViewRepo) GetViewCount(ctx context.Context, postID uint) (int64, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockPostViewRepo) HasViewed(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}
func (m *mockPostViewRepo) GetMostViewed(ctx context.Context, limit int) ([]model.ViewCount, error) {
	args := m.Called(ctx, limit)
	vc, _ := args.Get(0).([]model.ViewCount)
	return vc, args.Error(1)
}

// setupPostViewHandler は本物の usecase + port モックで PostViewHandler を組む。
func setupPostViewHandler() (*PostViewHandler, *mockPostViewRepo) {
	views := new(mockPostViewRepo)
	h := NewPostViewHandler(
		usecase.NewRecordPostViewUseCase(views),
		usecase.NewGetPostViewCountUseCase(views),
		usecase.NewGetMostViewedPostsUseCase(views),
	)
	return h, views
}

// --- RecordView ---

func TestPostViewRecordView_Success(t *testing.T) {
	h, views := setupPostViewHandler()
	views.On("HasViewed", mock.Anything, uint(1), uint(5)).Return(false, nil)
	views.On("RecordView", mock.Anything, mock.AnythingOfType("*model.PostView")).Return(nil)

	r := newRouter(1)
	r.POST("/posts/:postId/views", h.RecordView)

	w := doRequest(r, http.MethodPost, "/posts/5/views", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, "記録しました", body["message"])
	views.AssertExpectations(t)
}

func TestPostViewRecordView_InvalidID(t *testing.T) {
	h, _ := setupPostViewHandler()

	r := newRouter(1)
	r.POST("/posts/:postId/views", h.RecordView)

	w := doRequest(r, http.MethodPost, "/posts/abc/views", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostViewRecordView_ServiceError(t *testing.T) {
	h, views := setupPostViewHandler()
	views.On("HasViewed", mock.Anything, uint(1), uint(5)).Return(false, nil)
	views.On("RecordView", mock.Anything, mock.AnythingOfType("*model.PostView")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/posts/:postId/views", h.RecordView)

	w := doRequest(r, http.MethodPost, "/posts/5/views", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	views.AssertExpectations(t)
}

// --- GetViewCount ---

func TestPostViewGetViewCount_Success(t *testing.T) {
	h, views := setupPostViewHandler()
	views.On("GetViewCount", mock.Anything, uint(5)).Return(int64(42), nil)

	r := newRouter(1)
	r.GET("/posts/:postId/view-count", h.GetViewCount)

	w := doRequest(r, http.MethodGet, "/posts/5/view-count", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, float64(42), body["view_count"])
	views.AssertExpectations(t)
}

func TestPostViewGetViewCount_InvalidID(t *testing.T) {
	h, _ := setupPostViewHandler()

	r := newRouter(1)
	r.GET("/posts/:postId/view-count", h.GetViewCount)

	w := doRequest(r, http.MethodGet, "/posts/abc/view-count", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostViewGetViewCount_ServiceError(t *testing.T) {
	h, views := setupPostViewHandler()
	views.On("GetViewCount", mock.Anything, uint(5)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/:postId/view-count", h.GetViewCount)

	w := doRequest(r, http.MethodGet, "/posts/5/view-count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	views.AssertExpectations(t)
}

// --- GetMostViewed ---

func TestPostViewGetMostViewed_Success(t *testing.T) {
	h, views := setupPostViewHandler()
	viewCounts := []model.ViewCount{
		{PostID: 1, Count: 100},
		{PostID: 2, Count: 50},
	}
	views.On("GetMostViewed", mock.Anything, 20).Return(viewCounts, nil)

	r := newRouter(1)
	r.GET("/posts/trending/most-viewed", h.GetMostViewed)

	w := doRequest(r, http.MethodGet, "/posts/trending/most-viewed", nil)
	assertStatus(t, w, http.StatusOK)
	views.AssertExpectations(t)
}

func TestPostViewGetMostViewed_ServiceError(t *testing.T) {
	h, views := setupPostViewHandler()
	views.On("GetMostViewed", mock.Anything, 20).Return([]model.ViewCount(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/trending/most-viewed", h.GetMostViewed)

	w := doRequest(r, http.MethodGet, "/posts/trending/most-viewed", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	views.AssertExpectations(t)
}
