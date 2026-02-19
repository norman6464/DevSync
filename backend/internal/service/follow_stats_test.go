package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newFollowStatsTestService() (*FollowStatsService, *MockFollowStatsRepository) {
	repo := new(MockFollowStatsRepository)
	svc := NewFollowStatsService(repo)
	return svc, repo
}

func TestFollowStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newFollowStatsTestService()
	expected := &model.FollowStats{
		FollowerCount:  42,
		FollowingCount: 15,
	}
	repo.On("GetFollowStats", uint(1)).Return(expected, nil)

	result, err := svc.GetFollowStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestFollowStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newFollowStatsTestService()

	_, err := svc.GetFollowStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestFollowStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newFollowStatsTestService()
	repo.On("GetFollowStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetFollowStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestFollowStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newFollowStatsTestService()
	expected := &model.FollowStats{}
	repo.On("GetFollowStats", uint(99)).Return(expected, nil)

	result, err := svc.GetFollowStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.FollowerCount)
	assert.Equal(t, int64(0), result.FollowingCount)
	repo.AssertExpectations(t)
}
