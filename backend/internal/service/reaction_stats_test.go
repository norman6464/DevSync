package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newReactionStatsTestService() (*ReactionStatsService, *MockReactionStatsRepository) {
	repo := new(MockReactionStatsRepository)
	svc := NewReactionStatsService(repo)
	return svc, repo
}

func TestReactionStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newReactionStatsTestService()
	expected := &model.ReactionStats{
		TotalReactionsReceived: 120,
		UniqueReactors:         35,
		ReactionsThisMonth:     18,
	}
	repo.On("GetReactionStats", uint(1)).Return(expected, nil)

	result, err := svc.GetReactionStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestReactionStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newReactionStatsTestService()

	_, err := svc.GetReactionStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestReactionStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newReactionStatsTestService()
	repo.On("GetReactionStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetReactionStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestReactionStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newReactionStatsTestService()
	expected := &model.ReactionStats{}
	repo.On("GetReactionStats", uint(99)).Return(expected, nil)

	result, err := svc.GetReactionStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalReactionsReceived)
	assert.Equal(t, int64(0), result.UniqueReactors)
	assert.Equal(t, int64(0), result.ReactionsThisMonth)
	repo.AssertExpectations(t)
}
