package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newRoadmapStatsTestService() (*RoadmapStatsService, *MockRoadmapStatsRepository) {
	repo := new(MockRoadmapStatsRepository)
	svc := NewRoadmapStatsService(repo)
	return svc, repo
}

func TestRoadmapStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newRoadmapStatsTestService()
	expected := &model.RoadmapStats{
		TotalRoadmaps:     5,
		ActiveRoadmaps:    3,
		CompletedRoadmaps: 2,
		TotalSteps:        20,
		CompletedSteps:    12,
	}
	repo.On("GetRoadmapStats", uint(1)).Return(expected, nil)

	result, err := svc.GetRoadmapStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestRoadmapStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newRoadmapStatsTestService()

	_, err := svc.GetRoadmapStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestRoadmapStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newRoadmapStatsTestService()
	repo.On("GetRoadmapStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetRoadmapStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newRoadmapStatsTestService()
	expected := &model.RoadmapStats{}
	repo.On("GetRoadmapStats", uint(99)).Return(expected, nil)

	result, err := svc.GetRoadmapStats(99)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.TotalRoadmaps)
	assert.Equal(t, 0, result.CompletedSteps)
	repo.AssertExpectations(t)
}
