package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// --- port モック（ctx 付き） ---

// mockResourceReviewRepo は usecase/repository.ResourceReviewRepository のモック。
type mockResourceReviewRepo struct{ mock.Mock }

func (m *mockResourceReviewRepo) Create(ctx context.Context, review *model.ResourceReview) error {
	return m.Called(ctx, review).Error(0)
}

func (m *mockResourceReviewRepo) FindByID(ctx context.Context, id uint) (*model.ResourceReview, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.ResourceReview)
	return r, args.Error(1)
}

func (m *mockResourceReviewRepo) FindByResourceID(ctx context.Context, resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	args := m.Called(ctx, resourceID, limit, offset)
	reviews, _ := args.Get(0).([]model.ResourceReview)
	return reviews, args.Get(1).(int64), args.Error(2)
}

func (m *mockResourceReviewRepo) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceReview, error) {
	args := m.Called(ctx, userID, resourceID)
	r, _ := args.Get(0).(*model.ResourceReview)
	return r, args.Error(1)
}

func (m *mockResourceReviewRepo) Update(ctx context.Context, review *model.ResourceReview) error {
	return m.Called(ctx, review).Error(0)
}

func (m *mockResourceReviewRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// mockLearningResourceReader は usecase/repository.LearningResourceReader のモック。
type mockLearningResourceReader struct{ mock.Mock }

func (m *mockLearningResourceReader) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.LearningResource)
	return r, args.Error(1)
}

// setupResourceReviewHandler は本物の usecase + port モックで ResourceReviewHandler を組む。
func setupResourceReviewHandler() (*ResourceReviewHandler, *mockResourceReviewRepo, *mockLearningResourceReader) {
	reviews := new(mockResourceReviewRepo)
	resources := new(mockLearningResourceReader)
	h := NewResourceReviewHandler(
		usecase.NewCreateResourceReviewUseCase(reviews, resources),
		usecase.NewListResourceReviewsUseCase(reviews),
		usecase.NewUpdateResourceReviewUseCase(reviews),
		usecase.NewDeleteResourceReviewUseCase(reviews),
	)
	return h, reviews, resources
}

// --- Create ---

func TestResourceReview_Create_Success(t *testing.T) {
	h, reviews, resources := setupResourceReviewHandler()
	resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
	reviews.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).
		Return((*model.ResourceReview)(nil), nil)
	reviews.On("Create", mock.Anything, mock.AnythingOfType("*model.ResourceReview")).Return(nil)

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequest(r, http.MethodPost, "/resources/10/reviews", map[string]interface{}{
		"rating":  5,
		"comment": "素晴らしいリソース",
	})

	assertStatus(t, w, http.StatusCreated)
	reviews.AssertExpectations(t)
}

