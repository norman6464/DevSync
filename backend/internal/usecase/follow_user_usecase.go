package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// FollowUserUseCase は指定ユーザーをフォローする。
type FollowUserUseCase struct {
	follows repository.FollowRepository
}

// NewFollowUserUseCase は FollowUserUseCase を生成する。
func NewFollowUserUseCase(follows repository.FollowRepository) *FollowUserUseCase {
	return &FollowUserUseCase{follows: follows}
}

// Execute は followerID が followeeID をフォローする。自分自身のフォローは許可しない。
func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID, followeeID uint) error {
	if followerID == followeeID {
		return domain.ErrBadRequest
	}
	return uc.follows.Follow(ctx, followerID, followeeID)
}
