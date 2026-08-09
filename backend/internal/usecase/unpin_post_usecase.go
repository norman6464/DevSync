package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// UnpinPostUseCase は投稿のピン留めを解除する。
type UnpinPostUseCase struct {
	pins repository.PostPinRepository
}

// NewUnpinPostUseCase は UnpinPostUseCase を生成する。
func NewUnpinPostUseCase(pins repository.PostPinRepository) *UnpinPostUseCase {
	return &UnpinPostUseCase{pins: pins}
}

// Execute はピン留めを解除する。未ピン留めなら ErrNotFound。
func (uc *UnpinPostUseCase) Execute(ctx context.Context, userID, postID uint) error {
	pinned, err := uc.pins.IsPinned(ctx, userID, postID)
	if err != nil {
		return err
	}
	if !pinned {
		return domain.NewError(domain.ErrCodeNotFound, "ピン留めされていません", nil)
	}
	return uc.pins.Unpin(ctx, userID, postID)
}
