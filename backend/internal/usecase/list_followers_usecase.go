package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListFollowersUseCase は指定ユーザーのフォロワー一覧をページネーション付きで取得する。
type ListFollowersUseCase struct {
	follows repository.FollowRepository
}

// NewListFollowersUseCase は ListFollowersUseCase を生成する。
func NewListFollowersUseCase(follows repository.FollowRepository) *ListFollowersUseCase {
	return &ListFollowersUseCase{follows: follows}
}

// Execute は userID のフォロワーを取得する。戻り値は (ユーザー一覧, 総件数, error)。
func (uc *ListFollowersUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	return uc.follows.GetFollowers(ctx, userID, limit, offset)
}
