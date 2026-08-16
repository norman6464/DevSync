package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateCheckinStreak(t *testing.T) {
	// 実装が time.Now() 基準で判定するため、テストも実行時点からの相対日付で組み立てる。
	day := func(offset int) string {
		return time.Now().AddDate(0, 0, offset).Format("2006-01-02")
	}

	t.Run("チェックインが無ければ 0", func(t *testing.T) {
		assert.Equal(t, 0, calculateCheckinStreak(nil))
	})

	t.Run("今日から連続 3 日", func(t *testing.T) {
		assert.Equal(t, 3, calculateCheckinStreak([]string{day(0), day(-1), day(-2)}))
	})

	t.Run("昨日始まりも連続として数える", func(t *testing.T) {
		assert.Equal(t, 2, calculateCheckinStreak([]string{day(-1), day(-2)}))
	})

	t.Run("最新が一昨日以前なら途切れたとみなして 0", func(t *testing.T) {
		assert.Equal(t, 0, calculateCheckinStreak([]string{day(-2), day(-3)}))
	})

	t.Run("途中に空きがあればそこで打ち切る", func(t *testing.T) {
		assert.Equal(t, 2, calculateCheckinStreak([]string{day(0), day(-1), day(-3)}))
	})
}
