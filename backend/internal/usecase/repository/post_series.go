package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostSeriesRepository は投稿シリーズの永続化に対する、usecase 側が要求する契約。
type PostSeriesRepository interface {
	Create(ctx context.Context, series *model.PostSeries) error
	FindByID(ctx context.Context, id uint) (*model.PostSeries, error)
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]model.PostSeries, error)
	CountByUser(ctx context.Context, userID uint) (int64, error)
	Update(ctx context.Context, series *model.PostSeries) error
	Delete(ctx context.Context, id uint) error
	AddPost(ctx context.Context, item *model.PostSeriesItem) error
	RemovePost(ctx context.Context, seriesID, postID uint) error
	HasPost(ctx context.Context, seriesID, postID uint) (bool, error)
	GetPostsBySeriesID(ctx context.Context, seriesID uint) ([]model.PostSeriesItem, error)
}
