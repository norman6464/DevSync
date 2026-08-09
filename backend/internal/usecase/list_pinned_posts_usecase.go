package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListPinnedPostsUseCase は指定ユーザーのピン留め投稿一覧を取得する。
type ListPinnedPostsUseCase struct {
	pins repository.PostPinRepository
}

// NewListPinnedPostsUseCase は ListPinnedPostsUseCase を生成する。
func NewListPinnedPostsUseCase(pins repository.PostPinRepository) *ListPinnedPostsUseCase {
	return &ListPinnedPostsUseCase{pins: pins}
}

// Execute はユーザーのピン留め投稿一覧を pin_order 昇順で返す。
func (uc *ListPinnedPostsUseCase) Execute(ctx context.Context, userID uint) ([]model.PostPin, error) {
	return uc.pins.GetByUserID(ctx, userID)
}
