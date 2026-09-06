package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ReconcilePostMetricsUseCase はpost_metrics（like_count/comment_count/view_count）を
// 実件数と突き合わせて補正する。夜次バッチのスケジューラから呼ばれる想定
// （infra/scheduler.MetricsReconcilerを満たす）。
type ReconcilePostMetricsUseCase struct {
	metrics repository.PostMetricsRepository
}

// NewReconcilePostMetricsUseCase は ReconcilePostMetricsUseCase を生成する。
func NewReconcilePostMetricsUseCase(metrics repository.PostMetricsRepository) *ReconcilePostMetricsUseCase {
	return &ReconcilePostMetricsUseCase{metrics: metrics}
}

// Execute はpost_metrics全件を補正する。
func (uc *ReconcilePostMetricsUseCase) Execute(ctx context.Context) error {
	return uc.metrics.Reconcile(ctx)
}
