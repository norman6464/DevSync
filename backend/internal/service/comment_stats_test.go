package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newCommentStatsTestService() (*CommentStatsService, *MockCommentStatsRepository) {
	repo := new(MockCommentStatsRepository)
	svc := NewCommentStatsService(repo)
	return svc, repo
}

func TestCommentStatsService_GetStats_Success(t *testing.T) {
	svc, repo := newCommentStatsTestService()
	expected := &model.CommentStats{
		TotalComments:     42,
		TotalReplies:      15,
		CommentsReceived:  28,
		CommentsThisMonth: 8,
	}
	repo.On("GetCommentStats", uint(1)).Return(expected, nil)

	result, err := svc.GetCommentStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCommentStatsService_GetStats_InvalidUserID(t *testing.T) {
	svc, _ := newCommentStatsTestService()

	_, err := svc.GetCommentStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestCommentStatsService_GetStats_RepoError(t *testing.T) {
	svc, repo := newCommentStatsTestService()
	repo.On("GetCommentStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetCommentStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCommentStatsService_GetStats_NoActivity(t *testing.T) {
	svc, repo := newCommentStatsTestService()
	expected := &model.CommentStats{}
	repo.On("GetCommentStats", uint(99)).Return(expected, nil)

	result, err := svc.GetCommentStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalComments)
	assert.Equal(t, int64(0), result.TotalReplies)
	assert.Equal(t, int64(0), result.CommentsReceived)
	assert.Equal(t, int64(0), result.CommentsThisMonth)
	repo.AssertExpectations(t)
}
