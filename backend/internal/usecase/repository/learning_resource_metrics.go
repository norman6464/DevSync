package repository

import "context"

// LearningResourceMetricsRepository は learning_resource_metrics（DEVSYNC-159で
// like_count/save_countを分離した側テーブル）に対する、usecase 側が要求する契約。
type LearningResourceMetricsRepository interface {
	// Reconcile はresource_likes/resource_savesの実件数からlearning_resource_metrics
	// 全件を補正する。CASCADE削除等でIncrement/Decrement系クエリを経由しない変化を
	// 吸収するための夜次バッチ処理として使う。
	Reconcile(ctx context.Context) error
}
