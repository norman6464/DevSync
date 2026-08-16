package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeToCalendarDay(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)

	t.Run("UTC の時刻はその暦日の 00:00 になる", func(t *testing.T) {
		got := normalizeToCalendarDay(time.Date(2026, 8, 16, 23, 59, 59, 0, time.UTC))
		assert.Equal(t, time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC), got)
	})

	t.Run("JST の時刻はロケーションの暦日で切り捨てる（Truncate と違い UTC 境界に丸めない）", func(t *testing.T) {
		// JST 8/17 08:00 = UTC 8/16 23:00。Truncate(24h) だと 8/16 になるが、暦日は 8/17。
		v := time.Date(2026, 8, 17, 8, 0, 0, 0, jst)
		got := normalizeToCalendarDay(v)
		assert.Equal(t, time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC), got)
	})

	t.Run("DB の DATE() が返す 00:00 は同じ暦日に写る（セッションのタイムゾーンに依存しない）", func(t *testing.T) {
		utcMidnight := time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC)
		jstMidnight := time.Date(2026, 8, 16, 0, 0, 0, 0, jst)
		assert.Equal(t, normalizeToCalendarDay(utcMidnight), normalizeToCalendarDay(jstMidnight))
	})
}

func TestIsNextCalendarDay(t *testing.T) {
	day := func(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

	t.Run("翌日は true・同日と 2 日後は false", func(t *testing.T) {
		assert.True(t, isNextCalendarDay(day(2026, 8, 17), day(2026, 8, 16)))
		assert.False(t, isNextCalendarDay(day(2026, 8, 16), day(2026, 8, 16)))
		assert.False(t, isNextCalendarDay(day(2026, 8, 18), day(2026, 8, 16)))
	})

	t.Run("月またぎ・年またぎも暦日で判定する", func(t *testing.T) {
		assert.True(t, isNextCalendarDay(day(2026, 9, 1), day(2026, 8, 31)))
		assert.True(t, isNextCalendarDay(day(2027, 1, 1), day(2026, 12, 31)))
	})

	t.Run("サマータイム切替日（23 時間の日）でも連続と判定する", func(t *testing.T) {
		// 米国太平洋時間の 2026-03-08 は DST 開始で 23 時間しかない。
		// 経過時間の割り算（int(hours/24) == 1）だと 0 になり連続が切れていた。
		loc, err := time.LoadLocation("America/Los_Angeles")
		require.NoError(t, err)
		before := normalizeToCalendarDay(time.Date(2026, 3, 7, 12, 0, 0, 0, loc))
		after := normalizeToCalendarDay(time.Date(2026, 3, 8, 12, 0, 0, 0, loc))
		assert.True(t, isNextCalendarDay(after, before))
	})
}

func TestIsTodayOrYesterday(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 8, d, 0, 0, 0, 0, time.UTC) }
	today := day(16)

	assert.True(t, isTodayOrYesterday(day(16), today))
	assert.True(t, isTodayOrYesterday(day(15), today))
	assert.False(t, isTodayOrYesterday(day(14), today))
	assert.False(t, isTodayOrYesterday(day(17), today))
}
