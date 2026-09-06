package usecase

import (
	"context"
	"errors"
)

// metricsReconciler は各ドメインのreconcile usecase（ReconcilePostMetricsUseCase等）が
// 共通して満たす最小契約。ReconcileAllMetricsUseCase はこれの一覧を受け取るだけでよく、
// ドメインが増えるたびにこの型自体を変更する必要はない。
type metricsReconciler interface {
	Execute(ctx context.Context) error
}

// ReconcileAllMetricsUseCase は複数ドメインのメトリクスreconcilerを1つにまとめ、
// スケジューラからの登録ポイントを1つに保つためのオーケストレータ。
// ドメインが増えても scheduler.New への配線（cronエントリ）を増やさずに済む。
type ReconcileAllMetricsUseCase struct {
	reconcilers []metricsReconciler
}

// NewReconcileAllMetricsUseCase は ReconcileAllMetricsUseCase を生成する。
func NewReconcileAllMetricsUseCase(reconcilers ...metricsReconciler) *ReconcileAllMetricsUseCase {
	return &ReconcileAllMetricsUseCase{reconcilers: reconcilers}
}

// Execute は登録された全ドメインのreconcileを順に実行する。1つが失敗しても残りは
// 実行し続け、発生したエラーはすべてまとめて返す（1ドメインの失敗が他ドメインの
// reconcileを妨げないようにするため）。
func (uc *ReconcileAllMetricsUseCase) Execute(ctx context.Context) error {
	var errs []error
	for _, reconciler := range uc.reconcilers {
		if err := reconciler.Execute(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
