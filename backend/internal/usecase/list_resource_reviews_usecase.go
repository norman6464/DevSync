package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListResourceReviewsUseCase は指定リソースのレビュー一覧を取得する。
type ListResourceReviewsUseCase struct {
	reviews repository.ResourceReviewRepository
}

// NewListResourceReviewsUseCase は ListResourceReviewsUseCase を生成する。
func NewListResourceReviewsUseCase(reviews repository.ResourceReviewRepository) *ListResourceReviewsUseCase {
	return &ListResourceReviewsUseCase{reviews: reviews}
}

// Execute はリソースのレビュー一覧と総件数をページネーション付きで返す。
func (uc *ListResourceReviewsUseCase) Execute(ctx context.Context, resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	return uc.reviews.FindByResourceID(ctx, resourceID, limit, offset)
}
