package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newProjectStatsTestService() (*ProjectStatsService, *MockProjectStatsRepository) {
	repo := new(MockProjectStatsRepository)
	svc := NewProjectStatsService(repo)
	return svc, repo
}

func TestProjectStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newProjectStatsTestService()
	expected := &model.ProjectStats{
		TotalProjects:     8,
		FeaturedProjects:  2,
		OngoingProjects:   5,
		CompletedProjects: 3,
	}
	repo.On("GetProjectStats", uint(1)).Return(expected, nil)

	result, err := svc.GetProjectStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestProjectStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newProjectStatsTestService()

	_, err := svc.GetProjectStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestProjectStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newProjectStatsTestService()
	repo.On("GetProjectStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetProjectStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestProjectStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newProjectStatsTestService()
	expected := &model.ProjectStats{}
	repo.On("GetProjectStats", uint(99)).Return(expected, nil)

	result, err := svc.GetProjectStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalProjects)
	assert.Equal(t, int64(0), result.FeaturedProjects)
	repo.AssertExpectations(t)
}
