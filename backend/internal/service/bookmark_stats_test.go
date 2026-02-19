package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newBookmarkStatsTestService() (*BookmarkStatsService, *MockBookmarkStatsRepository) {
	repo := new(MockBookmarkStatsRepository)
	svc := NewBookmarkStatsService(repo)
	return svc, repo
}

func TestBookmarkStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newBookmarkStatsTestService()
	expected := &model.BookmarkStats{
		TotalBookmarksMade:     50,
		TotalBookmarksReceived: 200,
		BookmarksThisMonth:     12,
	}
	repo.On("GetBookmarkStats", uint(1)).Return(expected, nil)

	result, err := svc.GetBookmarkStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestBookmarkStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newBookmarkStatsTestService()

	_, err := svc.GetBookmarkStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestBookmarkStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newBookmarkStatsTestService()
	repo.On("GetBookmarkStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetBookmarkStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestBookmarkStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newBookmarkStatsTestService()
	expected := &model.BookmarkStats{}
	repo.On("GetBookmarkStats", uint(99)).Return(expected, nil)

	result, err := svc.GetBookmarkStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalBookmarksMade)
	assert.Equal(t, int64(0), result.TotalBookmarksReceived)
	assert.Equal(t, int64(0), result.BookmarksThisMonth)
	repo.AssertExpectations(t)
}
