package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// FollowService はフォロー・フォロワー関係のビジネスロジックを提供する。
type FollowService struct {
	repo repository.FollowRepositoryInterface
}

// NewFollowService は新しいFollowServiceインスタンスを生成する。
func NewFollowService(repo repository.FollowRepositoryInterface) *FollowService {
	return &FollowService{repo: repo}
}

// Follow は指定ユーザーをフォローする。自分自身のフォローは許可しない。
func (s *FollowService) Follow(followerID, followeeID uint) error {
	if followerID == followeeID {
		return ErrBadRequest
	}
	return s.repo.Follow(followerID, followeeID)
}

// Unfollow は指定ユーザーのフォローを解除する。自分自身のアンフォローは許可しない。
func (s *FollowService) Unfollow(followerID, followeeID uint) error {
	if followerID == followeeID {
		return ErrBadRequest
	}
	return s.repo.Unfollow(followerID, followeeID)
}

// IsFollowing は指定ユーザーをフォロー中かどうかを判定する。
func (s *FollowService) IsFollowing(followerID, followeeID uint) bool {
	return s.repo.IsFollowing(followerID, followeeID)
}

// GetFollowers は指定ユーザーの全フォロワーを取得する。
func (s *FollowService) GetFollowers(userID uint) ([]model.User, error) {
	return s.repo.GetFollowers(userID)
}

// GetFollowing は指定ユーザーがフォロー中の全ユーザーを取得する。
func (s *FollowService) GetFollowing(userID uint) ([]model.User, error) {
	return s.repo.GetFollowing(userID)
}
