package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ReconcileLearningResourceMetricsUseCase はlearning_resource_metrics（like_count/
// save_count）を実件数と突き合わせて補正する。ReconcileAllMetricsUseCase から
// ドメイン別reconcilerの1つとして呼ばれる想定。
type ReconcileLearningResourceMetricsUseCase struct {
	metrics repository.LearningResourceMetricsRepository
}

// NewReconcileLearningResourceMetricsUseCase は ReconcileLearningResourceMetricsUseCase を生成する。
func NewReconcileLearningResourceMetricsUseCase(metrics repository.LearningResourceMetricsRepository) *ReconcileLearningResourceMetricsUseCase {
	return &ReconcileLearningResourceMetricsUseCase{metrics: metrics}
}

// Execute はlearning_resource_metrics全件を補正する。
func (uc *ReconcileLearningResourceMetricsUseCase) Execute(ctx context.Context) error {
	return uc.metrics.Reconcile(ctx)
}
