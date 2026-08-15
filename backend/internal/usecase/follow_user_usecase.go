package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// FollowUserUseCase は指定ユーザーをフォローする。
type FollowUserUseCase struct {
	follows       repository.FollowRepository
	notifications repository.NotificationCreator
}

// NewFollowUserUseCase は FollowUserUseCase を生成する。
func NewFollowUserUseCase(follows repository.FollowRepository, notifications repository.NotificationCreator) *FollowUserUseCase {
	return &FollowUserUseCase{follows: follows, notifications: notifications}
}

// Execute は followerID が followeeID をフォローし、相手に通知する。
// 自分自身のフォローは許可しない。通知の失敗はフォローの成否に影響させない。
func (uc *FollowUserUseCase) Execute(ctx context.Context, followerID, followeeID uint) error {
	if followerID == followeeID {
		return domain.ErrBadRequest
	}
	if err := uc.follows.Follow(ctx, followerID, followeeID); err != nil {
		return err
	}

	_ = uc.notifications.Create(ctx, &model.Notification{
		UserID:  followeeID,
		ActorID: followerID,
		Type:    model.NotificationTypeFollow,
	})
	return nil
}
