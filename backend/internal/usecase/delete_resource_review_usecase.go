package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// DeleteResourceReviewUseCase はレビューを削除する。所有者のみ削除できる。
type DeleteResourceReviewUseCase struct {
	reviews repository.ResourceReviewRepository
}

// NewDeleteResourceReviewUseCase は DeleteResourceReviewUseCase を生成する。
func NewDeleteResourceReviewUseCase(reviews repository.ResourceReviewRepository) *DeleteResourceReviewUseCase {
	return &DeleteResourceReviewUseCase{reviews: reviews}
}

// Execute は所有権を確認したうえでレビューを削除する。
func (uc *DeleteResourceReviewUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, func(r *model.ResourceReview) uint { return r.UserID }); err != nil {
		return err
	}
	return uc.reviews.Delete(ctx, id)
}
