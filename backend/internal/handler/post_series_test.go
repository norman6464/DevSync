package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockPostSeriesService は PostSeriesServiceInterface のモック実装。
type MockPostSeriesService struct{ mock.Mock }

func (m *MockPostSeriesService) Create(series *model.PostSeries) error {
	return m.Called(series).Error(0)
}
func (m *MockPostSeriesService) GetByID(id uint) (*model.PostSeries, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*model.PostSeries), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostSeriesService) GetByUserID(userID uint) ([]model.PostSeries, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.([]model.PostSeries), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostSeriesService) Update(id, userID uint, updates *model.PostSeries) (*model.PostSeries, error) {
	args := m.Called(id, userID, updates)
	if v := args.Get(0); v != nil {
		return v.(*model.PostSeries), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostSeriesService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockPostSeriesService) AddPost(seriesID, postID uint, orderIndex int, userID uint) error {
	return m.Called(seriesID, postID, orderIndex, userID).Error(0)
}
func (m *MockPostSeriesService) RemovePost(seriesID, postID, userID uint) error {
	return m.Called(seriesID, postID, userID).Error(0)
}
func (m *MockPostSeriesService) GetPosts(seriesID uint) ([]model.PostSeriesItem, error) {
	args := m.Called(seriesID)
	if v := args.Get(0); v != nil {
		return v.([]model.PostSeriesItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupPostSeriesHandler() (*PostSeriesHandler, *MockPostSeriesService) {
	svc := new(MockPostSeriesService)
	h := NewPostSeriesHandler(svc)
	return h, svc
}

// ============================================================
// Create テスト
// ============================================================

func TestPostSeriesCreate_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("Create", mock.AnythingOfType("*model.PostSeries")).Return(nil)

	r := newRouter(1)
	r.POST("/series", h.Create)

	w := doRequest(r, http.MethodPost, "/series", map[string]interface{}{
		"title":       "Go入門シリーズ",
		"description": "Go言語の基礎を学ぶシリーズ",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestPostSeriesCreate_InvalidBody(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.POST("/series", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/series", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesCreate_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("Create", mock.AnythingOfType("*model.PostSeries")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/series", h.Create)

	w := doRequest(r, http.MethodPost, "/series", map[string]interface{}{
		"title": "Go入門シリーズ",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostSeriesGetByID_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	series := &model.PostSeries{Title: "Go入門シリーズ", UserID: 1}
	svc.On("GetByID", uint(5)).Return(series, nil)

	r := newRouter(1)
	r.GET("/series/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/series/5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostSeriesGetByID_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.GET("/series/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/series/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesGetByID_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("GetByID", uint(5)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/series/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/series/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestPostSeriesGetByUserID_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	seriesList := []model.PostSeries{{Title: "Go入門シリーズ", UserID: 1}}
	svc.On("GetByUserID", uint(1)).Return(seriesList, nil)

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/series", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostSeriesGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/series", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesGetByUserID_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("GetByUserID", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/series", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestPostSeriesUpdate_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	updated := &model.PostSeries{Title: "Updated", UserID: 1}
	svc.On("Update", uint(5), uint(1), mock.AnythingOfType("*model.PostSeries")).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/series/5", map[string]interface{}{
		"title":       "Updated",
		"description": "新しい説明",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostSeriesUpdate_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/series/abc", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesUpdate_InvalidBody(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/series/5", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesUpdate_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("Update", uint(5), uint(1), mock.AnythingOfType("*model.PostSeries")).Return(nil, errors.New("forbidden"))

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/series/5", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestPostSeriesDelete_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("Delete", uint(5), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostSeriesDelete_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesDelete_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("Delete", uint(5), uint(1)).Return(errors.New("forbidden"))

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetPosts テスト
// ============================================================

func TestPostSeriesGetPosts_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	items := []model.PostSeriesItem{{SeriesID: 5, PostID: 10, OrderIndex: 1}}
	svc.On("GetPosts", uint(5)).Return(items, nil)

	r := newRouter(1)
	r.GET("/series/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/series/5/posts", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostSeriesGetPosts_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.GET("/series/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/series/abc/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesGetPosts_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("GetPosts", uint(5)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/series/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/series/5/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// AddPost テスト
// ============================================================

func TestPostSeriesAddPost_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("AddPost", uint(5), uint(10), 1, uint(1)).Return(nil)

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/series/5/posts", map[string]interface{}{
		"post_id":     10,
		"order_index": 1,
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestPostSeriesAddPost_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/series/abc/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesAddPost_InvalidBody(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequestRaw(r, http.MethodPost, "/series/5/posts", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesAddPost_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("AddPost", uint(5), uint(10), 0, uint(1)).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/series/5/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// RemovePost テスト
// ============================================================

func TestPostSeriesRemovePost_Success(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("RemovePost", uint(5), uint(10), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/series/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/series/5/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostSeriesRemovePost_InvalidSeriesID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.DELETE("/series/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/series/abc/posts/10", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesRemovePost_InvalidPostID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.DELETE("/series/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/series/5/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesRemovePost_ServiceError(t *testing.T) {
	h, svc := setupPostSeriesHandler()
	svc.On("RemovePost", uint(5), uint(10), uint(1)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/series/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/series/5/posts/10", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
