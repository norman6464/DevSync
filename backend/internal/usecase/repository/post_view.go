package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostViewRepository は投稿閲覧数の永続化に対する、usecase 側が要求する契約。
type PostViewRepository interface {
	// RecordViewIfAbsent は未閲覧なら閲覧を記録し view_count を加算する原子的操作。
	// 実際に記録したとき true を返す（既に閲覧済みなら false, nil）。
	RecordViewIfAbsent(ctx context.Context, view *model.PostView) (bool, error)
	GetViewCount(ctx context.Context, postID uint) (int64, error)
	GetMostViewed(ctx context.Context, limit int) ([]model.ViewCount, error)
}
