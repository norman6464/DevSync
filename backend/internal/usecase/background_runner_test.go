package usecase

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackgroundRunner(t *testing.T) {
	t.Run("同時実行数がワーカー数を超えない", func(t *testing.T) {
		r := NewBackgroundRunner(2, 64, time.Second)

		var current, max int32
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			r.Submit("job", func(ctx context.Context) error {
				defer wg.Done()
				n := atomic.AddInt32(&current, 1)
				for {
					m := atomic.LoadInt32(&max)
					if n <= m || atomic.CompareAndSwapInt32(&max, m, n) {
						break
					}
				}
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&current, -1)
				return nil
			})
		}
		wg.Wait()

		assert.LessOrEqual(t, atomic.LoadInt32(&max), int32(2))
	})

	t.Run("ジョブに期限付き ctx が渡り、超過で打ち切られる", func(t *testing.T) {
		r := NewBackgroundRunner(1, 4, 20*time.Millisecond)

		done := make(chan error, 1)
		r.Submit("job", func(ctx context.Context) error {
			<-ctx.Done()
			done <- ctx.Err()
			return ctx.Err()
		})

		select {
		case err := <-done:
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		case <-time.After(time.Second):
			t.Fatal("期限が伝播していない（ジョブが打ち切られなかった）")
		}
	})

	t.Run("キューが満杯のときは呼び出し元をブロックせず破棄する", func(t *testing.T) {
		r := NewBackgroundRunner(1, 1, time.Second)

		block := make(chan struct{})
		executed := make(chan string, 8)
		// ワーカーを塞ぐ
		r.Submit("blocker", func(ctx context.Context) error {
			executed <- "blocker"
			<-block
			return nil
		})
		// blocker がワーカーを掴んだのを確認してからキューを埋める
		assert.Equal(t, "blocker", <-executed)
		r.Submit("queued", func(ctx context.Context) error {
			executed <- "queued"
			return nil
		})

		// キュー満杯での Submit がブロックせずに戻ること（破棄される）
		returned := make(chan struct{})
		go func() {
			r.Submit("dropped", func(ctx context.Context) error {
				executed <- "dropped"
				return nil
			})
			close(returned)
		}()
		select {
		case <-returned:
		case <-time.After(time.Second):
			t.Fatal("キュー満杯の Submit がブロックした")
		}

		close(block)
		assert.Equal(t, "queued", <-executed)
		select {
		case name := <-executed:
			t.Fatalf("破棄されたはずのジョブが実行された: %s", name)
		case <-time.After(50 * time.Millisecond):
		}
	})
}
