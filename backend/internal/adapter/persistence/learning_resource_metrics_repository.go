package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningResourceMetricsRepository は [repository.LearningResourceMetricsRepository] の
// sqlc(pgx) 実装。
type learningResourceMetricsRepository struct {
	q *sqlcgen.Queries
}

// NewLearningResourceMetricsRepository は LearningResourceMetricsRepository の
// sqlc(pgx) 実装を返す。
func NewLearningResourceMetricsRepository(q *sqlcgen.Queries) repository.LearningResourceMetricsRepository {
	return &learningResourceMetricsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningResourceMetricsRepository = (*learningResourceMetricsRepository)(nil)

// Reconcile はresource_likes/resource_savesの実件数からlearning_resource_metrics
// 全件を1文で補正する。
func (r *learningResourceMetricsRepository) Reconcile(ctx context.Context) error {
	return r.q.ReconcileAllLearningResourceMetrics(ctx)
}
