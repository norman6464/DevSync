// Package scheduler は cron による定期実行を担う。
// 何を実行するかは呼び出し側から注入し、ここでは実行の仕組みだけを持つ。
package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

// CronScheduler はcronスケジューラの抽象インターフェース。
type CronScheduler interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	Start()
	Stop() context.Context
}

// WeeklyReportSender はウィークリーレポート送信の抽象インターフェース。
// 定期実行の入り口であり、リクエスト ctx を持たないためここが ctx の起点になる。
type WeeklyReportSender interface {
	Execute(ctx context.Context) error
}

// MetricsReconciler はカウンタのreconcile（実件数との乖離補正）の抽象インターフェース。
// 定期実行の入り口であり、リクエスト ctx を持たないためここが ctx の起点になる。
type MetricsReconciler interface {
	Execute(ctx context.Context) error
}

// Scheduler は cron ベースの定期実行。
// ウィークリーレポートメール・カウンタのreconcileなどの定期タスクをスケジュール実行する。
type Scheduler struct {
	cron              CronScheduler
	emailSvc          WeeklyReportSender
	metricsReconciler MetricsReconciler
}

// New は新しい Scheduler インスタンスを生成する。weeklyReport は SMTP 未設定時など
// ウィークリーレポート機能自体を無効化したい場合に nil を渡せる（その場合はcronへ
// 登録しない）。metricsReconciler は SMTP 設定の有無に関わらず常に登録する。
func New(weeklyReport WeeklyReportSender, metricsReconciler MetricsReconciler) *Scheduler {
	return &Scheduler{
		cron:              cron.New(),
		emailSvc:          weeklyReport,
		metricsReconciler: metricsReconciler,
	}
}

// Start はスケジューラを起動し、定期タスクを登録する。
// 毎週月曜日 9:00（サーバーローカル時間）にウィークリーレポートメールを送信する
// （emailSvcがnilなら登録しない）。毎日 3:00 に投稿カウンタのreconcileを行う。
func (s *Scheduler) Start() {
	if s.emailSvc != nil {
		// 毎週月曜日 9:00 に実行
		_, err := s.cron.AddFunc("0 9 * * 1", func() {
			log.Println("スケジューラ: ウィークリーレポートメール送信開始")
			if err := s.emailSvc.Execute(context.Background()); err != nil {
				log.Printf("スケジューラ: ウィークリーレポートメール送信エラー: %v", err)
			} else {
				log.Println("スケジューラ: ウィークリーレポートメール送信完了")
			}
		})
		if err != nil {
			log.Printf("スケジューラ: ウィークリーレポートcronジョブ登録失敗: %v", err)
		}
	}

	if s.metricsReconciler != nil {
		// 毎日 3:00 に実行（アクセスが少ない深夜帯）
		_, err := s.cron.AddFunc("0 3 * * *", func() {
			log.Println("スケジューラ: 投稿カウンタreconcile開始")
			if err := s.metricsReconciler.Execute(context.Background()); err != nil {
				log.Printf("スケジューラ: 投稿カウンタreconcileエラー: %v", err)
			} else {
				log.Println("スケジューラ: 投稿カウンタreconcile完了")
			}
		})
		if err != nil {
			log.Printf("スケジューラ: reconcileジョブ登録失敗: %v", err)
		}
	}

	s.cron.Start()
	log.Println("スケジューラ: 起動完了")
}

// Stop はスケジューラを停止する。
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("スケジューラ: 停止完了")
}
