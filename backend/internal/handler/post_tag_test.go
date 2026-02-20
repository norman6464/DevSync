package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostTagService は PostTagServiceInterface のモック実装。
type MockPostTagService struct{ mock.Mock }

func (m *MockPostTagService) SetTags(postID, userID uint, tags []string) error {
	return m.Called(postID, userID, tags).Error(0)
}
func (m *MockPostTagService) SetAutoTags(postID, userID uint, content string) error {
	return m.Called(postID, userID, content).Error(0)
}
func (m *MockPostTagService) GetByPostID(postID uint) ([]string, error) {
	args := m.Called(postID)
	if v := args.Get(0); v != nil {
		return v.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostTagService) FindPostsByTag(tag string, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(tag, limit, offset)
	if v := args.Get(0); v != nil {
		return v.([]model.Post), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}
func (m *MockPostTagService) GetPopularTags(limit int) ([]model.TagCount, error) {
	args := m.Called(limit)
	if v := args.Get(0); v != nil {
		return v.([]model.TagCount), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupPostTagHandler() (*PostTagHandler, *MockPostTagService) {
	svc := new(MockPostTagService)
	h := NewPostTagHandler(svc)
	return h, svc
}

// ============================================================
// SetTags テスト
// ============================================================

func TestPostTagSetTags_Success(t *testing.T) {
	h, svc := setupPostTagHandler()
	svc.On("SetTags", uint(5), uint(1), []string{"Go", "React"}).Return(nil)

	r := newRouter(1)
	r.POST("/posts/:postId/tags", h.SetTags)

	w := doRequest(r, http.MethodPost, "/posts/5/tags", map[string]interface{}{
		"tags": []string{"Go", "React"},
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTagSetTags_InvalidID(t *testing.T) {
	h, _ := setupPostTagHandler()

	r := newRouter(1)
	r.POST("/posts/:postId/tags", h.SetTags)

	w := doRequest(r, http.MethodPost, "/posts/abc/tags", map[string]interface{}{
		"tags": []string{"Go"},
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTagSetTags_InvalidBody(t *testing.T) {
	h, _ := setupPostTagHandler()

	r := newRouter(1)
	r.POST("/posts/:postId/tags", h.SetTags)

	w := doRequestRaw(r, http.MethodPost, "/posts/5/tags", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTagSetTags_ServiceError(t *testing.T) {
	h, svc := setupPostTagHandler()
	svc.On("SetTags", uint(5), uint(1), []string{"Go"}).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/posts/:postId/tags", h.SetTags)

	w := doRequest(r, http.MethodPost, "/posts/5/tags", map[string]interface{}{
		"tags": []string{"Go"},
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByPostID テスト
// ============================================================

func TestPostTagGetByPostID_Success(t *testing.T) {
	h, svc := setupPostTagHandler()
	svc.On("GetByPostID", uint(5)).Return([]string{"Go", "React"}, nil)

	r := newRouter(1)
	r.GET("/posts/:postId/tags", h.GetByPostID)

	w := doRequest(r, http.MethodGet, "/posts/5/tags", nil)
	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	assert.NotNil(t, body["tags"])
	svc.AssertExpectations(t)
}

func TestPostTagGetByPostID_InvalidID(t *testing.T) {
	h, _ := setupPostTagHandler()

	r := newRouter(1)
	r.GET("/posts/:postId/tags", h.GetByPostID)

	w := doRequest(r, http.MethodGet, "/posts/abc/tags", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTagGetByPostID_ServiceError(t *testing.T) {
	h, svc := setupPostTagHandler()
	svc.On("GetByPostID", uint(5)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/:postId/tags", h.GetByPostID)

	w := doRequest(r, http.MethodGet, "/posts/5/tags", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// FindPostsByTag テスト
// ============================================================

func TestPostTagFindPostsByTag_Success(t *testing.T) {
	h, svc := setupPostTagHandler()
	posts := []model.Post{{Title: "Go入門", UserID: 1}}
	svc.On("FindPostsByTag", "Go", 20, 0).Return(posts, int64(1), nil)

	r := newRouter(1)
	r.GET("/posts/tags/search", h.FindPostsByTag)

	w := doRequest(r, http.MethodGet, "/posts/tags/search?tag=Go", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTagFindPostsByTag_MissingTag(t *testing.T) {
	h, _ := setupPostTagHandler()

	r := newRouter(1)
	r.GET("/posts/tags/search", h.FindPostsByTag)

	w := doRequest(r, http.MethodGet, "/posts/tags/search", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostTagFindPostsByTag_ServiceError(t *testing.T) {
	h, svc := setupPostTagHandler()
	svc.On("FindPostsByTag", "Go", 20, 0).Return(nil, int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/tags/search", h.FindPostsByTag)

	w := doRequest(r, http.MethodGet, "/posts/tags/search?tag=Go", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetPopularTags テスト
// ============================================================

func TestPostTagGetPopularTags_Success(t *testing.T) {
	h, svc := setupPostTagHandler()
	tags := []model.TagCount{{Tag: "Go", Count: 100}, {Tag: "React", Count: 50}}
	svc.On("GetPopularTags", 20).Return(tags, nil)

	r := newRouter(1)
	r.GET("/posts/tags/popular", h.GetPopularTags)

	w := doRequest(r, http.MethodGet, "/posts/tags/popular", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestPostTagGetPopularTags_ServiceError(t *testing.T) {
	h, svc := setupPostTagHandler()
	svc.On("GetPopularTags", 20).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/posts/tags/popular", h.GetPopularTags)

	w := doRequest(r, http.MethodGet, "/posts/tags/popular", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
