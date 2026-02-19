package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newLearningLogStatsTestService() (*LearningLogStatsService, *MockLearningLogStatsRepository) {
	repo := new(MockLearningLogStatsRepository)
	svc := NewLearningLogStatsService(repo)
	return svc, repo
}

func TestLearningLogStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newLearningLogStatsTestService()
	expected := &model.LearningLogStats{
		TotalLogs:     120,
		TotalDuration: 3600,
		CategoryCount: 4,
		LogsThisMonth: 15,
	}
	repo.On("GetLearningLogStats", uint(1)).Return(expected, nil)

	result, err := svc.GetLearningLogStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestLearningLogStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newLearningLogStatsTestService()

	_, err := svc.GetLearningLogStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestLearningLogStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newLearningLogStatsTestService()
	repo.On("GetLearningLogStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetLearningLogStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestLearningLogStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newLearningLogStatsTestService()
	expected := &model.LearningLogStats{}
	repo.On("GetLearningLogStats", uint(99)).Return(expected, nil)

	result, err := svc.GetLearningLogStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalLogs)
	assert.Equal(t, int64(0), result.TotalDuration)
	repo.AssertExpectations(t)
}
