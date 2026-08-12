package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newBookmarkRouter はブックマーク系のルートを登録したルーターを返す。
func newBookmarkRouter(h *PostHandler, userID uint) *gin.Engine {
	r := newRouter(userID)
	r.POST("/posts/:id/bookmark", h.Bookmark)
	r.DELETE("/posts/:id/bookmark", h.Unbookmark)
	r.GET("/posts/bookmarks", h.GetBookmarks)
	r.GET("/posts/bookmarks/count", h.GetBookmarksCount)
	return r
}

// ---------- Bookmark ----------

func TestPostBookmark_Success(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	bookmarks.On("Bookmark", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusOK)
	bookmarks.AssertExpectations(t)
	authors.AssertExpectations(t)
}

func TestPostBookmark_OwnPostForbidden(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

	w := doRequest(r, http.MethodPost, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusForbidden)
	bookmarks.AssertNotCalled(t, "Bookmark", mock.Anything, mock.Anything, mock.Anything)
	authors.AssertExpectations(t)
}

func TestPostBookmark_PostNotFound(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), nil)

	w := doRequest(r, http.MethodPost, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusNotFound)
	bookmarks.AssertNotCalled(t, "Bookmark", mock.Anything, mock.Anything, mock.Anything)
	authors.AssertExpectations(t)
}

func TestPostBookmark_InvalidID(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	w := doRequest(r, http.MethodPost, "/posts/abc/bookmark", nil)
	assertStatus(t, w, http.StatusBadRequest)
	bookmarks.AssertNotCalled(t, "Bookmark", mock.Anything, mock.Anything, mock.Anything)
}

func TestPostBookmark_RepositoryError(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	bookmarks.On("Bookmark", mock.Anything, uint(1), uint(5)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	bookmarks.AssertExpectations(t)
	authors.AssertExpectations(t)
}

// ---------- Unbookmark ----------

func TestPostUnbookmark_Success(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	bookmarks.On("Unbookmark", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusOK)
	bookmarks.AssertExpectations(t)
	authors.AssertExpectations(t)
}

func TestPostUnbookmark_OwnPostForbidden(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusForbidden)
	bookmarks.AssertNotCalled(t, "Unbookmark", mock.Anything, mock.Anything, mock.Anything)
	authors.AssertExpectations(t)
}

func TestPostUnbookmark_PostNotFound(t *testing.T) {
	h, bookmarks, authors := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusNotFound)
	bookmarks.AssertNotCalled(t, "Unbookmark", mock.Anything, mock.Anything, mock.Anything)
	authors.AssertExpectations(t)
}

func TestPostUnbookmark_InvalidID(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	w := doRequest(r, http.MethodDelete, "/posts/abc/bookmark", nil)
	assertStatus(t, w, http.StatusBadRequest)
	bookmarks.AssertNotCalled(t, "Unbookmark", mock.Anything, mock.Anything, mock.Anything)
}

// ---------- GetBookmarks ----------

func TestPostGetBookmarks_Success(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	bookmarks.On("FindBookmarkedByUserID", mock.Anything, uint(1), 1, 20).
		Return([]model.Post{{ID: 3, Title: "Bookmarked 1"}}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/posts/bookmarks", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"total":1`)
	bookmarks.AssertExpectations(t)
}

// ページ指定は offset の算出に反映される。
func TestPostGetBookmarks_Pagination(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	bookmarks.On("FindBookmarkedByUserID", mock.Anything, uint(1), 2, 5).Return([]model.Post{}, int64(7), nil)

	w := doRequest(r, http.MethodGet, "/posts/bookmarks?page=2&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"offset":5`)
	bookmarks.AssertExpectations(t)
}

// 1 件も無ければ null ではなく空配列を返す。
func TestPostGetBookmarks_Empty(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	bookmarks.On("FindBookmarkedByUserID", mock.Anything, uint(1), 1, 20).Return(nil, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/posts/bookmarks", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"posts":[]`)
}

func TestPostGetBookmarks_RepositoryError(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	bookmarks.On("FindBookmarkedByUserID", mock.Anything, uint(1), 1, 20).
		Return(nil, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/bookmarks", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	bookmarks.AssertExpectations(t)
}

// ---------- GetBookmarksCount ----------

func TestPostGetBookmarksCount_Success(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	bookmarks.On("CountBookmarkedByUserID", mock.Anything, uint(1)).Return(int64(7), nil)

	w := doRequest(r, http.MethodGet, "/posts/bookmarks/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":7`)
	bookmarks.AssertExpectations(t)
}

func TestPostGetBookmarksCount_RepositoryError(t *testing.T) {
	h, bookmarks, _ := setupPostHandlerWithBookmarkPorts()
	r := newBookmarkRouter(h, 1)

	bookmarks.On("CountBookmarkedByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/bookmarks/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	bookmarks.AssertExpectations(t)
}
