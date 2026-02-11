package handler

import (
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestResourceCreate_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources", h.Create)

	repo.On("Create", mock.AnythingOfType("*model.LearningResource")).Return(nil)

	w := doRequest(r, http.MethodPost, "/resources", map[string]string{
		"title": "Go Tutorial", "category": "programming",
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
	repo.On("FindByID", uint(10)).Return(resource, nil)
	repo.On("HasLiked", uint(1), uint(10)).Return(true, nil)
	repo.On("HasSaved", uint(1), uint(10)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/resources/10", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, true, body["has_liked"])
	assert.Equal(t, false, body["has_saved"])
}

func TestResourceGetByID_NotFound(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/:id", h.GetByID)

	repo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/resources/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestResourceGetByID_ForbiddenPrivate(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/:id", h.GetByID)

	// 非公開で他ユーザーのリソース
	resource := &model.LearningResource{Title: "Private", IsPublic: false}
	resource.ID = 10
	resource.UserID = 999
	repo.On("FindByID", uint(10)).Return(resource, nil)

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

	repo.On("FindPublic", 20, 0, "", "").Return(
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

	repo.On("FindPublic", 10, 5, "programming", "beginner").Return(
		[]model.LearningResource{}, int64(0), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources?limit=10&offset=5&category=programming&difficulty=beginner", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceGetPublic_LimitCap(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources", h.GetPublic)

	repo.On("FindPublic", 100, 0, "", "").Return(
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

	repo.On("Search", "go", 20, 0).Return(
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
	repo.On("FindByID", uint(10)).Return(resource, nil)
	repo.On("Update", mock.AnythingOfType("*model.LearningResource")).Return(nil)

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
	repo.On("FindByID", uint(10)).Return(resource, nil)

	w := doRequest(r, http.MethodPut, "/resources/10", map[string]string{
		"title": "Hacked",
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestResourceUpdate_NotFound(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.PUT("/resources/:id", h.Update)

	repo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

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
	repo.On("FindByID", uint(10)).Return(resource, nil)
	repo.On("Delete", uint(10)).Return(nil)

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
	repo.On("FindByID", uint(10)).Return(resource, nil)

	w := doRequest(r, http.MethodDelete, "/resources/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- Like / Unlike ----------

func TestResourceLike_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/like", h.Like)

	repo.On("Like", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/resources/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceUnlike_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/like", h.Unlike)

	repo.On("Unlike", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/resources/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Save / Unsave ----------

func TestResourceSave_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.POST("/resources/:id/save", h.SaveResource)

	repo.On("Save", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/resources/5/save", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestResourceUnsave_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.DELETE("/resources/:id/save", h.UnsaveResource)

	repo.On("Unsave", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/resources/5/save", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetSaved ----------

func TestResourceGetSaved_Success(t *testing.T) {
	h, repo := setupLearningResourceHandler()
	r := newRouter(1)
	r.GET("/resources/saved", h.GetSaved)

	repo.On("FindSavedByUserID", uint(1), 20, 0).Return(
		[]model.LearningResource{{Title: "Saved"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/resources/saved", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["total"])
}