func TestResourceReview_Create_InvalidID(t *testing.T) {
	h, _, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequest(r, http.MethodPost, "/resources/abc/reviews", map[string]interface{}{
		"rating":  5,
		"comment": "テスト",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Create_InvalidJSON(t *testing.T) {
	h, _, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/resources/10/reviews", "invalid json")

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Create_ServiceError(t *testing.T) {
	h, reviews, resources := setupResourceReviewHandler()
	resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
	reviews.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).
		Return((*model.ResourceReview)(nil), nil)
	reviews.On("Create", mock.Anything, mock.AnythingOfType("*model.ResourceReview")).
		Return(errors.New("db error"))

	r := newRouter(1)
	r.POST("/resources/:id/reviews", h.Create)

	w := doRequest(r, http.MethodPost, "/resources/10/reviews", map[string]interface{}{
		"rating":  3,
		"comment": "普通",
	})

	assertStatus(t, w, http.StatusInternalServerError)
	reviews.AssertExpectations(t)
}

// --- GetByResourceID ---

func TestResourceReview_GetByResourceID_Success(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	list := []model.ResourceReview{
		{Rating: 5, Comment: "最高"},
		{Rating: 3, Comment: "普通"},
	}
	reviews.On("FindByResourceID", mock.Anything, uint(10), 20, 0).Return(list, int64(2), nil)

	r := newRouter(1)
	r.GET("/resources/:id/reviews", h.GetByResourceID)

	w := doRequest(r, http.MethodGet, "/resources/10/reviews", nil)

	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	reviewList := data["reviews"].([]interface{})
	if len(reviewList) != 2 {
		t.Errorf("expected 2 reviews, got %d", len(reviewList))
	}
	reviews.AssertExpectations(t)
}

func TestResourceReview_GetByResourceID_Empty(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	reviews.On("FindByResourceID", mock.Anything, uint(10), 20, 0).Return([]model.ResourceReview{}, int64(0), nil)

	r := newRouter(1)
	r.GET("/resources/:id/reviews", h.GetByResourceID)

	w := doRequest(r, http.MethodGet, "/resources/10/reviews", nil)

	assertStatus(t, w, http.StatusOK)
	data := parseJSON(t, w)
	reviewList := data["reviews"].([]interface{})
	if len(reviewList) != 0 {
		t.Errorf("expected 0 reviews, got %d", len(reviewList))
	}
	reviews.AssertExpectations(t)
}

func TestResourceReview_GetByResourceID_ServiceError(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	reviews.On("FindByResourceID", mock.Anything, uint(10), 20, 0).
		Return([]model.ResourceReview{}, int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/resources/:id/reviews", h.GetByResourceID)

	w := doRequest(r, http.MethodGet, "/resources/10/reviews", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	reviews.AssertExpectations(t)
}

// --- Update ---

func TestResourceReview_Update_Success(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	reviews.On("FindByID", mock.Anything, uint(5)).Return(&model.ResourceReview{UserID: 1, Rating: 3}, nil)
	reviews.On("Update", mock.Anything, mock.AnythingOfType("*model.ResourceReview")).Return(nil)

	r := newRouter(1)
	r.PUT("/reviews/:reviewId", h.Update)

	w := doRequest(r, http.MethodPut, "/reviews/5", map[string]interface{}{
		"rating":  4,
		"comment": "更新後",
	})

	assertStatus(t, w, http.StatusOK)
	reviews.AssertExpectations(t)
}

func TestResourceReview_Update_InvalidID(t *testing.T) {
	h, _, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.PUT("/reviews/:reviewId", h.Update)

	w := doRequest(r, http.MethodPut, "/reviews/abc", map[string]interface{}{
		"rating":  4,
		"comment": "テスト",
	})

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Update_ServiceError(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	reviews.On("FindByID", mock.Anything, uint(5)).Return(&model.ResourceReview{UserID: 1, Rating: 3}, nil)
	reviews.On("Update", mock.Anything, mock.AnythingOfType("*model.ResourceReview")).
		Return(errors.New("db error"))

	r := newRouter(1)
	r.PUT("/reviews/:reviewId", h.Update)

	w := doRequest(r, http.MethodPut, "/reviews/5", map[string]interface{}{
		"rating":  4,
		"comment": "テスト",
	})

	assertStatus(t, w, http.StatusInternalServerError)
	reviews.AssertExpectations(t)
}

// --- Delete ---

func TestResourceReview_Delete_Success(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	reviews.On("FindByID", mock.Anything, uint(5)).Return(&model.ResourceReview{UserID: 1}, nil)
	reviews.On("Delete", mock.Anything, uint(5)).Return(nil)

	r := newRouter(1)
	r.DELETE("/reviews/:reviewId", h.Delete)

	w := doRequest(r, http.MethodDelete, "/reviews/5", nil)

	assertStatus(t, w, http.StatusOK)
	reviews.AssertExpectations(t)
}

func TestResourceReview_Delete_InvalidID(t *testing.T) {
	h, _, _ := setupResourceReviewHandler()

	r := newRouter(1)
	r.DELETE("/reviews/:reviewId", h.Delete)

	w := doRequest(r, http.MethodDelete, "/reviews/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestResourceReview_Delete_ServiceError(t *testing.T) {
	h, reviews, _ := setupResourceReviewHandler()

	reviews.On("FindByID", mock.Anything, uint(5)).Return(&model.ResourceReview{UserID: 1}, nil)
	reviews.On("Delete", mock.Anything, uint(5)).Return(errors.New("db error"))

	r := newRouter(1)
	r.DELETE("/reviews/:reviewId", h.Delete)

	w := doRequest(r, http.MethodDelete, "/reviews/5", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	reviews.AssertExpectations(t)
}
