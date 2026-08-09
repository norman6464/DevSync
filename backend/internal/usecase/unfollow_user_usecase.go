package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// UnfollowUserUseCase は指定ユーザーのフォローを解除する。
type UnfollowUserUseCase struct {
	follows repository.FollowRepository
}

// NewUnfollowUserUseCase は UnfollowUserUseCase を生成する。
func NewUnfollowUserUseCase(follows repository.FollowRepository) *UnfollowUserUseCase {
	return &UnfollowUserUseCase{follows: follows}
}

// Execute は followerID が followeeID のフォローを解除する。自分自身は許可しない。
func (uc *UnfollowUserUseCase) Execute(ctx context.Context, followerID, followeeID uint) error {
	if followerID == followeeID {
		return domain.ErrBadRequest
	}
	return uc.follows.Unfollow(ctx, followerID, followeeID)
}
