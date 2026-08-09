package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// UpdateResourceReviewUseCase は既存のレビューを更新する。所有者のみ更新できる。
type UpdateResourceReviewUseCase struct {
	reviews repository.ResourceReviewRepository
}

// NewUpdateResourceReviewUseCase は UpdateResourceReviewUseCase を生成する。
func NewUpdateResourceReviewUseCase(reviews repository.ResourceReviewRepository) *UpdateResourceReviewUseCase {
	return &UpdateResourceReviewUseCase{reviews: reviews}
}

// Execute は所有権を確認し、rating（0 は変更なし）・comment（空は変更なし）を更新する。
func (uc *UpdateResourceReviewUseCase) Execute(ctx context.Context, id, userID uint, rating int, comment string) (*model.ResourceReview, error) {
	review, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, func(r *model.ResourceReview) uint { return r.UserID })
	if err != nil {
		return nil, err
	}

	if rating != 0 {
		if err := domain.ValidateRating(rating); err != nil {
			return nil, err
		}
		review.Rating = rating
	}
	if c := strings.TrimSpace(comment); c != "" {
		if err := domain.ValidateStringLength(c, 1, 5000, "コメント"); err != nil {
			return nil, err
		}
		review.Comment = c
	}

	if err := uc.reviews.Update(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}
