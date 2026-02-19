package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newPostStatsTestService() (*PostStatsService, *MockPostStatsRepository) {
	repo := new(MockPostStatsRepository)
	svc := NewPostStatsService(repo)
	return svc, repo
}

func TestPostStatsService_GetPostStats_Success(t *testing.T) {
	svc, repo := newPostStatsTestService()
	expected := &model.PostStats{
		TotalPosts:         10,
		PublishedPosts:     7,
		DraftPosts:         3,
		TotalLikesReceived: 42,
		TotalComments:      15,
		PostsThisMonth:     2,
	}
	repo.On("GetPostStats", uint(1)).Return(expected, nil)

	result, err := svc.GetPostStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestPostStatsService_GetPostStats_InvalidUserID(t *testing.T) {
	svc, _ := newPostStatsTestService()

	_, err := svc.GetPostStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestPostStatsService_GetPostStats_RepoError(t *testing.T) {
	svc, repo := newPostStatsTestService()
	repo.On("GetPostStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetPostStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestPostStatsService_GetPostStats_NoPosts(t *testing.T) {
	svc, repo := newPostStatsTestService()
	expected := &model.PostStats{}
	repo.On("GetPostStats", uint(99)).Return(expected, nil)

	result, err := svc.GetPostStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalPosts)
	repo.AssertExpectations(t)
}
