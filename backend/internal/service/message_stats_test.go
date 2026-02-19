package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newMessageStatsTestService() (*MessageStatsService, *MockMessageStatsRepository) {
	repo := new(MockMessageStatsRepository)
	svc := NewMessageStatsService(repo)
	return svc, repo
}

func TestMessageStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newMessageStatsTestService()
	expected := &model.MessageStats{
		TotalSent:          150,
		TotalReceived:      200,
		ConversationsCount: 12,
		MessagesThisMonth:  45,
	}
	repo.On("GetMessageStats", uint(1)).Return(expected, nil)

	result, err := svc.GetMessageStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestMessageStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newMessageStatsTestService()

	_, err := svc.GetMessageStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestMessageStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newMessageStatsTestService()
	repo.On("GetMessageStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetMessageStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestMessageStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newMessageStatsTestService()
	expected := &model.MessageStats{}
	repo.On("GetMessageStats", uint(99)).Return(expected, nil)

	result, err := svc.GetMessageStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalSent)
	assert.Equal(t, int64(0), result.TotalReceived)
	assert.Equal(t, int64(0), result.ConversationsCount)
	assert.Equal(t, int64(0), result.MessagesThisMonth)
	repo.AssertExpectations(t)
}
