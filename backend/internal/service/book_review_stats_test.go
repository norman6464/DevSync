package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newBookReviewStatsTestService() (*BookReviewStatsService, *MockBookReviewStatsRepository) {
	repo := new(MockBookReviewStatsRepository)
	svc := NewBookReviewStatsService(repo)
	return svc, repo
}

func TestBookReviewStatsService_GetBookReviewStats_Success(t *testing.T) {
	svc, repo := newBookReviewStatsTestService()
	expected := &model.BookReviewStats{
		TotalReviews:  5,
		AverageRating: 3.8,
		MaxRating:     5,
		MinRating:     2,
		FiveStarCount: 2,
	}
	repo.On("GetBookReviewStats", uint(1)).Return(expected, nil)

	result, err := svc.GetBookReviewStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestBookReviewStatsService_GetBookReviewStats_InvalidUserID(t *testing.T) {
	svc, _ := newBookReviewStatsTestService()

	_, err := svc.GetBookReviewStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestBookReviewStatsService_GetBookReviewStats_RepoError(t *testing.T) {
	svc, repo := newBookReviewStatsTestService()
	repo.On("GetBookReviewStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetBookReviewStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestBookReviewStatsService_GetBookReviewStats_NoReviews(t *testing.T) {
	svc, repo := newBookReviewStatsTestService()
	expected := &model.BookReviewStats{}
	repo.On("GetBookReviewStats", uint(99)).Return(expected, nil)

	result, err := svc.GetBookReviewStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalReviews)
	repo.AssertExpectations(t)
}
