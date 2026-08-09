package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// otherComment は「他人（UserID=2）のコメント」。自己操作チェックを通す用。
func otherComment() *model.Comment { return &model.Comment{UserID: 2} }

// ---------- Like ----------

func TestCommentLike_Like_Success(t *testing.T) {
	h, likes, reader := setupCommentLikeHandler()
	r := newRouter(1)
	r.POST("/comments/:id/like", h.Like)

	reader.On("FindCommentByID", mock.Anything, uint(5)).Return(otherComment(), nil)
	likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(false, nil)
	likes.On("Like", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	likes.AssertExpectations(t)
}

func TestCommentLike_Like_InvalidID(t *testing.T) {
	h, _, _ := setupCommentLikeHandler()
	r := newRouter(1)
	r.POST("/comments/:id/like", h.Like)

	w := doRequest(r, http.MethodPost, "/comments/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentLike_Like_AlreadyLiked_BadRequest(t *testing.T) {
	h, likes, reader := setupCommentLikeHandler()
	r := newRouter(1)
	r.POST("/comments/:id/like", h.Like)

	reader.On("FindCommentByID", mock.Anything, uint(5)).Return(otherComment(), nil)
	likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)

	w := doRequest(r, http.MethodPost, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Unlike ----------

func TestCommentLike_Unlike_Success(t *testing.T) {
	h, likes, reader := setupCommentLikeHandler()
	r := newRouter(1)
	r.DELETE("/comments/:id/like", h.Unlike)

	reader.On("FindCommentByID", mock.Anything, uint(5)).Return(otherComment(), nil)
	likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)
	likes.On("Unlike", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	likes.AssertExpectations(t)
}

func TestCommentLike_Unlike_InvalidID(t *testing.T) {
	h, _, _ := setupCommentLikeHandler()
	r := newRouter(1)
	r.DELETE("/comments/:id/like", h.Unlike)

	w := doRequest(r, http.MethodDelete, "/comments/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentLike_Unlike_CommentNotFound(t *testing.T) {
	h, _, reader := setupCommentLikeHandler()
	r := newRouter(1)
	r.DELETE("/comments/:id/like", h.Unlike)

	reader.On("FindCommentByID", mock.Anything, uint(5)).Return((*model.Comment)(nil), errors.New("not found"))

	w := doRequest(r, http.MethodDelete, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- GetStatus ----------

func TestCommentLike_GetStatus_Success(t *testing.T) {
	h, likes, reader := setupCommentLikeHandler()
	r := newRouter(1)
	r.GET("/comments/:id/like", h.GetStatus)

	reader.On("FindCommentByID", mock.Anything, uint(5)).Return(otherComment(), nil)
	likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)
	likes.On("CountByCommentID", mock.Anything, uint(5)).Return(int64(3), nil)

	w := doRequest(r, http.MethodGet, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	likes.AssertExpectations(t)
}

func TestCommentLike_GetStatus_InvalidID(t *testing.T) {
	h, _, _ := setupCommentLikeHandler()
	r := newRouter(1)
	r.GET("/comments/:id/like", h.GetStatus)

	w := doRequest(r, http.MethodGet, "/comments/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestCommentLike_GetStatus_CommentNotFound(t *testing.T) {
	h, _, reader := setupCommentLikeHandler()
	r := newRouter(1)
	r.GET("/comments/:id/like", h.GetStatus)

	reader.On("FindCommentByID", mock.Anything, uint(5)).Return((*model.Comment)(nil), errors.New("not found"))

	w := doRequest(r, http.MethodGet, "/comments/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
}
