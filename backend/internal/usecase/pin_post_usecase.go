package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// maxPinsPerUser は 1 ユーザーがピン留めできる投稿の上限。
const maxPinsPerUser = 3

// PinPostUseCase は投稿をプロフィールにピン留めする。
type PinPostUseCase struct {
	pins  repository.PostPinRepository
	posts repository.PostReader
}

// NewPinPostUseCase は PinPostUseCase を生成する。
func NewPinPostUseCase(pins repository.PostPinRepository, posts repository.PostReader) *PinPostUseCase {
	return &PinPostUseCase{pins: pins, posts: posts}
}

// Execute は自分の投稿を最大 maxPinsPerUser 件までピン留めする。
func (uc *PinPostUseCase) Execute(ctx context.Context, userID, postID uint) error {
	post, err := uc.posts.FindByID(ctx, postID)
	if err != nil {
		return domain.ErrNotFound
	}
	if post.UserID != userID {
		return domain.NewError(domain.ErrCodeForbidden, "自分の投稿のみピン留めできます", nil)
	}

	count, err := uc.pins.CountByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if count >= maxPinsPerUser {
		return domain.NewError(domain.ErrCodeBadRequest, "ピン留めは最大3件までです", nil)
	}

	return uc.pins.Pin(ctx, &model.PostPin{
		UserID:   userID,
		PostID:   postID,
		PinOrder: int(count),
	})
}
