package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockBookmarkCollectionService
// ============================================================

type MockBookmarkCollectionService struct{ mock.Mock }

func (m *MockBookmarkCollectionService) Create(collection *model.BookmarkCollection) error {
	return m.Called(collection).Error(0)
}

func (m *MockBookmarkCollectionService) GetByUserID(userID uint) ([]model.BookmarkCollection, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.BookmarkCollection), args.Error(1)
}

func (m *MockBookmarkCollectionService) Update(id, userID uint, updates *model.BookmarkCollection) (*model.BookmarkCollection, error) {
	args := m.Called(id, userID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BookmarkCollection), args.Error(1)
}

func (m *MockBookmarkCollectionService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func (m *MockBookmarkCollectionService) AddPost(collectionID, postID, userID uint) error {
	return m.Called(collectionID, postID, userID).Error(0)
}

func (m *MockBookmarkCollectionService) RemovePost(collectionID, postID, userID uint) error {
	return m.Called(collectionID, postID, userID).Error(0)
}

func (m *MockBookmarkCollectionService) GetPosts(collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(collectionID, limit, offset)
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}
func (m *MockBookmarkCollectionService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// ============================================================
// BookmarkCollectionHandler テスト
// ============================================================

func TestBookmarkCollection_Create_Success(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("Create", mock.MatchedBy(func(c *model.BookmarkCollection) bool {
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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	collections := []model.BookmarkCollection{
		{ID: 1, UserID: 1, Name: "Go"},
		{ID: 2, UserID: 1, Name: "React"},
	}
	mockSvc.On("GetByUserID", uint(1)).Return(collections, nil)

	r := gin.New()
	r.GET("/bookmark-collections", authMiddleware(1), h.GetMyCollections)

	req, _ := http.NewRequest("GET", "/bookmark-collections", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_Update_Success(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	updated := &model.BookmarkCollection{ID: 1, UserID: 1, Name: "新名"}
	mockSvc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.BookmarkCollection")).Return(updated, nil)

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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.BookmarkCollection")).
		Return(nil, domain.ErrForbidden)

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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("Delete", uint(1), uint(1)).Return(nil)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id", authMiddleware(1), h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
}

func TestBookmarkCollection_AddPost_Success(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("AddPost", uint(1), uint(10), uint(1)).Return(nil)

	r := gin.New()
	r.POST("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.AddPost)

	req, _ := http.NewRequest("POST", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
}

func TestBookmarkCollection_AddPost_Conflict(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("AddPost", uint(1), uint(10), uint(1)).Return(domain.ErrConflict)

	r := gin.New()
	r.POST("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.AddPost)

	req, _ := http.NewRequest("POST", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusConflict)
}

func TestBookmarkCollection_GetPosts_Success(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	posts := []model.Post{
		{ID: 1, Title: "投稿1"},
		{ID: 2, Title: "投稿2"},
	}
	mockSvc.On("GetPosts", uint(1), 20, 0).Return(posts, int64(2), nil)

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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("RemovePost", uint(1), uint(10), uint(1)).Return(nil)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_RemovePost_InvalidCollectionID(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/abc/posts/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkCollection_RemovePost_InvalidPostID(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id/posts/:postId", authMiddleware(1), h.RemovePost)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1/posts/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkCollection_RemovePost_ServiceError(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("RemovePost", uint(1), uint(10), uint(1)).Return(domain.ErrForbidden)

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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("Create", mock.AnythingOfType("*model.BookmarkCollection")).Return(errors.New("db error"))

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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("GetByUserID", uint(1)).Return([]model.BookmarkCollection(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/bookmark-collections", authMiddleware(1), h.GetMyCollections)

	req, _ := http.NewRequest("GET", "/bookmark-collections", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestBookmarkCollection_Delete_ServiceError(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("Delete", uint(1), uint(1)).Return(domain.ErrForbidden)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id", authMiddleware(1), h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusForbidden)
}

func TestBookmarkCollection_Delete_InvalidID(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	r := gin.New()
	r.DELETE("/bookmark-collections/:id", authMiddleware(1), h.Delete)

	req, _ := http.NewRequest("DELETE", "/bookmark-collections/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookmarkCollection_GetPosts_ServiceError(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("GetPosts", uint(1), 20, 0).Return([]model.Post(nil), int64(0), errors.New("db error"))

	r := gin.New()
	r.GET("/bookmark-collections/:id/posts", h.GetPosts)

	req, _ := http.NewRequest("GET", "/bookmark-collections/1/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestBookmarkCollection_GetPosts_InvalidID(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

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
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("CountByUserID", uint(1)).Return(int64(3), nil)

	r := newRouter(1)
	r.GET("/bookmark-collections/my/count", h.GetMyCount)

	w := doRequest(r, http.MethodGet, "/bookmark-collections/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":3`)
	mockSvc.AssertExpectations(t)
}

func TestBookmarkCollection_GetMyCount_ServiceError(t *testing.T) {
	mockSvc := new(MockBookmarkCollectionService)
	h := NewBookmarkCollectionHandler(mockSvc)

	mockSvc.On("CountByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/bookmark-collections/my/count", h.GetMyCount)

	w := doRequest(r, http.MethodGet, "/bookmark-collections/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	mockSvc.AssertExpectations(t)
}
