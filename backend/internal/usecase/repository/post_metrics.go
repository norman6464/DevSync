package repository

import "context"

// PostMetricsRepository は post_metrics（DEVSYNC-159でlike_count/comment_count/
// view_countを分離した側テーブル）に対する、usecase 側が要求する契約。
type PostMetricsRepository interface {
	// Reconcile はlikes/comments/post_viewsの実件数からpost_metrics全件を補正する。
	// CASCADE削除等でIncrement/Decrement系クエリを経由しない変化を吸収するための
	// 夜次バッチ処理として使う。
	Reconcile(ctx context.Context) error
}
