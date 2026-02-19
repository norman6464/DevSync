package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newQAStatsTestService() (*QAStatsService, *MockQAStatsRepository) {
	repo := new(MockQAStatsRepository)
	svc := NewQAStatsService(repo)
	return svc, repo
}

func TestQAStatsService_GetQAStats_Success(t *testing.T) {
	svc, repo := newQAStatsTestService()
	expected := &model.QAStats{
		TotalQuestions:     5,
		TotalAnswers:       12,
		BestAnswerCount:    3,
		TotalVotesReceived: 28,
	}
	repo.On("GetQAStats", uint(1)).Return(expected, nil)

	result, err := svc.GetQAStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestQAStatsService_GetQAStats_InvalidUserID(t *testing.T) {
	svc, _ := newQAStatsTestService()

	_, err := svc.GetQAStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestQAStatsService_GetQAStats_RepoError(t *testing.T) {
	svc, repo := newQAStatsTestService()
	repo.On("GetQAStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetQAStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestQAStatsService_GetQAStats_NoActivity(t *testing.T) {
	svc, repo := newQAStatsTestService()
	expected := &model.QAStats{}
	repo.On("GetQAStats", uint(99)).Return(expected, nil)

	result, err := svc.GetQAStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalQuestions)
	assert.Equal(t, int64(0), result.TotalAnswers)
	repo.AssertExpectations(t)
}
