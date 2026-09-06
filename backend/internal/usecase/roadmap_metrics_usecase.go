package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ReconcileRoadmapMetricsUseCase はroadmap_metrics（step_count/completed_step_count）を
// 実件数と突き合わせて補正する。ReconcileAllMetricsUseCase からドメイン別reconcilerの
// 1つとして呼ばれる想定。
type ReconcileRoadmapMetricsUseCase struct {
	metrics repository.RoadmapMetricsRepository
}

// NewReconcileRoadmapMetricsUseCase は ReconcileRoadmapMetricsUseCase を生成する。
func NewReconcileRoadmapMetricsUseCase(metrics repository.RoadmapMetricsRepository) *ReconcileRoadmapMetricsUseCase {
	return &ReconcileRoadmapMetricsUseCase{metrics: metrics}
}

// Execute はroadmap_metrics全件を補正する。
func (uc *ReconcileRoadmapMetricsUseCase) Execute(ctx context.Context) error {
	return uc.metrics.Reconcile(ctx)
}
