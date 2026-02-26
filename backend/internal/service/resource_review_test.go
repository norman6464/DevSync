package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestResourceReviewService はResourceReviewServiceのテスト用インスタンスを生成するヘルパー。
func newTestResourceReviewService() (*ResourceReviewService, *MockResourceReviewRepository, *MockLearningResourceRepository) {
	reviewRepo := new(MockResourceReviewRepository)
	resourceRepo := new(MockLearningResourceRepository)
	svc := NewResourceReviewService(reviewRepo, resourceRepo)
	return svc, reviewRepo, resourceRepo
}

// ============================================================
// レビュー作成テスト
// ============================================================

func TestResourceReviewCreate_Success(t *testing.T) {
	svc, reviewRepo, resourceRepo := newTestResourceReviewService()

	resource := &model.LearningResource{UserID: 2}
	resource.ID = 10

	resourceRepo.On("FindByID", uint(10)).Return(resource, nil)
	reviewRepo.On("FindByUserAndResource", uint(1), uint(10)).Return(nil, errors.New("not found"))
	reviewRepo.On("Create", mock.AnythingOfType("*model.ResourceReview")).Return(nil)

	review := &model.ResourceReview{
		UserID:     1,
		ResourceID: 10,
		Rating:     4,
		Comment:    "とても良いリソースです",
	}

	err := svc.Create(review)
	assert.NoError(t, err)
	reviewRepo.AssertExpectations(t)
}

func TestResourceReviewCreate_ResourceNotFound(t *testing.T) {
	svc, _, resourceRepo := newTestResourceReviewService()

	resourceRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	review := &model.ResourceReview{
		UserID:     1,
		ResourceID: 99,
		Rating:     4,
	}

	err := svc.Create(review)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestResourceReviewCreate_InvalidRating_TooLow(t *testing.T) {
	svc, _, resourceRepo := newTestResourceReviewService()

	resource := &model.LearningResource{UserID: 2}
	resource.ID = 10
	resourceRepo.On("FindByID", uint(10)).Return(resource, nil)

	review := &model.ResourceReview{
		UserID:     1,
		ResourceID: 10,
		Rating:     0,
	}

	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "1〜5")
}

func TestResourceReviewCreate_InvalidRating_TooHigh(t *testing.T) {
	svc, _, resourceRepo := newTestResourceReviewService()

	resource := &model.LearningResource{UserID: 2}
	resource.ID = 10
	resourceRepo.On("FindByID", uint(10)).Return(resource, nil)

	review := &model.ResourceReview{
		UserID:     1,
		ResourceID: 10,
		Rating:     6,
	}

	err := svc.Create(review)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "1〜5")
}

func TestResourceReviewCreate_Duplicate(t *testing.T) {
	svc, reviewRepo, resourceRepo := newTestResourceReviewService()

	resource := &model.LearningResource{UserID: 2}
	resource.ID = 10

	resourceRepo.On("FindByID", uint(10)).Return(resource, nil)
	existing := &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: 3}
	reviewRepo.On("FindByUserAndResource", uint(1), uint(10)).Return(existing, nil)

	review := &model.ResourceReview{
		UserID:     1,
		ResourceID: 10,
		Rating:     4,
	}

	err := svc.Create(review)
	assert.ErrorIs(t, err, ErrConflict)
}

// ============================================================
// レビュー取得テスト
// ============================================================

func TestResourceReviewGetByResourceID_Success(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	reviews := []model.ResourceReview{
		{Rating: 5, Comment: "素晴らしい"},
		{Rating: 3, Comment: "普通"},
	}
	reviewRepo.On("FindByResourceID", uint(10), 20, 0).Return(reviews, int64(2), nil)

	result, total, err := svc.GetByResourceID(10, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	reviewRepo.AssertExpectations(t)
}

func TestResourceReviewGetByResourceID_Empty(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	reviewRepo.On("FindByResourceID", uint(10), 20, 0).Return([]model.ResourceReview{}, int64(0), nil)

	result, total, err := svc.GetByResourceID(10, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	reviewRepo.AssertExpectations(t)
}

// ============================================================
// レビュー更新テスト
// ============================================================

func TestResourceReviewUpdate_Success(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	existing := &model.ResourceReview{UserID: 1, ResourceID: 10, Rating: 3, Comment: "普通"}
	existing.ID = 1

	reviewRepo.On("FindByID", uint(1)).Return(existing, nil)
	reviewRepo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, 5, "最高でした")
	assert.NoError(t, err)
	assert.Equal(t, 5, result.Rating)
	assert.Equal(t, "最高でした", result.Comment)
	reviewRepo.AssertExpectations(t)
}

func TestResourceReviewUpdate_Forbidden(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	existing := &model.ResourceReview{UserID: 1}
	existing.ID = 1

	reviewRepo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 999, 5, "")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestResourceReviewUpdate_NotFound(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	reviewRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	_, err := svc.Update(99, 1, 5, "")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestResourceReviewUpdate_InvalidRating(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	existing := &model.ResourceReview{UserID: 1, Rating: 3}
	existing.ID = 1

	reviewRepo.On("FindByID", uint(1)).Return(existing, nil)

	_, err := svc.Update(1, 1, 6, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "1〜5")
}

// ============================================================
// レビュー削除テスト
// ============================================================

func TestResourceReviewDelete_Success(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	existing := &model.ResourceReview{UserID: 1}
	existing.ID = 1

	reviewRepo.On("FindByID", uint(1)).Return(existing, nil)
	reviewRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	reviewRepo.AssertExpectations(t)
}

func TestResourceReviewDelete_Forbidden(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	existing := &model.ResourceReview{UserID: 1}
	existing.ID = 1

	reviewRepo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestResourceReviewDelete_NotFound(t *testing.T) {
	svc, reviewRepo, _ := newTestResourceReviewService()

	reviewRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.Delete(99, 1)
	assert.ErrorIs(t, err, ErrNotFound)
}
