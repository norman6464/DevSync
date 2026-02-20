package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostCollectionService は PostCollectionServiceInterface のモック実装。
type MockPostCollectionService struct{ mock.Mock }

func (m *MockPostCollectionService) Create(collection *model.PostCollection) (*model.PostCollection, error) {
	args := m.Called(collection)
	if v := args.Get(0); v != nil {
		return v.(*model.PostCollection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostCollectionService) GetByID(id uint) (*model.PostCollection, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*model.PostCollection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostCollectionService) GetByUserID(userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	args := m.Called(userID, limit, offset)
	if v := args.Get(0); v != nil {
		return v.([]model.PostCollection), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}
func (m *MockPostCollectionService) GetPublicByUserID(userID uint) ([]model.PostCollection, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.([]model.PostCollection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostCollectionService) Update(id, userID uint, title, description string, isPublic bool) (*model.PostCollection, error) {
	args := m.Called(id, userID, title, description, isPublic)
	if v := args.Get(0); v != nil {
		return v.(*model.PostCollection), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostCollectionService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockPostCollectionService) AddPost(collectionID, userID, postID uint, note string) error {
	return m.Called(collectionID, userID, postID, note).Error(0)
}
func (m *MockPostCollectionService) RemovePost(collectionID, userID, postID uint) error {
	return m.Called(collectionID, userID, postID).Error(0)
}
func (m *MockPostCollectionService) GetPosts(collectionID uint) ([]model.PostCollectionItem, error) {
	args := m.Called(collectionID)
	if v := args.Get(0); v != nil {
		return v.([]model.PostCollectionItem), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupPostCollectionHandler() (*PostCollectionHandler, *MockPostCollectionService) {
	svc := new(MockPostCollectionService)
	h := NewPostCollectionHandler(svc)
	return h, svc
}

// ============================================================
// Create テスト
// ============================================================

func TestPostCollectionCreate_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	result := &model.PostCollection{Title: "My Collection", UserID: 1}
	svc.On("Create", mock.AnythingOfType("*model.PostCollection")).Return(result, nil)

	r := newRouter(1)
	r.POST("/collections", h.Create)

	w := doRequest(r, http.MethodPost, "/collections", map[string]interface{}{
		"title":       "My Collection",
		"description": "テスト用コレクション",
		"is_public":   true,
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestPostCollectionCreate_InvalidBody(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.POST("/collections", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/collections", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionCreate_ServiceError(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("Create", mock.AnythingOfType("*model.PostCollection")).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.POST("/collections", h.Create)

	w := doRequest(r, http.MethodPost, "/collections", map[string]interface{}{
		"title": "My Collection",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestPostCollectionGetByID_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	collection := &model.PostCollection{Title: "My Collection", UserID: 1}
	svc.On("GetByID", uint(5)).Return(collection, nil)

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostCollectionGetByID_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionGetByID_ServiceError(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("GetByID", uint(5)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestPostCollectionGetByUserID_OwnCollections(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	collections := []model.PostCollection{{Title: "My Collection", UserID: 1}}
	svc.On("GetByUserID", uint(1), 20, 0).Return(collections, int64(1), nil)

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/collections", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostCollectionGetByUserID_OtherUserPublicOnly(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	collections := []model.PostCollection{{Title: "Public Collection", UserID: 2, IsPublic: true}}
	svc.On("GetPublicByUserID", uint(2)).Return(collections, nil)

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/2/collections", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostCollectionGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/collections", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionGetByUserID_ServiceError(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("GetByUserID", uint(1), 20, 0).Return(nil, int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/collections", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/1/collections", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestPostCollectionUpdate_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	updated := &model.PostCollection{Title: "Updated", UserID: 1}
	svc.On("Update", uint(5), uint(1), "Updated", "新しい説明", true).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/collections/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/collections/5", map[string]interface{}{
		"title":       "Updated",
		"description": "新しい説明",
		"is_public":   true,
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
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
	h, svc := setupPostCollectionHandler()
	svc.On("Update", uint(5), uint(1), "Updated", "", false).Return(nil, errors.New("forbidden"))

	r := newRouter(1)
	r.PUT("/collections/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/collections/5", map[string]interface{}{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestPostCollectionDelete_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("Delete", uint(5), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/collections/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/collections/5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostCollectionDelete_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.DELETE("/collections/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/collections/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionDelete_ServiceError(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("Delete", uint(5), uint(1)).Return(errors.New("forbidden"))

	r := newRouter(1)
	r.DELETE("/collections/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/collections/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetPosts テスト
// ============================================================

func TestPostCollectionGetPosts_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	items := []model.PostCollectionItem{{CollectionID: 5, PostID: 10}}
	svc.On("GetPosts", uint(5)).Return(items, nil)

	r := newRouter(1)
	r.GET("/collections/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/collections/5/posts", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostCollectionGetPosts_InvalidID(t *testing.T) {
	h, _ := setupPostCollectionHandler()

	r := newRouter(1)
	r.GET("/collections/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/collections/abc/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCollectionGetPosts_ServiceError(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("GetPosts", uint(5)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/collections/:id/posts", h.GetPosts)

	w := doRequest(r, http.MethodGet, "/collections/5/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// AddPost テスト
// ============================================================

func TestPostCollectionAddPost_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("AddPost", uint(5), uint(1), uint(10), "メモ").Return(nil)

	r := newRouter(1)
	r.POST("/collections/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/collections/5/posts", map[string]interface{}{
		"post_id": 10,
		"note":    "メモ",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
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
	h, svc := setupPostCollectionHandler()
	svc.On("AddPost", uint(5), uint(1), uint(10), "").Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/collections/:id/posts", h.AddPost)

	w := doRequest(r, http.MethodPost, "/collections/5/posts", map[string]interface{}{
		"post_id": 10,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// RemovePost テスト
// ============================================================

func TestPostCollectionRemovePost_Success(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	svc.On("RemovePost", uint(5), uint(1), uint(10)).Return(nil)

	r := newRouter(1)
	r.DELETE("/collections/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/collections/5/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
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
	h, svc := setupPostCollectionHandler()
	svc.On("RemovePost", uint(5), uint(1), uint(10)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/collections/:id/posts/:postId", h.RemovePost)

	w := doRequest(r, http.MethodDelete, "/collections/5/posts/10", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// レスポンスボディ検証テスト
// ============================================================

func TestPostCollectionGetByID_ResponseBody(t *testing.T) {
	h, svc := setupPostCollectionHandler()
	collection := &model.PostCollection{Title: "Test Collection", UserID: 1}
	svc.On("GetByID", uint(1)).Return(collection, nil)

	r := newRouter(1)
	r.GET("/collections/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/collections/1", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.Equal(t, "Test Collection", body["title"])
	svc.AssertExpectations(t)
}
