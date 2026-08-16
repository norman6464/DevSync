package usecase

import (
	"context"
	"log"
	"sync"
	"time"
)

// バックグラウンド通知の既定値。通知は best-effort のため、処理能力を超えた分は
// 破棄してリクエスト処理側の資源（goroutine・DB 接続）を守る。
const (
	backgroundWorkers    = 4
	backgroundQueueSize  = 256
	backgroundJobTimeout = 30 * time.Second
)

// backgroundJob は上限付きワーカーで実行する 1 件のジョブ。
type backgroundJob struct {
	name string
	run  func(ctx context.Context) error
}

// BackgroundRunner は同時実行数と 1 件あたりの実行期限を制限してジョブを実行する。
// 呼び出しごとに goroutine を起動する方式と違い、ワーカー数以上の並行実行は起きず、
// ジョブには期限付き ctx が渡るため DB 操作が積み上がらない。
// キューが満杯のときはジョブを破棄してログへ残す（通知向けの best-effort 設計）。
type BackgroundRunner struct {
	queue   chan backgroundJob
	timeout time.Duration
}

// NewBackgroundRunner は workers 本のワーカーと queueSize のキューでランナーを起動する。
func NewBackgroundRunner(workers, queueSize int, timeout time.Duration) *BackgroundRunner {
	r := &BackgroundRunner{queue: make(chan backgroundJob, queueSize), timeout: timeout}
	for i := 0; i < workers; i++ {
		go r.worker()
	}
	return r
}

func (r *BackgroundRunner) worker() {
	for job := range r.queue {
		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		if err := job.run(ctx); err != nil {
			log.Printf("[WARN] バックグラウンドジョブ %s が失敗: %v", job.name, err)
		}
		cancel()
	}
}

// Submit はジョブを投入する。キューが満杯なら待たずに破棄してログへ残す。
// 呼び出し元（リクエスト処理）をブロックしないことを優先する。
func (r *BackgroundRunner) Submit(name string, run func(ctx context.Context) error) {
	select {
	case r.queue <- backgroundJob{name: name, run: run}:
	default:
		log.Printf("[WARN] バックグラウンドジョブ %s を破棄: キューが満杯", name)
	}
}

var (
	defaultRunnerOnce sync.Once
	defaultRunner     *BackgroundRunner
)

// defaultBackgroundRunner は通知系ジョブが共有する既定のランナーを返す。
// 横断的な資源制限のため 1 プロセスに 1 つでよく、コンストラクタ注入にすると
// 全 usecase の生成シグネチャへ波及するためパッケージ共有にしている。
func defaultBackgroundRunner() *BackgroundRunner {
	defaultRunnerOnce.Do(func() {
		defaultRunner = NewBackgroundRunner(backgroundWorkers, backgroundQueueSize, backgroundJobTimeout)
	})
	return defaultRunner
}
