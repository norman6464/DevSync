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

// Scheduler は cron ベースの定期実行。
// ウィークリーレポートメールなどの定期タスクをスケジュール実行する。
type Scheduler struct {
	cron     CronScheduler
	emailSvc WeeklyReportSender
}

// New は新しい Scheduler インスタンスを生成する。
func New(weeklyReport WeeklyReportSender) *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		emailSvc: weeklyReport,
	}
}

// Start はスケジューラを起動し、定期タスクを登録する。
// 毎週月曜日 9:00（サーバーローカル時間）にウィークリーレポートメールを送信する。
func (s *Scheduler) Start() {
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
		log.Printf("スケジューラ: cronジョブ登録失敗: %v", err)
		return
	}

	s.cron.Start()
	log.Println("スケジューラ: 起動完了（毎週月曜 9:00 にウィークリーレポートメール送信）")
}

// Stop はスケジューラを停止する。
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("スケジューラ: 停止完了")
}
