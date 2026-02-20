package handler

import (
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/service"
)

// ---------- Like ----------

func TestCommentLike_Like_Success(t *testing.T) {
	h, svc := setupCommentLikeHandler()
	r := newRouter(1)
	r.POST("/comments/:id/like", h.Like)

	svc.On("Like", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCommentLike_Like_InvalidID(t *testing.T) {
	h, _ := setupCommentLikeHandler()
	r := newRouter(1)
	r.POST("/comments/:id/like", h.Like)

	w := doRequest(r, http.MethodPost, "/comments/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentLike_Like_ServiceError(t *testing.T) {
	h, svc := setupCommentLikeHandler()
	r := newRouter(1)
	r.POST("/comments/:id/like", h.Like)

	svc.On("Like", uint(1), uint(5)).Return(service.ErrBadRequest)

	w := doRequest(r, http.MethodPost, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
	svc.AssertExpectations(t)
}

// ---------- Unlike ----------

func TestCommentLike_Unlike_Success(t *testing.T) {
	h, svc := setupCommentLikeHandler()
	r := newRouter(1)
	r.DELETE("/comments/:id/like", h.Unlike)

	svc.On("Unlike", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCommentLike_Unlike_InvalidID(t *testing.T) {
	h, _ := setupCommentLikeHandler()
	r := newRouter(1)
	r.DELETE("/comments/:id/like", h.Unlike)

	w := doRequest(r, http.MethodDelete, "/comments/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentLike_Unlike_ServiceError(t *testing.T) {
	h, svc := setupCommentLikeHandler()
	r := newRouter(1)
	r.DELETE("/comments/:id/like", h.Unlike)

	svc.On("Unlike", uint(1), uint(5)).Return(service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

// ---------- GetStatus ----------

func TestCommentLike_GetStatus_Success(t *testing.T) {
	h, svc := setupCommentLikeHandler()
	r := newRouter(1)
	r.GET("/comments/:id/like", h.GetStatus)

	svc.On("GetStatus", uint(1), uint(5)).Return(true, int64(3), nil)

	w := doRequest(r, http.MethodGet, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestCommentLike_GetStatus_InvalidID(t *testing.T) {
	h, _ := setupCommentLikeHandler()
	r := newRouter(1)
	r.GET("/comments/:id/like", h.GetStatus)

	w := doRequest(r, http.MethodGet, "/comments/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentLike_GetStatus_ServiceError(t *testing.T) {
	h, svc := setupCommentLikeHandler()
	r := newRouter(1)
	r.GET("/comments/:id/like", h.GetStatus)

	svc.On("GetStatus", uint(1), uint(5)).Return(false, int64(0), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}
