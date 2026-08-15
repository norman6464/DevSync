package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ResourceReviewRepository は学習リソースレビューの永続化に対する、usecase 側が要求する契約。
type ResourceReviewRepository interface {
	Create(ctx context.Context, review *model.ResourceReview) error
	// FindByID は指定 ID のレビューを返す。不在の場合は (nil, nil) を返す。
	FindByID(ctx context.Context, id uint) (*model.ResourceReview, error)
	FindByResourceID(ctx context.Context, resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error)
	// FindByUserAndResource は指定ユーザーの指定リソースへのレビューを返す。不在の場合は (nil, nil) を返す。
	FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceReview, error)
	Update(ctx context.Context, review *model.ResourceReview) error
	Delete(ctx context.Context, id uint) error
}

// LearningResourceReader は存在確認に必要なリソース読み取りだけを切り出した最小 port（-er）。
type LearningResourceReader interface {
	// FindByID は指定 ID の学習リソースを返す。不在の場合は (nil, nil) を返す。
	FindByID(ctx context.Context, id uint) (*model.LearningResource, error)
}
