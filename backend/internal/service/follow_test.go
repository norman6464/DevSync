package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTestFollowService() (*FollowService, *MockFollowRepository) {
	repo := new(MockFollowRepository)
	svc := NewFollowService(repo)
	return svc, repo
}

// ============================================================
// Follow
// ============================================================

func TestFollow_Success(t *testing.T) {
	svc, repo := newTestFollowService()

	repo.On("Follow", uint(1), uint(2)).Return(nil)

	err := svc.Follow(1, 2)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestFollow_SelfFollow(t *testing.T) {
	svc, _ := newTestFollowService()

	err := svc.Follow(1, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
}

// ============================================================
// Unfollow
// ============================================================

func TestUnfollow_Success(t *testing.T) {
	svc, repo := newTestFollowService()

	repo.On("Unfollow", uint(1), uint(2)).Return(nil)

	err := svc.Unfollow(1, 2)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// IsFollowing
// ============================================================

func TestIsFollowing_True(t *testing.T) {
	svc, repo := newTestFollowService()

	repo.On("IsFollowing", uint(1), uint(2)).Return(true)

	result := svc.IsFollowing(1, 2)
	assert.True(t, result)
	repo.AssertExpectations(t)
}

func TestIsFollowing_False(t *testing.T) {
	svc, repo := newTestFollowService()

	repo.On("IsFollowing", uint(1), uint(2)).Return(false)

	result := svc.IsFollowing(1, 2)
	assert.False(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetFollowers / GetFollowing
// ============================================================

func TestGetFollowers_Success(t *testing.T) {
	svc, repo := newTestFollowService()

	followers := []model.User{{Name: "Follower1"}, {Name: "Follower2"}}
	repo.On("GetFollowers", uint(1)).Return(followers, nil)

	result, err := svc.GetFollowers(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestGetFollowing_Success(t *testing.T) {
	svc, repo := newTestFollowService()

	following := []model.User{{Name: "Following1"}}
	repo.On("GetFollowing", uint(1)).Return(following, nil)

	result, err := svc.GetFollowing(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}
