package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newMentionStatsTestService() (*MentionStatsService, *MockMentionStatsRepository) {
	repo := new(MockMentionStatsRepository)
	svc := NewMentionStatsService(repo)
	return svc, repo
}

func TestMentionStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newMentionStatsTestService()
	expected := &model.MentionStats{
		MentionsReceived:  42,
		MentionsMade:      15,
		MentionsThisMonth: 8,
	}
	repo.On("GetMentionStats", uint(1)).Return(expected, nil)

	result, err := svc.GetMentionStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestMentionStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newMentionStatsTestService()

	_, err := svc.GetMentionStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestMentionStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newMentionStatsTestService()
	repo.On("GetMentionStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetMentionStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestMentionStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newMentionStatsTestService()
	expected := &model.MentionStats{}
	repo.On("GetMentionStats", uint(99)).Return(expected, nil)

	result, err := svc.GetMentionStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.MentionsReceived)
	assert.Equal(t, int64(0), result.MentionsMade)
	assert.Equal(t, int64(0), result.MentionsThisMonth)
	repo.AssertExpectations(t)
}
