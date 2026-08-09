package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetPostViewCountUseCase は投稿のユニーク閲覧数を取得する。
type GetPostViewCountUseCase struct {
	views repository.PostViewRepository
}

// NewGetPostViewCountUseCase は GetPostViewCountUseCase を生成する。
func NewGetPostViewCountUseCase(views repository.PostViewRepository) *GetPostViewCountUseCase {
	return &GetPostViewCountUseCase{views: views}
}

// Execute は投稿のユニーク閲覧数を返す。
func (uc *GetPostViewCountUseCase) Execute(ctx context.Context, postID uint) (int64, error) {
	if err := domain.ValidateRequiredID(postID, "postID"); err != nil {
		return 0, err
	}
	return uc.views.GetViewCount(ctx, postID)
}
