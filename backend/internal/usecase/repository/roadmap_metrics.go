package repository

import "context"

// RoadmapMetricsRepository は roadmap_metrics（DEVSYNC-159でstep_count/
// completed_step_countを分離した側テーブル）に対する、usecase 側が要求する契約。
type RoadmapMetricsRepository interface {
	// Reconcile はroadmap_stepsの実件数からroadmap_metrics全件を補正する。
	// CASCADE削除等でCTEベースの加減算を経由しない変化を吸収するための
	// 夜次バッチ処理として使う。
	Reconcile(ctx context.Context) error
}
