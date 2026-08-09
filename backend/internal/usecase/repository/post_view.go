package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostViewRepository は投稿閲覧数の永続化に対する、usecase 側が要求する契約。
type PostViewRepository interface {
	RecordView(ctx context.Context, view *model.PostView) error
	GetViewCount(ctx context.Context, postID uint) (int64, error)
	HasViewed(ctx context.Context, userID, postID uint) (bool, error)
	GetMostViewed(ctx context.Context, limit int) ([]model.ViewCount, error)
}
