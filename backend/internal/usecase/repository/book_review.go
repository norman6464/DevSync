package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// BookReviewRepository は書籍レビューの永続化に対する、usecase 側が要求する契約。
type BookReviewRepository interface {
	Create(ctx context.Context, review *model.BookReview) error
	// FindByID は指定 ID のレビューを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.BookReview, error)
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.BookReview, int64, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.BookReview, int64, error)
	FindByRating(ctx context.Context, userID uint, minRating, maxRating int) ([]model.BookReview, error)
	Search(ctx context.Context, query string, limit, offset int) ([]model.BookReview, int64, error)
	Update(ctx context.Context, review *model.BookReview) error
	Delete(ctx context.Context, id uint) error
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}
