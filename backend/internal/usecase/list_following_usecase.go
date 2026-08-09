package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListFollowingUseCase は指定ユーザーがフォロー中のユーザー一覧をページネーション付きで取得する。
type ListFollowingUseCase struct {
	follows repository.FollowRepository
}

// NewListFollowingUseCase は ListFollowingUseCase を生成する。
func NewListFollowingUseCase(follows repository.FollowRepository) *ListFollowingUseCase {
	return &ListFollowingUseCase{follows: follows}
}

// Execute は userID がフォロー中のユーザーを取得する。戻り値は (ユーザー一覧, 総件数, error)。
func (uc *ListFollowingUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	return uc.follows.GetFollowing(ctx, userID, limit, offset)
}
