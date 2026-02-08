package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// FollowService handles follow business logic.
type FollowService struct {
	repo repository.FollowRepositoryInterface
}

// NewFollowService creates a new FollowService.
func NewFollowService(repo repository.FollowRepositoryInterface) *FollowService {
	return &FollowService{repo: repo}
}

// Follow follows a user.
func (s *FollowService) Follow(followerID, followeeID uint) error {
	if followerID == followeeID {
		return ErrBadRequest
	}
	return s.repo.Follow(followerID, followeeID)
}

// Unfollow unfollows a user.
func (s *FollowService) Unfollow(followerID, followeeID uint) error {
	return s.repo.Unfollow(followerID, followeeID)
}

// IsFollowing checks if a user is following another user.
func (s *FollowService) IsFollowing(followerID, followeeID uint) bool {
	return s.repo.IsFollowing(followerID, followeeID)
}

// GetFollowers returns all followers of a user.
func (s *FollowService) GetFollowers(userID uint) ([]model.User, error) {
	return s.repo.GetFollowers(userID)
}

// GetFollowing returns all users a user is following.
func (s *FollowService) GetFollowing(userID uint) ([]model.User, error) {
	return s.repo.GetFollowing(userID)
}
