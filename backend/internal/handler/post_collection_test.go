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

// mockPostCollectionRepo は usecase/repository.PostCollectionRepository のモック（ctx 付き）。
type mockPostCollectionRepo struct{ mock.Mock }

func (m *mockPostCollectionRepo) Create(ctx context.Context, collection *model.PostCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockPostCollectionRepo) FindByID(ctx context.Context, id uint) (*model.PostCollection, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.PostCollection)
	return c, args.Error(1)
}

func (m *mockPostCollectionRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	c, _ := args.Get(0).([]model.PostCollection)
	return c, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostCollectionRepo) FindPublicByUserID(ctx context.Context, userID uint) ([]model.PostCollection, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.PostCollection)
	return c, args.Error(1)
}

func (m *mockPostCollectionRepo) Update(ctx context.Context, collection *model.PostCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockPostCollectionRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostCollectionRepo) AddPost(ctx context.Context, item *model.PostCollectionItem) error {
	return m.Called(ctx, item).Error(0)
}

func (m *mockPostCollectionRepo) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return m.Called(ctx, collectionID, postID).Error(0)
}

func (m *mockPostCollectionRepo) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	args := m.Called(ctx, collectionID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPostCollectionRepo) GetPostsByCollectionID(ctx context.Context, collectionID uint) ([]model.PostCollectionItem, error) {
	args := m.Called(ctx, collectionID)
	i, _ := args.Get(0).([]model.PostCollectionItem)
	return i, args.Error(1)
}

