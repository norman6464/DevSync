package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Create テスト
// ============================================================

func TestBookReviewCreate_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.POST("/book-reviews", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.BookReview")).Return(nil)

	w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
		"title":  "Go言語入門",
		"author": "テスト著者",
		"rating": 5,
		"review": "とても良い本です",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestBookReviewCreate_ValidationError(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.POST("/book-reviews", h.Create)

	// title と rating は required
	w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
		"review": "レビューのみ",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewCreate_InvalidJSON(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.POST("/book-reviews", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/book-reviews", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewCreate_ServiceError(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.POST("/book-reviews", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.BookReview")).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/book-reviews", map[string]interface{}{
		"title":  "Go言語入門",
		"rating": 5,
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByID テスト
// ============================================================

func TestBookReviewGetByID_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/:id", h.GetByID)

	review := &model.BookReview{Title: "Go言語入門", Rating: 5}
	review.ID = 1
	svc.On("GetByID", uint(1)).Return(review, nil)

	w := doRequest(r, http.MethodGet, "/book-reviews/1", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, "Go言語入門", body["title"])
	svc.AssertExpectations(t)
}

func TestBookReviewGetByID_NotFound(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/:id", h.GetByID)

	svc.On("GetByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/book-reviews/999", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestBookReviewGetByID_InvalidID(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/book-reviews/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetAll テスト
// ============================================================

func TestBookReviewGetAll_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews", h.GetAll)

	reviews := []model.BookReview{
		{Title: "Go言語入門", Rating: 5},
		{Title: "Rust入門", Rating: 4},
	}
	svc.On("GetAll", 20, 0).Return(reviews, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/book-reviews", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(2), body["total"])
	svc.AssertExpectations(t)
}

func TestBookReviewGetAll_WithPagination(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews", h.GetAll)

	svc.On("GetAll", 10, 5).Return([]model.BookReview{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/book-reviews?limit=10&offset=5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestBookReviewGetAll_LimitCap(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews", h.GetAll)

	// limit=200 は 100 に制限される
	svc.On("GetAll", 100, 0).Return([]model.BookReview{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/book-reviews?limit=200", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestBookReviewGetByUserID_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/users/:userId/book-reviews", h.GetByUserID)

	reviews := []model.BookReview{
		{Title: "Go言語入門", Rating: 5},
	}
	svc.On("GetByUserID", uint(1), 20, 0).Return(reviews, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/users/1/book-reviews", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestBookReviewGetByUserID_ServiceError(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/users/:userId/book-reviews", h.GetByUserID)

	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.BookReview(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/1/book-reviews", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestBookReviewGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/users/:userId/book-reviews", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/book-reviews", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetAll テスト（追加分）
// ============================================================

func TestBookReviewGetAll_ServiceError(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews", h.GetAll)

	svc.On("GetAll", 20, 0).Return([]model.BookReview(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/book-reviews", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestBookReviewUpdate_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)

	updated := &model.BookReview{Title: "更新後タイトル", Rating: 4}
	updated.ID = 1
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.BookReview")).Return(updated, nil)

	w := doRequest(r, http.MethodPut, "/book-reviews/1", map[string]interface{}{
		"title": "更新後タイトル",
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestBookReviewUpdate_Forbidden(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)

	svc.On("Update", uint(5), uint(1), mock.AnythingOfType("*model.BookReview")).Return(nil, service.ErrForbidden)

	w := doRequest(r, http.MethodPut, "/book-reviews/5", map[string]interface{}{
		"title": "不正更新",
	})
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestBookReviewUpdate_InvalidJSON(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/book-reviews/1", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewUpdate_InvalidID(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/book-reviews/abc", map[string]interface{}{
		"title": "テスト",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewUpdate_NotFound(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id", h.Update)

	svc.On("Update", uint(999), uint(1), mock.AnythingOfType("*model.BookReview")).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/book-reviews/999", map[string]interface{}{
		"title": "存在しない",
	})
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

// ============================================================
// Delete テスト
// ============================================================

func TestBookReviewDelete_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.DELETE("/book-reviews/:id", h.Delete)

	svc.On("Delete", uint(1), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/book-reviews/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestBookReviewDelete_Forbidden(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.DELETE("/book-reviews/:id", h.Delete)

	svc.On("Delete", uint(5), uint(1)).Return(service.ErrForbidden)

	w := doRequest(r, http.MethodDelete, "/book-reviews/5", nil)
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestBookReviewDelete_NotFound(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.DELETE("/book-reviews/:id", h.Delete)

	svc.On("Delete", uint(999), uint(1)).Return(service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/book-reviews/999", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByRating テスト
// ============================================================

func TestBookReviewGetByRating_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/rating", h.GetByRating)

	reviews := []model.BookReview{
		{Title: "良書", Rating: 4},
	}
	svc.On("GetByRating", uint(1), 4, 5).Return(reviews, nil)

	w := doRequest(r, http.MethodGet, "/book-reviews/rating?min_rating=4&max_rating=5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestBookReviewGetByRating_InvalidMinRating(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/rating", h.GetByRating)

	w := doRequest(r, http.MethodGet, "/book-reviews/rating?min_rating=abc&max_rating=5", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewGetByRating_InvalidMaxRating(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/rating", h.GetByRating)

	w := doRequest(r, http.MethodGet, "/book-reviews/rating?min_rating=1&max_rating=abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestBookReviewGetByRating_ServiceError(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.GET("/book-reviews/rating", h.GetByRating)

	svc.On("GetByRating", uint(1), 1, 5).Return([]model.BookReview(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/book-reviews/rating?min_rating=1&max_rating=5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Archive テスト
// ============================================================

func TestBookReviewArchive_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id/archive", h.Archive)

	svc.On("ArchiveReview", uint(1), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/book-reviews/1/archive", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "アーカイブしました")
	svc.AssertExpectations(t)
}

func TestBookReviewArchive_Forbidden(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id/archive", h.Archive)

	svc.On("ArchiveReview", uint(5), uint(1)).Return(service.ErrForbidden)

	w := doRequest(r, http.MethodPut, "/book-reviews/5/archive", nil)
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestBookReviewArchive_InvalidID(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id/archive", h.Archive)

	w := doRequest(r, http.MethodPut, "/book-reviews/abc/archive", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// Unarchive テスト
// ============================================================

func TestBookReviewUnarchive_Success(t *testing.T) {
	h, svc := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id/unarchive", h.Unarchive)

	svc.On("UnarchiveReview", uint(1), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/book-reviews/1/unarchive", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "アーカイブを解除しました")
	svc.AssertExpectations(t)
}

func TestBookReviewUnarchive_InvalidID(t *testing.T) {
	h, _ := setupBookReviewHandler()
	r := newRouter(1)
	r.PUT("/book-reviews/:id/unarchive", h.Unarchive)

	w := doRequest(r, http.MethodPut, "/book-reviews/abc/unarchive", nil)
	assertStatus(t, w, http.StatusBadRequest)
}
