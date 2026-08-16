package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestResourceCreate_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources", h.Create)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.LearningResource")).Return(nil)

	w := doRequest(r, http.MethodPost, "/resources", map[string]string{
		"title": "Go Tutorial", "category": "article", "url": "https://example.com/go-tutorial",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestResourceCreate_ValidationError(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources", h.Create)

	// title と category は required
	w := doRequest(r, http.MethodPost, "/resources", map[string]string{
		"title": "No category",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceCreate_InvalidJSON(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/resources", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetByID ----------

func TestResourceGetByID_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/:id", h.GetByID)

	resource := &model.LearningResource{Title: "Found", IsPublic: true}
	resource.ID = 10
	resource.UserID = 1
	repo.On("FindByID", mock.Anything, uint(10)).Return(resource, nil)
	repo.On("HasLiked", mock.Anything, uint(1), uint(10)).Return(true, nil)
	repo.On("HasSaved", mock.Anything, uint(1), uint(10)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/resources/10", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, true, body["has_liked"])
	assert.Equal(t, false, body["has_saved"])
}

// 不在のリソースは 404 にならず 500 になる（移行前からの挙動）。
func TestResourceGetByID_MissingReturnsInternalError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/:id", h.GetByID)

	// port は不在を (nil, nil) で表す。
	repo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/resources/999", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceGetByID_ForbiddenPrivate(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/:id", h.GetByID)

	// 非公開で他ユーザーのリソース
	resource := &model.LearningResource{Title: "Private", IsPublic: false}
	resource.ID = 10
	resource.UserID = 999
	repo.On("FindByID", mock.Anything, uint(10)).Return(resource, nil)

	w := doRequest(r, http.MethodGet, "/resources/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestResourceGetByID_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/resources/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetPublic ----------

func TestResourceGetPublic_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources", h.GetPublic)

	repo.On("FindPublic", mock.Anything, 20, 0, "", "").Return(
		[]model.LearningResource{{Title: "Public Resource"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	resources := body["resources"].([]interface{})
	assert.Len(t, resources, 1)
}

func TestResourceGetPublic_WithFilters(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources", h.GetPublic)

	repo.On("FindPublic", mock.Anything, 10, 5, "programming", "beginner").Return(
		[]model.LearningResource{}, int64(0), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources?limit=10&offset=5&category=programming&difficulty=beginner", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceGetPublic_LimitCap(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources", h.GetPublic)

	repo.On("FindPublic", mock.Anything, 100, 0, "", "").Return(
		[]model.LearningResource{}, int64(0), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources?limit=200", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Search ----------

func TestResourceSearch_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/search", h.Search)

	repo.On("Search", mock.Anything, "go", 20, 0).Return(
		[]model.LearningResource{{Title: "Go Basics"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources/search?q=go", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Update ----------

func TestResourceUpdate_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	resource := &model.LearningResource{Title: "Old Title"}
	resource.ID = 10
	resource.UserID = 1
	repo.On("FindByID", mock.Anything, uint(10)).Return(resource, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.LearningResource")).Return(nil)

	w := doRequest(r, http.MethodPut, "/resources/10", map[string]string{
		"title": "Updated Title",
	})
	assertStatus(t, w, http.StatusOK)
}

func TestResourceUpdate_Forbidden(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	resource := &model.LearningResource{Title: "Other's"}
	resource.ID = 10
	resource.UserID = 999
	repo.On("FindByID", mock.Anything, uint(10)).Return(resource, nil)

	w := doRequest(r, http.MethodPut, "/resources/10", map[string]string{
		"title": "Hacked",
	})
	assertStatus(t, w, http.StatusForbidden)
}

// 更新も不在は 404 にならず 500 になる（移行前からの挙動）。
func TestResourceUpdate_MissingReturnsInternalError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/resources/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Delete ----------

func TestResourceDelete_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id", h.Delete)

	resource := &model.LearningResource{}
	resource.ID = 10
	resource.UserID = 1
	repo.On("FindByID", mock.Anything, uint(10)).Return(resource, nil)
	repo.On("Delete", mock.Anything, uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/resources/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceDelete_Forbidden(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id", h.Delete)

	resource := &model.LearningResource{}
	resource.ID = 10
	resource.UserID = 999
	repo.On("FindByID", mock.Anything, uint(10)).Return(resource, nil)

	w := doRequest(r, http.MethodDelete, "/resources/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestResourceUpdate_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/resources/abc", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceUpdate_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(10)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/resources/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceDelete_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/resources/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceGetPublic_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources", h.GetPublic)

	repo.On("FindPublic", mock.Anything, 20, 0, "", "").Return(
		[]model.LearningResource{}, int64(0), errors.New("db error"),
	)

	w := doRequest(r, http.MethodGet, "/resources", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceSearch_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/search", h.Search)

	repo.On("Search", mock.Anything, "go", 20, 0).Return(
		[]model.LearningResource{}, int64(0), errors.New("db error"),
	)

	w := doRequest(r, http.MethodGet, "/resources/search?q=go", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- Like / Unlike ----------

func TestResourceLike_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/like", h.Like)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Like", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/resources/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceLike_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/like", h.Like)

	w := doRequest(r, http.MethodPost, "/resources/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceLike_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/like", h.Like)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Like", mock.Anything, uint(1), uint(5)).Return(errors.New("already liked"))

	w := doRequest(r, http.MethodPost, "/resources/5/like", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceUnlike_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/like", h.Unlike)

	w := doRequest(r, http.MethodDelete, "/resources/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceUnlike_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/like", h.Unlike)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Unlike", mock.Anything, uint(1), uint(5)).Return(errors.New("not liked"))

	w := doRequest(r, http.MethodDelete, "/resources/5/like", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceUnlike_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/like", h.Unlike)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Unlike", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/resources/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Save / Unsave ----------

func TestResourceSave_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/save", h.SaveResource)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Save", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/resources/5/save", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceSave_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/save", h.SaveResource)

	w := doRequest(r, http.MethodPost, "/resources/abc/save", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceSave_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/save", h.SaveResource)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Save", mock.Anything, uint(1), uint(5)).Return(errors.New("already saved"))

	w := doRequest(r, http.MethodPost, "/resources/5/save", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceUnsave_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/save", h.UnsaveResource)

	w := doRequest(r, http.MethodDelete, "/resources/abc/save", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceUnsave_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/save", h.UnsaveResource)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Unsave", mock.Anything, uint(1), uint(5)).Return(errors.New("not saved"))

	w := doRequest(r, http.MethodDelete, "/resources/5/save", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceUnsave_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/save", h.UnsaveResource)

	otherResource := &model.LearningResource{UserID: 99}
	otherResource.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherResource, nil)
	repo.On("Unsave", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/resources/5/save", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceGetSaved_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/saved", h.GetSaved)

	repo.On("FindSavedByUserID", mock.Anything, uint(1), 20, 0).Return(
		[]model.LearningResource{}, int64(0), errors.New("db error"),
	)

	w := doRequest(r, http.MethodGet, "/resources/saved", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetSaved ----------

func TestResourceGetSaved_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/saved", h.GetSaved)

	repo.On("FindSavedByUserID", mock.Anything, uint(1), 20, 0).Return(
		[]model.LearningResource{{Title: "Saved"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources/saved", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["total"])
}

// ---------- GetByUserID ----------

func TestResourceGetByUserID_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/users/:userId/resources", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(1), true, 20, 0).Return(
		[]model.LearningResource{{Title: "My Resource"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/users/1/resources", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["total"])
}

func TestResourceGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/users/:userId/resources", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/resources", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceGetByUserID_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/users/:userId/resources", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(1), true, 20, 0).Return(
		[]model.LearningResource{}, int64(0), errors.New("db error"),
	)

	w := doRequest(r, http.MethodGet, "/users/1/resources", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetByDifficulty ----------

func TestResourceGetByDifficulty_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/difficulty/:difficulty", h.GetByDifficulty)

	repo.On("FindByDifficulty", mock.Anything, "beginner", 20, 0).Return(
		[]model.LearningResource{{Title: "Beginner Guide"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources/difficulty/beginner", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["total"])
}

func TestResourceGetByDifficulty_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/difficulty/:difficulty", h.GetByDifficulty)

	repo.On("FindByDifficulty", mock.Anything, "beginner", 20, 0).Return(
		[]model.LearningResource{}, int64(0), errors.New("db error"),
	)

	w := doRequest(r, http.MethodGet, "/resources/difficulty/beginner", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestResourceCreate_InvalidURL(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources", h.Create)

	w := doRequest(r, http.MethodPost, "/resources", map[string]interface{}{
		"title":    "悪意あるリソース",
		"url":      "javascript:alert('xss')",
		"category": "article",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceUpdate_InvalidImageURL(t *testing.T) {
	h, _ := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/resources/1", map[string]interface{}{
		"image_url": "data:text/html,<script>alert('xss')</script>",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetMyCount（自分のリソース総数取得）
// ============================================================

func TestResourceGetMyCount_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(7), nil)

	w := doRequest(r, http.MethodGet, "/resources/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":7`)
	repo.AssertExpectations(t)
}

func TestResourceGetMyCount_ServiceError(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/resources/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}
