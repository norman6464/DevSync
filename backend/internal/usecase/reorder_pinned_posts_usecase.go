package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ReorderPinnedPostsUseCase はピン留め投稿の表示順序を変更する。
type ReorderPinnedPostsUseCase struct {
	pins repository.PostPinRepository
}

// NewReorderPinnedPostsUseCase は ReorderPinnedPostsUseCase を生成する。
func NewReorderPinnedPostsUseCase(pins repository.PostPinRepository) *ReorderPinnedPostsUseCase {
	return &ReorderPinnedPostsUseCase{pins: pins}
}

// Execute は postIDs の順にピン留めを並べ替える。
// postIDs はすべて userID のピン留め済み投稿でなければならない（上限・所有権を検証）。
func (uc *ReorderPinnedPostsUseCase) Execute(ctx context.Context, userID uint, postIDs []uint) error {
	if len(postIDs) > maxPinsPerUser {
		return domain.NewError(domain.ErrCodeBadRequest, "ピン留めは最大3件までです", nil)
	}

	pins, err := uc.pins.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	pinnedSet := make(map[uint]bool, len(pins))
	for _, pin := range pins {
		pinnedSet[pin.PostID] = true
	}
	for _, postID := range postIDs {
		if !pinnedSet[postID] {
			return domain.NewError(domain.ErrCodeForbidden, "自分のピン留め投稿のみ順序変更できます", nil)
		}
	}

	return uc.pins.UpdateOrder(ctx, userID, postIDs)
}
