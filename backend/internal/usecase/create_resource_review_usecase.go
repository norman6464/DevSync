package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// CreateResourceReviewUseCase は学習リソースへのレビューを新規作成する。
type CreateResourceReviewUseCase struct {
	reviews   repository.ResourceReviewRepository
	resources repository.LearningResourceReader
}

// NewCreateResourceReviewUseCase は CreateResourceReviewUseCase を生成する。
func NewCreateResourceReviewUseCase(
	reviews repository.ResourceReviewRepository,
	resources repository.LearningResourceReader,
) *CreateResourceReviewUseCase {
	return &CreateResourceReviewUseCase{reviews: reviews, resources: resources}
}

// Execute はリソースの存在確認・評価値/コメント検証・重複チェックを行いレビューを作成する。
func (uc *CreateResourceReviewUseCase) Execute(ctx context.Context, review *model.ResourceReview) error {
	if _, err := uc.resources.FindByID(ctx, review.ResourceID); err != nil {
		return domain.ErrNotFound
	}

	if err := domain.ValidateRating(review.Rating); err != nil {
		return err
	}

	review.Comment = strings.TrimSpace(review.Comment)
	if err := domain.ValidateStringLength(review.Comment, 0, 5000, "コメント"); err != nil {
		return err
	}

	existing, _ := uc.reviews.FindByUserAndResource(ctx, review.UserID, review.ResourceID)
	if existing != nil {
		return domain.ErrConflict
	}

	return uc.reviews.Create(ctx, review)
}
