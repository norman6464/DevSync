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
	// 不在（nil）と DB 障害（err）を区別する。障害を 404 に変換すると原因が隠れる。
	resource, err := uc.resources.FindByID(ctx, review.ResourceID)
	if err != nil {
		return err
	}
	if resource == nil {
		return domain.ErrNotFound
	}

	if err := domain.ValidateRating(review.Rating); err != nil {
		return err
	}

	review.Comment = strings.TrimSpace(review.Comment)
	if err := domain.ValidateStringLength(review.Comment, 0, 5000, "コメント"); err != nil {
		return err
	}

	// 重複チェックの DB 障害を握り潰すと、障害時に重複作成を許してしまう。
	existing, err := uc.reviews.FindByUserAndResource(ctx, review.UserID, review.ResourceID)
	if err != nil {
		return err
	}
	if existing != nil {
		return domain.ErrConflict
	}

	return uc.reviews.Create(ctx, review)
}
