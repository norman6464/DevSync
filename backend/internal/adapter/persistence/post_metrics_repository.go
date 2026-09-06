package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postMetricsRepository は [repository.PostMetricsRepository] の sqlc(pgx) 実装。
type postMetricsRepository struct {
	q *sqlcgen.Queries
}

// NewPostMetricsRepository は PostMetricsRepository の sqlc(pgx) 実装を返す。
func NewPostMetricsRepository(q *sqlcgen.Queries) repository.PostMetricsRepository {
	return &postMetricsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostMetricsRepository = (*postMetricsRepository)(nil)

// Reconcile はlikes/comments/post_viewsの実件数からpost_metrics全件を1文で補正する。
func (r *postMetricsRepository) Reconcile(ctx context.Context) error {
	return r.q.ReconcileAllPostMetrics(ctx)
}
