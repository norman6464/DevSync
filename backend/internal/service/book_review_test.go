package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTestBookReviewService() (*BookReviewService, *MockBookReviewRepository) {
	repo := new(MockBookReviewRepository)
	svc := NewBookReviewService(repo)
	return svc, repo
}

// ============================================================
// Update
// ============================================================

func TestBookReviewUpdate_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{Title: "Old", Author: "Author A", Rating: 3, UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.BookReview{Title: "New Title", Rating: 5}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, 5, result.Rating)
	assert.Equal(t, "Author A", result.Author) // 変更なし
	repo.AssertExpectations(t)
}

func TestBookReviewUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.BookReview{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestBookReviewUpdate_NotFound(t *testing.T) {
	svc, repo := newTestBookReviewService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.BookReview{Title: "New"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete
// ============================================================

func TestBookReviewDelete_Success(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBookReviewDelete_Forbidden(t *testing.T) {
	svc, repo := newTestBookReviewService()

	existing := &model.BookReview{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}