func (m *mockPostCollectionRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupPostCollectionHandler は本物の usecase と port モックで PostCollectionHandler を組む。
func setupPostCollectionHandler() (*PostCollectionHandler, *mockPostCollectionRepo) {
	repo := new(mockPostCollectionRepo)
	h := NewPostCollectionHandler(
		usecase.NewCreatePostCollectionUseCase(repo),
		usecase.NewGetPostCollectionUseCase(repo),
		usecase.NewListPostCollectionsForViewerUseCase(repo),
		usecase.NewCountPostCollectionsUseCase(repo),
		usecase.NewUpdatePostCollectionUseCase(repo),
		usecase.NewDeletePostCollectionUseCase(repo),
		usecase.NewAddPostToCollectionUseCase(repo),
		usecase.NewRemovePostFromCollectionUseCase(repo),
		usecase.NewListPostCollectionPostsUseCase(repo),
	)
	return h, repo
}

// ownedCollection は認証ユーザー（userID=1）が所有するコレクションを返す。
func ownedCollection() *model.PostCollection {
	return &model.PostCollection{Title: "My Collection", UserID: 1}
}

// ============================================================
// Create テスト
// ============================================================

func TestPostCollectionCreate_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostCollection")).Return(nil)

	r := newRouter(1)
	r.POST("/collections", h.Create)

	w := doRequest(r, http.MethodPost, "/collections", map[string]interface{}{
		"title":       "My Collection",
		"description": "テスト用コレクション",
		"is_public":   true,
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

func TestPostCollectionCreate_InvalidBody(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.POST("/collections", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/collections", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionCreate_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.PostCollection")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/collections", h.Create)

	w := doRequest(r, http.MethodPost, "/collections", map[string]interface{}{
		"title": "My Collection",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostCollectionGetByID_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	collection := &model.PostCollection{Title: "My Collection", UserID: 1}
	repo.On("FindByID", mock.Anything, uint(5)).Return(collection, nil)

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/5", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetByID_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionGetByID_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return((*model.PostCollection)(nil), errors.New("not found"))

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestPostCollectionGetByUserID_OwnCollections(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	collections := []model.PostCollection{{Title: "My Collection", UserID: 1}}
	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return(collections, int64(1), nil)

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/collections", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetByUserID_OtherUserPublicOnly(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	collections := []model.PostCollection{{Title: "Public Collection", UserID: 2, IsPublic: true}}
	repo.On("FindPublicByUserID", mock.Anything, uint(2)).Return(collections, nil)

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/2/collections", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/collections", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionGetByUserID_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return([]model.PostCollection(nil), int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/collections", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestPostCollectionUpdate_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.PostCollection")).Return(nil)

	r := newRouter(1)
	r.PUT("/collections/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/collections/5", map[string]interface{}{
		"title":       "Updated",
		"description": "新しい説明",
		"is_public":   true,
	})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionUpdate_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.PUT("/collections/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/collections/abc", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionUpdate_InvalidBody(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.PUT("/collections/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/collections/5", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionUpdate_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.PostCollection")).Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/collections/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/collections/5", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestPostCollectionDelete_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("Delete", mock.Anything, uint(5)).Return(nil)

	r := newRouter(1)
	r.DELETE("/collections/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/collections/5", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionDelete_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.DELETE("/collections/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/collections/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionDelete_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("Delete", mock.Anything, uint(5)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/collections/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/collections/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// GetPosts テスト
// ============================================================

func TestPostCollectionGetPosts_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	items := []model.PostCollectionItem{{CollectionID: 5, PostID: 10}}
	repo.On("GetPostsByCollectionID", mock.Anything, uint(5)).Return(items, nil)

	r := newRouter(1)
	r.GET("/collections/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/collections/5/posts", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetPosts_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.GET("/collections/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/collections/abc/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionGetPosts_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("GetPostsByCollectionID", mock.Anything, uint(5)).Return([]model.PostCollectionItem(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/collections/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/collections/5/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// AddPost テスト
// ============================================================

func TestPostCollectionAddPost_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
	repo.On("AddPost", mock.Anything, mock.AnythingOfType("*model.PostCollectionItem")).Return(nil)

	r := newRouter(1)
	r.POST("/collections/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/collections/5/posts", map[string]interface{}{
		"post_id": 10,
		"note":    "メモ",
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

func TestPostCollectionAddPost_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.POST("/collections/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/collections/abc/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionAddPost_InvalidBody(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.POST("/collections/:id/posts", h.AddPost)

	w := doRequestRaw(r, http.MethodPost, "/collections/5/posts", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionAddPost_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("HasPost", mock.Anything, uint(5), uint(10)).Return(false, nil)
	repo.On("AddPost", mock.Anything, mock.AnythingOfType("*model.PostCollectionItem")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/collections/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/collections/5/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// RemovePost テスト
// ============================================================

func TestPostCollectionRemovePost_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("RemovePost", mock.Anything, uint(5), uint(10)).Return(nil)

	r := newRouter(1)
	r.DELETE("/collections/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/collections/5/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestPostCollectionRemovePost_InvalidCollectionID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.DELETE("/collections/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/collections/abc/posts/10", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionRemovePost_InvalidPostID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.DELETE("/collections/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/collections/5/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionRemovePost_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("FindByID", mock.Anything, uint(5)).Return(ownedCollection(), nil)
	repo.On("RemovePost", mock.Anything, uint(5), uint(10)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/collections/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/collections/5/posts/10", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

// ============================================================
// レスポンスボディ検証テスト
// ============================================================

func TestPostCollectionGetByID_ResponseBody(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	collection := &model.PostCollection{Title: "Test Collection", UserID: 1}
	repo.On("FindByID", mock.Anything, uint(1)).Return(collection, nil)

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/1", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, "Test Collection", body["title"])
	repo.AssertExpectations(t)
}

// ============================================================
// GetMyCount テスト
// ============================================================

func TestPostCollectionGetMyCount_Success(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(5), nil)

	r := newRouter(1)
	r.GET("/collections/my/count", h.GetMyCount)

	w := doRequest(r, http.MethodGet, "/collections/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":5`)
	repo.AssertExpectations(t)
}

func TestPostCollectionGetMyCount_ServiceError(t *testing.T) {
	h, repo := setupPostCollectionHandler()
	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/collections/my/count", h.GetMyCount)

	w := doRequest(r, http.MethodGet, "/collections/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
