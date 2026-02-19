package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newLearningResourceStatsTestService() (*LearningResourceStatsService, *MockLearningResourceStatsRepository) {
	repo := new(MockLearningResourceStatsRepository)
	svc := NewLearningResourceStatsService(repo)
	return svc, repo
}

func TestLearningResourceStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newLearningResourceStatsTestService()
	expected := &model.LearningResourceStats{
		TotalResources: 20,
		TotalLikes:     85,
		TotalSaves:     30,
		CategoryCount:  4,
	}
	repo.On("GetLearningResourceStats", uint(1)).Return(expected, nil)

	result, err := svc.GetLearningResourceStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestLearningResourceStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newLearningResourceStatsTestService()

	_, err := svc.GetLearningResourceStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestLearningResourceStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newLearningResourceStatsTestService()
	repo.On("GetLearningResourceStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetLearningResourceStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newLearningResourceStatsTestService()
	expected := &model.LearningResourceStats{}
	repo.On("GetLearningResourceStats", uint(99)).Return(expected, nil)

	result, err := svc.GetLearningResourceStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalResources)
	assert.Equal(t, int64(0), result.TotalLikes)
	assert.Equal(t, int64(0), result.TotalSaves)
	assert.Equal(t, int64(0), result.CategoryCount)
	repo.AssertExpectations(t)
}
