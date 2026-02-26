package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// --- Mock ---

type MockResourceReviewService struct{ mock.Mock }

func (m *MockResourceReviewService) Create(review *model.ResourceReview) error {
	return m.Called(review).Error(0)
}

func (m *MockResourceReviewService) GetByResourceID(resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	args := m.Called(resourceID, limit, offset)
	return args.Get(0).([]model.ResourceReview), args.Get(1).(int64), args.Error(2)
}

func (m *MockResourceReviewService) Update(id, userID uint, rating int, comment string) (*model.ResourceReview, error) {
	args := m.Called(id, userID, rating, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceReview), args.Error(1)
}

func (m *MockResourceReviewService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

func setupResourceReviewHandler() (*ResourceReviewHandler, *MockResourceReviewService) {
	svc := new(MockResourceReviewService)
	h := NewResourceReviewHandler(svc)
	return h, svc
}

// --- Create ---

func TestResourceReview_Create_Success(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("Create", mock.AnythingOfType("*model.ResourceReview")).Return(nil)

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequest(r, http.MethodPost, "/resources/10/reviews", map[string]interface{}{
		"rating":  5,
		"comment": "素晴らしいリソース",
	})

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestResourceReview_Create_InvalidID(t *testing.T) {
	h, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequest(r, http.MethodPost, "/resources/abc/reviews", map[string]interface{}{
		"rating":  5,
		"comment": "テスト",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Create_InvalidJSON(t *testing.T) {
	h, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/resources/10/reviews", "invalid json")

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Create_ServiceError(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("Create", mock.AnythingOfType("*model.ResourceReview")).Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequest(r, http.MethodPost, "/resources/10/reviews", map[string]interface{}{
		"rating":  3,
		"comment": "普通",
	})

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// --- GetByResourceID ---

func TestResourceReview_GetByResourceID_Success(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	reviews := []model.ResourceReview{
		{Rating: 5, Comment: "最高"},
		{Rating: 3, Comment: "普通"},
	}
	svc.On("GetByResourceID", uint(10), 20, 0).Return(reviews, int64(2), nil)

	r := newRouter(1)
	r.GET("/resources/:id/reviews", h.GetByResourceID)

	w := doRequest(r, http.MethodGet, "/resources/10/reviews", nil)

	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	reviewList := data["reviews"].([]interface{})
	if len(reviewList) != 2 {
		t.Errorf("expected 2 reviews, got %d", len(reviewList))
	}
	svc.AssertExpectations(t)
}

func TestResourceReview_GetByResourceID_Empty(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("GetByResourceID", uint(10), 20, 0).Return([]model.ResourceReview{}, int64(0), nil)

	r := newRouter(1)
	r.GET("/resources/:id/reviews", h.GetByResourceID)

	w := doRequest(r, http.MethodGet, "/resources/10/reviews", nil)

	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	reviewList := data["reviews"].([]interface{})
	if len(reviewList) != 0 {
		t.Errorf("expected 0 reviews, got %d", len(reviewList))
	}
	svc.AssertExpectations(t)
}

func TestResourceReview_GetByResourceID_ServiceError(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("GetByResourceID", uint(10), 20, 0).Return([]model.ResourceReview{}, int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/resources/:id/reviews", h.GetByResourceID)

	w := doRequest(r, http.MethodGet, "/resources/10/reviews", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// --- Update ---

func TestResourceReview_Update_Success(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	updated := &model.ResourceReview{Rating: 4, Comment: "更新後"}
	svc.On("Update", uint(5), uint(1), 4, "更新後").Return(updated, nil)

	r := newRouter(1)
	r.PUT("/reviews/:reviewId", h.Update)

	w := doRequest(r, http.MethodPut, "/reviews/5", map[string]interface{}{
		"rating":  4,
		"comment": "更新後",
	})

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestResourceReview_Update_InvalidID(t *testing.T) {
	h, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.PUT("/reviews/:reviewId", h.Update)

	w := doRequest(r, http.MethodPut, "/reviews/abc", map[string]interface{}{
		"rating":  4,
		"comment": "テスト",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Update_ServiceError(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("Update", uint(5), uint(1), 4, "テスト").Return(nil, errors.New("forbidden"))

	r := newRouter(1)
	r.PUT("/reviews/:reviewId", h.Update)

	w := doRequest(r, http.MethodPut, "/reviews/5", map[string]interface{}{
		"rating":  4,
		"comment": "テスト",
	})

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// --- Delete ---

func TestResourceReview_Delete_Success(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("Delete", uint(5), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/reviews/:reviewId", h.Delete)

	w := doRequest(r, http.MethodDelete, "/reviews/5", nil)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestResourceReview_Delete_InvalidID(t *testing.T) {
	h, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.DELETE("/reviews/:reviewId", h.Delete)

	w := doRequest(r, http.MethodDelete, "/reviews/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Delete_ServiceError(t *testing.T) {
	h, svc := setupResourceReviewHandler()

	svc.On("Delete", uint(5), uint(1)).Return(errors.New("not found"))

	r := newRouter(1)
	r.DELETE("/reviews/:reviewId", h.Delete)

	w := doRequest(r, http.MethodDelete, "/reviews/5", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}
