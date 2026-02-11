package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestPostCreate_Success(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil)
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)
	postRepo.On("FindByID", mock.AnythingOfType("uint")).Return(&model.Post{
		Title: "Test Post", Content: "Hello",
	}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Test Post", "content": "Hello",
	})

	assertStatus(t, w, http.StatusCreated)
}

func TestPostCreate_ValidationError(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	// title と content は required
	w := doRequest(r, http.MethodPost, "/posts", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreate_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/posts", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetAll ----------

func TestPostGetAll_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	postRepo.On("FindAll", 1, 20).Return([]model.Post{
		{Title: "A"}, {Title: "B"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusOK)

	var posts []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &posts)
	assert.Len(t, posts, 2)
}

func TestPostGetAll_WithPagination(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	postRepo.On("FindAll", 2, 5).Return([]model.Post{}, nil)

	w := doRequest(r, http.MethodGet, "/posts?page=2&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetByID ----------

func TestPostGetByID_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	postRepo.On("FindByID", uint(10)).Return(&model.Post{Title: "Found"}, nil)
	postRepo.On("HasLiked", uint(1), uint(0)).Return(false)

	w := doRequest(r, http.MethodGet, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetByID_NotFound(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	postRepo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetByID_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestPostUpdate_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	post := &model.Post{Title: "Old", Content: "Old Content"}
	post.ID = 10
	post.UserID = 1
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{
		"title": "New Title",
	})
	assertStatus(t, w, http.StatusOK)
}

func TestPostUpdate_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1) // userID=1
	r.PUT("/posts/:id", h.Update)

	post := &model.Post{Title: "Other's post"}
	post.ID = 10
	post.UserID = 999 // 別ユーザーの投稿
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{
		"title": "Hacked",
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestPostUpdate_NotFound(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	postRepo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Delete ----------

func TestPostDelete_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	post := &model.Post{}
	post.ID = 10
	post.UserID = 1
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Delete", uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostDelete_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	post := &model.Post{}
	post.ID = 10
	post.UserID = 999
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- Timeline ----------

func TestPostTimeline_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/timeline", h.Timeline)

	postRepo.On("Timeline", uint(1), 1, 20).Return([]model.Post{
		{Title: "Timeline Post"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/timeline", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Like / Unlike ----------

func TestPostLike_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	postRepo.On("Like", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostUnlike_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	postRepo.On("Unlike", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Comments ----------

func TestPostGetComments_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	postRepo.On("GetComments", uint(5)).Return([]model.Comment{
		{Content: "Nice!"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/comments", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostCreateComment_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	postRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{
		"content": "Great post!",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestPostCreateComment_ValidationError(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	// content は required
	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDeleteComment_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	postRepo.On("DeleteComment", uint(3), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/3", nil)
	assertStatus(t, w, http.StatusOK)
}
