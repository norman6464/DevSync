package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// CountPinnedPostsUseCase は指定ユーザーのピン留め投稿数を返す。
type CountPinnedPostsUseCase struct {
	pins repository.PostPinRepository
}

// NewCountPinnedPostsUseCase は CountPinnedPostsUseCase を生成する。
func NewCountPinnedPostsUseCase(pins repository.PostPinRepository) *CountPinnedPostsUseCase {
	return &CountPinnedPostsUseCase{pins: pins}
}

// Execute はユーザーのピン留め投稿数を返す。
func (uc *CountPinnedPostsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.pins.CountByUserID(ctx, userID)
}
