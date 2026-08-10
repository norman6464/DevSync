package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockBookmarkCollectionRepo は usecase/repository.BookmarkCollectionRepository のモック（ctx 付き）。
type mockBookmarkCollectionRepo struct{ mock.Mock }

func (m *mockBookmarkCollectionRepo) Create(ctx context.Context, collection *model.BookmarkCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockBookmarkCollectionRepo) FindByID(ctx context.Context, id uint) (*model.BookmarkCollection, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.BookmarkCollection)
	return c, args.Error(1)
}

func (m *mockBookmarkCollectionRepo) FindByUserID(ctx context.Context, userID uint) ([]model.BookmarkCollection, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.BookmarkCollection)
	return c, args.Error(1)
}

func (m *mockBookmarkCollectionRepo) Update(ctx context.Context, collection *model.BookmarkCollection) error {
	return m.Called(ctx, collection).Error(0)
}

func (m *mockBookmarkCollectionRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockBookmarkCollectionRepo) AddPost(ctx context.Context, item *model.BookmarkCollectionItem) error {
	return m.Called(ctx, item).Error(0)
}

func (m *mockBookmarkCollectionRepo) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return m.Called(ctx, collectionID, postID).Error(0)
}

func (m *mockBookmarkCollectionRepo) GetPosts(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(ctx, collectionID, limit, offset)
	p, _ := args.Get(0).([]model.Post)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockBookmarkCollectionRepo) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	args := m.Called(ctx, collectionID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *mockBookmarkCollectionRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// newBookmarkCollectionHandlerWithRepo は本物の usecase と port モックで handler を組む。
func newBookmarkCollectionHandlerWithRepo(repo *mockBookmarkCollectionRepo) *BookmarkCollectionHandler {
	return NewBookmarkCollectionHandler(
		usecase.NewCreateBookmarkCollectionUseCase(repo),
		usecase.NewListBookmarkCollectionsUseCase(repo),
		usecase.NewUpdateBookmarkCollectionUseCase(repo),
		usecase.NewDeleteBookmarkCollectionUseCase(repo),
		usecase.NewAddPostToBookmarkCollectionUseCase(repo),
		usecase.NewRemovePostFromBookmarkCollectionUseCase(repo),
		usecase.NewListBookmarkCollectionPostsUseCase(repo),
		usecase.NewCountBookmarkCollectionsUseCase(repo),
	)
}

// ownedBookmarkCollection は認証ユーザー（userID=1）が所有するコレクションを返す。
func ownedBookmarkCollection() *model.BookmarkCollection {
	return &model.BookmarkCollection{Name: "Go学習", UserID: 1}
}

// ============================================================
// BookmarkCollectionHandler テスト
// ============================================================

func TestBookmarkCollection_Create_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model.BookmarkCollection) bool {
		return c.UserID == 1 && c.Name == "Go学習"
	})).Return(nil)

	r := gin.New()
	r.POST("/bookmark-collections", authMiddleware(1), h.Create)

	body := jsonBody(map[string]string{"name": "Go学習", "color": "blue"})
	req, _ := http.NewRequest("POST", "/bookmark-collections", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_GetMyCollections_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	collections := []model.BookmarkCollection{
		{ID: 1, UserID: 1, Name: "Go"},
		{ID: 2, UserID: 1, Name: "React"},
	}
	mockSvc.On("FindByUserID", mock.Anything, uint(1)).Return(collections, nil)

	r := gin.New()
	r.GET("/bookmark-collections", authMiddleware(1), h.GetMyCollections)

	req, _ := http.NewRequest("GET", "/bookmark-collections", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_Update_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(ownedBookmarkCollection(), nil)
	mockSvc.On("Update", mock.Anything, mock.AnythingOfType("*model.BookmarkCollection")).Return(nil)

	r := gin.New()
	r.PUT("/bookmark-collections/:id", authMiddleware(1), h.Update)

	body := jsonBody(map[string]string{"name": "新名"})
	req, _ := http.NewRequest("PUT", "/bookmark-collections/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	result := parseJSON(t, w)
	if result["name"] != "新名" {
		t.Errorf("expected name=新名, got %v", result["name"])
	}
}

func TestBookmarkCollection_Update_Forbidden(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	// 他人のコレクションなので所有権チェックで弾かれ、更新は呼ばれない
	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(&model.BookmarkCollection{ID: 1, UserID: 99, Name: "他人の"}, nil)

	r := gin.New()
	r.PUT("/bookmark-collections/:id", authMiddleware(1), h.Update)

	body := jsonBody(map[string]string{"name": "新名"})
	req, _ := http.NewRequest("PUT", "/bookmark-collections/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusForbidden)
}

func TestBookmarkCollection_Delete_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(ownedBookmarkCollection(), nil)
	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id", authMiddleware(1), h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
}

func TestBookmarkCollection_AddPost_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(ownedBookmarkCollection(), nil)
	mockSvc.On("HasPost", mock.Anything, uint(1), uint(10)).Return(false, nil)
	mockSvc.On("AddPost", mock.Anything, mock.AnythingOfType("*model.BookmarkCollectionItem")).Return(nil)

	r := gin.New()
	r.POST("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.AddPost)

	req, _ := http.NewRequest("POST", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
}

func TestBookmarkCollection_AddPost_Conflict(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(ownedBookmarkCollection(), nil)
	mockSvc.On("HasPost", mock.Anything, uint(1), uint(10)).Return(true, nil)

	r := gin.New()
	r.POST("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.AddPost)

	req, _ := http.NewRequest("POST", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusConflict)
}

func TestBookmarkCollection_GetPosts_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	posts := []model.Post{
		{ID: 1, Title: "投稿1"},
		{ID: 2, Title: "投稿2"},
	}
	mockSvc.On("GetPosts", mock.Anything, uint(1), 20, 0).Return(posts, int64(2), nil)

	r := gin.New()
	r.GET("/bookmark-collections/:id/posts", h.GetPosts)

	req, _ := http.NewRequest("GET", "/bookmark-collections/1/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	result := parseJSON(t, w)
	if result["total"].(float64) != 2 {
		t.Errorf("expected total=2, got %v", result["total"])
	}
}

// ============================================================
// RemovePost テスト
// ============================================================

func TestBookmarkCollection_RemovePost_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(ownedBookmarkCollection(), nil)
	mockSvc.On("RemovePost", mock.Anything, uint(1), uint(10)).Return(nil)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_RemovePost_InvalidCollectionID(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/abc/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkCollection_RemovePost_InvalidPostID(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1/posts/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkCollection_RemovePost_ServiceError(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(&model.BookmarkCollection{Name: "他人の", UserID: 99}, nil)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusForbidden)
}

// ============================================================
// エラーケース追加テスト
// ============================================================

func TestBookmarkCollection_Create_ServiceError(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*model.BookmarkCollection")).Return(errors.New("db error"))

	r := gin.New()
	r.POST("/bookmark-collections", authMiddleware(1), h.Create)

	body := jsonBody(map[string]string{"name": "Test"})
	req, _ := http.NewRequest("POST", "/bookmark-collections", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestBookmarkCollection_GetMyCollections_ServiceError(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByUserID", mock.Anything, uint(1)).Return([]model.BookmarkCollection(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/bookmark-collections", authMiddleware(1), h.GetMyCollections)

	req, _ := http.NewRequest("GET", "/bookmark-collections", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestBookmarkCollection_Delete_ServiceError(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("FindByID", mock.Anything, uint(1)).Return(&model.BookmarkCollection{Name: "他人の", UserID: 99}, nil)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id", authMiddleware(1), h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusForbidden)
}

func TestBookmarkCollection_Delete_InvalidID(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id", authMiddleware(1), h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkCollection_GetPosts_ServiceError(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("GetPosts", mock.Anything, uint(1), 20, 0).Return([]model.Post(nil), int64(0), errors.New("db error"))

	r := gin.New()
	r.GET("/bookmark-collections/:id/posts", h.GetPosts)

	req, _ := http.NewRequest("GET", "/bookmark-collections/1/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestBookmarkCollection_GetPosts_InvalidID(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	r := gin.New()
	r.GET("/bookmark-collections/:id/posts", h.GetPosts)

	req, _ := http.NewRequest("GET", "/bookmark-collections/abc/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetMyCount テスト
// ============================================================

func TestBookmarkCollection_GetMyCount_Success(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)

	r := newRouter(1)
	r.GET("/bookmark-collections/my/count", h.GetMyCount)

	w := doRequest(r, http.MethodGet, "/bookmark-collections/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":3`)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_GetMyCount_ServiceError(t *testing.T) {
	mockSvc := new(mockBookmarkCollectionRepo)
	h := newBookmarkCollectionHandlerWithRepo(mockSvc)

	mockSvc.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/bookmark-collections/my/count", h.GetMyCount)

	w := doRequest(r, http.MethodGet, "/bookmark-collections/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	mockSvc.AssertExpectations(t)
}
