package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postViewRepository は [repository.PostViewRepository] の sqlc(pgx) 実装。
type postViewRepository struct {
	q *sqlcgen.Queries
}

// NewPostViewRepository は PostViewRepository の sqlc(pgx) 実装を返す。
func NewPostViewRepository(pool *pgxpool.Pool) repository.PostViewRepository {
	return &postViewRepository{q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostViewRepository = (*postViewRepository)(nil)

// RecordViewIfAbsent は (user_id, post_id) のユニーク制約と ON CONFLICT DO NOTHING を用いて
// 未閲覧のときだけ閲覧を記録し、実際に挿入できた場合のみ post_metrics.view_count を加算する。
// 記録と加算は同一SQL文（queries/post_view.sqlのCreatePostView）で行うため、
// 並行リクエストによる二重記録・二重加算・重複エラーは起こらない。
func (r *postViewRepository) RecordViewIfAbsent(ctx context.Context, view *model.PostView) (bool, error) {
	rowsAffected, err := r.q.CreatePostView(ctx, sqlcgen.CreatePostViewParams{
		UserID: int64(view.UserID),
		PostID: int64(view.PostID),
	})
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

// GetViewCount は投稿の閲覧数を返す。
func (r *postViewRepository) GetViewCount(ctx context.Context, postID uint) (int64, error) {
	return r.q.CountPostViewsByPost(ctx, int64(postID))
}

// GetMostViewed は閲覧数上位の投稿を返す。
func (r *postViewRepository) GetMostViewed(ctx context.Context, limit int) ([]model.ViewCount, error) {
	rows, err := r.q.ListMostViewedPosts(ctx, int32Param(limit))
	if err != nil {
		return nil, err
	}
	results := make([]model.ViewCount, len(rows))
	for i, row := range rows {
		results[i] = model.ViewCount{PostID: uint(row.PostID), Count: int(row.ViewCount)}
	}
	return results, nil
}
