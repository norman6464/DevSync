package service

import (
	"log"

	"github.com/robfig/cron/v3"
)

// Scheduler はcronベースの定期実行サービス。
// ウィークリーレポートメールなどの定期タスクをスケジュール実行する。
type Scheduler struct {
	cron     *cron.Cron
	emailSvc *WeeklyReportEmailService
}

// NewScheduler は新しいSchedulerインスタンスを生成する。
func NewScheduler(emailSvc *WeeklyReportEmailService) *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		emailSvc: emailSvc,
	}
}

// Start はスケジューラを起動し、定期タスクを登録する。
// 毎週月曜日 9:00（サーバーローカル時間）にウィークリーレポートメールを送信する。
func (s *Scheduler) Start() {
	// 毎週月曜日 9:00 に実行
	_, err := s.cron.AddFunc("0 9 * * 1", func() {
		log.Println("スケジューラ: ウィークリーレポートメール送信開始")
		if err := s.emailSvc.SendAllWeeklyReports(); err != nil {
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
