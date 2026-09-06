package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// roadmapMetricsRepository は [repository.RoadmapMetricsRepository] の sqlc(pgx) 実装。
type roadmapMetricsRepository struct {
	q *sqlcgen.Queries
}

// NewRoadmapMetricsRepository は RoadmapMetricsRepository の sqlc(pgx) 実装を返す。
func NewRoadmapMetricsRepository(q *sqlcgen.Queries) repository.RoadmapMetricsRepository {
	return &roadmapMetricsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RoadmapMetricsRepository = (*roadmapMetricsRepository)(nil)

// Reconcile はroadmap_stepsの実件数からroadmap_metrics全件を1文で補正する。
func (r *roadmapMetricsRepository) Reconcile(ctx context.Context) error {
	return r.q.ReconcileAllRoadmapMetrics(ctx)
}
