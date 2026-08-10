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

// mockPostSeriesRepo は usecase/repository.PostSeriesRepository のモック（ctx 付き）。
type mockPostSeriesRepo struct{ mock.Mock }

func (m *mockPostSeriesRepo) Create(ctx context.Context, series *model.PostSeries) error {
	return m.Called(ctx, series).Error(0)
}

func (m *mockPostSeriesRepo) FindByID(ctx context.Context, id uint) (*model.PostSeries, error) {
	args := m.Called(ctx, id)
	s, _ := args.Get(0).(*model.PostSeries)
	return s, args.Error(1)
}

func (m *mockPostSeriesRepo) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]model.PostSeries, error) {
	args := m.Called(ctx, userID, offset, limit)
	s, _ := args.Get(0).([]model.PostSeries)
	return s, args.Error(1)
}

func (m *mockPostSeriesRepo) CountByUser(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostSeriesRepo) Update(ctx context.Context, series *model.PostSeries) error {
	return m.Called(ctx, series).Error(0)
}

func (m *mockPostSeriesRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostSeriesRepo) AddPost(ctx context.Context, item *model.PostSeriesItem) error {
	return m.Called(ctx, item).Error(0)
}

func (m *mockPostSeriesRepo) RemovePost(ctx context.Context, seriesID, postID uint) error {
	return m.Called(ctx, seriesID, postID).Error(0)
}

func (m *mockPostSeriesRepo) HasPost(ctx context.Context, seriesID, postID uint) (bool, error) {
	args := m.Called(ctx, seriesID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostSeriesRepo) GetPostsBySeriesID(ctx context.Context, seriesID uint) ([]model.PostSeriesItem, error) {
	args := m.Called(ctx, seriesID)
	i, _ := args.Get(0).([]model.PostSeriesItem)
	return i, args.Error(1)
}

// setupPostSeriesHandler は本物の usecase と port モックで PostSeriesHandler を組む。
func setupPostSeriesHandler() (*PostSeriesHandler, *mockPostSeriesRepo) {
	repo := new(mockPostSeriesRepo)
	h := NewPostSeriesHandler(
		usecase.NewCreatePostSeriesUseCase(repo),
		usecase.NewGetPostSeriesUseCase(repo),
		usecase.NewListPostSeriesUseCase(repo),
		usecase.NewCountPostSeriesUseCase(repo),
		usecase.NewUpdatePostSeriesUseCase(repo),
		usecase.NewDeletePostSeriesUseCase(repo),
		usecase.NewAddPostToSeriesUseCase(repo),
		usecase.NewRemovePostFromSeriesUseCase(repo),
		usecase.NewListPostSeriesPostsUseCase(repo),
	)
	return h, repo
}

// ownedSeries は認証ユーザー（userID=1）が所有するシリーズを返す。
func ownedSeries() *model.PostSeries {
	return &model.PostSeries{Title: "Go入門シリーズ", UserID: 1}
}

// ============================================================
// Create テスト
// ============================================================

func TestPostSeriesCreate_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostSeries")).Return(nil)

	r := newRouter(1)
	r.POST("/series", h.Create)

	w := doRequest(r, http.MethodPost, "/series", map[string]interface{}{
		"title":       "Go入門シリーズ",
		"description": "Go言語の基礎を学ぶシリーズ",
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

func TestPostSeriesCreate_InvalidBody(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.POST("/series", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/series", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesCreate_ServiceError(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostSeries")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/series", h.Create)

	w := doRequest(r, http.MethodPost, "/series", map[string]interface{}{
		"title": "Go入門シリーズ",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostSeriesGetByID_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)

	r := newRouter(1)
	r.GET("/series/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/series/5", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByID_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.GET("/series/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/series/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesGetByID_ServiceError(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return((*model.PostSeries)(nil), errors.New("not found"))

	r := newRouter(1)
	r.GET("/series/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/series/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestPostSeriesGetByUserID_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	seriesList := []model.PostSeries{*ownedSeries()}
	// page=1 / limit=20 は offset=0 に変換される
	repo.On("FindByUserID", mock.Anything, uint(1), 0, 20).Return(seriesList, nil)
	repo.On("CountByUser", mock.Anything, uint(1)).Return(int64(1), nil)

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/series", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/series", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesGetByUserID_ServiceError(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByUserID", mock.Anything, uint(1), 0, 20).Return([]model.PostSeries(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/series", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetByUserID_CountError(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	seriesList := []model.PostSeries{*ownedSeries()}
	repo.On("FindByUserID", mock.Anything, uint(1), 0, 20).Return(seriesList, nil)
	repo.On("CountByUser", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/series", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/series", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestPostSeriesUpdate_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.PostSeries")).Return(nil)

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/series/5", map[string]interface{}{
		"title":       "Updated",
		"description": "新しい説明",
	})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 他人のシリーズは更新できない（403）。
func TestPostSeriesUpdate_Forbidden(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{Title: "他人の", UserID: 99}, nil)

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/series/5", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Update")
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
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.PostSeries")).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/series/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/series/5", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestPostSeriesDelete_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("Delete", mock.Anything, uint(5)).Return(nil)

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/5", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 他人のシリーズは削除できない（403）。
func TestPostSeriesDelete_Forbidden(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(&model.PostSeries{Title: "他人の", UserID: 99}, nil)

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/5", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Delete")
}

func TestPostSeriesDelete_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesDelete_ServiceError(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("Delete", mock.Anything, uint(5)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/series/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/series/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetPosts テスト
// ============================================================

func TestPostSeriesGetPosts_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	items := []model.PostSeriesItem{{SeriesID: 5, PostID: 10, OrderIndex: 1}}
	repo.On("GetPostsBySeriesID", mock.Anything, uint(5)).Return(items, nil)

	r := newRouter(1)
	r.GET("/series/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/series/5/posts", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostSeriesGetPosts_InvalidID(t *testing.T) {
	h, _ := setupPostSeriesHandler()

	r := newRouter(1)
	r.GET("/series/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/series/abc/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostSeriesGetPosts_ServiceError(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("GetPostsBySeriesID", mock.Anything, uint(5)).Return([]model.PostSeriesItem(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/series/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/series/5/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// AddPost テスト
// ============================================================

func TestPostSeriesAddPost_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
	repo.On("AddPost", mock.Anything, mock.AnythingOfType("*model.PostSeriesItem")).Return(nil)

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/series/5/posts", map[string]interface{}{
		"post_id":     10,
		"order_index": 1,
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

// すでに追加済みの投稿は 400 で拒否する。
func TestPostSeriesAddPost_Duplicate(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(true, nil)

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/series/5/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "AddPost")
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
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
	repo.On("AddPost", mock.Anything, mock.AnythingOfType("*model.PostSeriesItem")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/series/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/series/5/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// RemovePost テスト
// ============================================================

func TestPostSeriesRemovePost_Success(t *testing.T) {
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("RemovePost", mock.Anything, uint(5), uint(10)).Return(nil)

	r := newRouter(1)
	r.DELETE("/series/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/series/5/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
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
	h, repo := setupPostSeriesHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedSeries(), nil)
	repo.On("RemovePost", mock.Anything, uint(5), uint(10)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/series/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/series/5/posts/10", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
